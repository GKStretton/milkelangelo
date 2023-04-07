import { AppBar, Toolbar, Typography } from "@mui/material";
import { PrecisionManufacturing, Videocam, VideocamOff, Podcasts, PausePresentation } from "@mui/icons-material";
import { useSessionStatus, useStreamStatus } from "../util/hooks";
export default function Header() {
	const streamStatus = useStreamStatus();
	const sessionStatus = useSessionStatus();

	const sessionId = !sessionStatus ? "UNKOWN" :
		sessionStatus.getProduction() ? sessionStatus.getProductionId() : `dev.${sessionStatus.getId()}`

	return (
	<AppBar position="relative">
		<Toolbar>
			<PrecisionManufacturing/>
			<Typography variant="h6">
				A Study of Light
			</Typography>

			<div style={{margin: "20px"}}></div>

			{streamStatus?.getLive() ? <>
				<Podcasts sx={{margin: 1}}/>
				<Typography variant="h6" flexGrow={1}>
					LIVE
				</Typography>
			</> : <>
				<Podcasts sx={{margin: 1, opacity: 0.2}}/>
				<Typography variant="h6" flexGrow={1}>
					OFF
				</Typography>
			</>}

			{sessionStatus?.getComplete() === false && sessionStatus?.getPaused() === false && <>
				<Typography variant="h6">
					RECORDING {sessionId}
				</Typography>
				<Videocam sx={{margin: 1}}/>
			</>}
			{sessionStatus?.getComplete() === false && sessionStatus?.getPaused() === true && <>
				<Typography variant="h6">
					PAUSED {sessionId}
				</Typography>
				<PausePresentation sx={{margin: 1}}/>
			</>}
			{sessionStatus?.getComplete() === true && <>
				<Typography variant="h6">
					ENDED {sessionId}
				</Typography>
				<VideocamOff sx={{margin: 1}}/>
			</>}
		</Toolbar>
	</AppBar>
	);
}