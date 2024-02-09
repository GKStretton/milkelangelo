import { useContext, useEffect } from "react";
import MqttContext from "./mqttContext";
import {
  StateReport,
  SessionStatus,
  StreamStatus,
  VialProfileCollection,
  SystemVialConfiguration,
} from "../machinepb/machine";
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

export function useBoolTopic(reqTopic: string, respTopic: string): boolean {
  const { client: c, messages } = useContext(MqttContext);

  // Subscribe
  useEffect(() => {
    if (!c || !c.connected) {
      return;
    }
    c?.subscribe(respTopic, (m) => {
      console.log(`subscribed to '${respTopic}': ${m}`);
      c?.publish(reqTopic, "")
    });
  }, [c?.connected]);

  const msg = messages[respTopic];
  return (msg?.toString() ?? "").toLowerCase() === "true"
}

export function useProtoTopic(topic: string): Buffer | null {
  const { client: c, messages } = useContext(MqttContext);

  // Subscribe
  useEffect(() => {
    if (!c || !c.connected) {
      return;
    }
    c?.subscribe(topic, (m) => {
      console.log(`subscribed to '${topic}': ${m}`);
    });
  }, [c?.connected]);

  return messages[topic] ?? null;
}

export function useStateReport(): StateReport | null {
  const msg: Buffer | null = useProtoTopic(TOPIC_STATE_REPORT_JSON);
  if (!msg) return null;

  return StateReport.fromJSON(JSON.parse(msg.toString()));
}

// Could add an alternate useXRef that returns a ref wrapper so the latest
// value can be accessed from within closures like the key listeners.

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
  const { client: c } = useContext(MqttContext);

  // initial get request
  useEffect(() => {
    if (!c || !c.connected) {
      console.error(`cannot get vial profiles, client not connected: ${c}`);
      return;
    }
    console.log("initial getting vial profiles");
    c?.publish(TOPIC_KV_GET + KV_KEY_ALL_VIAL_PROFILES, "");
  }, [c?.connected]);

  // publish function to be returned
  const pub = (collection: VialProfileCollection) => {
    if (!c || !c.connected) {
      console.error("cannot publish vial profile collection because mqtt client is null");
      return;
    }
    const json = JSON.stringify(collection, null, 2);
    c.publish(TOPIC_KV_SET + KV_KEY_ALL_VIAL_PROFILES, json);
  };

  const msg: Buffer | null = useProtoTopic(TOPIC_KV_GET_RESP + KV_KEY_ALL_VIAL_PROFILES);
  if (!msg || msg.toString() == "") return [null, pub];

  const obj = JSON.parse(msg.toString());
  const profile = obj ? VialProfileCollection.fromJSON(obj) : null;

  return [profile, pub];
}

export function useSystemVialProfiles(): [
  SystemVialConfiguration | null,
  (collection: SystemVialConfiguration) => void
] {
  const { client: c } = useContext(MqttContext);

  // initial get request
  useEffect(() => {
    if (!c || !c.connected) {
      console.error(`cannot get system profiles, client not connected: ${c}`);
      return;
    }
    console.log("initial getting system vial profiles");
    c?.publish(TOPIC_KV_GET + KV_KEY_SYSTEM_VIAL_PROFILES, "");
  }, [c?.connected]);

  // publish function to be returned
  const pub = (collection: SystemVialConfiguration) => {
    if (!c || !c.connected) {
      console.error("cannot publish system vial profile collection because mqtt client is null");
      return;
    }
    const json = JSON.stringify(collection, null, 2);
    c.publish(TOPIC_KV_SET + KV_KEY_SYSTEM_VIAL_PROFILES, json);
  };

  const msg: Buffer | null = useProtoTopic(TOPIC_KV_GET_RESP + KV_KEY_SYSTEM_VIAL_PROFILES);
  if (!msg || msg.toString() == "") return [null, pub];

  const obj = JSON.parse(msg.toString());
  const profile = obj ? SystemVialConfiguration.fromJSON(obj) : null;

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
