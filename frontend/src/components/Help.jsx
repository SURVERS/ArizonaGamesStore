import React, { useEffect } from 'react';
import BottomNavigation from './BottomNavigation';
import '../styles/EmptyPage.css';

function Help() {

  useEffect(() => {
    document.title = 'Arz Store | Помощь';
  }, []);

  return (
    <div className="empty-page">
      <div className="empty-content">
        <h1>Помощь</h1>
        <p>Раздел в разработке</p>
      </div>
      <BottomNavigation />
    </div>
  );
}

export default Help;
