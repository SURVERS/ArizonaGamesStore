import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { useRateLimit } from '../hooks/useRateLimit';
import SnowEffect from './SnowEffect';
import '../styles/Auth.css';

function Login() {
  const [formData, setFormData] = useState({
    nickname: '',
    password: ''
  });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const { login } = useAuth();
  const navigate = useNavigate();
  const { executeWithRateLimit, isOnCooldown, remainingTime } = useRateLimit(2);


  useEffect(() => {
    document.title = 'Arz Store | Вход';
  }, []);

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
    setError('');
  };

  const handleLogin = async () => {
    if (!formData.nickname || !formData.password) {
      setError('Заполните все поля');
      return;
    }

    if (isOnCooldown) {
      setError(`Подождите ${remainingTime} сек. перед следующей попыткой`);
      return;
    }

    setLoading(true);
    setError('');

    const result = await executeWithRateLimit(async () => {
      try {
        const loginResult = await login(formData.nickname, formData.password, '');

        if (loginResult.success) {
          navigate('/feed');
        } else {
          setError(loginResult.error);
        }
        return loginResult;
      } catch (err) {
        setError('Ошибка соединения с сервером');
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

        <h1 className="auth-title">АВТОРИЗАЦИЯ...</h1>

        <p className="auth-description">
          Введите данные ниже, чтобы погрузиться в игровой магазин!
          Если ваш аккаунт не зарегистрирован, нажмите ниже кнопку
          «Нету Аккаунта? Зарегистрироваться».
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
              onKeyPress={(e) => e.key === 'Enter' && handleLogin()}
            />
          </div>
        </div>

        <button
          className="btn-primary"
          onClick={handleLogin}
          disabled={loading || isOnCooldown}
        >
          {loading ? 'ВХОД...' : isOnCooldown ? `ПОДОЖДИТЕ (${remainingTime}с)` : 'ВОЙТИ'}
        </button>

        <p className="toggle-text">Нету аккаунта?</p>

        <Link to="/register" className="btn-secondary">
          ЗАРЕГИСТРИРОВАТЬСЯ
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

export default Login;
