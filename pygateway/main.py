import serial
import time
import paho.mqtt.client as mqtt

def on_connect(client, userdata, flags, rc):
    print("MQTT connected")

def on_disconnect(client, userdata, rc):
    print("MQTT disconnected")

def handleInput(client: mqtt.Client, line: str):
    if line[0] == '>':
        data = line[1:].strip("\r\n\t ").split('\t')
        if len(data) == 2:
            client.publish(data[0], data[1])
            # print(data[0], ":", data[1])
        else:
            print("malformed data, expected 2 parameters after splitting on string:", line)
    elif line == "" or line == "\r\n" or line == "\n":
        return
    else:
        print(line)

# subscribe to all topics, or whitelist?
#   Write [topic];[payload]\n to serial

# subscribe to /flash
#   spawn avrdude script to flash the message payload to the board

if __name__ == "__main__":
    client = mqtt.Client()
    client.on_connect = on_connect
    client.on_disconnect = on_disconnect

    client.connect("mosquitto", 1883, 10)

    print("Attempting to open serial interface...")
    with serial.Serial('/dev/ttyACM0', 1000000, timeout=1) as ser:
        print("Opened.")

        while True:
            if ser.inWaiting() > 0:
                handleInput(client, ser.readline().decode("utf-8"))