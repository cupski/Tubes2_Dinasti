import React from 'react';
import './styles.css';
import chopper from './assets/bfs.png';

const ChopperButton = () => {
    return (
        <button className="chopper-button">
            <img src={chopper} alt="Button 2" />
            BFS
        </button>
    );
}

export default ChopperButton;
