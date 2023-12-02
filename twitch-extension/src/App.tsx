import { useState, useEffect } from 'react';
import './App.css';
import { createCollectionVote, createLocationVote } from './ebs/helpers';
import useWindowDimensions from './hooks/WindowSize';
import _ from 'lodash';

function App() {
  const ext = window?.Twitch?.ext
  const [auth, setAuth] = useState<Twitch.ext.Authorized>();

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

    if (!auth) return;
    console.log("voting location ", x, y);
    fetch(`http://localhost:8080/vote`, {
      method: 'POST',
      body: JSON.stringify(createLocationVote(x, y)),
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

  return (
    <div className="App">
      <header className="App-header">
        {ext && auth ?
        <>
          <div className="debug-text">{`${width}, ${height}. chat: ${ext.features.isChatEnabled}`}</div>
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
