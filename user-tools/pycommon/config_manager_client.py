import paho.mqtt.publish as mqttpub
import paho.mqtt.subscribe as mqttsub
import yaml

KV_ROOT_TOPIC = "asol/kv/"
TOPIC_SET = KV_ROOT_TOPIC+"set/"
TOPIC_SET_RESP = KV_ROOT_TOPIC+"set-resp/"
TOPIC_GET = KV_ROOT_TOPIC+"get/"
TOPIC_GET_RESP = KV_ROOT_TOPIC+"get-resp/"

TOPIC_TRIGGER_DSLR = "asol/dslr-crop-capture"
TOPIC_TRIGGER_DSLR_RESP = "asol/dslr-crop-capture-resp"

CLIENT_ID="pythoninterfaces"
HOST="DEPTH"

def write_remote_crop_config(name, cfg):
    yml = yaml.dump(cfg)
    mqttpub.single(TOPIC_SET+name, payload=yml, hostname=HOST, port=1883, client_id=CLIENT_ID)
    print("sent conf to", TOPIC_SET+name)
    
    print("waiting for response...")
    resp = mqttsub.simple(TOPIC_SET_RESP+name, hostname=HOST, port=1883, client_id=CLIENT_ID, keepalive=1)
    print("got response")
    print(resp.payload.decode("utf-8"))

def read_remote_crop_config(name):
    print("sending config request to", TOPIC_GET+name)
    print("and waiting for response on", TOPIC_GET_RESP+name, "...")
    mqttpub.single(TOPIC_GET+name, payload="", hostname=HOST, port=1883, client_id=CLIENT_ID)
    resp = mqttsub.simple(TOPIC_GET_RESP+name, hostname=HOST, port=1883, client_id=CLIENT_ID, keepalive=1)
    if resp.payload.decode("utf-8") == "404":
        print("no config yet, using 0 values")
        return None
    else:
        print("got response")
        print(resp.payload.decode("utf-8"))
        cfg = yaml.load(resp.payload, Loader=yaml.FullLoader)
        return cfg

# returns True if successful, otherwise False
def trigger_dslr_capture():
    print("triggering dslr capture...")
    mqttpub.single(TOPIC_TRIGGER_DSLR, payload="", hostname=HOST, port=1883, client_id=CLIENT_ID)
    resp = mqttsub.simple(TOPIC_TRIGGER_DSLR_RESP, hostname=HOST, port=1883, client_id=CLIENT_ID, keepalive=1)
    msg = resp.payload.decode("utf-8")
    if msg == "ack":
        return True
    
    print("noack:", msg)
    return False