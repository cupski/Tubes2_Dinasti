import React, { useState, useEffect } from 'react';
import './styles.css';
import bfs_title from './assets/bfs-text.png';
import { useNavigate } from 'react-router-dom';
import Graph from 'react-vis-network-graph';
import chopperRunning from './assets/crying-chopper-while-running.gif'

const BFSPage = () => {
    const [startArticle, setStartArticle] = useState('');
    const [targetArticle, setTargetArticle] = useState('');
    const [result, setResult] = useState(null);
    const [isLoading, setIsLoading] = useState(false);
    const [loadingMessage, setLoadingMessage] = useState('');
    const [showMessage, setShowMessage] = useState(true);
    const [graphData, setGraphData] = useState(null);
    const [startSuggestions, setStartSuggestions] = useState([]);
    const [targetSuggestions, setTargetSuggestions] = useState([]);
    const [clickedEdge, setClickedEdge] = useState(null);
    const [error, setError] = useState(null);

    useEffect(() => {
        let interval;
        if (isLoading) {
            let messages = [
                'Tunggu ya... rutemu lagi dicari Chopper nih!',
                'Duh.. kayaknya agak jauh...',
                'Hmm.. kamu masih sabar kannn?'
            ];
            let index = 0;
            setLoadingMessage(messages[index]);

            interval = setInterval(() => {
                setShowMessage(false); // mulai pudar
                setTimeout(() => {
                    index = (index + 1) % messages.length; // Ganti index pesan
                    setLoadingMessage(messages[index]); // Setel pesan baru
                    setShowMessage(true); // Munculkan kembali pesan
                }, 500); // Waktu memudar
            }, 12000);
        }
        return () => {
            clearInterval(interval);
        };
    }, [isLoading]);
    
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
        setError(null);
        try {
            const formattedStartArticle = startArticle.replace(/ /g, '_');
            const formattedTargetArticle = targetArticle.replace(/ /g, '_');
            const fullStartArticleURL = `https://en.wikipedia.org/wiki/${formattedStartArticle}`;
            const fullTargetArticleURL = `https://en.wikipedia.org/wiki/${formattedTargetArticle}`;
    
            const response = await fetch(`http://localhost:8080/shortestpath?algorithm=bfs&start=${encodeURIComponent(fullStartArticleURL)}&target=${encodeURIComponent(fullTargetArticleURL)}`);
            if (!response.ok) {
                throw new Error('Articles invalid or no route found!');
            }
            const data = await response.json();

            // Check if result is empty
            if (!data.path || data.path.length === 0) {
                throw new Error('No path found');
            }

            const nodes = data.path.map((url, index) => {
                const pageName = url.split('/').pop(); // ambil nama aja
                return {
                    id: index,
                    label: pageName,
                    shape: 'star',
                    color: index === 0 ? 'red' : index === data.path.length - 1 ? 'rgba(133, 225, 4, 0.946)' : undefined // Set color for start and end points
                };
            });
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
            setError(error.message);
            setIsLoading(false);
        }
    };

    const handleEdgeClick = (event) => {
        setClickedEdge(event.edges[0]);
    };

    const navigate = useNavigate();

    const handleBack = () => {
        navigate('/');
    };

    const handleNodeClick = (event) => {
        const nodeId = event.nodes[0]; // Ambil id node yang diklik
        const clickedNode = graphData.nodes.find(node => node.id === nodeId); // Cari node yang sesuai dengan id
        if (clickedNode) {
            const wikiLink = `https://en.wikipedia.org/wiki/${clickedNode.label}`;
            window.open(wikiLink, '_blank'); // Buka link Wikipedia pada tab baru
        }
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
                    <div>
                        <img src={chopperRunning} alt="Chopper running" className="loading-gif"/>
                        <p className={`loading-message ${showMessage ? 'fade-in' : 'fade-out'}`}>{loadingMessage}</p>
                    </div>
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
                                    selectEdge: handleEdgeClick,
                                    selectNode: handleNodeClick
                                }}
                            />
                        </div>
                    </div>
                ) : null}
                {error && (
                    <div style={{display: 'flex', flexDirection: 'column', justifyContent: 'center', alignItems: 'center'}}>
                        <p>{error}</p>
                        <button className="close-button" onClick={() => setError(null)}>Close</button>
                    </div>
                )}
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
