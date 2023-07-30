// MqttProvider.tsx
import React, { useState, useEffect, ReactNode } from "react";
import mqtt, { MqttClient } from "precompiled-mqtt";
import MqttContext, { MqttContextValue } from "./mqttContext";

interface MqttProviderProps {
  url: string;
  children: ReactNode;
}

interface Credentials {
  username: string;
  password: string;
}

const MqttProvider: React.FC<MqttProviderProps> = ({ url, children }) => {
  const [client, setClient] = useState<MqttClient | null>(null);
  const [messages, setMessages] = useState<{ [topic: string]: Buffer }>({});
  const [credentials, setCredentials] = useState<Credentials>({ username: "", password: "" });

  const do_auth: boolean = process.env.REACT_APP_MQTT_AUTHENTICATE === "true";

  useEffect(() => {
    if (!do_auth || (credentials.username && credentials.password)) return;

    // Ask for username and password using prompts
    const username = prompt("Username") ?? "";
    const password = prompt("Password") ?? "";

    setCredentials({ username, password });
  }, []);

  useEffect(() => {
    if (do_auth && (!credentials.username || !credentials.password)) return;

    const options = {
      username: credentials.username,
      password: credentials.password,
    };

    const mqttClient = mqtt.connect(url, do_auth ? options : undefined);

    mqttClient.on("connect", () => {
      console.log("connected");
      setClient(mqttClient);
    });

    mqttClient.on("message", (topic, message) => {
      console.log("got message %s %s", topic, message.toString());
      setMessages((prevMessages) => ({
        ...prevMessages,
        [topic]: message,
      }));
    });

    return () => {
      mqttClient.end();
    };
  }, [credentials]);

  const contextValue: MqttContextValue = { client, messages };

  return <MqttContext.Provider value={contextValue}>{client ? children : <div>Loading...</div>}</MqttContext.Provider>;
};

export default MqttProvider;
