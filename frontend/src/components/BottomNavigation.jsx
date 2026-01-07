import React from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import '../styles/BottomNavigation.css';

function BottomNavigation() {
  const navigate = useNavigate();
  const location = useLocation();


  if (['/auth', '/register', '/verify-email'].includes(location.pathname)) {
    return null;
  }

  const isActive = (path) => location.pathname === path;

  return (
    <div className="bottom-navigation">
      <button
        className={`nav-item ${isActive('/feed') ? 'active' : ''}`}
        onClick={() => navigate('/feed')}
      >
        <img
          src="/src/images/icons/home.png"
          alt="Home"
          className="nav-icon"
        />
      </button>

      <button
        className={`nav-item ${isActive('/settings') ? 'active' : ''}`}
        onClick={() => navigate('/settings')}
      >
        <img
          src="/src/images/icons/settings.png"
          alt="Settings"
          className="nav-icon"
        />
      </button>

      <button
        className={`nav-item ${isActive('/profile') ? 'active' : ''}`}
        onClick={() => navigate('/profile')}
      >
        <img
          src="/src/images/icons/user.png"
          alt="Profile"
          className="nav-icon"
        />
      </button>
    </div>
  );
}

export default BottomNavigation;
