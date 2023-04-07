import { useContext } from 'react';
import MqttContext from './mqttContext';
import { StateReport, SessionStatus, StreamStatus } from '../machinepb/machine_pb';
import { TOPIC_STATE_REPORT_RAW } from '../topics_firmware/topics_firmware';
import { TOPIC_STREAM_STATUS_RESP_RAW, TOPIC_SESSION_STATUS_RESP_RAW } from '../topics_backend/topics_backend';

function useProtoTopic(topic: string): Buffer | null {
	const { client: c, messages } = useContext(MqttContext);
	const message = messages[topic];

	if (message) {
		return message;
	}

	return null;
}

export function useStateReport(): StateReport | null {
	const msg: Buffer | null = useProtoTopic(TOPIC_STATE_REPORT_RAW);
	if (!msg) return null;

	return StateReport.deserializeBinary(msg);
}

export function useSessionStatus(): SessionStatus | null {
	const msg: Buffer | null = useProtoTopic(TOPIC_SESSION_STATUS_RESP_RAW);
	if (!msg) return null;

	return SessionStatus.deserializeBinary(msg);
}

export function useStreamStatus(): StreamStatus | null {
	const msg: Buffer | null = useProtoTopic(TOPIC_STREAM_STATUS_RESP_RAW);
	if (!msg) return null;

	return StreamStatus.deserializeBinary(msg);
}
