import { useEffect, useRef, useState } from "react";
import { WebRTCReceiver } from "./webrtc";

// Signalling goes through goo on this page's origin (/video), so it works
// wherever the page does. VITE_WEBRTC_URL points straight at MediaMTX instead.
const VIDEO_BASE =
	import.meta.env.VITE_WEBRTC_URL ??
	`${window.location.protocol === "https:" ? "wss" : "ws"}://${window.location.host}/video`;

// Plays a MediaMTX stream into the returned video ref. hasVideo turns true
// once the first frame is decoded.
export function useVideoStream(stream: string) {
	const ref = useRef<HTMLVideoElement>(null);
	const [hasVideo, setHasVideo] = useState(false);

	useEffect(() => {
		const video = ref.current;
		if (!video) return;
		const receiver = new WebRTCReceiver(`${VIDEO_BASE}/${stream}/ws`, video);
		const onFrame = () => setHasVideo(true);
		const onEmptied = () => setHasVideo(false);
		video.addEventListener("loadeddata", onFrame);
		video.addEventListener("emptied", onEmptied);
		return () => {
			receiver.stop();
			video.removeEventListener("loadeddata", onFrame);
			video.removeEventListener("emptied", onEmptied);
		};
	}, [stream]);

	return { ref, hasVideo };
}
