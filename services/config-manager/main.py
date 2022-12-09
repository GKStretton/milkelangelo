import paho.mqtt.client as mqtt
import signal
import yaml
import time
import os

YML_FILE = os.getenv("CROP_YML_FILE")

def on_connect(client, userdata, flags, rc):
	print("MQTT connected")

def on_disconnect(client, userdata, rc):
	print("MQTT disconnected")

def on_message(client: mqtt.Client, userdata, msg: mqtt.MQTTMessage):
	print("got misc message", msg.payload, "on", msg.topic)

def set_config_listener(client: mqtt.Client, userdata, msg: mqtt.MQTTMessage):
	print("received set config request")
	# write it to the file at /crop/crop.yml
	yml = yaml.load(msg.payload, Loader=yaml.FullLoader)
	print(yml)
	time.sleep(2)
	with open(YML_FILE, 'w') as f:
		yaml.dump(yml, f)
		print("config written to", YML_FILE)
	client.publish("crop-config/set-resp", "ACK")
	print("ACK sent")

def get_config_listener(client: mqtt.Client, userdata, msg: mqtt.MQTTMessage):
	print("received get config request")
	# Add delay because it was returning faster than the client could subscribe
	time.sleep(2)
	if os.path.isfile(YML_FILE):
		with open(YML_FILE, 'r') as f:
			conf = f.read()
			print("returning", conf)
			client.publish("crop-config/get-resp", conf)
	else:
		print("file not found", YML_FILE)
		client.publish("crop-config/get-resp", "404")

	# wait for client to reconnect
	print("sent response")

def term_handler(signum, frame):
	exit()

print("launch")
if __name__ == "__main__":
	print("Starting")

	signal.signal(signal.SIGTERM, term_handler)
	signal.signal(signal.SIGINT, term_handler)

	print("Sleeping to ensure mqtt broker is running after computer restart")
	time.sleep(2)

	# MQTT
	client = mqtt.Client()
	client.on_connect = on_connect
	client.on_disconnect = on_disconnect
	client.on_message = on_message

	client.connect("mosquitto", 1883, 10)
	print("Connected to broker")
	client.message_callback_add("crop-config/set", set_config_listener)
	client.message_callback_add("crop-config/get", get_config_listener)

	client.subscribe([
		("crop-config/set", 1),
		("crop-config/get", 1),
	])

	print("Starting broker network loop")
	client.loop_forever()