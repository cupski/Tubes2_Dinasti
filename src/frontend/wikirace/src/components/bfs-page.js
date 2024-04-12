import React, { useState } from 'react';
import './styles.css';
import chopper_img from './assets/chopper.png'
import bfs_title from './assets/bfs-text.png'

const BFSPage = () => {
    const [startPoint, setStartPoint] = useState("");
    const [endPoint, setEndPoint] = useState("");
    const [result, setResult] = useState("");

    const handleSearch = () => {
        // seach logic disini, masih bingung.
        setResult(`Results for BFS with start point "${startPoint}" and end point "${endPoint}"`);
    };

    return (
        <div className="logic-container">
            <img src={bfs_title} alt="BFS TITLE" className='header-bfs'/>
            <div className="start-container">
                <p style={{ fontFamily: 'Comic Sans MS', fontSize: '16px', marginBottom: '5px' }}>Enter Start Point:</p>
                <input type="text" value={startPoint} onChange={(e) => setStartPoint(e.target.value)} placeholder="Start Point" />
            </div>
            <div className="end-container">
                <p style={{ fontFamily: 'Comic Sans MS', fontSize: '16px', marginBottom: '5px' }}>Enter End Point:</p>
                <input type="text" value={endPoint} onChange={(e) => setEndPoint(e.target.value)} placeholder="End Point" />
            </div>
            <div className="search-container">
                <button onClick={handleSearch}>Search</button>
            </div>
            <div>
                <p>{result}</p>
            </div>
            <img src={chopper_img} alt="Chopper" className="bottom-corner-image" />
        </div>
    );
}

export default BFSPage;