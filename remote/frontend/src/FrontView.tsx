import { useVideoStream } from "./useVideoStream";
import VideoStatus from "./VideoStatus";

// The cropped front camera, looking at the bowl from the side. View only.
export default function FrontView() {
	const { ref, status } = useVideoStream("front-cam-crop");

	return (
		<div className={`front ${status.phase === "playing" ? "has-video" : ""}`}>
			<video ref={ref} muted autoPlay playsInline />
			<VideoStatus status={status} />
			<span className="front-label">Front</span>
		</div>
	);
}
