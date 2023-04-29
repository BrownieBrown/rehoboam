// src/pages/SignIn.tsx
import React, { FormEvent, useState } from 'react'
import axios from 'axios'
import { useNavigate } from 'react-router-dom'
import './SignIn.css'

const SignIn: React.FC = () => {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const navigate = useNavigate()

  const submitHandler = async (e: FormEvent) => {
    e.preventDefault()

    try {
      const response = await axios.post(
        'http://localhost:8080/api/v1/auth/signIn',
        {
          email,
          password,
        }
      )

      if (response.status === 200) {
        console.log('Authentication successful')
        // Handle successful authentication, e.g., save the token or update the state.

        // Redirect to the Homepage
        navigate('/')
      }
    } catch (error) {
      console.error('Authentication failed')
      // Handle authentication failure, e.g., show an error message.
      navigate('/error')
    }
  }

  return (
    <div className="signin-container">
      {/* ... */}
      <form className="signin-form" onSubmit={submitHandler}>
        <input
          type="email"
          className="signin-input"
          placeholder="Email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
        />
        <input
          type="password"
          className="signin-input"
          placeholder="Password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
        />
        <button type="submit" className="signin-button">
          Sign In
        </button>
      </form>
    </div>
  )
}

export default SignIn
