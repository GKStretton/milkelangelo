import paho.mqtt.client as mqtt
from pycommon.const import HOST
import topics_backend.topics_backend as tb
import topics_firmware.topics_firmware as tf
import yaml

CLIENT_ID="py-interfaces"
DEBUG = False

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
    pub(tf.TOPIC_GOTO_XY, pl)

    debug("wrote goto_xy payload.")# Listening for response...")
    #! commented out because there's no timeout supported, so it hangs if
    #! there's no client responding
    # resp = sub(GOTO_RESP_TOPIC)
    # debug("got goto_xy response payload '{}'".format(resp.payload))

def dispense(ul):
    pl = "{:.1f}".format(ul)
    debug("writing dispense payload '{}'".format(pl))
    pub(tf.TOPIC_DISPENSE, pl)

    debug("wrote dispense payload")#. Listening for response...")
    # resp = sub(DISPENSE_RESP_TOPIC)
    # debug("got dispense response payload '{}'".format(resp.payload))

def collect(pos, ul):
    debug("writing collect payload '{}'".format(pos))
    pl = f"{pos},{ul:.1f}"
    pub(tf.TOPIC_COLLECT, pl)

    debug("wrote collect payload")

def sleep():
    pub(tf.TOPIC_SLEEP, "")

def shutdown():
    pub(tf.TOPIC_SHUTDOWN, "")

def wake():
    pub(tf.TOPIC_WAKE, "")

def uncalibrate():
    pub(tf.TOPIC_UNCALIBRATE, "")

def goto_node(node):
    pub(tf.TOPIC_GOTO_NODE, node)

def toggle_manual():
    pub(tf.TOPIC_TOGGLE_MANUAL, "")

def begin_session(production: bool = False):
    pl = ""
    if production:
        pl = "PRODUCTION"
    pub(tb.TOPIC_SESSION_BEGIN, pl)

def end_session():
    pub(tb.TOPIC_SESSION_END, "")

def pause_session():
    pub(tb.TOPIC_SESSION_PAUSE, "")

def resume_session():
    pub(tb.TOPIC_SESSION_RESUME, "")

def fluid_req(req_type, ml, open_drain=False):
    msg = f"{req_type},{ml:.2f},{open_drain}"
    pub(tf.TOPIC_FLUID, msg)

def start_stream():
    pub(tb.TOPIC_STREAM_START, "")

def end_stream():
    pub(tb.TOPIC_STREAM_END, "")