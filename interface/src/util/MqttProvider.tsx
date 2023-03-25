// MqttProvider.tsx
import React, { useState, useEffect, ReactNode } from 'react';
import mqtt, { MqttClient } from 'precompiled-mqtt';
import MqttContext, { MqttContextValue } from './mqttContext';

interface MqttProviderProps {
	children: ReactNode;
}

const MqttProvider: React.FC<MqttProviderProps> = ({ children }) => {
	const [client, setClient] = useState<MqttClient | null>(null);
	const [messages, setMessages] = useState<{ [topic: string]: string }>({});

	useEffect(() => {
		const mqttClient = mqtt.connect('ws://DEPTH:9001');

		mqttClient.on('connect', () => {
			console.log('connected');
			setClient(mqttClient);
		});

		mqttClient.on('message', (topic, message) => {
			console.log("got message %s %s", topic, message.toString());
			setMessages((prevMessages) => ({
				...prevMessages,
				[topic]: message.toString(),
			}));
		});

		return () => {
			mqttClient.end();
		};
	}, []);

	const contextValue: MqttContextValue = { client, messages };

	return (
		<MqttContext.Provider value={contextValue}>
			{client ? children : <div>Loading...</div>}
		</MqttContext.Provider>
	);
};

export default MqttProvider;
