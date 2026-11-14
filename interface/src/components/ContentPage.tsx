import React, { useCallback, useContext, useEffect, useRef, useState } from 'react'
import { useSessionStatus } from '../util/hooks';
import StateReport from './StateReport';
import { Button, MenuItem, Select, Typography } from '@mui/material';
import MqttContext from '../util/mqttContext';
import { TOPIC_GENERATE_CONTENT, TRIGGER_UPLOAD_FROM_CONTENT_PLAN } from '../topics_backend/topics_backend';

const SERVER =
  process.env.REACT_APP_GOO_URL ?? `http://${window.location.hostname}:8089`;

function ContentPage() {
  const { client: c } = useContext(MqttContext);
  const connected = c?.connected;

  const sessionStatus = useSessionStatus()

  const videoRef = useRef<HTMLVideoElement>(null);
  const [videoNames, setVideoNames] = useState<string[]>();
  const [selectedVideo, setSelectedVideo] = useState("");
  const [currentTime, setCurrentTime] = useState(0);
  // changes whenever a still is selected, to refetch the image
  const [stillVersion, setStillVersion] = useState(Date.now());
  const [stillSelected, setStillSelected] = useState(false);
  const [err, setErr] = useState("");

  const getList = useCallback(() => {
    if (!sessionStatus?.id) return;

    setErr("");
    fetch(`${SERVER}/list-session-videos?session_id=${sessionStatus?.id}`)
    .then(async (resp) => {
      if (!resp.ok) {
        throw new Error(await resp.text());
      }
      const data: string[] = await resp.json();
      setVideoNames(data);
      // default to the last recording, which contains the end of the session
      setSelectedVideo(data[data.length - 1] ?? "");
    })
    .catch((e) => {
      console.error(e);
      setErr(e.toString());
    })
  }, [sessionStatus?.id])

  useEffect(() => {
    getList()
  }, [getList])

  const step = (s: number) => {
    const v = videoRef.current;
    if (!v) return;
    v.pause();
    v.currentTime = Math.min(Math.max(v.currentTime + s, 0), v.duration || 0);
  }

  const select = () => {
    if (!sessionStatus?.id || !selectedVideo) return;

    setErr("");
    const t = videoRef.current?.currentTime ?? 0;
    fetch(`${SERVER}/select-still?session_id=${sessionStatus.id}&name=${selectedVideo}&t=${t.toFixed(3)}`)
    .then(async (resp) => {
      if (!resp.ok) {
        throw new Error(await resp.text());
      }
      setStillVersion(Date.now());
    })
    .catch((e) => {
      console.error(e);
      setErr(e.toString());
    })
  }

  const [generateMessage, setGenerateMessage] = useState("");
  const [generateRequested, setGenerateRequested] = useState(false);
  const generateContent = () => {
    if (!connected || !c) {
      setGenerateMessage("mqtt broker not connected");
      return;
    }
    if (!sessionStatus?.id) {
      setGenerateMessage("session id blank");
      return;
    }
    c.publish(TOPIC_GENERATE_CONTENT, `${sessionStatus.id}`)
    setGenerateMessage("requested generation of content");
    setGenerateRequested(true);
  }

  return (
    <div style={{padding: "1rem"}}>
      {err ?
        <Typography variant='h4' color='error'>{err}</Typography>
      : null}
      {!sessionStatus?.id ?
        <Typography>No session status received yet, waiting for goo over mqtt...</Typography>
      : null}
      {videoNames ?
      <>
        <Typography variant="h5">Piece selection session {sessionStatus?.productionId} ({sessionStatus?.id})</Typography>
        {videoNames.length > 1 ?
          <Select size="small" value={selectedVideo} onChange={(e) => setSelectedVideo(e.target.value)}>
            {videoNames.map((n) => <MenuItem key={n} value={n}>{n}</MenuItem>)}
          </Select>
        : null}
        {selectedVideo ?
          <video
            ref={videoRef}
            key={selectedVideo}
            width={600}
            controls
            muted
            preload="auto"
            src={`${SERVER}/get-session-video?session_id=${sessionStatus?.id}&name=${selectedVideo}`}
            onTimeUpdate={(e) => setCurrentTime(e.currentTarget.currentTime)}
            style={{display: "block"}}
          />
        : <div>No top-cam recordings found for this session</div>}
        <div>{selectedVideo} @ {currentTime.toFixed(2)}s</div>
        <Button variant="outlined" onClick={()=>{step(-300)}}>{"- 5m"}</Button>
        <Button variant="outlined" onClick={()=>{step(-30)}}>{"- 30s"}</Button>
        <Button variant="contained" disabled={!selectedVideo} onClick={()=>{select()}}>{"Select"}</Button>
        <Button variant="outlined" onClick={()=>{step(30)}}>{"+ 30s"}</Button>
        <Button variant="outlined" onClick={()=>{step(300)}}>{"+ 5m"}</Button>
        <br/>
        <br/>
        {!connected ? <div>ERROR: not connected to mqtt broker</div>: null}
        {!stillSelected ? <div>Cannot generate because a still has not been selected</div>: null}
        <Button variant="contained" color="success" disabled={!connected || !stillSelected || generateRequested} onClick={()=>{generateContent()}}>Generate Content</Button>
        <div>{generateMessage}</div>
        <br/>
        <hr />
        <br/>
        <Button variant="outlined" color="warning" disabled={!connected} onClick={()=>{c?.publish(TRIGGER_UPLOAD_FROM_CONTENT_PLAN, String(sessionStatus?.id))}}>Trigger upload</Button>
        <br/>
        <br/>
        <div style={{display: stillSelected ? "block" : "none"}}>
          <div>Selected:</div>
          <img
            alt="selected still"
            width={600}
            src={`${SERVER}/get-still?session_id=${sessionStatus?.id}&v=${stillVersion}`}
            onLoad={() => setStillSelected(true)}
            onError={() => setStillSelected(false)}
          />
        </div>
      </>
      : null }
      <br/>
      <br/>
      <hr/>
      <br/>
      <br/>
      <StateReport/>
    </div>
  )
}

export default ContentPage;
