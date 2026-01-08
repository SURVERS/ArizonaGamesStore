import React from 'react';
import { Navigate, useLocation } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

function VerifyEmailRoute({ children }) {
  const { pendingVerificationEmail, loading } = useAuth();
  const location = useLocation();

  if (loading) {
    return <div>Checking authorization...</div>;
  }

  if (!pendingVerificationEmail) {
    return <Navigate to="/auth" state={{ from: location }} replace />;
  }

  return children;
}

export default VerifyEmailRoute;
