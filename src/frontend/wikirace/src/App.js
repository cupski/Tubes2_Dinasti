import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import Header from "./components/header";
import Button from "./components/button";
import IDSPage from "./components/ids-page";

function App() {
  return (
    <Router>
      <div className="app-container">
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/ids-page/*" element={<IDSPageWOHeader />} />
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
    <>
      <IDSPage />
    </>
  );
};

export default App;