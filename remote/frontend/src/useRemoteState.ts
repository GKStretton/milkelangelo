import { useEffect, useState } from "react";
import type { RemoteState } from "./types";

// Follows /api/events. EventSource reconnects by itself after errors.
export function useRemoteState() {
	const [state, setState] = useState<RemoteState | null>(null);
	const [connected, setConnected] = useState(false);

	useEffect(() => {
		const es = new EventSource("/api/events");
		es.addEventListener("state", (e) => {
			setState(JSON.parse((e as MessageEvent).data));
			setConnected(true);
		});
		es.onerror = () => setConnected(false);
		return () => es.close();
	}, []);

	return { state, connected };
}
