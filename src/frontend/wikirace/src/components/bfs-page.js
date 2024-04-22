import React, { useState, useEffect } from 'react';
import './styles.css';
import bfs_title from './assets/bfs-text.png';
import { useNavigate } from 'react-router-dom';
import Graph from 'react-vis-network-graph';

const BFSPage = () => {
    const [startArticle, setStartArticle] = useState('');
    const [targetArticle, setTargetArticle] = useState('');
    const [result, setResult] = useState(null);
    const [isLoading, setIsLoading] = useState(false);
    const [graphData, setGraphData] = useState(null);
    const [startSuggestions, setStartSuggestions] = useState([]);
    const [targetSuggestions, setTargetSuggestions] = useState([]);
    const [clickedEdge, setClickedEdge] = useState(null);

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
            const formattedStartArticle = startArticle.replace(/ /g, '_');
            const formattedTargetArticle = targetArticle.replace(/ /g, '_');
            const fullStartArticleURL = `https://en.wikipedia.org/wiki/${formattedStartArticle}`;
            const fullTargetArticleURL = `https://en.wikipedia.org/wiki/${formattedTargetArticle}`;
    
            const response = await fetch(`http://localhost:8080/shortestpath?algorithm=bfs&start=${encodeURIComponent(fullStartArticleURL)}&target=${encodeURIComponent(fullTargetArticleURL)}`);
            const data = await response.json();

            const nodes = data.path.map((url, index) => ({
                id: index,
                label: url,
                shape: 'star',
                color: index === 0 ? 'red' : index === data.path.length - 1 ? 'green' : undefined // Set color for start and end points
            }));
            const edges = [];
            for (let i = 0; i < nodes.length - 1; i++) {
                edges.push({ id: `edge${i}`, from: i, to: i + 1, arrows: "to" });
            }
            const graph = { nodes, edges };
            setGraphData(graph);

            setResult(data);
            setIsLoading(false);
        } catch (error) {
            console.error('Error fetching data:', error);
            setIsLoading(false);
        }
    };

    const handleEdgeClick = (event) => {
        setClickedEdge(event.edges[0]); // Store the clicked edge id
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
                        list="startSuggestions"
                        placeholder="Enter Start Point"
                    />
                    <datalist id="startSuggestions">
                        {startSuggestions.map((suggestion, index) => (
                            <option key={index} value={suggestion} />
                        ))}
                    </datalist>
                </div>
            </div>
            <div className="end-container">
                <p style={{ fontFamily: 'Comic Sans MS', fontSize: '16px', marginBottom: '5px' }}>Enter End Point:</p>
                <div className="input-select-container">
                    <input
                        type="text"
                        value={targetArticle}
                        onChange={(e) => setTargetArticle(e.target.value)}
                        list="targetSuggestions"
                        placeholder="Enter End Point"
                    />
                    <datalist id="targetSuggestions">
                        {targetSuggestions.map((suggestion, index) => (
                            <option key={index} value={suggestion} />
                        ))}
                    </datalist>
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
                        <div style={{ height: '400px' }}>
                            <Graph
                                graph={graphData}
                                options={{
                                    nodes: {
                                        shape: 'star',
                                        size: 20, // Set the size of the nodes
                                        font: {
                                            color: 'white' // Set the font color of the labels
                                        }
                                    },
                                    edges: {
                                        font: {
                                            align: 'horizontal'
                                        },
                                        color: {
                                            color: clickedEdge ? 'white' : 'white', // Set the color of the edges
                                            highlight: clickedEdge ? 'yellow' : 'white' // Set the color of the edges when clicked
                                        }
                                    }
                                }}
                                events={{
                                    selectEdge: handleEdgeClick // Handle edge click event
                                }}
                            />
                        </div>
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
