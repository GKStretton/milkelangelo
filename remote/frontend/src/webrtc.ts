// Relay (TURN) servers from goo, for viewers who can't reach MediaMTX
// directly. null when goo has none, so MediaMTX's own list is used.
async function relayIceServers(): Promise<RTCIceServer[] | null> {
	try {
		const res = await fetch("/api/ice-servers", { cache: "no-store" });
		if (!res.ok) return null;
		const data = await res.json();
		return Array.isArray(data?.iceServers) ? data.iceServers : null;
	} catch {
		return null;
	}
}

// Plays a MediaMTX WebRTC stream (websocket signalling) into a video element.
// Adapted from interface/src/util/WebRTCReceiver.ts, with a way to stop it.
export class WebRTCReceiver {
	private stopped = false;
	private ws: WebSocket | null = null;
	private pc: RTCPeerConnection | null = null;
	private restartTimeout: number | null = null;

	constructor(
		private url: string,
		private video: HTMLVideoElement,
	) {
		this.start();
	}

	stop() {
		this.stopped = true;
		if (this.restartTimeout !== null) window.clearTimeout(this.restartTimeout);
		this.teardown();
	}

	private async start() {
		// fetched per connection, since the credentials expire
		const relay = await relayIceServers();
		if (this.stopped) return;

		const ws = new WebSocket(this.url);
		this.ws = ws;
		ws.onerror = () => ws.close();
		ws.onclose = () => this.scheduleRestart();
		ws.onmessage = (msg) => this.onIceServers(msg, relay);
	}

	private onIceServers(msg: MessageEvent, relay: RTCIceServer[] | null) {
		const ws = this.ws;
		if (!ws) return;

		const pc = new RTCPeerConnection({ iceServers: relay ?? JSON.parse(msg.data) });
		this.pc = pc;

		ws.onmessage = (msg) => {
			pc.setRemoteDescription(new RTCSessionDescription(JSON.parse(msg.data)));
			ws.onmessage = (msg) => pc.addIceCandidate(JSON.parse(msg.data));
		};
		pc.onicecandidate = (evt) => {
			if (evt.candidate && evt.candidate.candidate !== "") {
				ws.send(JSON.stringify(evt.candidate));
			}
		};
		pc.oniceconnectionstatechange = () => {
			if (pc.iceConnectionState === "failed" || pc.iceConnectionState === "closed") {
				this.scheduleRestart();
			}
		};
		pc.ontrack = (evt) => {
			if (this.video.srcObject !== evt.streams[0]) {
				this.video.srcObject = evt.streams[0];
				this.video.play().catch(() => {});
			}
		};

		pc.addTransceiver("video", { direction: "sendrecv" });
		pc.addTransceiver("audio", { direction: "sendrecv" });
		pc.createOffer().then((desc) => {
			pc.setLocalDescription(desc);
			ws.send(JSON.stringify(desc));
		});
	}

	private teardown() {
		if (this.ws) {
			this.ws.onclose = null;
			this.ws.close();
			this.ws = null;
		}
		if (this.pc) {
			this.pc.close();
			this.pc = null;
		}
	}

	private scheduleRestart() {
		if (this.stopped || this.restartTimeout !== null) return;
		this.teardown();
		this.restartTimeout = window.setTimeout(() => {
			this.restartTimeout = null;
			this.start();
		}, 2000);
	}
}
