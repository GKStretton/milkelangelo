import type { VideoStatus as Status } from "./useVideoStream";

// Loading and error state shown over a camera view until video plays.
export default function VideoStatus({
	status,
	className = "",
}: {
	status: Status;
	className?: string;
}) {
	if (status.phase === "playing") return null;

	if (status.phase === "error") {
		return (
			<div className={`video-status error ${className}`} role="status">
				<span className="video-status-title">No video</span>
				<span className="video-status-detail">{status.message}</span>
				<span className="video-status-retry">
					<span className="spinner" aria-hidden="true" />
					Retrying{status.step ? ` · ${status.step}` : ""}
				</span>
			</div>
		);
	}

	return (
		<div className={`video-status ${className}`} role="status">
			<span className="spinner" aria-hidden="true" />
			<span className="video-status-detail">{status.step}…</span>
		</div>
	);
}
