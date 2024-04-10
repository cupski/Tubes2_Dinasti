import React from 'react';
import './styles.css';
import luffy from './assets/ids.png';
import chopper from './assets/bfs.png';

const Button = () => {
    return (
        <div className="button-container">
            <button className="luffy-button">
                <img src={luffy} alt="Button 1" />
            </button>

            <button className="chopper-button">
                <img src={chopper} alt="Button 2" />
            </button>
        </div>
    );
}

export default Button;
