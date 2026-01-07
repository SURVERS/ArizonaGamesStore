import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import SnowEffect from './SnowEffect';
import BottomNavigation from './BottomNavigation';
import Toast from './Toast';
import FilterModal from './FilterModal';
import '../styles/MobileCategory.css';

const SERVERS = [
  'ViceCity', 'Phoenix', 'Tucson', 'Scottdale', 'Winslow', 'Brainburg',
  'BumbleBee', 'CasaGrande', 'Chandler', 'Christmas', 'Faraway', 'Gilbert',
  'Glendale', 'Holiday', 'Kingman', 'Mesa', 'Page', 'Payson', 'Prescott',
  'QueenCreek', 'RedRock', 'SaintRose', 'Sedona', 'ShowLow', 'SunCity',
  'Surprise', 'Wednesday', 'Yava', 'Yuma', 'Love', 'Mirage', 'Drake', 'Space'
];


const LISTING_TYPES = ['Продать', 'Купить', 'Услуги'];


const PRICE_CONFIG = {
  'Продать': {
    'VC': { min: 1000, max: 50000000000, period: null },
    '$': { min: 10000, max: 1500000000000, period: null },
    'BTC': { min: 1, max: 15000000, period: null },
    'EURO': { min: 100, max: 650000000, period: null }
  },
  'Купить': {
    'VC': { min: 1000, max: 50000000000, period: null },
    '$': { min: 10000, max: 1500000000000, period: null },
    'BTC': { min: 1, max: 15000000, period: null },
    'EURO': { min: 100, max: 650000000, period: null }
  },
  'Услуги': {
    'VC': { min: 1000, max: 10000000, period: 'час' },
    '$': { min: 10000, max: 1000000000, period: 'час' },
    'BTC': { min: 1, max: 100000, period: 'час' },
    'EURO': { min: 100, max: 500000, period: 'час' }
  }
};

function Others() {
  const navigate = useNavigate();
  const { user } = useAuth();
  const [isCreating, setIsCreating] = useState(false);
  const [selectedServer, setSelectedServer] = useState('ViceCity');
  const [formData, setFormData] = useState({
    server: 'ViceCity',
    title: '',
    description: '',
    type: 'Продать',
    currency: null,
    price: '',
    rentalHoursLimit: '',
    image: null,
    imagePreview: null
  });

  const [listings, setListings] = useState([]);
  const [isLoading, setIsLoading] = useState(false);
  const [adCount, setAdCount] = useState(0);
  const [offset, setOffset] = useState(0);
  const [hasMore, setHasMore] = useState(true);
  const [isLoadingMore, setIsLoadingMore] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [cooldownTime, setCooldownTime] = useState(0);
  const [toast, setToast] = useState(null);
  const [rentalModal, setRentalModal] = useState({ isOpen: false, ad: null });
  const [rentalHours, setRentalHours] = useState('');


  const [sortBy, setSortBy] = useState('date_desc');
  const [showSortMenu, setShowSortMenu] = useState(false);
  const [filterModalOpen, setFilterModalOpen] = useState(false);
  const [filters, setFilters] = useState({
    type: '',
    currency: '',
    priceMin: '',
    priceMax: ''
  });


  useEffect(() => {
    if (cooldownTime > 0) {
      const timer = setTimeout(() => {
        setCooldownTime(cooldownTime - 1);
      }, 1000);
      return () => clearTimeout(timer);
    }
  }, [cooldownTime]);


  useEffect(() => {
    document.title = 'Arz Store | Другое';
  }, []);


  useEffect(() => {
    setListings([]);
    setOffset(0);
    setHasMore(true);
    fetchAds(0);
    fetchAdCount();

  }, [selectedServer, sortBy, filters]);


  useEffect(() => {
    const handleScroll = () => {
      if (window.innerHeight + document.documentElement.scrollTop >= document.documentElement.offsetHeight - 100) {
        if (!isLoading && !isLoadingMore && hasMore) {
          loadMore();
        }
      }
    };

    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  }, [isLoading, isLoadingMore, hasMore, offset]);

  const fetchAds = async (currentOffset) => {
    if (currentOffset === 0) {
      setIsLoading(true);
    } else {
      setIsLoadingMore(true);
    }

    try {

      let url = `http://localhost:8080/api/ads?category=others&server=${selectedServer}&limit=20&offset=${currentOffset}`;

      if (sortBy) {
        url += `&sort=${sortBy}`;
      }

      if (filters.type) {
        url += `&type=${encodeURIComponent(filters.type)}`;
      }
      if (filters.currency) {
        url += `&currency=${filters.currency}`;
      }
      if (filters.priceMin) {
        url += `&price_min=${filters.priceMin}`;
      }
      if (filters.priceMax) {
        url += `&price_max=${filters.priceMax}`;
      }

      const response = await fetch(url, { credentials: 'include' });

      if (response.ok) {
        const data = await response.json();
        const newAds = data.ads || [];

        if (currentOffset === 0) {
          setListings(newAds);
        } else {
          setListings(prev => [...prev, ...newAds]);
        }

        if (newAds.length < 20) {
          setHasMore(false);
        }
      }
    } catch (error) {
      console.error('Ошибка загрузки объявлений:', error);
    } finally {
      setIsLoading(false);
      setIsLoadingMore(false);
    }
  };

  const loadMore = () => {
    const newOffset = offset + 20;
    setOffset(newOffset);
    fetchAds(newOffset);
  };

  const fetchAdCount = async () => {
    try {
      const response = await fetch('http://localhost:8080/api/getadcount?CategoryName=others');
      const data = await response.json();
      setAdCount(data.count || 0);
    } catch (error) {
      console.error('Error fetching ad count:', error);
    }
  };

  const handleImageChange = (e) => {
    const file = e.target.files[0];
    if (file) {

      if (file.size > 10 * 1024 * 1024) {
        alert('Размер изображения не должен превышать 10 МБ');
        return;
      }


      if (!file.type.startsWith('image/')) {
        alert('Пожалуйста, загрузите изображение');
        return;
      }


      const img = new Image();
      const objectUrl = URL.createObjectURL(file);

      img.onload = () => {
        URL.revokeObjectURL(objectUrl);


        if (img.width > 1920 || img.height > 1080) {
          alert(`Максимальное разрешение: 1920x1080 пикселей.\nВаше изображение: ${img.width}x${img.height}`);
          return;
        }


        if (img.width < 300 || img.height < 200) {
          alert(`Минимальное разрешение: 300x200 пикселей.\nВаше изображение: ${img.width}x${img.height}`);
          return;
        }


        const reader = new FileReader();
        reader.onloadend = () => {
          setFormData({
            ...formData,
            image: file,
            imagePreview: reader.result
          });
        };
        reader.readAsDataURL(file);
      };

      img.onerror = () => {
        URL.revokeObjectURL(objectUrl);
        alert('Ошибка загрузки изображения');
      };

      img.src = objectUrl;
    }
  };


  const getAvailableCurrencies = () => {
    const baseCurrencies = formData.server === 'ViceCity'
      ? ['VC', 'BTC', 'EURO']
      : ['$', 'BTC', 'EURO'];
    return [...baseCurrencies, 'Договорная'];
  };


  const getPriceConfig = () => {
    if (!formData.type || !formData.currency) return null;
    return PRICE_CONFIG[formData.type]?.[formData.currency];
  };


  const formatNumber = (num) => {
    return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, '.');
  };


  const formatPrice = (price) => {
    if (!price) return '0';
    return price.toString().replace(/\B(?=(\d{3})+(?!\d))/g, '.');
  };


  const getCurrencySymbol = (currency) => {
    switch(currency) {
      case 'VC': return 'VC$';
      case '$': return '$';
      case 'BTC': return 'BTC';
      case 'EURO': return 'EURO';
      default: return currency;
    }
  };


  const getActionText = (type) => {
    switch(type) {
      case 'Продать': return 'КУПИТЬ';
      case 'Купить': return 'ПРОДАТЬ';
      case 'Сдать в аренду': return 'АРЕНДОВАТЬ';
      case 'Поиск Заместителя': return 'СВЯЗАТЬСЯ';
      case 'Услуги': return 'СВЯЗАТЬСЯ';
      default: return 'ДЕЙСТВИЕ';
    }
  };


  const getTypeCost = (type) => {
    switch(type) {
      case 'Услуги': return '/ час';
      default: return '';
    }
  };


  const handleOpenRental = (ad) => {
    setRentalModal({ isOpen: true, ad });
    setRentalHours('');
  };


  const handleCloseRental = () => {
    setRentalModal({ isOpen: false, ad: null });
    setRentalHours('');
  };


  const calculateRentalCost = () => {
    if (!rentalHours || !rentalModal.ad) return 0;
    const hours = parseInt(rentalHours);
    const pricePerHour = rentalModal.ad.price || 0;
    return hours * pricePerHour;
  };


  const sanitizeInput = (input) => {
    if (!input) return '';
    return input
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')
      .replace(/"/g, '&quot;')
      .replace(/'/g, '&#x27;')
      .replace(/\
  };

  const handleInputChange = (e) => {
    const { name, value } = e.target;


    if (name === 'type') {
      setFormData({
        ...formData,
        [name]: value,
        currency: null,
        price: '',
        rentalHoursLimit: ''
      });
    }

    else if (name === 'server') {
      const availableCurrencies = value === 'ViceCity' ? ['VC', 'BTC', 'EURO'] : ['$', 'BTC', 'EURO'];
      const newCurrency = availableCurrencies.includes(formData.currency) ? formData.currency : null;
      setFormData({
        ...formData,
        [name]: value,
        currency: newCurrency,
        price: newCurrency ? formData.price : ''
      });
    }
    else {

      const sanitizedValue = (name === 'title' || name === 'description') ? sanitizeInput(value) : value;
      setFormData({
        ...formData,
        [name]: sanitizedValue
      });
    }
  };

  const handleCurrencyChange = (currency) => {
    setFormData({
      ...formData,
      currency,
      price: ''
    });
  };

  const handlePriceChange = (e) => {
    const value = e.target.value.replace(/\D/g, '');
    setFormData({
      ...formData,
      price: value
    });
  };

  const handleRentalHoursChange = (e) => {
    const value = e.target.value.replace(/\D/g, '');
    const numValue = parseInt(value || '0');
    if (numValue <= 180) {
      setFormData({
        ...formData,
        rentalHoursLimit: value
      });
    }
  };

  const handleCreateListing = async () => {

    if (!formData.image) {
      setToast({ message: 'Пожалуйста, загрузите изображение', type: 'error' });
      return;
    }
    if (!formData.title.trim()) {
      setToast({ message: 'Введите название товара', type: 'error' });
      return;
    }
    if (!formData.description.trim()) {
      setToast({ message: 'Введите описание товара', type: 'error' });
      return;
    }


    if (!formData.currency) {
      setToast({ message: 'Выберите валюту', type: 'error' });
      return;
    }


    if (formData.currency !== 'Договорная') {
      if (!formData.price || formData.price === '0') {
        const priceConfig = getPriceConfig();
        if (priceConfig && priceConfig.min > 0) {
          setToast({ message: 'Введите цену', type: 'error' });
          return;
        }
      }


      const priceConfig = getPriceConfig();
      if (priceConfig) {
        const priceNum = parseInt(formData.price || '0');
        if (priceNum < priceConfig.min || priceNum > priceConfig.max) {
          const currencySymbol = formData.currency === 'VC' ? 'VC$' : formData.currency === '$' ? '$' : formData.currency;
          setToast({ message: `Цена должна быть от ${formatNumber(priceConfig.min)} до ${formatNumber(priceConfig.max)} ${currencySymbol}`, type: 'error' });
          return;
        }
      }
    }


    if (formData.type === 'Услуги') {
      if (!formData.rentalHoursLimit || formData.rentalHoursLimit === '0') {
        setToast({ message: 'Укажите лимит часов услуги (от 1 до 180)', type: 'error' });
        return;
      }
      const hoursLimit = parseInt(formData.rentalHoursLimit);
      if (hoursLimit < 1 || hoursLimit > 180) {
        setToast({ message: 'Лимит часов должен быть от 1 до 180', type: 'error' });
        return;
      }
    }

    setIsSubmitting(true);
    setToast({ message: 'Идёт создание объявления... Подождите немного', type: 'loading' });


    const timestamp = Date.now();
    const fileExtension = formData.image.name.split('.').pop();
    const imagePath = `ads/others/${timestamp}_${user?.nickname || 'user'}.${fileExtension}`;


    const formDataToSend = new FormData();
    formDataToSend.append('server', formData.server);
    formDataToSend.append('title', formData.title);
    formDataToSend.append('description', formData.description);
    formDataToSend.append('type', formData.type);
    formDataToSend.append('currency', formData.currency);


    if (formData.currency === 'Договорная') {
      formDataToSend.append('price', '0');
    } else {
      const priceConfig = getPriceConfig();
      formDataToSend.append('price', formData.price || '0');
      if (priceConfig?.period) {
        formDataToSend.append('pricePeriod', priceConfig.period);
      }
    }
    formDataToSend.append('image', formData.image);
    formDataToSend.append('imagePath', imagePath);
    formDataToSend.append('category', 'others');
    formDataToSend.append('nickname', user?.nickname || 'Unknown');


    if (formData.type === 'Услуги' && formData.rentalHoursLimit) {
      formDataToSend.append('rentalHoursLimit', formData.rentalHoursLimit);
    }

    try {
      const response = await fetch('http://localhost:8080/api/createnewads', {
        method: 'POST',
        credentials: 'include',
        body: formDataToSend
      });

      const result = await response.json();

      if (response.ok) {
        setToast({ message: 'Объявление успешно создано!', type: 'success' });


        setFormData({
          server: 'ViceCity',
          title: '',
          description: '',
          type: 'Продать',
          currency: null,
          price: '',
          rentalHoursLimit: '',
          image: null,
          imagePreview: null
        });


        setListings([]);
        setOffset(0);
        setHasMore(true);
        await fetchAds(0);
        await fetchAdCount();


        setCooldownTime(60);


        setTimeout(() => {
          setIsCreating(false);
        }, 1500);
      } else {
        setToast({ message: `Ошибка: ${result.error || 'Не удалось создать объявление'}`, type: 'error' });
      }
    } catch (error) {
      console.error('Ошибка при создании объявления:', error);
      setToast({ message: 'Ошибка соединения с сервером', type: 'error' });
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleCancel = () => {
    setFormData({
      server: 'ViceCity',
      title: '',
      description: '',
      type: 'Продать',
      currency: null,
      price: '',
      rentalHoursLimit: '',
      image: null,
      imagePreview: null
    });
    setIsCreating(false);
  };

  if (isCreating) {
    return (
      <div className="mobile-category-page">
        <SnowEffect />
        {toast && <Toast message={toast.message} type={toast.type} onClose={() => setToast(null)} />}

        <div className="mobile-category-content create-listing-content">
          <div className="mobile-category-header">
            <button
              className="back-button"
              onClick={() => setIsCreating(false)}
              disabled={isSubmitting}
            >
              ← Назад
            </button>
            <h1 className="mobile-category-title">Создать объявление</h1>
          </div>

          <div className="mobile-divider"></div>

          <div className="listing-form">
            <div className="form-group">
              <label>Сервер *</label>
              <select
                name="server"
                value={formData.server}
                onChange={handleInputChange}
                className="form-select"
              >
                {SERVERS.map(server => (
                  <option key={server} value={server}>{server}</option>
                ))}
              </select>
            </div>

            <div className="form-group image-upload-group">
              <label>Изображение товара * (макс. 10 МБ)</label>
              <div
                className="image-upload-area"
                onClick={() => document.getElementById('image-input').click()}
              >
                {formData.imagePreview ? (
                  <img src={formData.imagePreview} alt="Preview" className="image-preview" />
                ) : (
                  <div className="upload-placeholder">
                    <div className="upload-icon">📷</div>
                    <p>Нажмите для загрузки изображения</p>
                    <span className="upload-hint">JPG, PNG, GIF (макс. 10 МБ)</span>
                  </div>
                )}
              </div>
              <input
                id="image-input"
                type="file"
                accept="image/*"
                onChange={handleImageChange}
                style={{ display: 'none' }}
              />
            </div>

            <div className="form-group">
              <label>Название товара * (макс. 25 символов)</label>
              <input
                type="text"
                name="title"
                value={formData.title}
                onChange={handleInputChange}
                maxLength={25}
                placeholder="Например: Реклама услуг"
                className="form-input"
              />
              <span className="char-counter">{formData.title.length}/25</span>
            </div>

            <div className="form-group">
              <label>Описание товара * (макс. 500 символов)</label>
              <textarea
                name="description"
                value={formData.description}
                onChange={handleInputChange}
                maxLength={500}
                placeholder="Опишите ваше предложение подробнее..."
                className="form-textarea"
                rows={6}
              />
              <span className="char-counter">{formData.description.length}/500</span>
            </div>

            <div className="form-group">
              <label>Тип объявления *</label>
              <div className="type-selector">
                {LISTING_TYPES.map(type => (
                  <button
                    key={type}
                    type="button"
                    className={`type-btn ${formData.type === type ? 'active' : ''}`}
                    onClick={() => handleInputChange({ target: { name: 'type', value: type } })}
                  >
                    {type}
                  </button>
                ))}
              </div>
            </div>

            
            <div className="form-group">
              <label>Валюта *</label>
              <div className="type-selector">
                {getAvailableCurrencies().map(currency => (
                  <button
                    key={currency}
                    type="button"
                    className={`type-btn ${formData.currency === currency ? 'active' : ''}`}
                    onClick={() => handleCurrencyChange(currency)}
                  >
                    {currency === 'VC' ? 'VC$' : currency}
                  </button>
                ))}
              </div>
            </div>

            
            {formData.currency && formData.currency !== 'Договорная' && (
              <div className="form-group">
                <label>
                  Цена *
                  {(() => {
                    const config = getPriceConfig();
                    if (config) {
                      const currencySymbol = formData.currency === 'VC' ? 'VC$' : formData.currency === '$' ? '$' : formData.currency;
                      const periodText = config.period ? ` / ${config.period}` : '';
                      return ` (от ${formatNumber(config.min)} до ${formatNumber(config.max)} ${currencySymbol}${periodText})`;
                    }
                    return '';
                  })()}
                </label>
                <input
                  type="text"
                  value={formData.price}
                  onChange={handlePriceChange}
                  placeholder={(() => {
                    const config = getPriceConfig();
                    if (config && config.min === 0) return '0 (бесплатно)';
                    return 'Введите цену';
                  })()}
                  className="form-input"
                />
              </div>
            )}

            {formData.currency === 'Договорная' && (
              <div className="form-group">
                <div className="negotiable-notice">
                  💬 Цена договорная - обсуждается с покупателем
                </div>
              </div>
            )}

            
            {formData.type === 'Услуги' && (
              <div className="form-group">
                <label>Лимит часов услуги * (от 1 до 180 часов)</label>
                <input
                  type="text"
                  value={formData.rentalHoursLimit}
                  onChange={handleRentalHoursChange}
                  placeholder="Например: 24"
                  className="form-input"
                />
                <span className="char-counter">
                  {formData.rentalHoursLimit ? `${formData.rentalHoursLimit} час${formData.rentalHoursLimit === '1' ? '' : formData.rentalHoursLimit > 1 && formData.rentalHoursLimit < 5 ? 'а' : 'ов'}` : 'Не указано'}
                </span>
              </div>
            )}

            <div className="form-actions">
              <button
                className="btn-create"
                onClick={handleCreateListing}
                disabled={isSubmitting || cooldownTime > 0}
              >
                {cooldownTime > 0
                  ? `Подождите ${cooldownTime} сек`
                  : isSubmitting
                  ? '⏳ Создание...'
                  : '✅ Создать объявление'
                }
              </button>
              <button
                className="btn-cancel"
                onClick={handleCancel}
                disabled={isSubmitting}
              >
                ← Вернуться назад
              </button>
            </div>
          </div>
        </div>

        <BottomNavigation />
      </div>
    );
  }

  return (
    <div className="mobile-category-page">
      <SnowEffect />
      {toast && <Toast message={toast.message} type={toast.type} onClose={() => setToast(null)} />}

      <div className="mobile-category-content">
        <div className="mobile-category-header">
          <button className="back-button" onClick={() => navigate('/feed')}>
            ← Назад
          </button>
          <h1 className="mobile-category-title">Реклама / Остальное</h1>
        </div>

        <div className="mobile-controls">
          <div className="sort-dropdown">
            <button className="control-btn" onClick={() => setShowSortMenu(!showSortMenu)}>
              <img src="/src/images/icons/sort.png" alt="Sort" />
              <span>Сортировка</span>
            </button>
            {showSortMenu && (
              <div className="sort-menu">
                <button onClick={() => { setSortBy('date_desc'); setShowSortMenu(false); }}>
                  Новые
                </button>
                <button onClick={() => { setSortBy('date_asc'); setShowSortMenu(false); }}>
                  Старые
                </button>
                <button onClick={() => { setSortBy('price_desc'); setShowSortMenu(false); }}>
                  Дорогие
                </button>
                <button onClick={() => { setSortBy('price_asc'); setShowSortMenu(false); }}>
                  Дешевые
                </button>
                <button onClick={() => { setSortBy('views_desc'); setShowSortMenu(false); }}>
                  Популярные
                </button>
              </div>
            )}
          </div>
          <button className="control-btn" onClick={() => setFilterModalOpen(true)}>
            <img src="/src/images/icons/filter.png" alt="Filter" />
            <span>Фильтры</span>
          </button>
        </div>

        <div className="mobile-divider"></div>

        <div className="mobile-server-selector">
          <label>Сервер:</label>
          <select value={selectedServer} onChange={(e) => setSelectedServer(e.target.value)}>
            {SERVERS.map(server => (
              <option key={server} value={server}>{server}</option>
            ))}
          </select>
        </div>

        <div className="mobile-ad-count">
          Всего объявлений: {adCount}
        </div>

        <button className="mobile-create-btn" onClick={() => setIsCreating(true)}>
          + Создать объявление
        </button>

        <div className="mobile-listings-container">
          {isLoading ? (
            <div className="mobile-no-listings">
              <div className="no-listings-icon">⏳</div>
              <h3>Загрузка объявлений...</h3>
            </div>
          ) : listings.length === 0 ? (
            <div className="mobile-no-listings">
              <div className="no-listings-icon">📭</div>
              <h3>В данный момент нет объявлений</h3>
              <p>Станьте первым, кто создаст объявление!</p>
            </div>
          ) : (
            <div className="mobile-listings-grid">
              {listings.map((ad) => (
                <div key={ad.id} className="mobile-listing-card">
                  <div className="mobile-listing-image">
                    <img src={ad.image} alt={ad.title} />
                    <div className="mobile-listing-id">ID: {ad.id}</div>
                  </div>
                  <div className="mobile-listing-content">
                    <div className="mobile-listing-price">
                      {ad.currency === 'Договорная' ? 'Договорная' : `${formatPrice(ad.price)} ${getCurrencySymbol(ad.currency)} ${getTypeCost(ad.type)}`}
                    </div>
                    <div className="mobile-listing-seller">@{ad.nickname}</div>
                    <div className="mobile-listing-title">{ad.title}</div>
                    <div className="mobile-listing-description">{ad.description}</div>
                    <div className="mobile-listing-rating">
                      <img src="/src/images/icons/star.png" alt="Rating" className="star-icon" />
                      <span>{(ad.author_rating || 0).toFixed(1)}</span>
                    </div>
                    <button
                      className="mobile-listing-action"
                      onClick={() => ad.type === 'Услуги' ? handleOpenRental(ad) : null}
                    >
                      <img src="/src/images/icons/store-icon.png" alt="Action" className="action-icon" />
                      <span>{getActionText(ad.type)}</span>
                    </button>
                  </div>
                </div>
              ))}
            </div>
          )}

          {isLoadingMore && (
            <div className="mobile-loading-more">
              <div className="loading-spinner">⏳</div>
              <p>Загрузка...</p>
            </div>
          )}
          {!hasMore && listings.length > 0 && (
            <div className="mobile-no-more">
              <p>Все объявления загружены</p>
            </div>
          )}
        </div>
      </div>

      
      {rentalModal.isOpen && (
        <div className="modal-overlay" onClick={handleCloseRental}>
          <div className="modal-content rental-modal" onClick={(e) => e.stopPropagation()}>
            <h2>Объявление #{rentalModal.ad.id} от {rentalModal.ad.nickname}</h2>
            <div className="rental-ad-preview">
              <img src={rentalModal.ad.image} alt={rentalModal.ad.title} />
              <div className="rental-details">
                <h3>{rentalModal.ad.title}</h3>
                <p>{rentalModal.ad.description}</p>
                <p><strong>Создатель объявления:</strong> @{rentalModal.ad.nickname}</p>
                <p><strong>Цена:</strong> {formatPrice(rentalModal.ad.price)} {getCurrencySymbol(rentalModal.ad.currency)} / час</p>
              </div>
            </div>
            <div className="rental-calculator">
              <label>
                Введите количество часов на которое готовы взять услугу
                {rentalModal.ad.rental_hours_limit && ` (от 1 до ${rentalModal.ad.rental_hours_limit} часов)`}
              </label>
              <input
                type="number"
                min="1"
                max={rentalModal.ad.rental_hours_limit || 180}
                value={rentalHours}
                onChange={(e) => {
                  const value = parseInt(e.target.value);
                  const maxHours = rentalModal.ad.rental_hours_limit || 180;
                  if (value <= maxHours) {
                    setRentalHours(e.target.value);
                  }
                }}
                placeholder={`От 1 до ${rentalModal.ad.rental_hours_limit || 180} часов`}
              />
              {rentalHours && (
                <div className="rental-result">
                  <strong>Итог:</strong> {formatPrice(calculateRentalCost())} {getCurrencySymbol(rentalModal.ad.currency)}
                </div>
              )}
            </div>
            <div className="modal-actions">
              <button className="btn-cancel" onClick={handleCloseRental}>
                Закрыть
              </button>
            </div>
          </div>
        </div>
      )}

      <FilterModal
        isOpen={filterModalOpen}
        onClose={() => setFilterModalOpen(false)}
        onApply={(newFilters) => setFilters(newFilters)}
        currentFilters={filters}
      />

      <BottomNavigation />
    </div>
  );
}

export default Others;
