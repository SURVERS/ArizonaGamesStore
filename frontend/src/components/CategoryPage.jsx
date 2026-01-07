import React from 'react';
import { useNavigate } from 'react-router-dom';
import SnowEffect from './SnowEffect';
import '../styles/Category.css';

function CategoryPage({ title, description, icon }) {
  const navigate = useNavigate();

  return (
    <div className="category-page">
      <SnowEffect />

      <div className="category-content">
        <button className="back-btn" onClick={() => navigate('/feed')}>
          ← Назад к меню
        </button>

        <div className="category-header">
          <div className="category-icon">{icon}</div>
          <h1 className="category-page-title">{title}</h1>
          <p className="category-page-description">{description}</p>
        </div>

        <div className="coming-soon">
          <div className="coming-soon-icon">🚧</div>
          <h2>Раздел в разработке</h2>
          <p>Скоро здесь появится функционал для работы с {title.toLowerCase()}</p>
        </div>
      </div>
    </div>
  );
}

export default CategoryPage;
