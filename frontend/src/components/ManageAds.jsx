import { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import SnowEffect from './SnowEffect';
import BottomNavigation from './BottomNavigation';
import Toast from './Toast';
import '../styles/ManageAds.css';

function ManageAds() {
  const navigate = useNavigate();
  const { id } = useParams();
  const { user } = useAuth();

  const [ad, setAd] = useState(null);
  const [isLoading, setIsLoading] = useState(true);
  const [toast, setToast] = useState(null);
  const [isEditing, setIsEditing] = useState(false);

  const [formData, setFormData] = useState({
    title: '',
    type: '',
    description: '',
    price: '',
    currency: '',
    category: '',
    server_name: '',
    rental_hours_limit: ''
  });

  const [selectedImage, setSelectedImage] = useState(null);
  const [imagePreview, setImagePreview] = useState(null);

  useEffect(() => {
    document.title = 'Arz Store | Управление объявлением';
    fetchAd();
  }, [id]);

  const fetchAd = async () => {
    try {
      const response = await fetch(`http://localhost:8080/api/listings/user/${user.nickname}`);
      const data = await response.json();

      if (data.success && data.listings) {
        const foundAd = data.listings.find(listing => listing.id === parseInt(id));
        if (foundAd) {
          setAd(foundAd);
          setFormData({
            title: foundAd.title || '',
            type: foundAd.type || '',
            description: foundAd.description || '',
            price: foundAd.price || '',
            currency: foundAd.currency || '',
            category: foundAd.category || '',
            server_name: foundAd.server_name || '',
            rental_hours_limit: foundAd.rental_hours_limit || ''
          });
        } else {
          setToast({ message: 'Объявление не найдено', type: 'error' });
          setTimeout(() => navigate('/profile'), 2000);
        }
      }
    } catch (error) {
      console.error('Ошибка загрузки объявления:', error);
      setToast({ message: 'Ошибка загрузки объявления', type: 'error' });
    } finally {
      setIsLoading(false);
    }
  };

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
  };

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

    setSelectedImage(file);
    const reader = new FileReader();
    reader.onloadend = () => {
      setImagePreview(reader.result);
    };
    reader.readAsDataURL(file);
  };

  const handleSave = async () => {
    if (!formData.title || !formData.type || !formData.description) {
      setToast({ message: 'Заполните все обязательные поля', type: 'error' });
      return;
    }

    setToast({ message: 'Сохраняем изменения...', type: 'loading' });

    try {
      const formDataToSend = new FormData();
      formDataToSend.append('title', formData.title);
      formDataToSend.append('type', formData.type);
      formDataToSend.append('description', formData.description);
      formDataToSend.append('price', formData.price || '0');
      formDataToSend.append('currency', formData.currency || 'Договорная');
      formDataToSend.append('category', formData.category);
      formDataToSend.append('server_name', formData.server_name);
      if (formData.rental_hours_limit) {
        formDataToSend.append('rental_hours_limit', formData.rental_hours_limit);
      }

      if (selectedImage) {
        formDataToSend.append('image', selectedImage);
      }

      const response = await fetch(`http://localhost:8080/api/ads/${id}`, {
        method: 'PUT',
        credentials: 'include',
        body: formDataToSend
      });

      const data = await response.json();

      if (response.ok) {
        setToast({ message: 'Объявление успешно обновлено!', type: 'success' });
        setIsEditing(false);
        setTimeout(() => navigate('/profile'), 1500);
      } else {
        setToast({ message: data.error || 'Ошибка при обновлении объявления', type: 'error' });
      }
    } catch (error) {
      console.error('Ошибка обновления объявления:', error);
      setToast({ message: 'Ошибка при обновлении объявления', type: 'error' });
    }
  };

  const handleDelete = async () => {
    if (!window.confirm('Вы уверены, что хотите удалить это объявление? Это действие нельзя отменить.')) {
      return;
    }

    setToast({ message: 'Удаляем объявление...', type: 'loading' });

    try {
      const response = await fetch(`http://localhost:8080/api/ads/${id}`, {
        method: 'DELETE',
        credentials: 'include'
      });

      const data = await response.json();

      if (response.ok) {
        setToast({ message: 'Объявление успешно удалено!', type: 'success' });
        setTimeout(() => navigate('/profile'), 1500);
      } else {
        setToast({ message: data.error || 'Ошибка при удалении объявления', type: 'error' });
      }
    } catch (error) {
      console.error('Ошибка удаления объявления:', error);
      setToast({ message: 'Ошибка при удалении объявления', type: 'error' });
    }
  };

  if (isLoading) {
    return (
      <div className="manage-ads-page">
        <SnowEffect />
        <div className="manage-ads-content">
          <div className="loading-spinner">⏳</div>
          <p>Загрузка...</p>
        </div>
      </div>
    );
  }

  if (!ad) {
    return null;
  }

  return (
    <div className="manage-ads-page">
      <SnowEffect />

      <div className="manage-ads-content">
        <div className="manage-ads-header">
          <button className="back-button" onClick={() => navigate('/profile')}>
            ← Назад
          </button>
          <h1 className="manage-ads-title">Управление объявлением</h1>
        </div>

        <div className="manage-ads-card">
          <div className="manage-ads-image-section">
            <img
              src={imagePreview || ad.image}
              alt={formData.title || ad.title}
              className="manage-ads-image"
            />
            {isEditing && (
              <div className="manage-ads-image-upload">
                <label htmlFor="image-upload" className="image-upload-label">
                  Изменить изображение
                </label>
                <input
                  id="image-upload"
                  type="file"
                  accept="image/*"
                  onChange={handleImageSelect}
                  style={{ display: 'none' }}
                />
              </div>
            )}
          </div>

          <div className="manage-ads-form">
            <div className="form-group">
              <label>Название</label>
              {isEditing ? (
                <input
                  type="text"
                  name="title"
                  value={formData.title}
                  onChange={handleInputChange}
                  className="form-input"
                />
              ) : (
                <div className="form-value">{ad.title}</div>
              )}
            </div>

            <div className="form-group">
              <label>Тип объявления</label>
              {isEditing ? (
                <select
                  name="type"
                  value={formData.type}
                  onChange={handleInputChange}
                  className="form-input"
                >
                  <option value="Продать">Продать</option>
                  <option value="Купить">Купить</option>
                  <option value="Сдать в аренду">Сдать в аренду</option>
                  <option value="Услуги">Услуги</option>
                  <option value="Поиск Заместителя">Поиск Заместителя</option>
                </select>
              ) : (
                <div className="form-value">{ad.type}</div>
              )}
            </div>

            <div className="form-group">
              <label>Описание</label>
              {isEditing ? (
                <textarea
                  name="description"
                  value={formData.description}
                  onChange={handleInputChange}
                  className="form-input form-textarea"
                  rows={4}
                />
              ) : (
                <div className="form-value">{ad.description}</div>
              )}
            </div>

            <div className="form-group">
              <label>Цена</label>
              {isEditing ? (
                <input
                  type="number"
                  name="price"
                  value={formData.price}
                  onChange={handleInputChange}
                  className="form-input"
                />
              ) : (
                <div className="form-value">{ad.price} {ad.currency}</div>
              )}
            </div>

            <div className="form-group">
              <label>Категория</label>
              {isEditing ? (
                <select
                  name="category"
                  value={formData.category}
                  onChange={handleInputChange}
                  className="form-input"
                >
                  <option value="business">Бизнесы</option>
                  <option value="accs">Аккаунты</option>
                  <option value="house">Дома</option>
                  <option value="security">Охранники</option>
                  <option value="vehicle">Транспорт</option>
                  <option value="others">Прочее</option>
                </select>
              ) : (
                <div className="form-value">{ad.category}</div>
              )}
            </div>

            <div className="form-group">
              <label>Сервер</label>
              {isEditing ? (
                <select
                  name="server_name"
                  value={formData.server_name}
                  onChange={handleInputChange}
                  className="form-input"
                >
                  <option value="Phoenix">Phoenix</option>
                  <option value="Tucson">Tucson</option>
                  <option value="Scottdale">Scottdale</option>
                  <option value="Chandler">Chandler</option>
                  <option value="Brainburg">Brainburg</option>
                  <option value="Saint Rose">Saint Rose</option>
                  <option value="Mesa">Mesa</option>
                  <option value="Red-Rock">Red-Rock</option>
                  <option value="Yuma">Yuma</option>
                  <option value="Surprise">Surprise</option>
                  <option value="Prescott">Prescott</option>
                  <option value="Glendale">Glendale</option>
                  <option value="Kingman">Kingman</option>
                  <option value="Winslow">Winslow</option>
                  <option value="Payson">Payson</option>
                  <option value="Gilbert">Gilbert</option>
                  <option value="Show Low">Show Low</option>
                  <option value="Casa-Grande">Casa-Grande</option>
                  <option value="Page">Page</option>
                  <option value="Sun-City">Sun-City</option>
                  <option value="Queen-Creek">Queen-Creek</option>
                  <option value="Sedona">Sedona</option>
                  <option value="Holiday">Holiday</option>
                  <option value="Wednesday">Wednesday</option>
                  <option value="Yava">Yava</option>
                  <option value="Faraway">Faraway</option>
                  <option value="Bumble Bee">Bumble Bee</option>
                  <option value="Christmas">Christmas</option>
                </select>
              ) : (
                <div className="form-value">{ad.server_name}</div>
              )}
            </div>

            {formData.type === 'Сдать в аренду' && (
              <div className="form-group">
                <label>Лимит аренды (часы)</label>
                {isEditing ? (
                  <input
                    type="number"
                    name="rental_hours_limit"
                    value={formData.rental_hours_limit}
                    onChange={handleInputChange}
                    className="form-input"
                  />
                ) : (
                  <div className="form-value">{ad.rental_hours_limit || 'Не указано'}</div>
                )}
              </div>
            )}

            <div className="manage-ads-actions">
              {isEditing ? (
                <>
                  <button className="btn-save" onClick={handleSave}>
                    Сохранить изменения
                  </button>
                  <button className="btn-cancel" onClick={() => {
                    setIsEditing(false);
                    setSelectedImage(null);
                    setImagePreview(null);
                    setFormData({
                      title: ad.title || '',
                      type: ad.type || '',
                      description: ad.description || '',
                      price: ad.price || '',
                      currency: ad.currency || '',
                      category: ad.category || '',
                      server_name: ad.server_name || '',
                      rental_hours_limit: ad.rental_hours_limit || ''
                    });
                  }}>
                    Отменить
                  </button>
                </>
              ) : (
                <>
                  <button className="btn-edit" onClick={() => setIsEditing(true)}>
                    Редактировать
                  </button>
                  <button className="btn-delete" onClick={handleDelete}>
                    Удалить объявление
                  </button>
                </>
              )}
            </div>
          </div>
        </div>
      </div>

      <BottomNavigation />

      {toast && <Toast message={toast.message} type={toast.type} onClose={() => setToast(null)} />}
    </div>
  );
}

export default ManageAds;
