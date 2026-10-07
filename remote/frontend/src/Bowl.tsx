import { useVideoStream } from "./useVideoStream";
import VideoStatus from "./VideoStatus";

// The cropped top camera frames the bowl exactly, so video coordinates map
// straight onto the robot's unit circle.

// keep aimed points just inside the rim
const MAX_RADIUS = 0.95;

interface BowlProps {
	x: number;
	y: number;
	canAim: boolean;
	onAim: (x: number, y: number) => void;
}

export default function Bowl({ x, y, canAim, onAim }: BowlProps) {
	const { ref: videoRef, status } = useVideoStream("top-cam-crop");

	const handleClick = (e: React.MouseEvent<HTMLDivElement>) => {
		if (!canAim) return;
		const rect = e.currentTarget.getBoundingClientRect();
		let nx = ((e.clientX - rect.left) / rect.width) * 2 - 1;
		let ny = -(((e.clientY - rect.top) / rect.height) * 2 - 1);
		const r = Math.hypot(nx, ny);
		if (r > MAX_RADIUS) {
			nx = (nx / r) * MAX_RADIUS;
			ny = (ny / r) * MAX_RADIUS;
		}
		onAim(nx, ny);
	};

	return (
		<div
			className={`bowl ${canAim ? "can-aim" : ""}`}
			onClick={handleClick}
			role="img"
			aria-label={`Bowl, pipette at ${x.toFixed(2)}, ${y.toFixed(2)}`}
		>
			<video ref={videoRef} muted autoPlay playsInline />
			<VideoStatus status={status} className="in-bowl" />
			<div className="bowl-rim" />
			<div
				className="crosshair"
				style={{ left: `${(x + 1) * 50}%`, top: `${(1 - y) * 50}%` }}
			/>
			<div className="coords">
				{x.toFixed(2)}, {y.toFixed(2)}
			</div>
		</div>
	);
}
