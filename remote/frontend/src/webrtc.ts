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

// how long to wait for video before giving up and retrying
const CONNECT_TIMEOUT_MS = 20_000;
const RETRY_DELAY_MS = 2_000;

export type ReceiverEvent =
	| { kind: "progress"; step: string }
	| { kind: "error"; message: string }
	// a connection that was already showing video recovered by itself
	| { kind: "resumed" };

// Plays a MediaMTX WebRTC stream (websocket signalling) into a video element,
// reporting progress and errors, and retrying after failures.
// Adapted from interface/src/util/WebRTCReceiver.ts.
export class WebRTCReceiver {
	private stopped = false;
	private ws: WebSocket | null = null;
	private pc: RTCPeerConnection | null = null;
	private restartTimeout: number | null = null;
	private connectTimeout: number | null = null;
	private connected = false;
	private hasPlayed = false;

	constructor(
		private url: string,
		private video: HTMLVideoElement,
		private report: (e: ReceiverEvent) => void,
	) {
		this.start();
	}

	stop() {
		this.stopped = true;
		if (this.restartTimeout !== null) window.clearTimeout(this.restartTimeout);
		this.teardown();
	}

	private async start() {
		this.connected = false;
		this.hasPlayed = false;
		this.report({ kind: "progress", step: "Connecting to the camera" });
		this.connectTimeout = window.setTimeout(
			() => this.fail("The camera took too long to start"),
			CONNECT_TIMEOUT_MS,
		);

		// fetched per connection, since the credentials expire
		const relay = await relayIceServers();
		if (this.stopped) return;

		const ws = new WebSocket(this.url);
		this.ws = ws;
		ws.onclose = (e) => {
			const reason = e.reason ? `: ${e.reason}` : "";
			this.fail(
				this.connected
					? `The video stream stopped${reason}`
					: `Couldn't reach the video server${reason}`,
			);
		};
		ws.onmessage = (msg) => this.onIceServers(msg, relay);
	}

	private onIceServers(msg: MessageEvent, relay: RTCIceServer[] | null) {
		const ws = this.ws;
		if (!ws) return;

		this.report({ kind: "progress", step: "Setting up video" });
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
			switch (pc.iceConnectionState) {
				case "checking":
					this.report({ kind: "progress", step: "Finding a route to the camera" });
					break;
				case "connected":
				case "completed":
					this.connected = true;
					this.report(
						this.hasPlayed
							? { kind: "resumed" }
							: { kind: "progress", step: "Starting video" },
					);
					break;
				case "failed":
					this.fail("Couldn't find a network route to the camera");
					break;
				case "disconnected":
					this.report({ kind: "progress", step: "Connection interrupted, waiting" });
					break;
			}
		};
		pc.ontrack = (evt) => {
			if (this.video.srcObject !== evt.streams[0]) {
				this.video.srcObject = evt.streams[0];
				this.video.play().catch((err: Error) => {
					this.report({ kind: "error", message: `The browser blocked playback: ${err.message}` });
				});
			}
		};

		pc.addTransceiver("video", { direction: "sendrecv" });
		pc.addTransceiver("audio", { direction: "sendrecv" });
		pc.createOffer()
			.then((desc) => {
				pc.setLocalDescription(desc);
				ws.send(JSON.stringify(desc));
			})
			.catch((err: Error) => this.fail(`Couldn't set up video: ${err.message}`));
	}

	// called by the owner once frames are showing
	playing() {
		this.hasPlayed = true;
		this.clearConnectTimeout();
	}

	private fail(message: string) {
		if (this.stopped || this.restartTimeout !== null) return;
		this.report({ kind: "error", message });
		this.teardown();
		this.restartTimeout = window.setTimeout(() => {
			this.restartTimeout = null;
			this.start();
		}, RETRY_DELAY_MS);
	}

	private clearConnectTimeout() {
		if (this.connectTimeout !== null) {
			window.clearTimeout(this.connectTimeout);
			this.connectTimeout = null;
		}
	}

	private teardown() {
		this.clearConnectTimeout();
		if (this.ws) {
			this.ws.onclose = null;
			this.ws.close();
			this.ws = null;
		}
		if (this.pc) {
			this.pc.oniceconnectionstatechange = null;
			this.pc.close();
			this.pc = null;
		}
	}
}
