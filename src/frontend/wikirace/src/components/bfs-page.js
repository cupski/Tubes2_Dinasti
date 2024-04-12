import React, { useState } from 'react';
import './styles.css';

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
            <h1>BFS</h1>
            <div>
                <input type="text" value={startPoint} onChange={(e) => setStartPoint(e.target.value)} placeholder="Start Point" />
            </div>
            <div>
                <input type="text" value={endPoint} onChange={(e) => setEndPoint(e.target.value)} placeholder="End Point" />
            </div>
            <div>
                <button onClick={handleSearch}>Search</button>
            </div>
            <div>
                <p>{result}</p>
            </div>
        </div>
    );
}

export default BFSPage;