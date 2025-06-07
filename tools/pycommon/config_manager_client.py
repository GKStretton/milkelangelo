import paho.mqtt.publish as mqttpub
import paho.mqtt.subscribe as mqttsub
import yaml
import topics_backend.topics_backend as tb

CLIENT_ID="pythoninterfaces"
HOST="localhost"

def write_remote_crop_config(name, cfg):
    yml = yaml.dump(cfg)
    mqttpub.single(tb.TOPIC_KV_SET+name, payload=yml, hostname=HOST, port=1883, client_id=CLIENT_ID)
    print("sent conf to", tb.TOPIC_KV_SET+name)
    
    print("waiting for response...")
    resp = mqttsub.simple(tb.TOPIC_KV_SET_RESP+name, hostname=HOST, port=1883, client_id=CLIENT_ID, keepalive=1)
    print("got response")
    print(resp.payload.decode("utf-8"))

def read_remote_crop_config(name):
    print("sending config request to", tb.TOPIC_KV_GET+name)
    print("and waiting for response on", tb.TOPIC_KV_GET_RESP+name, "...")
    mqttpub.single(tb.TOPIC_KV_GET+name, payload="", hostname=HOST, port=1883, client_id=CLIENT_ID)
    resp = mqttsub.simple(tb.TOPIC_KV_GET_RESP+name, hostname=HOST, port=1883, client_id=CLIENT_ID, keepalive=1)
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
    mqttpub.single(tb.TOPIC_TRIGGER_DSLR, payload="", hostname=HOST, port=1883, client_id=CLIENT_ID)
    resp = mqttsub.simple(tb.TOPIC_TRIGGER_DSLR_RESP, hostname=HOST, port=1883, client_id=CLIENT_ID, keepalive=1)
    msg = resp.payload.decode("utf-8")
    if msg == "ack":
        return True
    
    print("noack:", msg)
    return False