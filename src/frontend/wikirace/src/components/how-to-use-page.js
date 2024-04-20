import React from 'react';
import './styles.css';
import island1 from './assets/island1.png';
import island2 from './assets/island2.png';
import island3 from './assets/island3.png';
import island4 from './assets/island4.png';
import island5 from './assets/island5.png';
import step1 from './assets/step1.png';
import step2 from './assets/step2.png';
import step3 from './assets/step3.png';
import step4 from './assets/step4.png';
import step5 from './assets/step5.png';
import line1to2 from './assets/1to2.png';
import line2to3 from './assets/2to3.png';
import line3to4 from './assets/3to4.png';
import line4to5 from './assets/4to5.png';
import ship from './assets/pirate-ship.png';
import title from './assets/how-to-use-title.png';
import { useNavigate } from 'react-router-dom';

const HowToUsePage = () => {
    const navigate = useNavigate();

    const handleBack = () => {
        navigate('/');
    };
    return (
        <div className="container">
            <div>
                <img src={title} alt="Title" className="how-to-use-title"/>
                <button className="back-button" onClick={handleBack}>
                            Back
                </button>
            </div>

            <div>
                <img src={island1} alt="Island 1" className="island1"/>
            </div>
            <div>
                <img src={step1} alt="Step 1" className="step1"/>
            </div>
            <div>
                <img src={line1to2} alt="Line 1 to 2" className="line1to2"/>
            </div>

            <div>
                <img src={island2} alt="Island 2" className="island2"/>
            </div>
            <div>
                <img src={step2} alt="Step 2" className="step2"/>
            </div>
            <div>
                <img src={line2to3} alt="Line 2 to 3" className="line2to3"/>
            </div>

            <div>
                <img src={island3} alt="Island 3" className="island3"/>
            </div>
            <div>
                <img src={step3} alt="Step 3" className="step3"/>
            </div>
            <div>
                <img src={line3to4} alt="Line 3 to 4" className="line3to4"/>
            </div>

            <div>
                <img src={island4} alt="Island 4" className="island4"/>
            </div>
            <div>
                <img src={step4} alt="Step 4" className="step4"/>
            </div>
            <div>
                <img src={line4to5} alt="Line 4 to 5" className="line4to5"/>
            </div>

            <div>
                <img src={island5} alt="Island 5" className="island5"/>
            </div>
            <div>
                <img src={step5} alt="Step 5" className="step5"/>
            </div>

            <div>
                <img src={ship} alt="Ship" className="pirate-ship"/>
            </div>
        </div>
    );
}

export default HowToUsePage;
