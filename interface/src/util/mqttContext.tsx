// mqttContext.ts
import { createContext } from "react";
import { MqttClient } from "precompiled-mqtt";

export interface MqttContextValue {
  client: MqttClient | null;
  messages: { [topic: string]: Buffer };
}

const MqttContext = createContext<MqttContextValue>({
  client: null,
  messages: {},
});

export default MqttContext;
