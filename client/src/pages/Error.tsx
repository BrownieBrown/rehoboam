import React from 'react'
import './Error.css'

const Error: React.FC = () => {
  return (
    <div className="error-container">
      <div className="error-content">
        <h1 className="error-code">404</h1>
        <h2 className="error-message">Page Not Found</h2>
        <p className="error-description">
          The page you are looking for might have been removed, had
          its name changed, or is temporarily unavailable.
        </p>
        <button
          className="error-back-btn"
          onClick={() => window.history.back()}
        >
          Go Back
        </button>
      </div>
    </div>
  )
}

export default Error
