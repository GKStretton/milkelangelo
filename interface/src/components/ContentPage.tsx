import React, { useCallback, useContext, useEffect, useMemo, useState } from 'react'
import { useSessionStatus } from '../util/hooks';
import StateReport from './StateReport';
import { Button, Typography } from '@mui/material';
import { createIncrementalCompilerHost } from 'typescript';
import MqttContext from '../util/mqttContext';
import { TOPIC_GENERATE_CONTENT, TOPIC_STILLS_GENERATED } from '../topics_backend/topics_backend';

function ContentPage() {
  const { client: c } = useContext(MqttContext);
  const connected = c?.connected;

  const sessionStatus = useSessionStatus()

  const [imageNames, setImageNames] = useState<string[]>();
  const [selectedImgIndex, setSelectedImgIndex] = useState(0);
  const [err, setErr] = useState("");

  const getList = useCallback(() => {
    if (!sessionStatus?.id) return;

    setErr("");
    fetch(`http://depth:8089/list-dslr-post?session_id=${sessionStatus?.id}`)
    .then((resp) => {
      console.log("got resp", resp)
      if (!resp.ok) {
        resp.text().then((t) => {
          console.error(t);
        })
      }
      resp.json().then((data) => {
        setImageNames(data);
      })
    })
    .catch((e) => {
      console.error(e);
      setErr(e.toString());
    })
  }, [sessionStatus?.id])

  useEffect(() => {
    getList()
  }, [getList])

  const increment = (i: number) => {
    const n = (imageNames ?? []).length;
    setSelectedImgIndex((selectedImgIndex + i + n) % n);
  }

  const select = () => {
    if (!imageNames) return;
    if (imageNames[selectedImgIndex] === "selected.jpg") {
      alert("must selected numbered image, not selected.jpg")
      return;
    }

    setErr("");

    fetch(`http://depth:8089/select-dslr-post?session_id=${sessionStatus?.id}&image_name=${imageNames[selectedImgIndex]}`)
    .catch(r => {
      console.error(r);
      setErr(r);
    })
    .finally(() => {
      getList()
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

  const imageSelected = useMemo(() => imageNames?.find((v) => v === "selected.jpg") !== undefined, [imageNames])

  return (
    <div style={{padding: "1rem"}}>
      {err ? 
        <Typography variant='h4' color='error'>{err}</Typography>
      : null}
      {imageNames ?
      <>
        <Typography variant="h5">Piece selection session {sessionStatus?.productionId} ({sessionStatus?.id})</Typography>
        <img
          alt={imageNames[selectedImgIndex]}
          width={600}
          src={`http://depth:8089/get-dslr-post?session_id=${sessionStatus?.id}&image_name=${imageNames[selectedImgIndex]}`}
        />
        <br/>
        <div>{imageNames[selectedImgIndex]}</div>
        <Button variant="outlined" onClick={()=>{increment(-10)}}>{"- 10"}</Button>
        <Button variant="outlined" onClick={()=>{increment(-1)}}>{"- 1"}</Button>
        <Button variant="contained" onClick={()=>{select()}}>{"Select"}</Button>
        <Button variant="outlined" onClick={()=>{increment(1)}}>{"+ 1"}</Button>
        <Button variant="outlined" onClick={()=>{increment(10)}}>{"+ 10"}</Button>
        <br/>
        <br/>
        {!connected ? <div>ERROR: not connected to mqtt broker</div>: null}
        {!imageSelected ? <div>Cannot generate because image has not been selected</div>: null}
        <Button variant="contained" color="warning" disabled={!connected || !imageSelected || generateRequested} onClick={()=>{generateContent()}}>Generate Content</Button>
        <div>{generateMessage}</div>
        <br/>
        <br/>
        {imageSelected ?
        <>
          <div>Selected:</div>
          <img
            alt="selected.jpg"
            width={600}
            src={`http://depth:8089/get-dslr-post?session_id=${sessionStatus?.id}&image_name=selected.jpg`}
          />
        </>
        : null}
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