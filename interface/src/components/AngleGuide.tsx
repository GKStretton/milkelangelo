import React, { useRef, useState } from "react";
import { Button } from "@mui/material";

// Points are stored as fractions of the frame (0..1, y down) so the guide
// survives the video being resized.
interface Point {
	x: number;
	y: number;
}

const CENTRE: Point = { x: 0.5, y: 0.5 };
// Pointer movement (px) below which a press counts as a click, not a drag
const CLICK_THRESHOLD_PX = 5;

interface AngleGuideProps {
	width: number;
	height: number;
}

// Overlay for measuring angles on the top cam. While enabled, a line follows
// the mouse from the frame centre; click to lock/unlock it, or drag to draw a
// line between any two points.
//
// Angles use the same frame as TOPIC_GOTO_XY (x right, y up): 0° points
// right and angles increase anticlockwise. "Ring" is the angle plus 180°,
// which is what getRingAndYawFromXY in firmware ik_algorithm.cpp gives as the
// ring angle when the arm pivot sits at the end of a line drawn from the
// centre.
export default function AngleGuide({ width, height }: AngleGuideProps) {
	const [enabled, setEnabled] = useState(false);
	const [start, setStart] = useState<Point>(CENTRE);
	const [end, setEnd] = useState<Point | null>(null);
	const [locked, setLocked] = useState(false);
	const drag = useRef<{ down: Point; wasLocked: boolean } | null>(null);

	const toPoint = (e: React.PointerEvent<HTMLDivElement>): Point => {
		const rect = e.currentTarget.getBoundingClientRect();
		return {
			x: (e.clientX - rect.left) / rect.width,
			y: (e.clientY - rect.top) / rect.height,
		};
	};

	const onPointerDown = (e: React.PointerEvent<HTMLDivElement>) => {
		e.currentTarget.setPointerCapture(e.pointerId);
		const p = toPoint(e);
		drag.current = { down: p, wasLocked: locked };
		setLocked(false);
		setStart(p);
		setEnd(p);
	};

	const onPointerMove = (e: React.PointerEvent<HTMLDivElement>) => {
		const p = toPoint(e);
		if (drag.current) {
			setEnd(p);
		} else if (!locked) {
			setStart(CENTRE);
			setEnd(p);
		}
	};

	const onPointerUp = (e: React.PointerEvent<HTMLDivElement>) => {
		if (!drag.current) return;
		const p = toPoint(e);
		const { down, wasLocked } = drag.current;
		drag.current = null;
		const movedPx = Math.hypot((p.x - down.x) * width, (p.y - down.y) * height);
		if (movedPx < CLICK_THRESHOLD_PX) {
			// Click: line from the centre, toggling the lock
			setStart(CENTRE);
			setEnd(p);
			setLocked(!wasLocked);
		} else {
			setEnd(p);
			setLocked(true);
		}
	};

	const onPointerLeave = () => {
		if (!locked && !drag.current) setEnd(null);
	};

	let angle: number | null = null;
	let line: React.ReactNode = null;
	if (end) {
		const dx = (end.x - start.x) * width;
		const dy = (end.y - start.y) * height;
		const len = Math.hypot(dx, dy);
		if (len > 0) {
			angle = ((Math.atan2(-dy, dx) * 180) / Math.PI + 360) % 360;
			const sx = start.x * width;
			const sy = start.y * height;
			const ex = end.x * width;
			const ey = end.y * height;
			// Faint extension across the whole frame to line up with long edges
			const L = 2 * Math.max(width, height);
			const ux = dx / len;
			const uy = dy / len;
			line = (
				<>
					<line
						x1={sx - ux * L}
						y1={sy - uy * L}
						x2={sx + ux * L}
						y2={sy + uy * L}
						stroke="yellow"
						strokeOpacity={0.5}
						strokeDasharray="6 4"
						strokeWidth={1}
					/>
					<line x1={sx} y1={sy} x2={ex} y2={ey} stroke="yellow" strokeWidth={2} />
					<circle cx={sx} cy={sy} r={3} fill="yellow" />
				</>
			);
		}
	}

	return (
		<>
			{enabled && (
				<div
					onPointerDown={onPointerDown}
					onPointerMove={onPointerMove}
					onPointerUp={onPointerUp}
					onPointerLeave={onPointerLeave}
					style={{
						position: "absolute",
						top: 0,
						left: 0,
						width: `${width}px`,
						height: `${height}px`,
						cursor: "crosshair",
						touchAction: "none",
					}}
				>
					<svg
						width={width}
						height={height}
						style={{ position: "absolute", top: 0, left: 0, pointerEvents: "none" }}
					>
						{line}
					</svg>
				</div>
			)}
			<div
				style={{
					position: "absolute",
					top: 8,
					left: 8,
					display: "flex",
					alignItems: "center",
					gap: 8,
				}}
			>
				<Button
					size="small"
					variant={enabled ? "contained" : "outlined"}
					onClick={() => {
						setEnabled(!enabled);
						setLocked(false);
						setEnd(null);
					}}
					sx={{ bgcolor: enabled ? undefined : "rgba(0,0,0,0.5)" }}
				>
					Angle
				</Button>
				{enabled && (
					<span
						style={{
							background: "rgba(0,0,0,0.6)",
							color: "yellow",
							padding: "2px 6px",
							borderRadius: 4,
							fontFamily: "monospace",
							fontSize: 13,
							pointerEvents: "none",
						}}
					>
						{angle === null
							? "move / click / drag"
							: `${angle.toFixed(1)}° · ring ${((angle + 180) % 360).toFixed(1)}°${
									locked ? " 🔒" : ""
							  }`}
					</span>
				)}
			</div>
		</>
	);
}
