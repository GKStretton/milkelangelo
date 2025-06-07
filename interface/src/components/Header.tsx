import { AppBar, Link, Toolbar, Typography } from "@mui/material";
import { PrecisionManufacturing, Videocam, VideocamOff, Podcasts, PausePresentation } from "@mui/icons-material";
import { useSessionStatus, useStreamStatus } from "../util/hooks";
import { NavLink } from "react-router-dom";

const nav = (
  <div style={{ flexGrow: 1, marginLeft: "3rem", color:"#ffffff"}}>
    <NavLink color="inherit" to="/" style={{marginRight:"1rem"}}>
      <Typography display="inline" variant="h6" marginRight="1rem">Dev</Typography>
    </NavLink>
    <NavLink color="inherit" to="/cleaning" style={{marginRight:"1rem"}}>
      <Typography display="inline" variant="h6">Cleaning</Typography>
    </NavLink>
    <NavLink color="inherit" to="/content" style={{marginRight:"1rem"}}>
      <Typography display="inline" variant="h6">Content</Typography>
    </NavLink>
    <NavLink color="inherit" to="/config" style={{marginRight:"1rem"}}>
      <Typography display="inline" variant="h6">Config</Typography>
    </NavLink>
  </div>
);

export default function Header() {
  const streamStatus = useStreamStatus();
  const sessionStatus = useSessionStatus();

  const sessionId = !sessionStatus
    ? "UNKOWN"
    : sessionStatus.production
    ? sessionStatus.productionId
    : `dev.${sessionStatus.id}`;

  return (
    <AppBar position="relative">
      <Toolbar>
        <PrecisionManufacturing />
        <Typography variant="h6">Milkelangelo</Typography>

        <div style={{ margin: "20px" }}></div>

        {streamStatus?.live ? (
          <>
            <Podcasts sx={{ margin: 1 }} />
            <Typography variant="h6" flexGrow={1}>
              LIVE
            </Typography>
          </>
        ) : (
          <>
            <Podcasts sx={{ margin: 1, opacity: 0.2 }} />
            <Typography variant="h6">
              OFF
            </Typography>
          </>
        )}

        {nav}

        {sessionStatus?.complete === false && sessionStatus?.paused === false && (
          <>
            <Typography variant="h6">RECORDING {sessionId}</Typography>
            <Videocam sx={{ margin: 1 }} />
          </>
        )}
        {sessionStatus?.complete === false && sessionStatus?.paused === true && (
          <>
            <Typography variant="h6">PAUSED {sessionId}</Typography>
            <PausePresentation sx={{ margin: 1 }} />
          </>
        )}
        {sessionStatus?.complete === true && (
          <>
            <Typography variant="h6">ENDED {sessionId}</Typography>
            <VideocamOff sx={{ margin: 1 }} />
          </>
        )}
      </Toolbar>
    </AppBar>
  );
}
