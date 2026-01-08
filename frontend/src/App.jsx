import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider } from './context/AuthContext';
import Login from './components/Login';
import Register from './components/Register';
import VerifyEmail from './components/VerifyEmail';
import Feed from './components/Feed';
import Accs from './components/Accs';
import Business from './components/Business';
import House from './components/House';
import Security from './components/Security';
import Vehicle from './components/Vehicle';
import Others from './components/Others';
import Help from './components/Help';
import Profile from './components/Profile';
import Settings from './components/Settings';
import ManageAds from './components/ManageAds';
import ProtectedRoute from './components/ProtectedRoute';
import PublicRoute from './components/PublicRoute';
import VerifyEmailRoute from './components/VerifyEmailRoute';

function App() {
  return (
    <Router>
      <AuthProvider>
        <Routes>
          <Route path="/" element={<Navigate to="/feed" replace />} />
          <Route
            path="/auth"
            element={
              <PublicRoute>
                <Login />
              </PublicRoute>
            }
          />
          <Route
            path="/register"
            element={
              <PublicRoute>
                <Register />
              </PublicRoute>
            }
          />
          <Route
            path="/verify-email"
            element={
              <VerifyEmailRoute>
                <VerifyEmail />
              </VerifyEmailRoute>
            }
          />
          <Route
            path="/feed"
            element={
              <ProtectedRoute>
                <Feed />
              </ProtectedRoute>
            }
          />
          <Route
            path="/accs"
            element={
              <ProtectedRoute>
                <Accs />
              </ProtectedRoute>
            }
          />
          <Route
            path="/business"
            element={
              <ProtectedRoute>
                <Business />
              </ProtectedRoute>
            }
          />
          <Route
            path="/house"
            element={
              <ProtectedRoute>
                <House />
              </ProtectedRoute>
            }
          />
          <Route
            path="/security"
            element={
              <ProtectedRoute>
                <Security />
              </ProtectedRoute>
            }
          />
          <Route
            path="/vehicle"
            element={
              <ProtectedRoute>
                <Vehicle />
              </ProtectedRoute>
            }
          />
          <Route
            path="/others"
            element={
              <ProtectedRoute>
                <Others />
              </ProtectedRoute>
            }
          />
          <Route
            path="/help"
            element={
              <ProtectedRoute>
                <Help />
              </ProtectedRoute>
            }
          />
          <Route
            path="/profile"
            element={
              <ProtectedRoute>
                <Profile />
              </ProtectedRoute>
            }
          />
          <Route
            path="/settings"
            element={
              <ProtectedRoute>
                <Settings />
              </ProtectedRoute>
            }
          />
          <Route
            path="/manage-ads/:id"
            element={
              <ProtectedRoute>
                <ManageAds />
              </ProtectedRoute>
            }
          />
          <Route path="*" element={<Navigate to="/auth" replace />} />
        </Routes>
      </AuthProvider>
    </Router>
  );
}

export default App;
