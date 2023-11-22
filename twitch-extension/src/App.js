import React, { useState } from 'react';
import './App.css';
import useWindowDimensions from './hooks/WindowSize';

function App() {
  const { height, width } = useWindowDimensions();
  const [{x, y}, setCoords] = useState({x: 0, y: 0})

  const clickHandler = (e) => {
    const bounds = e.target.getBoundingClientRect();
    const x = e.clientX - bounds.left;
    const y = e.clientY - bounds.top;
    setCoords({x: x, y: y});
  }
  return (
    <div className="App">
      <header className="App-header">
        <div className="debug-text">{`${width}, ${height}. chat: ${window['Twitch']?.ext.features.isChatEnabled}`}</div>
        <div id="color-vote-area">
          {/* <div className="color-option" /> */}
          {/* <div className="color-option" /> */}
        </div>
        <div className='canvas' onClick={clickHandler}>
          <div className='cursor' style={{left: x, top: y}}/>
        </div>
      </header>
    </div>
  );
}

export default App;
