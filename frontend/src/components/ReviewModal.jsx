import { useState } from 'react';
import Toast from './Toast';
import '../styles/ReviewModal.css';

const ReviewModal = ({ ad, isOpen, onClose, onSuccess }) => {
  const [rating, setRating] = useState(1);
  const [hoveredRating, setHoveredRating] = useState(0);
  const [reviewText, setReviewText] = useState('');
  const [selectedImage, setSelectedImage] = useState(null);
  const [imagePreview, setImagePreview] = useState(null);
  const [toast, setToast] = useState(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  if (!isOpen || !ad) return null;

  const handleImageSelect = (e) => {
    const file = e.target.files?.[0];
    if (!file) return;


    if (file.size > 15 * 1024 * 1024) {
      setToast({ message: 'Размер изображения не должен превышать 15 МБ', type: 'error' });
      return;
    }


    if (!file.type.startsWith('image/')) {
      setToast({ message: 'Пожалуйста, выберите файл изображения', type: 'error' });
      return;
    }


    const img = new Image();
    img.onload = () => {
      const aspectRatio = img.width / img.height;

      if (aspectRatio < 1.2 || aspectRatio > 2.5) {
        setToast({
          message: 'Изображение должно быть прямоугольным (соотношение сторон от 1.2:1 до 2.5:1)',
          type: 'error'
        });
        return;
      }

      setSelectedImage(file);
      const reader = new FileReader();
      reader.onloadend = () => {
        setImagePreview(reader.result);
      };
      reader.readAsDataURL(file);
    };
    img.src = URL.createObjectURL(file);
  };

  const handleSubmit = async () => {
    if (!reviewText.trim()) {
      setToast({ message: 'Пожалуйста, напишите отзыв', type: 'error' });
      return;
    }

    if (!selectedImage) {
      setToast({ message: 'Пожалуйста, загрузите изображение-доказательство', type: 'error' });
      return;
    }

    setIsSubmitting(true);
    setToast({ message: 'Отправляем отзыв...', type: 'loading' });

    try {
      const formData = new FormData();
      formData.append('ad_id', ad.ID || ad.id);
      formData.append('rating', rating);
      formData.append('review_text', reviewText);
      formData.append('proof_image', selectedImage);

      const response = await fetch('http://localhost:8080/api/feedback', {
        method: 'POST',
        credentials: 'include',
        body: formData
      });

      const data = await response.json();

      if (response.ok) {
        setToast({ message: 'Отзыв успешно отправлен на модерацию!', type: 'success' });
        setTimeout(() => {
          onSuccess?.();
          onClose();
        }, 1500);
      } else {
        setToast({ message: data.error || 'Ошибка при отправке отзыва', type: 'error' });
        setIsSubmitting(false);
      }
    } catch (error) {
      console.error('Ошибка отправки отзыва:', error);
      setToast({ message: 'Ошибка при отправке отзыва', type: 'error' });
      setIsSubmitting(false);
    }
  };

  const renderStars = () => {
    const stars = [];
    for (let i = 1; i <= 5; i++) {
      const isFilled = i <= (hoveredRating || rating);
      stars.push(
        <span
          key={i}
          className={`review-star ${isFilled ? 'filled' : ''}`}
          onClick={() => setRating(i)}
          onMouseEnter={() => setHoveredRating(i)}
          onMouseLeave={() => setHoveredRating(0)}
        >
          ★
        </span>
      );
    }
    return stars;
  };

  return (
    <>
      <div className="review-modal-overlay" onClick={onClose}>
        <div className="review-modal" onClick={(e) => e.stopPropagation()}>
          <button className="review-modal-close" onClick={onClose}>✕</button>

          <h2 className="review-modal-title">Оставить отзыв</h2>

          
          <div className="review-ad-info">
            <img
              src={ad.image || ad.Image}
              alt={ad.title || ad.Title}
              className="review-ad-image"
            />
            <div className="review-ad-details">
              <h3>{ad.title || ad.Title}</h3>
              <p>{ad.description || ad.Description}</p>
              <div className="review-ad-meta">
                <span className="review-ad-id">ID: {ad.ID || ad.id}</span>
                <span className="review-ad-owner">Продавец: {ad.nickname || ad.Nickname}</span>
              </div>
            </div>
          </div>

          
          <div className="review-rating-section">
            <label className="review-label">Ваша оценка *</label>
            <div className="review-stars">
              {renderStars()}
            </div>
            <p className="review-rating-text">{rating} из 5 звезд</p>
          </div>

          
          <div className="review-text-section">
            <label className="review-label">Текст отзыва *</label>
            <textarea
              className="review-textarea"
              placeholder="Расскажите о вашем опыте взаимодействия с продавцом..."
              value={reviewText}
              onChange={(e) => setReviewText(e.target.value)}
              rows={5}
              maxLength={1000}
            />
            <div className="review-char-count">{reviewText.length}/1000</div>
          </div>

          
          <div className="review-image-section">
            <label className="review-label">Доказательство (скриншот) *</label>
            <p className="review-image-hint">
              Загрузите прямоугольное изображение до 15 МБ
            </p>

            {imagePreview ? (
              <div className="review-image-preview">
                <img src={imagePreview} alt="Preview" />
                <button
                  className="review-image-remove"
                  onClick={() => {
                    setSelectedImage(null);
                    setImagePreview(null);
                  }}
                >
                  Удалить
                </button>
              </div>
            ) : (
              <label className="review-image-upload">
                <input
                  type="file"
                  accept="image/*"
                  onChange={handleImageSelect}
                  style={{ display: 'none' }}
                />
                <div className="review-image-upload-content">
                  <span className="review-upload-icon">📷</span>
                  <span>Нажмите для загрузки изображения</span>
                </div>
              </label>
            )}
          </div>

          
          <div className="review-actions">
            <button
              className="review-btn review-btn-submit"
              onClick={handleSubmit}
              disabled={isSubmitting}
            >
              {isSubmitting ? 'Отправка...' : 'Отправить отзыв'}
            </button>
            <button
              className="review-btn review-btn-cancel"
              onClick={onClose}
              disabled={isSubmitting}
            >
              Закрыть
            </button>
          </div>
        </div>
      </div>

      {toast && <Toast message={toast.message} type={toast.type} onClose={() => setToast(null)} />}
    </>
  );
};

export default ReviewModal;
