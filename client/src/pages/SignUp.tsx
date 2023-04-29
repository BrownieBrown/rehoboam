// src/pages/SignUp.tsx
import axios from 'axios'
import React, { FormEvent, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import './SignUp.css'

const SignUp: React.FC = () => {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const navigate = useNavigate()

  // src/pages/SignUp.tsx
  const submitHandler = async (e: FormEvent) => {
    e.preventDefault()

    if (password !== confirmPassword) {
      // Show an error message or handle the password mismatch
      console.error('Passwords do not match')
      return
    }

    try {
      // Send a request to your Go backend
      const response = await axios.post(
        'http://localhost:8080/api/v1/auth/signUp',
        {
          email,
          password,
        }
      )

      if (response.status === 201) {
        console.log('Sign up successful')
        // Redirect to another page after successful sign-up, e.g., the Homepage
        navigate('/')
      }
    } catch (error) {
      console.error('Sign up failed')
      // Handle sign-up failure, e.g., show an error message or redirect to the error page
      navigate('/error')
    }
  }

  return (
    <div className="signup-container">
      <h1 className="signup-title">Sign Up</h1>
      <form className="signup-form" onSubmit={submitHandler}>
        <input
          type="email"
          className="signup-input"
          placeholder="Email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
        />
        <input
          type="password"
          className="signup-input"
          placeholder="Password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
        />
        <input
          type="password"
          className="signup-input"
          placeholder="Confirm Password"
          value={confirmPassword}
          onChange={(e) => setConfirmPassword(e.target.value)}
          required
        />
        <button type="submit" className="signup-button">
          Sign Up
        </button>
      </form>
    </div>
  )
}

export default SignUp
