import pycommon.machine_pb2
import paho.mqtt.subscribe as mqttsub

resp = mqttsub.simple("mega/state-report", hostname="depth", port=1883, client_id="meme", keepalive=1)
print("got state report:")
print(resp.payload)