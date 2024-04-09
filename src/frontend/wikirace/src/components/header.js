import React from 'react';
import './styles.css';
import text from './assets/title.png'; 

const Header = () => {
    return (
        <div className='header-container'>
            <div className='max-width-container'>
                <h1 className='header-title'>Wikirace by Dinasti</h1>
                <img src={text} alt="one piece text" className='header-image'/>
            </div>
        </div>
    );
}

export default Header;