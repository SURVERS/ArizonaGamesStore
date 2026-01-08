import React, { useState } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import '../styles/CategoryNav.css';

const CATEGORIES = [
  { name: 'Аксессуары', path: '/accs', description: 'Арендовать, купить, продать аксессуары' },
  { name: 'Бизнесы', path: '/business', description: 'Назначить заместителя, купить, продать бизнес' },
  { name: 'Дома', path: '/house', description: 'Арендовать, купить, продать дома' },
  { name: 'Охранники', path: '/security', description: 'Арендовать, купить, продать охранника' },
  { name: 'Транспорт', path: '/vehicle', description: 'Арендовать, купить, продать транспорт' },
  { name: 'Реклама / Остальное', path: '/others', description: 'Рекламировать услуги, товары и прочее' },
];

function CategoryNav() {
  const [isOpen, setIsOpen] = useState(true);
  const navigate = useNavigate();
  const location = useLocation();

  const toggleNav = () => {
    setIsOpen(!isOpen);
  };

  const handleCategoryClick = (path) => {
    navigate(path);
  };

  return (
    <div className={`category-nav ${isOpen ? 'open' : 'closed'}`}>
      <button className="category-nav-toggle" onClick={toggleNav}>
        {isOpen ? '‹' : '›'}
      </button>

      {isOpen && (
        <div className="category-nav-content">
          <h3 className="category-nav-title">Категории</h3>
          <div className="category-list">
            {CATEGORIES.map((category) => (
              <button
                key={category.path}
                className={`category-item ${location.pathname === category.path ? 'active' : ''}`}
                onClick={() => handleCategoryClick(category.path)}
                title={category.description}
              >
                <div className="category-content">
                  <span className="category-name">{category.name}</span>
                  <span className="category-description">{category.description}</span>
                </div>
              </button>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}

export default CategoryNav;
