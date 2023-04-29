import React from 'react'
import {
  BrowserRouter as Router,
  Route,
  Routes,
} from 'react-router-dom'
import Home from '../pages/Home'
import SignIn from '../pages/SignIn'
import Error from '../pages/Error'
import SignUp from '../pages/SignUp'

const AppRoutes: React.FC = () => {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/signIn" element={<SignIn />} />
        <Route path="/signUp" element={<SignUp />} />
        <Route path="/error" element={<Error />} />
      </Routes>
    </Router>
  )
}

export default AppRoutes
