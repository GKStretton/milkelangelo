import { useContext, useEffect } from "react";
import MqttContext from "./mqttContext";
import {
  StateReport,
  SessionStatus,
  StreamStatus,
  VialProfileCollection,
  SystemVialConfiguration,
} from "../machinepb/machine";
import { TOPIC_STATE_REPORT_RAW } from "../topics_firmware/topics_firmware";
import {
  TOPIC_STREAM_STATUS_RESP_RAW,
  TOPIC_SESSION_STATUS_RESP_RAW,
  TOPIC_KV_GET_RESP,
  KV_KEY_ALL_VIAL_PROFILES,
  TOPIC_KV_SET,
  KV_KEY_SYSTEM_VIAL_PROFILES,
  TOPIC_KV_GET,
  TOPIC_STATE_REPORT_JSON,
} from "../topics_backend/topics_backend";
import { Type } from "protobufjs";

export function useProtoTopic(topic: string): Buffer | null {
  const { client: c, messages } = useContext(MqttContext);

  // Subscribe
  useEffect(() => {
    c?.subscribe(topic, (m) => {
      console.log(`subscribed to '${topic}': ${m}`);
    });
  }, [c?.connected]);

  if (!c || !c.connected) {
    console.error("cannot useProtoTopic, client not connected:", c);
    return null;
  }

  return messages[topic] ?? null;
}

export function useStateReport(): StateReport | null {
  const msg: Buffer | null = useProtoTopic(TOPIC_STATE_REPORT_JSON);
  if (!msg) return null;

  return StateReport.fromJSON(JSON.parse(msg.toString()));
}

export function useSessionStatus(): SessionStatus | null {
  const msg: Buffer | null = useProtoTopic(TOPIC_SESSION_STATUS_RESP_RAW);
  if (!msg) return null;

  return SessionStatus.decode(msg);
}

export function useStreamStatus(): StreamStatus | null {
  const msg: Buffer | null = useProtoTopic(TOPIC_STREAM_STATUS_RESP_RAW);
  if (!msg) return null;

  return StreamStatus.decode(msg);
}

export function useVialProfiles(): [VialProfileCollection | null, (collection: VialProfileCollection) => void] {
  const msg: Buffer | null = useProtoTopic(TOPIC_KV_GET_RESP + KV_KEY_ALL_VIAL_PROFILES);
  const obj = msg ? JSON.parse(msg.toString()) : null;
  const profile = obj ? VialProfileCollection.fromJSON(obj) : null;

  const { client: c } = useContext(MqttContext);

  // initial get request
  useEffect(() => {
    console.log("initial getting vial profiles");
    c?.publish(TOPIC_KV_GET + KV_KEY_ALL_VIAL_PROFILES, "");
  }, []);

  // publish function to be returned
  const pub = (collection: VialProfileCollection) => {
    if (!c) {
      console.error("cannot publish vial profile collection because mqtt client is null");
      return;
    }
    const json = JSON.stringify(VialProfileCollection.toJSON(collection), null, 2);
    c.publish(TOPIC_KV_SET + KV_KEY_ALL_VIAL_PROFILES, json);
  };

  return [profile, pub];
}

interface ProtoMethods<T> {
  fromJSON(data: object): any;
  toJSON(message: any): any;
  [key: string]: any;
}

function bufferToJSON(b: Buffer): object {
  return JSON.parse(b.toString());
}

/*
// Something strange going on with the types, I don't understand right now
export function useKeyValueStore<T extends ProtoMethods<T>>(
  name: string,
  messageType: T
): [T | null, (newValue: T) => void] {
  // Get latest value
  const raw: Buffer | null = useProtoTopic(TOPIC_KV_GET_RESP + name);
  const json = raw ? bufferToJSON(raw) : null;
  const msg = json ? messageType.fromJSON(json) : null;

  const { client: c } = useContext(MqttContext);

  // initial get request
  useEffect(() => {
    if (!c || !c.connected) {
      console.error(`cannot get ${name}, client not connected: ${c}`);
      return;
    }
    console.log(`initial getting ${name}`);
    c.publish(TOPIC_KV_GET + name, "");
  }, []);

  // publish function to be returned
  const pub = (message: T) => {
    if (!c) {
      console.error("cannot publish vial profile collection because mqtt client is null");
      return;
    }
    const json = JSON.stringify(messageType.toJSON(message), null, 2);
    c.publish(TOPIC_KV_SET + KV_KEY_ALL_VIAL_PROFILES, json);
  };

  return [msg, pub];
}

*/
