import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import Header from "./components/header";
import Button from "./components/button";
import IDSPage from "./components/ids-page";
import BFSPage from "./components/bfs-page";
import HowToUsePage from "./components/how-to-use-page";
import AboutUsPage from "./components/about-us-page";

function App() {
  return (
    <Router>
      <div className="app-container">
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/ids-page/*" element={<IDSPageWOHeader />} />
          <Route path="/bfs-page/*" element={<BFSPageWOHeader />} />
          <Route path="/how-to-use-page/*" element={<HowToUse />} />
          <Route path="/about-us-page/*" element={<AboutUs />} />
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

const HowToUse = () => {
  return (
    <div className="how-to-use-container">
      <HowToUsePage />
    </div>
  );
};

const AboutUs = () => {
  return (
    <div className="page-container">
      <AboutUsPage />
    </div>
  );
};

export default App;