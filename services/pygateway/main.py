import signal
import os
import serial
import time
import paho.mqtt.client as mqtt
import subprocess

# global var for the serial connection
serialConn = None

FIRMWARE_LOCATION = "./light.hex"
flashing = False

exiting = False

debug = False
# If true, just print, don't handle serial
justPrint = False

mqtt_connected = False


def debugf(msg, *args):
    if debug:
        print(msg, *args)

#####################
### MQTT HANDLERS ###
#####################


def on_connect(client, userdata, flags, rc):
    mqtt_connected = True
    client.subscribe([
        ("mega/req/#", 1),
        ("mega/flash", 1),
        ("pygateway/justprint", 1),
    ])
    print("MQTT connected, subscribed to topics")
    # note: making serial open/close with mqtt because there's a strange
    # bug where on mqtt reconnect, the serial just reads garbage bytes.
    # maybe some strange threading thing?
    serialConn.open()
    print("opened serial")


def on_disconnect(client, userdata, rc):
    mqtt_connected = False
    print("MQTT disconnected")
    serialConn.close()
    print("closed serial")

# default handler for things we aren't handling explicitly


def on_message(client, userdata, msg: mqtt.MQTTMessage):
    print("Received mqtt topic" + msg.topic + "without handler. Doing nothing.")


def just_print_handler(client, userdata, msg: mqtt.MQTTMessage):
    global justPrint
    justPrint = not justPrint

# handler for mega requests


def mega_handler(client, userdata, msg: mqtt.MQTTMessage):
    # send to serial
    # [topic];[payload]\n
    serialConn.write(msg.topic.encode("utf-8"))
    serialConn.write(';'.encode("utf-8"))
    serialConn.write(msg.payload)
    serialConn.write('\n'.encode("utf-8"))
    print("forwarded msg on {} to serial".format(msg.topic))

# This flashes the arduino mega
# spawns avrdude script to flash the message payload to the board


def flash_mega(client: mqtt.Client, userdata, msg: mqtt.MQTTMessage):
    global flashing
    with open(FIRMWARE_LOCATION, 'wb') as f:
        f.write(msg.payload)

    flashing = True
    serialConn.close()
    time.sleep(0.2)

    res = subprocess.run(["/bin/sh", "/src/flash", FIRMWARE_LOCATION])
    if res.returncode != 0:
        print("Flash command failed.")
        client.publish("mega/flashresp", "Flash failed")
    else:
        print("Flash command successful.")
        client.publish("mega/flashresp", "Flash complete")

    print("waiting...")
    time.sleep(1)
    print("Opening serial port...")

    serialConn.open()
    flashing = False
    print("Serial open")
    client.publish("mega/flashresp", "Serial reopened")


#######################
### SERIAL HANDLERS ###
#######################

def term_handler(signum, frame):
    global exiting
    print("term received, closing serial")
    exiting = True
    serialConn.close()
    exit()


if __name__ == "__main__":
    debug = True if os.environ["DEBUG_PYGATEWAY"].lower() == "true" else False

    signal.signal(signal.SIGTERM, term_handler)
    signal.signal(signal.SIGINT, term_handler)

    # SERIAL
    print("Configuring serial...")
    serialConn = serial.Serial()
    serialConn.port = '/dev/ttyACM0'
    serialConn.baudrate = 1000000
    serialConn.timeout = 10

    # MQTT
    print("Sleeping to ensure mqtt broker is running after computer restart")
    time.sleep(2)

    client = mqtt.Client()
    client.on_connect = on_connect
    client.on_disconnect = on_disconnect

    # client.connect("localhost", 1883, 10)
    host = os.environ["BROKER_HOST"]
    print("connecting to", host)
    client.connect(host, 1883, 10)
    print("Connected to broker")

    client.message_callback_add("mega/flash", flash_mega)
    client.message_callback_add("mega/req/#", mega_handler)
    client.message_callback_add("pygateway/justprint", just_print_handler)
    client.on_message = on_message

    client.loop_start()
    print("Started broker network loop")

    START_SYMBOL = b'>'
    TOPIC_END = b';'
    PLAINTEXT_IDENTIFIER = b'#'
    PROTOBUF_IDENTIFIER = b'$'
    PAYLOAD_END = b'\n'
    MISC_TOPIC = "mega/l/misc"

    while not exiting:
        # Reconnect serial if down
        if mqtt_connected and not serialConn.isOpen():
            serialConn.open()
            time.sleep(0.5)

        while serialConn.isOpen():
            if justPrint:
                if serialConn.in_waiting > 0:
                    b = serialConn.read_all()
                    print(b.hex(':'))
            else:
                try:
                    # Read until start, handle undefined output
                    miscOutput = serialConn.read_until(START_SYMBOL)[:-1]
                    if len(miscOutput) > 0:
                        client.publish(MISC_TOPIC, miscOutput)

                    debugf("received start symbol, discarding the following:", miscOutput)
                    # now it's the topic
                    topic = serialConn.read_until(TOPIC_END)[:-1]
                    debugf(f"received topic end for '{topic}'")
                    payloadType = serialConn.read()
                    debugf(f"received payload type '{payloadType}'")

                    if payloadType == PLAINTEXT_IDENTIFIER:
                        payload = serialConn.read_until(PAYLOAD_END)[:-1]
                        debugf(f"received plaintext payload '{payload}'")
                    elif payloadType == PROTOBUF_IDENTIFIER:
                        # todo: support longer lengths than 255
                        payloadSizeRaw = serialConn.read(1)
                        payloadSize = int(payloadSizeRaw[0])
                        debugf("received protobuf payload size", payloadSize)
                        payload = serialConn.read(payloadSize)
                        debugf("received protobuf payload:,", payload)
                        end = serialConn.read()
                        if end != PAYLOAD_END:
                            print("error, payload_end not found after protobuf")
                            continue
                    else:
                        print("error, payloadType", payloadType, "is invalid")
                        print("Closing serial connection")
                        serialConn.close()
                        continue

                    print("topic:", topic, "; payload:", payload)
                    client.publish(topic.decode("utf-8"), payload)
                except TypeError as err:
                    print("TypeError:", err)

        time.sleep(0.1)
