import React, { useState, useEffect } from 'react';
import './styles.css';
import bfs_title from './assets/bfs-text.png';
import { useNavigate } from 'react-router-dom';
import { Graph } from 'react-d3-graph';

const BFSPage = () => {
    const [startArticle, setStartArticle] = useState('');
    const [targetArticle, setTargetArticle] = useState('');
    const [result, setResult] = useState(null);
    const [isLoading, setIsLoading] = useState(false);
    const [graphData, setGraphData] = useState(null);
    const [startSuggestions, setStartSuggestions] = useState([]);
    const [targetSuggestions, setTargetSuggestions] = useState([]);

    const fetchSuggestions = async (input, setSuggestions) => {
        try {
            const response = await fetch(
                `https://en.wikipedia.org/w/api.php?action=opensearch&limit=10&format=json&search=${input}&origin=*`
              );
            const data = await response.json();
            const suggestions = data[1] || [];
            setSuggestions(suggestions);
        } catch (error) {
            console.error('Error fetching suggestions:', error);
        }
    };

    useEffect(() => {
        if (startArticle.trim() !== '') {
            fetchSuggestions(startArticle, setStartSuggestions);
        }
    }, [startArticle]);

    useEffect(() => {
        if (targetArticle.trim() !== '') {
            fetchSuggestions(targetArticle, setTargetSuggestions);
        }
    }, [targetArticle]);
    
    const handleBFSClick = async (e) => {
        e.preventDefault();
        setIsLoading(true);
        try {
            // Construct full URLs

            const formattedStartArticle = startArticle.replace(/ /g, '_');
            const formattedTargetArticle = targetArticle.replace(/ /g, '_');

            const fullStartArticleURL = `https://en.wikipedia.org/wiki/${formattedStartArticle}`;
            const fullTargetArticleURL = `https://en.wikipedia.org/wiki/${formattedTargetArticle}`;
    
            // Make API request with full URLs
            const response = await fetch(`http://localhost:8080/shortestpath?algorithm=bfs&start=${encodeURIComponent(fullStartArticleURL)}&target=${encodeURIComponent(fullTargetArticleURL)}`);
            const data = await response.json();

            // Transform data into graph format
            const graphNodes = data.path.map((url) => ({ id: url, label: url }));
            const graphLinks = data.path.slice(0, -1).map((url, index) => ({ source: url, target: data.path[index + 1] }));

            const graphData = {
                nodes: graphNodes,
                links: graphLinks
            };

            setResult(data);
            setGraphData(graphData);
            setIsLoading(false);
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
                <div className="input-select-container">
                    <input
                        type="text"
                        value={startArticle}
                        onChange={(e) => setStartArticle(e.target.value)}
                        placeholder="Enter Start Point"
                    />
                    <select value={startArticle} onChange={(e) => setStartArticle(e.target.value)}>
                        <option value="">Select Start Point</option>
                        {startSuggestions
                            .filter(suggestion => suggestion.toLowerCase().includes(startArticle.toLowerCase()))
                            .map((suggestion, index) => (
                                <option key={index} value={suggestion}>{suggestion}</option>
                            ))}
                    </select>
                </div>
            </div>
            <div className="end-container">
                <p style={{ fontFamily: 'Comic Sans MS', fontSize: '16px', marginBottom: '5px' }}>Enter End Point:</p>
                <div className="input-select-container">
                    <input
                        type="text"
                        value={targetArticle}
                        onChange={(e) => setTargetArticle(e.target.value)}
                        placeholder="Enter End Point"
                    />
                    <select value={targetArticle} onChange={(e) => setTargetArticle(e.target.value)}>
                        <option value="">Select End Point</option>
                        {targetSuggestions
                            .filter(suggestion => suggestion.toLowerCase().includes(targetArticle.toLowerCase()))
                            .map((suggestion, index) => (
                                <option key={index} value={suggestion}>{suggestion}</option>
                            ))}
                    </select>
                </div>
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
                        <ol>
                            {result.path.map((url, index) => (
                                <li key={index}>
                                    <a href={url} target="_blank" rel="noopener noreferrer">{url}</a>
                                </li>
                            ))}
                        </ol>
                        <p>Articles Visited: {result.articlesVisited}</p>
                        <p>Articles Checked: {result.articlesChecked}</p>
                        <p>Execution Time: {result.executionTime} ms</p>
                        {graphData && (
                            <div className="graph-container">
                                <Graph
                                    id="graph-id" // id is mandatory
                                    data={graphData}
                                    config={{
                                        node: {
                                            size: 1000,
                                            highlightStrokeColor: 'blue'
                                        },
                                        link: { highlightColor: 'lightblue' }
                                    }}
                                />
                            </div>
                        )}
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
