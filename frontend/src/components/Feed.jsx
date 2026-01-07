import React, { useState, useRef, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import StarRating from './StarRating';
import BottomNavigation from './BottomNavigation';
import '../styles/Feed.css';

function Feed() {
  const [showProfileMenu, setShowProfileMenu] = useState(false);
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const menuRef = useRef(null);


  const [hotAds, setHotAds] = useState([]);
  const [hotAdsLoading, setHotAdsLoading] = useState(false);
  const [hotAdsOffset, setHotAdsOffset] = useState(0);
  const [hasMoreHotAds, setHasMoreHotAds] = useState(true);
  const hotSectionRef = useRef(null);

  const categories = [
    {
      title: 'Аксессуары',
      description: 'В этом разделе игроки могут арендовать, купить, продать аксессуары',
      path: '/accs',
      image: '/src/images/block_items/accs.png'
    },
    {
      title: 'Бизнесы',
      description: 'В этом разделе игроки могут назначить заместителя, купить, продать бизнес',
      path: '/business',
      image: '/src/images/block_items/business.png'
    },
    {
      title: 'Дома',
      description: 'В этом разделе игроки могут арендовать, купить, продать дома',
      path: '/house',
      image: '/src/images/block_items/house.png'
    },
    {
      title: 'Охранники',
      description: 'В этом разделе игроки могут арендовать, купить, продать охранника',
      path: '/security',
      image: '/src/images/block_items/security.png'
    },
    {
      title: 'Транспорт',
      description: 'В этом разделе игроки могут арендовать, купить, продать транспорт',
      path: '/vehicle',
      image: '/src/images/block_items/vehicle.png'
    },
    {
      title: 'Реклама / Остальное',
      description: 'В этом разделе игроки могут рекламировать свои услуги/товары и прочее, купить или продать у игроков',
      path: '/others',
      image: '/src/images/block_items/others.png'
    }
  ];

  useEffect(() => {
    function handleClickOutside(event) {
      if (menuRef.current && !menuRef.current.contains(event.target)) {
        setShowProfileMenu(false);
      }
    }

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);


  useEffect(() => {
    document.title = 'Arz Store | Лента';
  }, []);


  const fetchHotAds = async (offset = 0) => {
    if (hotAdsLoading || (!hasMoreHotAds && offset > 0)) return;

    setHotAdsLoading(true);
    try {
      const response = await fetch(`http://localhost:8080/api/ads/random?limit=15&offset=${offset}`);
      if (!response.ok) throw new Error('Ошибка загрузки объявлений');

      const data = await response.json();
      const newAds = data.ads || [];

      if (offset === 0) {
        setHotAds(newAds);
      } else {
        setHotAds(prev => [...prev, ...newAds]);
      }

      setHotAdsOffset(offset + 15);
      setHasMoreHotAds(newAds.length === 15);
    } catch (error) {
      console.error('Ошибка загрузки горячих объявлений:', error);
    } finally {
      setHotAdsLoading(false);
    }
  };


  useEffect(() => {
    fetchHotAds(0);
  }, []);


  useEffect(() => {
    if (!hotSectionRef.current || hotAds.length === 0) return;

    const observer = new IntersectionObserver(
      (entries) => {
        const lastEntry = entries[0];
        if (lastEntry.isIntersecting && hasMoreHotAds && !hotAdsLoading) {
          fetchHotAds(hotAdsOffset);
        }
      },
      { threshold: 0.5 }
    );


    const cards = hotSectionRef.current.querySelectorAll('.hot-ad-card');
    if (cards.length >= 14) {
      observer.observe(cards[14]);
    }

    return () => observer.disconnect();
  }, [hotAds, hasMoreHotAds, hotAdsLoading, hotAdsOffset]);

  const handleLogout = async () => {
    await logout();
    navigate('/auth');
  };

  return (
    <div className="feed-container">
      <header className="feed-header">
        <div className="header-content">
          <img src="/src/images/logo_arz.png" alt="Arizona Games Store" className="header-logo" />

          <div className="profile-section" ref={menuRef}>
            <div className="header-user-info" onClick={() => setShowProfileMenu(!showProfileMenu)}>
              <span className="header-nickname">{user?.nickname || 'User'}</span>
              <img
                src={user?.avatar || 'https://storage.yandexcloud.net/fotora.ru/uploads/2b0c131e8cfe54b1.jpeg'}
                alt="Avatar"
                className="header-avatar"
              />
            </div>

            {showProfileMenu && (
              <div className="profile-dropdown">
                <div className="profile-header">
                  <img
                    src={user?.avatar && user.avatar.trim() !== '' ? user.avatar : 'https://storage.yandexcloud.net/fotora.ru/uploads/2b0c131e8cfe54b1.jpeg'}
                    alt="Avatar"
                    className="profile-avatars"
                  />
                  <div className="profile-info">
                    <div className="profile-nickname">{user?.nickname || 'User'}</div>
                    <StarRating rating={user?.rating || 0} />
                  </div>
                </div>

                <div className="profile-actions">
                  <button className="profile-btn" onClick={() => navigate('/settings')}>
                    ⚙️ Настройки
                  </button>
                  <button className="profile-btn" onClick={() => navigate('/rules')}>
                    📋 Правила использования сайта
                  </button>
                  <button className="profile-btn logout-btn" onClick={handleLogout}>
                    🚪 Выйти с аккаунта
                  </button>
                </div>
              </div>
            )}
          </div>
        </div>
      </header>

      <main className="feed-main">
        <div className="feed-intro">
          <h1 className="feed-title">ARIZONA GAMES STORE</h1>
          <p className="feed-description">
            Добро пожаловать в игровой магазин Arizona Role Play! Здесь вы можете безопасно покупать,
            продавать и арендовать игровое имущество. Выберите категорию ниже и начните свой путь к успеху!
          </p>
        </div>

        <div className="categories-grid">
          {categories.map((category, index) => (
            <div
              key={index}
              className="category-card"
              onClick={() => navigate(category.path)}
            >
              <div className="category-image-wrapper">
                <img src={category.image} alt={category.title} className="category-image" />
              </div>
              <div className="category-content">
                <h3 className="category-title">{category.title}</h3>
                <p className="category-description">{category.description}</p>
              </div>
            </div>
          ))}
        </div>

        
        {hotAds.length > 0 && (
          <div className="hot-ads-section">
            <div className="hot-ads-header">
              <h2 className="hot-ads-title">🔥 ГОРЯЧЕЕ</h2>
              <p className="hot-ads-subtitle">Популярные объявления из разных категорий</p>
            </div>

            <div className="hot-ads-grid" ref={hotSectionRef}>
              {hotAds.map((ad, index) => (
                <div key={ad.ID || index} className="hot-ad-card">
                  <div className="hot-badge">🔥</div>
                  <div className="hot-ad-image-wrapper">
                    <img
                      src={ad.image || 'https://via.placeholder.com/300x200'}
                      alt={ad.title}
                      className="hot-ad-image"
                    />
                  </div>
                  <div className="hot-ad-content">
                    <div className="hot-ad-header">
                      <h3 className="hot-ad-title">{ad.title}</h3>
                      <span className="hot-ad-category">{ad.category}</span>
                    </div>
                    <p className="hot-ad-description">{ad.description}</p>
                    <div className="hot-ad-footer">
                      <div className="hot-ad-author">
                        <img
                          src={ad.author_avatar || 'https://storage.yandexcloud.net/fotora.ru/uploads/2b0c131e8cfe54b1.jpeg'}
                          alt={ad.nickname}
                          className="hot-ad-author-avatar"
                        />
                        <div className="hot-ad-author-info">
                          <span className="hot-ad-author-name">{ad.nickname}</span>
                          <div className="hot-ad-author-rating">
                            ⭐ {ad.author_rating?.toFixed(1) || '5.0'}
                          </div>
                        </div>
                      </div>
                      <div className="hot-ad-price">
                        {ad.price ? (
                          <>
                            <span className="hot-ad-price-value">{ad.price.toLocaleString()}</span>
                            <span className="hot-ad-price-currency">{ad.currency || '$'}</span>
                          </>
                        ) : (
                          <span className="hot-ad-price-free">Договорная</span>
                        )}
                      </div>
                    </div>
                  </div>
                </div>
              ))}
            </div>

            {hotAdsLoading && (
              <div className="hot-ads-loading">
                <div className="loading-spinner"></div>
                <p>Загрузка объявлений...</p>
              </div>
            )}
          </div>
        )}
      </main>

      <footer className="feed-footer">
        <div className="footer-content">
          <div className="footer-item">
            <span className="footer-icon">📱</span>
            <a href="https://t.me/survers_team" target="_blank" rel="noopener noreferrer">
              Telegram: @survers_team
            </a>
          </div>
          <div className="footer-item">
            <span className="footer-icon">✉️</span>
            <a href="mailto:arizonagamesstore@rambler.ru">
              E-Mail: arizonagamesstore@rambler.ru
            </a>
          </div>
        </div>
        <div className="footer-copyright">
          © 2025 Arizona Games Store. Все права защищены.
        </div>
      </footer>

      <BottomNavigation />
    </div>
  );
}

export default Feed;
