import serial
import time
import paho.mqtt.client as mqtt

# global var for the serial connection
serialConn = None

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

# handler for mega requests MEGA_REQUEST_PREFIX + #. only # gets passed to mega
def mega_handler(client, userdata, msg: mqtt.MQTTMessage):
    # send to serial
    # [topic];[payload]\n
    serialConn.write(msg.topic.encode("utf-8"))
    serialConn.write(';'.encode("utf-8"))
    serialConn.write(msg.payload)
    serialConn.write('\n'.encode("utf-8"))
    
# This flashes the arduino mega
def flash_mega(client, userdata, msg: mqtt.MQTTMessage):
    print("flash mega not implement")


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

# subscribe to all topics, or whitelist?
#   Write [topic];[payload]\n to serial

# subscribe to /flash
#   spawn avrdude script to flash the message payload to the board

if __name__ == "__main__":
    client = mqtt.Client()
    client.on_connect = on_connect
    client.on_disconnect = on_disconnect

    client.connect("mosquitto", 1883, 10)
    print("Connected to broker")
    

    client.message_callback_add("mega/flash", flash_mega)
    client.message_callback_add("mega/#", mega_handler)
    client.on_message = on_message
    # We don't want everything to prevent ard slowdown
    client.subscribe([
        ("mega/req/#", 1),
        ("mega/flash", 1),
    ])


    print("Attempting to open serial interface...")
    with serial.Serial('/dev/ttyACM0', 1000000, timeout=1) as ser:
        serialConn = ser
        print("Opened Serial.")

        client.loop_start()
        print("Started broker network loop")

        while True:
            if ser.inWaiting() > 0:
                handleSerialLine(client, ser.readline().decode("utf-8"))
            time.sleep(0.01)