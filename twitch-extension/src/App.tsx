import { useState, useEffect, useMemo } from 'react';
import './App.css';
import { createCollectionVote, createLocationVote } from './ebs/helpers';
import useWindowDimensions from './hooks/WindowSize';
import _ from 'lodash';

function App() {
  const ext = window?.Twitch?.ext
  const [auth, setAuth] = useState<Twitch.ext.Authorized>();
  const [robotState, setRobotState] = useState();

  useEffect(() => {
    if (!ext) {
      console.error("ext not defined, not running on twitch?");
      return;
    }
    ext.onAuthorized((auth) => {
      console.log("got auth: ", auth);
      setAuth(auth);
    })
    ext.listen("broadcast", (target, contentType, message) => {
      console.log("got broadcast: ", target, contentType, message);
      setRobotState(JSON.parse(message));
    })
  }, [ext])

  const { height, width } = useWindowDimensions();
  const [{x, y}, setCoords] = useState({x: 0, y: 0})

  const locationVoteHandler = (e: React.MouseEvent<HTMLElement>) => {
    const target = e.target as HTMLElement;
    const bounds = target.getBoundingClientRect();
    const x = e.clientX - bounds.left;
    const y = e.clientY - bounds.top;
    setCoords({x: x, y: y});
    const xMod = x / 100.0;
    const yMod = y / 100.0;

    if (!auth) return;
    console.log("voting location ", xMod, yMod);
    fetch(`http://localhost:8080/vote`, {
      method: 'POST',
      body: JSON.stringify(createLocationVote(xMod, yMod)),
      headers: {
        "Authorization": "Bearer " + auth.token,
        "X-Twitch-Extension-Client-Id": auth.clientId,
      },
    }).catch(e => console.error(e))
  }

  const collectionVoteHandler = (i: number) => () => {
    if (!auth) return;
    console.log("voting collection for ", 6-i);
    fetch(`http://localhost:8080/vote`, {
      method: 'POST',
      body: JSON.stringify(createCollectionVote(6-i)),
      headers: {
        "Authorization": "Bearer " + auth.token,
        "X-Twitch-Extension-Client-Id": auth.clientId,
      },
    }).catch(e => console.error(e))
  }

  const stateText = useMemo(() => JSON.stringify(robotState, null, " "), [robotState])

  return (
    <div className="App">
      <header className="App-header">
        {ext && auth ?
        <>
          <textarea readOnly={true} className="debug-text" value={stateText}/>
          <div id="color-vote-area">
            {_.times(5, (i) => <div className="color-option" onClick={collectionVoteHandler(i)}/>)}
          </div>
          <div className='canvas' onClick={locationVoteHandler}>
            <div className='cursor' style={{left: x, top: y}}/>
          </div>
        </>
        : <p style={{color: "#ff00ff"}}>Error: could not get auth from twitch!</p>}
      </header>
    </div>
  );
}

export default App;
