import React from 'react';
import './styles.css';
import chopper from './assets/chopper.png';
import luffy from './assets/luffy.png';
import zoro from './assets/zoro.png';
import title from './assets/about-us-title.png';
import { useNavigate } from 'react-router-dom';

const AboutUsPage = () => {
    const navigate = useNavigate();

    const handleBack = () => {
        navigate('/');
    };

    return (
        <div className="about-us-container">
            <div>
                <img src={title} alt="Title" className="about-us-title"/>
            </div>
            <div className="intro">
                <h1>Welcome to Wikirace by Dinasti!</h1>
                <p>Wikirace by Dinasti merupakan sebuah website Wikirace dengan algoritma BFS dan IDS. Ayo temukan rute tersingkat dari suatu keyword ke keyword lain!</p>
            </div>
            <div className="character-section">
                <h2>Meet Our Crew</h2>
                <div className="character-card">
                    <img src={chopper} alt="Tony Tony Chopper" className="chopper-img" />
                    <h3>Denise Felicia Tiowanni</h3>
                    <p>NIM: 13522013</p>
                </div>
                <div className="character-card">
                    <img src={zoro} alt="Roronoa Zoro" className="zoro-img" />
                    <h3>Muhammad Yusuf Rafi</h3>
                    <p>NIM: 13522009</p>
                </div>
                <div className="character-card">
                    <img src={luffy} alt="Monkey D. Luffy" className="luffy-img"/>
                    <h3>Rafii Ahmad Fahreza</h3>
                    <p>NIM: 10023570</p>
                </div>
            </div>
            <div>
                <button className="back-button2" onClick={handleBack}>
                            Back
                </button>
            </div>
        </div>
    );
}

export default AboutUsPage;
