import { useState, useEffect } from 'react';
import './App.css';
import { createLocationVote } from './ebs/helpers';
import useWindowDimensions from './hooks/WindowSize';

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

  const clickHandler = (e: React.MouseEvent<HTMLElement>) => {
    const target = e.target as HTMLElement;
    const bounds = target.getBoundingClientRect();
    const x = e.clientX - bounds.left;
    const y = e.clientY - bounds.top;
    setCoords({x: x, y: y});

    if (!auth) return;
    fetch(`http://localhost:8080/vote`, {
      method: 'POST',
      body: JSON.stringify(createLocationVote(x, y)),
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
            {/* <div className="color-option" /> */}
            {/* <div className="color-option" /> */}
          </div>
          <div className='canvas' onClick={clickHandler}>
            <div className='cursor' style={{left: x, top: y}}/>
          </div>
        </>
        : <p style={{color: "#ff00ff"}}>Error: could not get auth from twitch!</p>}
      </header>
    </div>
  );
}

export default App;
