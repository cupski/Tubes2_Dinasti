import React from 'react';
import './styles.css';
import text from './assets/title.png';
import title from './assets/wikirace2.png';

const Header = () => {
    return (
        <div className='header-container'>
            <div className='max-width-container'>
            <img src={title} alt="Wikirace by Dinasti" className='header-wikirace'/>
                <img src={text} alt="one piece text" className='header-image'/>
            </div>
        </div>
    );
}

export default Header;