import React from 'react';
import { useNavigate } from 'react-router-dom';
import './styles.css';
import chopper from './assets/bfs.png';

const ChopperButton = () => {
    const navigate = useNavigate();

    const handleClick = () => {
        navigate('/bfs-page');
    };
    
    return (
        <button className="chopper-button" onClick={handleClick}>
            <img src={chopper} alt="Button 2" />
            BFS
        </button>
    );
}

export default ChopperButton;
