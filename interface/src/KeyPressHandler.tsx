// KeyPressHandler.js for global events
import { useEffect, useContext } from 'react';
import MqttContext from './util/mqttContext';
import { MqttClient } from 'precompiled-mqtt';

const handleKeyPress = (mqttClient: MqttClient | null, event: KeyboardEvent): void => {
	// console.log(event);
}

const KeyPressHandler: React.FC = () => {
  const { client: mqttClient, messages } = useContext(MqttContext);

  useEffect(() => {
    const keyDownHandler = (event: KeyboardEvent) => {
      handleKeyPress(mqttClient, event);
    };

    window.addEventListener('keydown', keyDownHandler);
    return () => {
      window.removeEventListener('keydown', keyDownHandler);
    };
  }, [mqttClient]);

  return null;
};

export default KeyPressHandler;