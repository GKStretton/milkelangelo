import { FormControlLabel, Switch } from "@mui/material";
import { TOPIC_KV_GET, TOPIC_KV_GET_RESP, TOPIC_KV_SET } from "../topics_backend/topics_backend";
import { useBoolTopic } from "./hooks";
import { useContext } from "react";
import MqttContext from "./mqttContext";

const useKVBool = (kv_name: string) => {
	const get_req = TOPIC_KV_GET + kv_name
	const get_resp = TOPIC_KV_GET_RESP + kv_name
	const set_req = TOPIC_KV_SET + kv_name

	const value = useBoolTopic(get_req, get_resp)
	const { client: c } = useContext(MqttContext);

	const setter = (new_value: boolean) => {
		c?.publish(set_req, new_value.toString())
	}

	return { value, setter};
}

export default function KVBool({
	kv_name,
	desc,
}: {
	kv_name: string,
	desc: string,
}) {
	const { value, setter } = useKVBool(kv_name);

	const changeHandler = (e: React.ChangeEvent<HTMLInputElement>) => {
		console.log(`Switch ${kv_name} changed to ${e.target.checked}`);
		setter(e.target.checked);
	};

	return (
		<FormControlLabel control={<Switch onChange={changeHandler} checked={value}/>} label={desc} />
	);
}