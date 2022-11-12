import paho.mqtt.client as mqtt
from pycommon.const import HOST
import yaml
import paho.mqtt.publish as mqttpub
import paho.mqtt.subscribe as mqttsub

CLIENT_ID="py-interfaces"
DEBUG = True


GOTO_TOPIC = "mega/req/goto-xy"
GOTO_RESP_TOPIC = "mega/resp/goto-xy"
DISPENSE_TOPIC = "mega/req/dispense"
DISPENSE_RESP_TOPIC = "mega/resp/dispense"

# debug print
def debug(msg):
    if DEBUG:
        print("[MQTT]", msg)

def pub(topic, pl):
    mqttpub.single(topic, payload=pl, hostname=HOST, port=1883, client_id=CLIENT_ID, keepalive=2)

def sub(topic):
    return mqttsub.simple(topic, hostname=HOST, port=1883, client_id=CLIENT_ID, keepalive=2)


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