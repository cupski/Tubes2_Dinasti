import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import Header from "./components/header";
import Button from "./components/button";
import IDSPage from "./components/ids-page";
import BFSPage from "./components/bfs-page";

function App() {
  return (
    <Router>
      <div className="app-container">
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/ids-page/*" element={<IDSPageWOHeader />} />
          <Route path="/bfs-page/*" element={<BFSPageWOHeader />} />
        </Routes>
      </div>
    </Router>
  );
}

const Home = () => {
  return (
    <>
      <Header />
      <Button />
    </>
  );
};

const IDSPageWOHeader = () => {
  return (
    <div className="page-container">
      <IDSPage />
    </div>
  );
};

const BFSPageWOHeader = () => {
  return (
    <div className="page-container">
      <BFSPage />
    </div>
  );
};

export default App;