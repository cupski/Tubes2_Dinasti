import React from 'react';
import './styles.css';
import LuffyButton from './luffy-button';
import ChopperButton from './chopper-button';

const Button = () => {
    return (
        <div className="button-container">
            <LuffyButton/>
            <ChopperButton/>
        </div>
    );
}

export default Button;
