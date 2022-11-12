import paho.mqtt.publish as mqttpub
import paho.mqtt.subscribe as mqttsub
import yaml

GET_TOPIC = "crop-config/get"
GET_RESP_TOPIC = "crop-config/get-resp"
SET_TOPIC = "crop-config/set"
SET_RESP_TOPIC = "crop-config/set-resp"
CLIENT_ID="pythoninterfaces"
HOST="DEPTH"

def write_remote_crop_config(yml):
    conf = yaml.dump(yml)
    mqttpub.single(SET_TOPIC, payload=conf, hostname=HOST, port=1883, client_id=CLIENT_ID)
    print("sent conf to", SET_TOPIC)
    
    print("waiting for response...")
    resp = mqttsub.simple(SET_RESP_TOPIC, hostname=HOST, port=1883, client_id=CLIENT_ID, keepalive=1)
    print("got response")
    print(resp.payload.decode("utf-8"))

def read_remote_crop_config():
    mqttpub.single(GET_TOPIC, payload="get", hostname=HOST, port=1883, client_id=CLIENT_ID)
    print("sent config request to", GET_TOPIC)
    print("waiting for response on", GET_RESP_TOPIC, "...")
    resp = mqttsub.simple(GET_RESP_TOPIC, hostname=HOST, port=1883, client_id=CLIENT_ID, keepalive=1)
    if resp.payload.decode("utf-8") == "404":
        print("no config yet, using 0 values")
        return None
    else:
        print("got response")
        print(resp.payload.decode("utf-8"))
        yml = yaml.load(resp.payload, Loader=yaml.FullLoader)
        return yml