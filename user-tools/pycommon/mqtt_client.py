import paho.mqtt.client as mqtt
from pycommon.const import HOST
import yaml

CLIENT_ID="py-interfaces"
DEBUG = True

GOTO_NODE_TOPIC = "mega/req/goto-node"
GOTO_TOPIC = "mega/req/goto-xy"
GOTO_RESP_TOPIC = "mega/resp/goto-xy"
DISPENSE_TOPIC = "mega/req/dispense"
DISPENSE_RESP_TOPIC = "mega/resp/dispense"
COLLECT_TOPIC = "mega/req/collect"
SLEEP_TOPIC = "mega/req/sleep"
WAKE_TOPIC = "mega/req/wake"
UNCALIBRATE_TOPIC = "mega/req/uncalibrate"
OPEN_DRAIN_TOPIC = "mega/req/open-drain"
CLOSE_DRAIN_TOPIC = "mega/req/close-drain"

# mqtt client
client = None

# debug print
def debug(msg):
    if DEBUG:
        print("[MQTT]", msg)

def on_connect(client, userdata, flags, rc):
    # client.subscribe([
    #     ("topic", 1)
    # ])
    print("Connected to broker")

def on_disconnect(client, userdata, rc):
    print("Disconnected from broker")

def connect():
    global client
    client = mqtt.Client(reconnect_on_failure=True)
    client.on_connect = on_connect
    client.on_disconnect = on_disconnect
    client.connect(HOST, 1883, 10)
    client.loop_start()
    print("Starter broker network loop")
    

def pub(topic, payload):
    global client
    if client is None:
        print("client is None, call connect?")
    else:
        client.publish(topic, payload)

def goto_xy(x, y):
    pl = "{:.3f},{:.3f}".format(x, y)

    debug("writing goto_xy payload '{}'".format(pl))
    pub(GOTO_TOPIC, pl)

    debug("wrote goto_xy payload.")# Listening for response...")
    #! commented out because there's no timeout supported, so it hangs if
    #! there's no client responding
    # resp = sub(GOTO_RESP_TOPIC)
    # debug("got goto_xy response payload '{}'".format(resp.payload))

def dispense(ul):
    debug("writing dispense payload '{}'".format(ul))
    pub(DISPENSE_TOPIC, ul)

    debug("wrote dispense payload")#. Listening for response...")
    # resp = sub(DISPENSE_RESP_TOPIC)
    # debug("got dispense response payload '{}'".format(resp.payload))

def collect(pos, ul):
    debug("writing collect payload '{}'".format(pos))
    pl = "{},{:.1f}".format(pos, ul)
    pub(COLLECT_TOPIC, pl)

    debug("wrote collect payload")

def sleep():
    pub(SLEEP_TOPIC, "")

def wake():
    pub(WAKE_TOPIC, "")

def uncalibrate():
    pub(UNCALIBRATE_TOPIC, "")

def set_drain(b: bool):
    if b:
        pub(OPEN_DRAIN_TOPIC, "")
    else:
        pub(CLOSE_DRAIN_TOPIC, "")
    
def goto_node(node):
    pub(GOTO_NODE_TOPIC, node)