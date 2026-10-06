import { useCallback, useEffect, useRef, useState } from "react";
import { api, unclaimOnExit } from "./api";

// goo drops a claim after 8s without renewal
const RENEW_INTERVAL_MS = 3000;

export interface Notice {
	tone: "info" | "error";
	text: string;
}

// Holds this page's claim on the robot, and sends commands with it.
export function useControl() {
	const [token, setToken] = useState<string | null>(null);
	const [pending, setPending] = useState(false);
	const [notice, setNotice] = useState<Notice | null>(null);
	const tokenRef = useRef<string | null>(null);
	tokenRef.current = token;

	useEffect(() => {
		if (!notice) return;
		const t = window.setTimeout(() => setNotice(null), 5000);
		return () => window.clearTimeout(t);
	}, [notice]);

	// keep the claim alive
	useEffect(() => {
		if (!token) return;
		const interval = window.setInterval(async () => {
			const res = await api.claim(token);
			if (res.status === 403) {
				setToken(null);
				setNotice({ tone: "error", text: "You lost control of the robot" });
			}
		}, RENEW_INTERVAL_MS);
		return () => window.clearInterval(interval);
	}, [token]);

	useEffect(() => {
		const onExit = () => {
			if (tokenRef.current) unclaimOnExit(tokenRef.current);
		};
		window.addEventListener("pagehide", onExit);
		return () => window.removeEventListener("pagehide", onExit);
	}, []);

	const claim = useCallback(async () => {
		setPending(true);
		const res = await api.claim(null);
		setPending(false);
		if (res.ok && res.data?.token) {
			setToken(res.data.token);
		} else {
			setNotice({ tone: "error", text: res.message });
		}
	}, []);

	const release = useCallback(async () => {
		const t = tokenRef.current;
		setToken(null);
		if (t) await api.unclaim(t);
	}, []);

	// run a command with the current token, reporting a rejection
	const send = useCallback(
		async (
			command: (token: string) => ReturnType<typeof api.collect>,
			success?: string,
		) => {
			const t = tokenRef.current;
			if (!t) return;
			setPending(true);
			const res = await command(t);
			setPending(false);
			if (res.ok) {
				if (success) setNotice({ tone: "info", text: success });
			} else {
				if (res.status === 403) setToken(null);
				setNotice({ tone: "error", text: res.message });
			}
		},
		[],
	);

	return {
		hasControl: token !== null,
		pending,
		notice,
		claim,
		release,
		collect: (id: number, name: string) =>
			send((t) => api.collect(t, id), `Collecting ${name}`),
		dispense: (x: number, y: number) => send((t) => api.dispense(t, x, y)),
		goTo: (x: number, y: number) => send((t) => api.goTo(t, x, y)),
	};
}
