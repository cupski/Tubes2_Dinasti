import React, { useState } from 'react';
import './styles.css';
import bfs_title from './assets/bfs-text.png'
import { useNavigate } from 'react-router-dom';

const BFSPage = () => {
    const [startArticle, setStartArticle] = useState('');
    const [targetArticle, setTargetArticle] = useState('');
    const [result, setResult] = useState(null);
    const [isLoading, setIsLoading] = useState(false); // Add loading state

    // const handleSubmit = async (e) => {
    //     e.preventDefault();
    //     setIsLoading(true); // Set loading state to true
    //     try {
    //         const response = await fetch(`http://localhost:8080/shortestpath?start=${startArticle}&target=${targetArticle}`);
    //         const data = await response.json();
    //         setResult(data);
    //         setIsLoading(false); // Set loading state to false after data is fetched
    //         console.log('Data fetched successfully:', data); // Log success message
    //     } catch (error) {
    //         console.error('Error fetching data:', error); // Log error message
    //         setIsLoading(false); // Set loading state to false in case of error
    //     }
    // };

    const handleBFSClick = async (e) => {
        e.preventDefault();
        setIsLoading(true);
        try {
            // Construct full URLs
            const fullStartArticleURL = `https://en.wikipedia.org/wiki/${startArticle}`;
            const fullTargetArticleURL = `https://en.wikipedia.org/wiki/${targetArticle}`;
    
            // Make API request with full URLs
            // const response = await fetch(`http://localhost:8080/shortestpath?start=${encodeURIComponent(fullStartArticleURL)}&target=${encodeURIComponent(fullTargetArticleURL)}`);
            const response = await fetch(`http://localhost:8080/shortestpath?algorithm=bfs&start=${encodeURIComponent(fullStartArticleURL)}&target=${encodeURIComponent(fullTargetArticleURL)}`);
            const data = await response.json();
            setResult(data);
            setIsLoading(false);
            console.log('Data fetched successfully:', data);
        } catch (error) {
            console.error('Error fetching data:', error);
            setIsLoading(false);
        }
    };

    const navigate = useNavigate();

    const handleBack = () => {
        navigate('/');
    };

    return (
        <div className="logic-container">
            <img src={bfs_title} alt="BFS TITLE" className='header-bfs'/>
            <div className="start-container">
                <p style={{ fontFamily: 'Comic Sans MS', fontSize: '16px', marginBottom: '5px' }}>Enter Start Point:</p>
                <input type="text" value={startArticle} onChange={(e) => setStartArticle(e.target.value)} placeholder="Start Point" />
            </div>
            <div className="end-container">
                <p style={{ fontFamily: 'Comic Sans MS', fontSize: '16px', marginBottom: '5px' }}>Enter End Point:</p>
                <input type="text" value={targetArticle} onChange={(e) => setTargetArticle(e.target.value)} placeholder="End Point" />
            </div>
            <div className="search-container">
                <button onClick={handleBFSClick}>Search</button>
            </div>
            <div className="result-container">
                {isLoading ? (
                    <p>Loading...</p>
                ) : result ? (
                    <div>
                        <h2>Result</h2>
                        <ul>
                            {result.path.map((url, index) => (
                                <li key={index}>
                                    <a href={url} target="_blank" rel="noopener noreferrer">{url}</a>
                                </li>
                            ))}
                        </ul>
                        <p>Articles Visited: {result.articlesVisited}</p>
                        <p>Articles Checked: {result.articlesChecked}</p>
                        <p>Execution Time: {result.executionTime} ms</p>
                    </div>
                ) : null}
            </div>
            <div>
                <button className="back-button3" onClick={handleBack}>
                            Back
                </button>
            </div>
        </div>
    );
}

export default BFSPage;
