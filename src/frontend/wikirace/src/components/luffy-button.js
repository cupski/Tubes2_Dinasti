import React from 'react';
import './styles.css';
import luffy from './assets/ids.png';

const LuffyButton = () => {
    return (
        <button className="luffy-button">
            <img src={luffy} alt="Button 1" />
        </button>
    );
}

export default LuffyButton;