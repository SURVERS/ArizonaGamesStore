import { useState, useEffect, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import SnowEffect from './SnowEffect';
import BottomNavigation from './BottomNavigation';
import Toast from './Toast';
import ReviewModal from './ReviewModal';
import '../styles/Profile.css';

function Profile() {
  const navigate = useNavigate();
  const { user } = useAuth();
  const [activeTab, setActiveTab] = useState('listings');
  const [listings, setListings] = useState([]);
  const [viewedAds, setViewedAds] = useState([]);
  const [feedbacks, setFeedbacks] = useState([]);
  const [isLoading, setIsLoading] = useState(false);
  const [selectedBackground, setSelectedBackground] = useState(null);
  const [backgroundPreview, setBackgroundPreview] = useState(null);
  const [isUploading, setIsUploading] = useState(false);
  const [toast, setToast] = useState(null);
  const [reviewModalOpen, setReviewModalOpen] = useState(false);
  const [selectedAdForReview, setSelectedAdForReview] = useState(null);
  const fileInputRef = useRef(null);


  useEffect(() => {
    if (activeTab === 'listings') {
      fetchUserListings();
    } else if (activeTab === 'viewed') {
      fetchViewedAds();
    } else if (activeTab === 'reviews') {
      fetchFeedbacks();
    }
  }, [activeTab]);


  useEffect(() => {
    document.title = `Arz Store | ${user?.nickname || 'Профиль'}`;
  }, [user?.nickname]);

  const fetchUserListings = async () => {
    if (!user) return;
    setIsLoading(true);
    try {
      const response = await fetch(`http://localhost:8080/api/listings/user/${user.nickname}`);
      const data = await response.json();

      if (data.success) {
        setListings(data.listings || []);
      }
    } catch (error) {
      console.error('Ошибка загрузки объявлений:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const fetchViewedAds = async () => {
    if (!user) return;
    setIsLoading(true);
    try {
      const response = await fetch('http://localhost:8080/api/viewed-ads', {
        credentials: 'include'
      });
      const data = await response.json();

      if (response.ok) {
        setViewedAds(data.viewed_ads || []);
      }
    } catch (error) {
      console.error('Ошибка загрузки просмотренных объявлений:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const fetchFeedbacks = async () => {
    if (!user) return;
    setIsLoading(true);
    try {
      const response = await fetch(`http://localhost:8080/api/feedback/${user.nickname}`);
      const data = await response.json();

      if (response.ok) {
        setFeedbacks(data.feedbacks || []);
      }
    } catch (error) {
      console.error('Ошибка загрузки отзывов:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const formatPrice = (price) => {
    return price.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ' ');
  };

  const getCurrencySymbol = (currency) => {
    const symbols = {
      'VC': 'VC',
      '$': '$',
      'BTC': '฿',
      'EURO': '€'
    };
    return symbols[currency] || currency;
  };


  const getTypeCost = (type) => {
    switch(type) {
      case 'Сдать в аренду': return '/ час';
      default: return '';
    }
  };

  const getCategoryName = (category) => {
    const categories = {
      'business': 'Бизнесы',
      'accs': 'Аккаунты',
      'house': 'Дома',
      'security': 'ОХРАННИКИ',
      'vehicle': 'Транспорт',
      'others': 'Прочее'
    };
    return categories[category] || category;
  };


  const isUserOnline = (lastSeenAt) => {
    if (!lastSeenAt) return false;
    const now = new Date();
    const lastSeen = new Date(lastSeenAt);
    const diffMs = now - lastSeen;
    const diffMinutes = Math.floor(diffMs / 60000);
    return diffMinutes < 5;
  };


  const getRoleBadge = (role) => {
    const roles = {
      'user': { text: 'Пользователь', color: '#888888' },
      'vip': { text: 'VIP', color: '#FFD700' },
      'premium': { text: 'PREMIUM', color: '#FF69B4' },
      'moderator': { text: 'МОДЕРАТОР', color: '#87CEEB' },
      'developer': { text: 'РАЗРАБОТЧИК', color: '#DF3535' },
      'owner': { text: 'ОСНОВАТЕЛЬ', color: '#DF3535' }
    };
    const normalizedRole = role ? role.toLowerCase() : 'user';
    return roles[normalizedRole] || roles['user'];
  };


  const handleBackgroundClick = () => {
    if (!isUploading && !selectedBackground) {
      fileInputRef.current?.click();
    }
  };


  const handleFileSelect = (e) => {
    const file = e.target.files?.[0];
    if (!file) return;


    if (file.size > 20 * 1024 * 1024) {
      setToast({ message: 'Размер изображения не должен превышать 20 МБ', type: 'error' });
      return;
    }


    if (!file.type.startsWith('image/')) {
      setToast({ message: 'Пожалуйста, выберите файл изображения', type: 'error' });
      return;
    }

    setSelectedBackground(file);


    const reader = new FileReader();
    reader.onloadend = () => {
      setBackgroundPreview(reader.result);
    };
    reader.readAsDataURL(file);
  };


  const handleSaveBackground = async () => {
    if (!selectedBackground) return;

    setIsUploading(true);
    setToast({ message: 'Загружаем фон профиля...', type: 'loading' });

    try {
      const formData = new FormData();
      formData.append('background', selectedBackground);

      const response = await fetch('http://localhost:8080/api/profile/update-background', {
        method: 'POST',
        credentials: 'include',
        body: formData,
      });

      const data = await response.json();

      if (response.ok) {
        setToast({ message: 'Фон профиля успешно обновлен!', type: 'success' });
        setTimeout(() => window.location.reload(), 1500);
      } else {
        setToast({ message: data.error || 'Ошибка при загрузке изображения', type: 'error' });
        setIsUploading(false);
      }
    } catch (error) {
      console.error('Ошибка загрузки фона:', error);
      setToast({ message: 'Ошибка при загрузке изображения', type: 'error' });
      setIsUploading(false);
    }
  };


  const handleDeleteBackground = async () => {
    if (!user?.background_avatar_profile) {
      setToast({ message: 'Фон профиля уже отсутствует', type: 'error' });
      return;
    }

    if (!window.confirm('Вы уверены, что хотите удалить фон профиля?')) {
      return;
    }

    setIsUploading(true);
    setToast({ message: 'Удаляем фон профиля...', type: 'loading' });

    try {
      const response = await fetch('http://localhost:8080/api/profile/delete-background', {
        method: 'DELETE',
        credentials: 'include',
      });

      const data = await response.json();

      if (response.ok) {
        setToast({ message: 'Фон профиля успешно удален!', type: 'success' });
        setTimeout(() => window.location.reload(), 1500);
      } else {
        setToast({ message: data.error || 'Ошибка при удалении фона', type: 'error' });
        setIsUploading(false);
      }
    } catch (error) {
      console.error('Ошибка удаления фона:', error);
      setToast({ message: 'Ошибка при удалении фона', type: 'error' });
      setIsUploading(false);
    }
  };


  const handleCancelBackground = () => {
    setSelectedBackground(null);
    setBackgroundPreview(null);
    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }
  };

  return (
    <div className="profile-page">
      <SnowEffect />

      <div className="profile-content">
        <div className="profile-header">
          <button className="back-button" onClick={() => navigate('/feed')}>
            ← Назад
          </button>
          <h1 className="profile-title">Профиль</h1>
        </div>

        
        <div
          className="profile-info-card"
          onClick={handleBackgroundClick}
          style={{
            backgroundImage: backgroundPreview
              ? `url(${backgroundPreview})`
              : user?.background_avatar_profile
                ? `url(${user.background_avatar_profile})`
                : 'none',
            backgroundSize: 'cover',
            backgroundPosition: 'center',
            cursor: (!isUploading && !selectedBackground) ? 'pointer' : 'default',
            position: 'relative'
          }}
        >
          
          <input
            ref={fileInputRef}
            type="file"
            accept="image/*"
            onChange={handleFileSelect}
            style={{ display: 'none' }}
          />

          
          {user?.background_avatar_profile && !selectedBackground && (
            <button
              className="profile-delete-background-btn"
              onClick={(e) => {
                e.stopPropagation();
                handleDeleteBackground();
              }}
              disabled={isUploading}
              title="Удалить фон профиля"
            >
              🗑️
            </button>
          )}
          

          <button
            className="profile-edit-btn"
            onClick={(e) => {
              e.stopPropagation();
              window.location.href = '/settings';
            }}
            disabled={isUploading}
            title="Редактировать профиль"
          >
            🌣
          </button>

          
          <div className="profile-background-overlay">
            <div className="profile-avatar-section">
              <div className="profile-avatar">
                <img
                  src={user?.avatar || '/src/images/icons/user.png'}
                  alt="Avatar"
                  onError={(e) => {
                    e.target.src = '/src/images/icons/user.png';
                  }}
                />
                <span
                  className={`online-dot ${isUserOnline(user?.last_seen_at) ? 'online' : 'offline'}`}
                ></span>
              </div>
              <div className="profile-details">
                <div className="profile-nickname-row">
                  <span className="profile-nickname">{user?.nickname || 'NickName'} </span>
                </div>
                <div className="profile-description">
                  {user.user_description ? user.user_description : 'Обычный бродяга по сайту Arizona Games Store ^_^'}
                </div>
                <br></br>
                {user?.user_role && (
                  <div
                    className="profile-role-badge"
                    style={{ color: getRoleBadge(user.user_role).color }}
                  >
                    {getRoleBadge(user.user_role).text}
                  </div>
                )}<br></br>
                <div className="profile-stats-row">
                  <div className="profile-rating">
                    <img src="/src/images/icons/star.png" alt="Rating" className="star-icon" />
                    <span>{user?.rating !== undefined ? user.rating.toFixed(1) : '0.0'}</span>
                  </div>
                  <div className="profile-reviews">Отзывов: {user?.reviews_count || 0}</div>
                </div>
              </div>
            </div>
          </div>

          
          {selectedBackground && (
            <div className="profile-background-actions">
              <button
                className="profile-background-save"
                onClick={(e) => {
                  e.stopPropagation();
                  handleSaveBackground();
                }}
                disabled={isUploading}
              >
                {isUploading ? 'ЗАГРУЗКА...' : 'СОХРАНИТЬ ИЗМЕНЕНИЯ'}
              </button>
              <button
                className="profile-background-cancel"
                onClick={(e) => {
                  e.stopPropagation();
                  handleCancelBackground();
                }}
                disabled={isUploading}
              >
                ВЕРНУТЬ НАЗАД
              </button>
            </div>
          )}
        </div>

        
        <div className="profile-tabs">
          <button
            className={`profile-tab ${activeTab === 'listings' ? 'active' : ''}`}
            onClick={() => setActiveTab('listings')}
          >
            Объявления
          </button>
          <button
            className={`profile-tab ${activeTab === 'reviews' ? 'active' : ''}`}
            onClick={() => setActiveTab('reviews')}
          >
            Отзывы
          </button>
          <button
            className={`profile-tab ${activeTab === 'viewed' ? 'active' : ''}`}
            onClick={() => setActiveTab('viewed')}
          >
            Просмотренные
          </button>
        </div>

        
        <div className="profile-tab-content">
          {activeTab === 'listings' && (
            <>
              {isLoading ? (
                <div className="profile-loading">
                  <div className="loading-spinner">⏳</div>
                  <p>Загрузка объявлений...</p>
                </div>
              ) : listings.length === 0 ? (
                <div className="profile-empty">
                  <p>У вас пока нет объявлений</p>
                </div>
              ) : (
                <div className="profile-listings-grid">
                  {listings.map((ad) => (
                    <div key={ad.id} className="profile-listing-card">
                      <div className="profile-listing-image">
                        <img src={ad.image} alt={ad.title} />
                        <div className="profile-listing-id">ID: {ad.id}</div>
                      </div>
                      <div className="profile-listing-content">
                        <div className="profile-listing-price">
                          {ad.currency === 'Договорная' ? 'Договорная' : `${formatPrice(ad.price)} ${getCurrencySymbol(ad.currency)} ${getTypeCost(ad.type)}`}
                        </div>
                        <div className="profile-listing-category">
                          {getCategoryName(ad.category)}
                        </div>
                        <div className="profile-listing-title">{ad.title}</div>
                        <div className="profile-listing-description">{ad.description}</div>
                        <button
                          className="profile-listing-action"
                          onClick={() => navigate(`/manage-ads/${ad.id}`)}
                        >
                          <img src="/src/images/icons/store-icon.png" alt="Action" className="action-icon" />
                          <span>ПЕРЕЙТИ В УПРАВЛЕНИЕ</span>
                        </button>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </>
          )}

          {activeTab === 'reviews' && (
            <>
              {isLoading ? (
                <div className="profile-loading">
                  <div className="loading-spinner">⏳</div>
                  <p>Загрузка отзывов...</p>
                </div>
              ) : feedbacks.length === 0 ? (
                <div className="profile-empty">
                  <p>Отзывов пока нет</p>
                </div>
              ) : (
                <div className="profile-feedbacks-grid">
                  {feedbacks.map((feedback) => (
                    <div key={feedback.id} className="profile-feedback-card">
                      <div className="feedback-header">
                        <img
                          src={feedback.reviewer_avatar || '/src/images/icons/user.png'}
                          alt={feedback.reviewer_nickname}
                          className="feedback-avatar"
                        />
                        <div className="feedback-author-info">
                          <span className="feedback-author">{feedback.reviewer_nickname}</span>
                          <div className="feedback-stars">
                            {[...Array(5)].map((_, i) => (
                              <span key={i} className={`star ${i < feedback.rating ? 'filled' : ''}`}>
                                ★
                              </span>
                            ))}
                          </div>
                        </div>
                      </div>
                      <p className="feedback-text">{feedback.review_text}</p>
                      {feedback.proof_image && (
                        <img
                          src={feedback.proof_image}
                          alt="Доказательство"
                          className="feedback-proof"
                        />
                      )}
                      <div className="feedback-date">
                        {new Date(feedback.created_at).toLocaleDateString('ru-RU')}
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </>
          )}

          {activeTab === 'viewed' && (
            <>
              {isLoading ? (
                <div className="profile-loading">
                  <div className="loading-spinner">⏳</div>
                  <p>Загрузка просмотренных объявлений...</p>
                </div>
              ) : viewedAds.length === 0 ? (
                <div className="profile-empty">
                  <p>Вы еще не просматривали объявления</p>
                </div>
              ) : (
                <div className="profile-listings-grid">
                  {viewedAds.map((item) => {
                    const ad = item.Ad || item;
                    return (
                      <div key={item.id} className="profile-listing-card">
                        <div className="profile-listing-image">
                          <img src={ad.image || ad.Image} alt={ad.title || ad.Title} />
                          <div className="profile-listing-id">ID: {ad.id || ad.ID}</div>
                        </div>
                        <div className="profile-listing-content">
                          <div className="profile-listing-price">
                            {ad.currency === 'Договорная' || ad.Currency === 'Договорная'
                              ? 'Договорная'
                              : `${formatPrice(ad.price || ad.Price)} ${getCurrencySymbol(ad.currency || ad.Currency)} ${getTypeCost(ad.type || ad.Type)}`}
                          </div>
                          <div className="profile-listing-category">
                            {getCategoryName(ad.category || ad.Category)}
                          </div>
                          <div className="profile-listing-title">{ad.title || ad.Title}</div>
                          <div className="profile-listing-description">{ad.description || ad.Description}</div>
                          <div className="viewed-ad-actions">
                            <button
                              className="profile-listing-action review-btn"
                              onClick={() => {
                                setSelectedAdForReview(ad);
                                setReviewModalOpen(true);
                              }}
                            >
                              Оставить отзыв
                            </button>
                          </div>
                        </div>
                      </div>
                    );
                  })}
                </div>
              )}
            </>
          )}
        </div>
      </div>

      <BottomNavigation />

      {toast && <Toast message={toast.message} type={toast.type} onClose={() => setToast(null)} />}

      <ReviewModal
        ad={selectedAdForReview}
        isOpen={reviewModalOpen}
        onClose={() => {
          setReviewModalOpen(false);
          setSelectedAdForReview(null);
        }}
        onSuccess={() => {
          setToast({ message: 'Отзыв успешно отправлен!', type: 'success' });
        }}
      />
    </div>
  );
}

export default Profile;
