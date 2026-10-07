import { useEffect, useRef, useState } from "react";
import { type ReceiverEvent, WebRTCReceiver } from "./webrtc";

// Signalling goes through goo on this page's origin (/video), so it works
// wherever the page does. VITE_WEBRTC_URL points straight at MediaMTX instead.
const VIDEO_BASE =
	import.meta.env.VITE_WEBRTC_URL ??
	`${window.location.protocol === "https:" ? "wss" : "ws"}://${window.location.host}/video`;

export type VideoStatus =
	| { phase: "loading"; step: string }
	| { phase: "playing" }
	// the receiver keeps retrying; step shows how the retry is going
	| { phase: "error"; message: string; step?: string };

// Plays a MediaMTX stream into the returned video ref, with its status.
export function useVideoStream(stream: string) {
	const ref = useRef<HTMLVideoElement>(null);
	const [status, setStatus] = useState<VideoStatus>({
		phase: "loading",
		step: "Connecting to the camera",
	});

	useEffect(() => {
		const video = ref.current;
		if (!video) return;

		const onEvent = (e: ReceiverEvent) =>
			setStatus((prev) => {
				if (e.kind === "error") return { phase: "error", message: e.message };
				if (e.kind === "resumed") return { phase: "playing" };
				// while retrying after an error, keep the error visible
				if (prev.phase === "error") return { ...prev, step: e.step };
				return { phase: "loading", step: e.step };
			});
		const receiver = new WebRTCReceiver(`${VIDEO_BASE}/${stream}/ws`, video, onEvent);

		// loadeddata fires once the first frame is decoded
		const onFrame = () => {
			receiver.playing();
			setStatus({ phase: "playing" });
		};
		video.addEventListener("loadeddata", onFrame);
		return () => {
			receiver.stop();
			video.removeEventListener("loadeddata", onFrame);
		};
	}, [stream]);

	return { ref, status };
}
