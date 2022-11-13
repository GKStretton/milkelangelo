import signal
import serial
import time
import paho.mqtt.client as mqtt
import subprocess

# global var for the serial connection
serialConn = None

FIRMWARE_LOCATION = "./light.hex"
flashing = False

exiting = False

#####################
### MQTT HANDLERS ###
#####################


def on_connect(client, userdata, flags, rc):
    print("MQTT connected")

def on_disconnect(client, userdata, rc):
    print("MQTT disconnected")

# default handler for things we aren't handling explicitly
def on_message(client, userdata, msg: mqtt.MQTTMessage):
    print("Received mqtt topic" + msg.topic + "without handler. Doing nothing.")

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
    time.sleep(1)

    res = subprocess.run(["/bin/sh", "/src/flash.sh", FIRMWARE_LOCATION])
    if res.returncode != 0:
        print("Flash command failed.")
        client.publish("mega/flashresp", "Flash failed")
    else:
        print("Flash command successful.")
        client.publish("mega/flashresp", "Flash complete")
    
    print("waiting...")
    time.sleep(2)
    print("Opening serial port...")

    serialConn.open()
    flashing = False
    print("Serial open") 
    client.publish("mega/flashresp", "Serial reopened")


#######################
### SERIAL HANDLERS ###
#######################


# process a line received over serial from arduino.
# >[topic];[payload]\n will be published as (topic, payload)
# anything else will be published as (mega/log/misc)
def handleSerialLine(client: mqtt.Client, line: str):
    if line[0] == '>':
        data = line[1:].strip("\r\n\t ").split(';')
        if len(data) == 2:
            client.publish(data[0], data[1])
            # print(data[0], ":", data[1])
        else:
            print("malformed data, expected 2 parameters after splitting on string:", line)
    elif line == "" or line == "\r\n" or line == "\n":
        # Empty lines, just skip
        return
    else:
        # Print anything that isn't an mqtt pub > and publish it to special topic
        l = line.strip("\r\n\t")
        print("Misc output:", l)
        client.publish("mega/log/misc", l)


def term_handler(signum, frame):
    global exiting
    print("term received, closing serial")
    exiting = True
    serialConn.close()
    exit()

if __name__ == "__main__":
    signal.signal(signal.SIGTERM, term_handler)
    signal.signal(signal.SIGINT, term_handler)

    # MQTT

    client = mqtt.Client()
    client.on_connect = on_connect
    client.on_disconnect = on_disconnect

    client.connect("mosquitto", 1883, 10)
    print("Connected to broker")
    

    client.message_callback_add("mega/flash", flash_mega)
    client.message_callback_add("mega/req/#", mega_handler)
    client.on_message = on_message
    client.subscribe([
        ("mega/req/#", 1),
        ("mega/flash", 1),
    ])

    client.loop_start()
    print("Started broker network loop")

    # SERIAL

    print("Attempting to open serial interface...")
    serialConn = serial.Serial('/dev/ttyACM0', 1000000, timeout=1)
    print("Opened Serial.")

    while not exiting:
        if not flashing and serialConn.inWaiting() > 0:
            handleSerialLine(client, serialConn.readline().decode("utf-8"))
        time.sleep(0.01)