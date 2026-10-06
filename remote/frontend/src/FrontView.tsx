import { useVideoStream } from "./useVideoStream";

// The cropped front camera, looking at the bowl from the side. View only.
export default function FrontView() {
	const { ref, hasVideo } = useVideoStream("front-cam-crop");

	return (
		<div className={`front ${hasVideo ? "has-video" : ""}`}>
			<video ref={ref} muted autoPlay playsInline />
			{!hasVideo && <div className="front-blank">No front camera</div>}
			<span className="front-label">Front</span>
		</div>
	);
}
