import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { useRateLimit } from '../hooks/useRateLimit';
import SnowEffect from './SnowEffect';
import '../styles/Auth.css';

function Register() {
  const [formData, setFormData] = useState({
    nickname: '',
    email: '',
    password: '',
    confirmPassword: ''
  });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const { register } = useAuth();
  const navigate = useNavigate();
  const { executeWithRateLimit, isOnCooldown, remainingTime } = useRateLimit(3);


  useEffect(() => {
    document.title = 'Arz Store | Регистрация';
  }, []);

  const validateForm = () => {
    if (!formData.nickname || !formData.email || !formData.password || !formData.confirmPassword) {
      setError('Заполните все поля');
      return false;
    }

    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(formData.email)) {
      setError('Введите корректный email');
      return false;
    }

    if (formData.nickname.length < 3) {
      setError('Никнейм должен быть минимум 3 символа');
      return false;
    }

    if (formData.password.length < 6) {
      setError('Пароль должен быть минимум 6 символов');
      return false;
    }

    if (formData.password !== formData.confirmPassword) {
      setError('Пароли не совпадают');
      return false;
    }

    return true;
  };

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
    setError('');
  };

  const registerAccount = async () => {
    setError('');

    if (!validateForm()) {
      return;
    }

    if (isOnCooldown) {
      setError(`Подождите ${remainingTime} сек. перед следующей попыткой`);
      return;
    }

    setLoading(true);

    const result = await executeWithRateLimit(async () => {
      try {
        const regResult = await register(formData.nickname, formData.password, formData.email, '');

        console.log('Registration result:', regResult);
        console.log('requires_verify:', regResult.data?.requires_verify);

        if (regResult.success) {
          if (regResult.data.requires_verify) {
            console.log('Navigating to verify-email');
            navigate('/verify-email', { state: { email: formData.email } });
          } else {
            console.log('Navigating to feed');
            navigate('/feed');
          }
        } else {
          setError(regResult.error);
        }
        return regResult;
      } catch (err) {
        setError('Ошибка соединения с сервером. Причина: ' + err);
        return { success: false, error: 'Ошибка соединения с сервером' };
      } finally {
        setLoading(false);
      }
    });

    if (result.rateLimited) {
      setError(`Подождите ${remainingTime} сек. перед следующей попыткой`);
    }
  };

  return (
    <div className="auth-container">
      <SnowEffect />
      <div className="auth-box">
        <img
          src="/src/images/logo_arz.png"
          alt="Arizona Games Store Logo"
          className="logo"
        />

        <h1 className="auth-title">РЕГИСТРАЦИЯ...</h1>

        <p className="auth-description">
          Создайте новый аккаунт, чтобы начать использовать игровой магазин
          Arizona Role Play. Заполните данные ниже для регистрации.
        </p>

        {error && <div className="error-message">{error}</div>}

        <div className="input-group">
          <div className="input-wrapper">
            <span className="input-icon">
              <img
                src="/src/images/icons/6387915.png"
                className="icon-block-smail"
              />
            </span>
            <input
              type="text"
              name="nickname"
              placeholder="Nick_Name"
              className="auth-input"
              value={formData.nickname}
              onChange={handleInputChange}
              maxLength={50}
            />
          </div>

          <div className="input-wrapper">
            <span className="input-icon">
              <img
                src="/src/images/icons/6387915.png"
                className="icon-block-smail"
              />
            </span>
            <input
              type="email"
              name="email"
              placeholder="Email"
              className="auth-input"
              value={formData.email}
              onChange={handleInputChange}
              maxLength={255}
            />
          </div>

          <div className="input-wrapper">
            <span className="input-icon">
              <img
                src="/src/images/icons/8472244.png"
                className="icon-block-smail"
              />
            </span>
            <input
              type="password"
              name="password"
              placeholder="Password"
              className="auth-input"
              value={formData.password}
              onChange={handleInputChange}
              maxLength={100}
            />
          </div>

          <div className="input-wrapper">
            <span className="input-icon">
              <img
                src="/src/images/icons/8472244.png"
                className="icon-block-smail"
              />
            </span>
            <input
              type="password"
              name="confirmPassword"
              placeholder="Confirm Password"
              className="auth-input"
              value={formData.confirmPassword}
              onChange={handleInputChange}
              maxLength={100}
            />
          </div>
        </div>

        <button
          className="btn-primary"
          onClick={registerAccount}
          disabled={loading || isOnCooldown}
        >
          {loading ? 'РЕГИСТРАЦИЯ...' : isOnCooldown ? `ПОДОЖДИТЕ (${remainingTime}с)` : 'ЗАРЕГИСТРИРОВАТЬСЯ'}
        </button>

        <p className="toggle-text">Уже есть аккаунт?</p>

        <Link to="/auth" className="btn-secondary">
          ВОЙТИ
        </Link>

        <a
          href="https://t.me/your_channel"
          target="_blank"
          rel="noopener noreferrer"
          className="telegram-link"
        >
          <span className="telegram-icon">
              <img
                src="/src/images/icons/telegram.png"
                className="icon-block-smail"
              />
          </span>
          Мы в Telegram!
        </a>

        <p className="footer-text">
          Arizona Games Store - Игровой магазин созданный специально для проекта
          "Arizona Role Play". Здесь вы можете не только покупать/продавать имущество.
          Но и арендовывать у других игроков! Попробуй скорее :3
        </p>
      </div>
    </div>
  );
}

export default Register;
