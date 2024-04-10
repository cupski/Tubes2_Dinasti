import React from 'react';
import { useNavigate } from 'react-router-dom';
import './styles.css';
import luffy from './assets/ids.png';

const LuffyButton = () => {
    const navigate = useNavigate();

    const handleClick = () => {
        navigate('/ids-page');
    };

    return (
        <button className="luffy-button" onClick={handleClick}>
            <img src={luffy} alt="Button 1" />
            IDS
        </button>
    );
}

export default LuffyButton;