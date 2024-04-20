import React, { useState } from 'react';

function App() {
  const [startArticle, setStartArticle] = useState('');
  const [targetArticle, setTargetArticle] = useState('');
  const [result, setResult] = useState(null);

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const response = await fetch(`http://localhost:8080/shortestpath?start=${startArticle}&target=${targetArticle}`);
      const data = await response.json();
      setResult(data);
    } catch (error) {
      console.error('Error:', error);
    }
  };

  return (
    <div>
      <h1>Shortest Path Finder</h1>
      <form onSubmit={handleSubmit}>
        <label>
          Start Article:
          <input type="text" value={startArticle} onChange={(e) => setStartArticle(e.target.value)} />
        </label>
        <label>
          Target Article:
          <input type="text" value={targetArticle} onChange={(e) => setTargetArticle(e.target.value)} />
        </label>
        <button type="submit">Find Shortest Path</button>
      </form>
      {result && (
        <div>
          <h2>Result</h2>
          <p>Checked articles: {result.checked}</p>
          <p>Search time: {result.search_time_ms} ms</p>
          <p>Path to target article: {result.path.join(' -> ')}</p>
        </div>
      )}
    </div>
  );
}

export default App;
