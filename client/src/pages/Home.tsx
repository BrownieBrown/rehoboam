// src/pages/HomePage.tsx
import React from 'react'
import { Link } from 'react-router-dom'
import './Home.css'

const Home = () => {
  return (
    <div className="homepage-container">
      <div className="title-container">
        <h1 className="homepage-title">Your Name</h1>
        <h2 className="homepage-subtitle">Your Profession</h2>
      </div>
      <div className="navigation">
        <Link to="/about" className="nav-link">
          About
        </Link>
        <Link to="/work" className="nav-link">
          Work
        </Link>
        <Link to="/contact" className="nav-link">
          Contact
        </Link>
      </div>
    </div>
  )
}

export default Home
