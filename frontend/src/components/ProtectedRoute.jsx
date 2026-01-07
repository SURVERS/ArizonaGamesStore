import React from 'react';
import { Navigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

function ProtectedRoute({ children }) {
  const { user, loading } = useAuth();

  if (loading) {
    return <div>Загрузка...</div>;
  }

  if (!user) {
    return <Navigate to="/auth" replace />;
  }

  return children;
}

export default ProtectedRoute;
