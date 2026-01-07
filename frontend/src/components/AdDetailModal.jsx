import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import StarRating from './StarRating';
import Toast from './Toast';
import '../styles/AdDetailModal.css';

const AdDetailModal = ({ ad, isOpen, onClose, currentUser }) => {
  const [toast, setToast] = useState(null);
  const navigate = useNavigate();


  useEffect(() => {
    const adId = ad?.id || ad?.ID;
    if (isOpen && ad && adId) {

      fetch(`http://localhost:8080/api/ads/${adId}/view`, {
        method: 'POST'
      }).catch(error => {
        console.error('Ошибка увеличения просмотров:', error);
      });


      if (currentUser) {
        const formData = new FormData();
        formData.append('ad_id', adId);

        fetch('http://localhost:8080/api/viewed-ads', {
          method: 'POST',
          credentials: 'include',
          body: formData
        }).catch(error => {
          console.error('Ошибка сохранения в просмотренные:', error);
        });
      }
    }
  }, [isOpen, ad, currentUser]);

  if (!isOpen || !ad) return null;


  const getActionButtonText = () => {
    const type = ad.type || ad.Type;

    if (type === 'Продать' || type === 'Купить') {

      return type === 'Продать' ? 'Купить' : 'Продать';
    }

    if (type === 'Сдать в аренду' || type === 'Услуги') {
      return 'Арендовать';
    }

    if (type === 'Поиск Заместителя') {
      return 'Связаться';
    }

    return 'Связаться';
  };


  const handleAction = async () => {
    try {



      const telegram = ad.owner_telegram || ad.telegram || ad.Telegram;

      if (!telegram) {
        setToast({
          message: 'У владельца объявления не указан Telegram. Невозможно связаться.',
          type: 'error'
        });
        return;
      }


      let telegramLink = telegram;
      if (telegram.startsWith('@')) {
        telegramLink = `https://t.me/${telegram.substring(1)}`;
      } else if (!telegram.startsWith('https://t.me/')) {
        telegramLink = `https://t.me/${telegram}`;
      }


      window.open(telegramLink, '_blank');

      setToast({
        message: 'Переход в Telegram...',
        type: 'success'
      });

    } catch (error) {
      console.error('Ошибка при обработке действия:', error);
      setToast({
        message: 'Произошла ошибка. Попробуйте позже.',
        type: 'error'
      });
    }
  };


  const isOwner = currentUser && (currentUser.nickname === ad.nickname || currentUser.nickname === ad.Nickname);


  const handleManage = () => {
    const adId = ad?.id || ad?.ID;
    navigate(`/manage-ads/${adId}`);
    onClose();
  };

  return (
    <>
      <div className="ad-detail-overlay" onClick={onClose}>
        <div className="ad-detail-modal" onClick={(e) => e.stopPropagation()}>
          <button className="ad-detail-close" onClick={onClose}>✕</button>

          <div className="ad-detail-content">
            <div className="ad-detail-image-section">
              <img
                src={ad.image || ad.Image || 'https://via.placeholder.com/600x400'}
                alt={ad.title || ad.Title}
                className="ad-detail-image"
              />
            </div>

            <div className="ad-detail-info-section">
              <div className="ad-detail-header">
                <h2 className="ad-detail-title">{ad.title || ad.Title}</h2>
                <span className="ad-detail-type">{ad.type || ad.Type}</span>
              </div>

              <div className="ad-detail-price-block">
                {ad.price || ad.Price ? (
                  <div className="ad-detail-price">
                    <span className="ad-detail-price-value">
                      {(ad.price || ad.Price).toLocaleString()}
                    </span>
                    <span className="ad-detail-price-currency">
                      {ad.currency || ad.Currency || '$'}
                    </span>
                    {ad.price_period || ad.PricePeriod ? (
                      <span className="ad-detail-price-period">
                        {ad.price_period || ad.PricePeriod}
                      </span>
                    ) : null}
                  </div>
                ) : (
                  <div className="ad-detail-price-free">Договорная</div>
                )}
              </div>

              <div className="ad-detail-description">
                <h3>Описание</h3>
                <p>{ad.description || ad.Description}</p>
              </div>

              <div className="ad-detail-meta">
                <div className="ad-detail-meta-item">
                  <span className="ad-detail-meta-label">Сервер:</span>
                  <span className="ad-detail-meta-value">{ad.server_name || ad.ServerName}</span>
                </div>
                <div className="ad-detail-meta-item">
                  <span className="ad-detail-meta-label">Категория:</span>
                  <span className="ad-detail-meta-value">{ad.category || ad.Category}</span>
                </div>
                <div className="ad-detail-meta-item">
                  <span className="ad-detail-meta-label">Просмотров:</span>
                  <span className="ad-detail-meta-value">{ad.views || ad.Views || 0}</span>
                </div>
                {(ad.rental_hours_limit || ad.RentalHoursLimit) && (
                  <div className="ad-detail-meta-item">
                    <span className="ad-detail-meta-label">Лимит аренды:</span>
                    <span className="ad-detail-meta-value">
                      {ad.rental_hours_limit || ad.RentalHoursLimit} часов
                    </span>
                  </div>
                )}
              </div>

              <div className="ad-detail-author">
                <img
                  src={ad.author_avatar || 'https://storage.yandexcloud.net/fotora.ru/uploads/2b0c131e8cfe54b1.jpeg'}
                  alt={ad.nickname || ad.Nickname}
                  className="ad-detail-author-avatar"
                />
                <div className="ad-detail-author-info">
                  <span className="ad-detail-author-name">{ad.nickname || ad.Nickname}</span>
                  <StarRating rating={ad.author_rating || ad.AuthorRating || 0} />
                </div>
              </div>

              <div className="ad-detail-actions">
                {isOwner ? (
                  <button className="ad-detail-action-btn manage-btn" onClick={handleManage}>
                    Перейти в управление
                  </button>
                ) : (
                  <button className="ad-detail-action-btn primary-btn" onClick={handleAction}>
                    {getActionButtonText()}
                  </button>
                )}
                <button className="ad-detail-action-btn cancel-btn" onClick={onClose}>
                  Закрыть
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      {toast && <Toast message={toast.message} type={toast.type} onClose={() => setToast(null)} />}
    </>
  );
};

export default AdDetailModal;
