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

const HowToUsePage = () => {
    return (
        <div className="container">
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
        </div>
    );
}

export default HowToUsePage;
