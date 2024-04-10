import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import Header from "./components/header";
import Button from "./components/button";
import IDSPage from "./components/ids-page";

function App() {
  return (
    <Router>
      <div className="app-container">
        <Header />
        <Routes>
          <Route path="/" element={<Button />} />
          <Route path="/ids-page" element={<IDSPage />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;