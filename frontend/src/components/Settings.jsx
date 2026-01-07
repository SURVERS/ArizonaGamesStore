import { useState, useRef, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import SnowEffect from './SnowEffect';
import BottomNavigation from './BottomNavigation';
import Toast from './Toast';
import '../styles/Settings.css';

function Settings() {
  const navigate = useNavigate();
  const { user, logout, checkAuth } = useAuth();
  const avatarInputRef = useRef(null);

  const [selectedAvatar, setSelectedAvatar] = useState(null);
  const [avatarPreview, setAvatarPreview] = useState(null);
  const [isAvatarUploading, setIsAvatarUploading] = useState(false);

  const [isEditingNickname, setIsEditingNickname] = useState(false);
  const [nickname, setNickname] = useState('');

  const [isEditingEmail, setIsEditingEmail] = useState(false);
  const [email, setEmail] = useState('');

  const [isEditingTelegram, setIsEditingTelegram] = useState(false);
  const [telegram, setTelegram] = useState('');

  const [isEditingPassword, setIsEditingPassword] = useState(false);
  const [isEditingDescriptionProfile, setIsEditingDescriptionProfile] = useState(false);
  const [description, setDescription] = useState('');
  const [passwordData, setPasswordData] = useState({
    oldPassword: '',
    newPassword: '',
    confirmPassword: ''
  });

  const [theme, setTheme] = useState('dark');
  const [toast, setToast] = useState(null);


  useEffect(() => {
    if (user) {
      setNickname(user.nickname || '');
      setEmail(user.email || '');
      setTelegram(user.telegram || '');
      setDescription(user.user_description || '');
      setTheme(user.theme || 'dark');
    }
  }, [user]);


  useEffect(() => {
    document.title = 'Arz Store | Настройки';
  }, []);


  const handleAvatarClick = () => {
    if (!isAvatarUploading) {
      avatarInputRef.current?.click();
    }
  };

  const handleAvatarSelect = (e) => {
    const file = e.target.files?.[0];
    if (!file) return;

    if (file.size > 5 * 1024 * 1024) {
      setToast({ message: 'Размер изображения не должен превышать 5 МБ', type: 'error' });
      return;
    }

    if (!file.type.startsWith('image/')) {
      setToast({ message: 'Пожалуйста, выберите файл изображения', type: 'error' });
      return;
    }

    setSelectedAvatar(file);

    const reader = new FileReader();
    reader.onloadend = () => {
      setAvatarPreview(reader.result);
    };
    reader.readAsDataURL(file);
  };

  const handleSaveAvatar = async () => {
    if (!selectedAvatar) return;

    setIsAvatarUploading(true);
    setToast({ message: 'Загружаем новую аватарку...', type: 'loading' });

    try {
      const formData = new FormData();
      formData.append('avatar', selectedAvatar);

      const response = await fetch('http://localhost:8080/api/profile/update-avatar', {
        method: 'POST',
        credentials: 'include',
        body: formData,
      });

      const data = await response.json();

      if (response.ok) {
        setToast({ message: 'Аватарка успешно обновлена!', type: 'success' });
        setTimeout(() => window.location.reload(), 1500);
      } else {
        setToast({ message: data.error || 'Ошибка при загрузке аватарки', type: 'error' });
        setIsAvatarUploading(false);
      }
    } catch (error) {
      console.error('Ошибка загрузки аватарки:', error);
      setToast({ message: 'Ошибка при загрузке аватарки', type: 'error' });
      setIsAvatarUploading(false);
    }
  };

  const handleCancelAvatar = () => {
    setSelectedAvatar(null);
    setAvatarPreview(null);
    if (avatarInputRef.current) {
      avatarInputRef.current.value = '';
    }
  };


  const handleSaveNickname = async () => {
    if (!nickname.trim()) {
      setToast({ message: 'Никнейм не может быть пустым', type: 'error' });
      return;
    }

    setToast({ message: 'Сохраняем никнейм...', type: 'loading' });

    try {
      const response = await fetch('http://localhost:8080/api/profile/update-nickname', {
        method: 'PUT',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ nickname }),
      });

      const data = await response.json();

      if (response.ok) {
        setToast({ message: 'Никнейм успешно обновлен!', type: 'success' });
        setIsEditingNickname(false);
        setTimeout(() => window.location.reload(), 1500);
      } else {
        setToast({ message: data.error || 'Ошибка при обновлении никнейма', type: 'error' });
      }
    } catch (error) {
      console.error('Ошибка обновления никнейма:', error);
      setToast({ message: 'Ошибка при обновлении никнейма', type: 'error' });
    }
  };


  const handleSaveEmail = async () => {
    if (!email.trim()) {
      setToast({ message: 'Email не может быть пустым', type: 'error' });
      return;
    }

    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(email)) {
      setToast({ message: 'Введите корректный email', type: 'error' });
      return;
    }

    setToast({ message: 'Сохраняем email...', type: 'loading' });

    try {
      const response = await fetch('http://localhost:8080/api/profile/update-email', {
        method: 'PUT',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email }),
      });

      const data = await response.json();

      if (response.ok) {
        setToast({ message: 'Email успешно обновлен!', type: 'success' });
        setIsEditingEmail(false);
        setTimeout(() => window.location.reload(), 1500);
      } else {
        setToast({ message: data.error || 'Ошибка при обновлении email', type: 'error' });
      }
    } catch (error) {
      console.error('Ошибка обновления email:', error);
      setToast({ message: 'Ошибка при обновлении email', type: 'error' });
    }
  };


  const handleSaveTelegram = async () => {
    if (!telegram.trim()) {
      setToast({ message: 'Telegram не может быть пустым', type: 'error' });
      return;
    }

    setToast({ message: 'Сохраняем Telegram...', type: 'loading' });

    try {
      const response = await fetch('http://localhost:8080/api/profile/update-telegram', {
        method: 'PUT',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ telegram }),
      });

      const data = await response.json();

      if (response.ok) {
        setToast({ message: 'Telegram успешно обновлен!', type: 'success' });
        setIsEditingTelegram(false);
        setTimeout(() => window.location.reload(), 1500);
      } else {
        setToast({ message: data.error || 'Ошибка при обновлении Telegram', type: 'error' });
      }
    } catch (error) {
      console.error('Ошибка обновления Telegram:', error);
      setToast({ message: 'Ошибка при обновлении Telegram', type: 'error' });
    }
  };


  const handleSaveDescription = async () => {
    if (!description.trim()) {
      setToast({ message: 'Описание не может быть пустым', type: 'error' });
      return;
    }

    if (description.length < 3) {
      setToast({ message: 'Описание должно содержать минимум 3 символа', type: 'error' });
      return;
    }

    if (description.length > 200) {
      setToast({ message: 'Описание не должно превышать 200 символов', type: 'error' });
      return;
    }

    setToast({ message: 'Сохраняем описание...', type: 'loading' });

    try {
      const response = await fetch('http://localhost:8080/api/profile/update-description', {
        method: 'PUT',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ description }),
      });

      const data = await response.json();

      if (response.ok) {
        setToast({ message: 'Описание успешно обновлено!', type: 'success' });
        setIsEditingDescriptionProfile(false);
        setTimeout(() => window.location.reload(), 1500);
      } else {
        setToast({ message: data.error || 'Ошибка при обновлении описания', type: 'error' });
      }
    } catch (error) {
      console.error('Ошибка обновления описания:', error);
      setToast({ message: 'Ошибка при обновлении описания', type: 'error' });
    }
  };


  const handleSavePassword = async () => {
    if (!passwordData.oldPassword || !passwordData.newPassword || !passwordData.confirmPassword) {
      setToast({ message: 'Заполните все поля', type: 'error' });
      return;
    }

    if (passwordData.newPassword !== passwordData.confirmPassword) {
      setToast({ message: 'Пароли не совпадают', type: 'error' });
      return;
    }

    if (passwordData.newPassword.length < 6) {
      setToast({ message: 'Пароль должен содержать минимум 6 символов', type: 'error' });
      return;
    }

    setToast({ message: 'Сохраняем пароль...', type: 'loading' });

    try {
      const response = await fetch('http://localhost:8080/api/profile/update-password', {
        method: 'PUT',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(passwordData),
      });

      const data = await response.json();

      if (response.ok) {
        setToast({ message: 'Пароль успешно обновлен!', type: 'success' });
        setIsEditingPassword(false);
        setPasswordData({ oldPassword: '', newPassword: '', confirmPassword: '' });
      } else {
        setToast({ message: data.error || 'Ошибка при обновлении пароля', type: 'error' });
      }
    } catch (error) {
      console.error('Ошибка обновления пароля:', error);
      setToast({ message: 'Ошибка при обновлении пароля', type: 'error' });
    }
  };


  const handleThemeChange = async (newTheme) => {
    const oldTheme = theme;
    setTheme(newTheme);


    if (newTheme === 'light') {
      document.body.classList.add('light-theme');
      document.body.classList.remove('dark-theme');
    } else {
      document.body.classList.add('dark-theme');
      document.body.classList.remove('light-theme');
    }

    try {
      const response = await fetch('http://localhost:8080/api/profile/update-theme', {
        method: 'PUT',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ theme: newTheme }),
      });

      if (response.ok) {
        setToast({ message: 'Тема успешно изменена!', type: 'success' });

        await checkAuth();
      } else {

        setTheme(oldTheme);
        if (oldTheme === 'light') {
          document.body.classList.add('light-theme');
          document.body.classList.remove('dark-theme');
        } else {
          document.body.classList.add('dark-theme');
          document.body.classList.remove('light-theme');
        }
        setToast({ message: 'Ошибка при изменении темы', type: 'error' });
      }
    } catch (error) {
      console.error('Ошибка изменения темы:', error);

      setTheme(oldTheme);
      if (oldTheme === 'light') {
        document.body.classList.add('light-theme');
        document.body.classList.remove('dark-theme');
      } else {
        document.body.classList.add('dark-theme');
        document.body.classList.remove('light-theme');
      }
      setToast({ message: 'Ошибка при изменении темы', type: 'error' });
    }
  };


  const handleLogout = async () => {
    if (!window.confirm('Вы уверены, что хотите выйти из профиля?')) {
      return;
    }

    await logout();
    navigate('/auth');
  };

  return (
    <div className="settings-page">
      <SnowEffect />

      <div className="settings-content">
        
        <div className="profile-header">
          <button className="back-button" onClick={() => navigate('/profile')}>
            ← Назад
          </button>
          <h1 className="profile-title">Настройки</h1>
        </div>

        
        <div className="settings-block">
          
          <div className="settings-item">
            <div className="settings-avatar-container">
              <div className="settings-avatar" onClick={handleAvatarClick}>
                <img
                  src={avatarPreview || user?.avatar || '/src/images/icons/user.png'}
                  alt="Avatar"
                  onError={(e) => {
                    e.target.src = '/src/images/icons/user.png';
                  }}
                />
                {!selectedAvatar && (
                  <div className="settings-avatar-overlay">
                    <span>📷</span>
                  </div>
                )}
              </div>
              <input
                ref={avatarInputRef}
                type="file"
                accept="image/*"
                onChange={handleAvatarSelect}
                style={{ display: 'none' }}
              />
            </div>

            {selectedAvatar && (
              <div className="settings-actions">
                <button
                  className="settings-save-btn"
                  onClick={handleSaveAvatar}
                  disabled={isAvatarUploading}
                >
                  {isAvatarUploading ? 'Загрузка...' : 'Сохранить изменения'}
                </button>
                <button
                  className="settings-cancel-btn"
                  onClick={handleCancelAvatar}
                  disabled={isAvatarUploading}
                >
                  Отменить изменения
                </button>
              </div>
            )}
          </div>
          <h1 className='nickname-settings'>Вы: {nickname}</h1>
          
          <div className="settings-item">
            <label className="settings-label">Никнейм</label>
            {isEditingNickname ? (
              <>
                <input
                  type="text"
                  className="settings-input"
                  value={nickname}
                  onChange={(e) => setNickname(e.target.value)}
                  placeholder="Введите новый никнейм"
                />
                <div className="settings-actions">
                  <button className="settings-save-btn" onClick={handleSaveNickname}>
                    Сохранить изменения
                  </button>
                  <button
                    className="settings-cancel-btn"
                    onClick={() => {
                      setIsEditingNickname(false);
                      setNickname(user?.nickname || '');
                    }}
                  >
                    Отменить
                  </button>
                </div>
              </>
            ) : (
              <div className="settings-value" onClick={() => setIsEditingNickname(true)}>
                {user?.nickname || 'Не указан'}
              </div>
            )}
          </div>
          
          <div className="settings-item">
            <label className="settings-label">Описание профиля</label>
            {isEditingDescriptionProfile ? (
              <>
                <textarea
                  className="settings-input settings-textarea"
                  value={description}
                  onChange={(e) => setDescription(e.target.value)}
                  placeholder="Введите описание профиля. Например: Люблю фармить вирты ^_^"
                  maxLength={200}
                  rows={3}
                />
                <div className="char-counter">
                  {description.length}/200 символов
                </div>
                <div className="settings-actions">
                  <button className="settings-save-btn" onClick={handleSaveDescription}>
                    Сохранить изменения
                  </button>
                  <button
                    className="settings-cancel-btn"
                    onClick={() => {
                      setIsEditingDescriptionProfile(false);
                      setDescription(user?.user_description || '');
                    }}
                  >
                    Отменить
                  </button>
                </div>
              </>
            ) : (
              <div className="settings-value" onClick={() => setIsEditingDescriptionProfile(true)}>
                {user?.user_description || 'Обычный бродяга по сайту Arizona Games Store ^_^'}
              </div>
            )}
          </div>
          
          <div className="settings-item">
            <label className="settings-label">Email</label>
            {isEditingEmail ? (
              <>
                <input
                  type="email"
                  className="settings-input"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  placeholder="Введите новый email"
                />
                <div className="settings-actions">
                  <button className="settings-save-btn" onClick={handleSaveEmail}>
                    Сохранить изменения
                  </button>
                  <button
                    className="settings-cancel-btn"
                    onClick={() => {
                      setIsEditingEmail(false);
                      setEmail(user?.email || '');
                    }}
                  >
                    Отменить
                  </button>
                </div>
              </>
            ) : (
              <div className="settings-value" onClick={() => setIsEditingEmail(true)}>
                {user?.email || 'Не указан'}
              </div>
            )}
          </div>

          
          <div className="settings-item">
            <label className="settings-label">Telegram</label>
            {isEditingTelegram ? (
              <>
                <input
                  type="text"
                  className="settings-input"
                  value={telegram}
                  onChange={(e) => setTelegram(e.target.value)}
                  placeholder="Введите Telegram (например: @username)"
                />
                <div className="settings-actions">
                  <button className="settings-save-btn" onClick={handleSaveTelegram}>
                    Сохранить изменения
                  </button>
                  <button
                    className="settings-cancel-btn"
                    onClick={() => {
                      setIsEditingTelegram(false);
                      setTelegram(user?.telegram || '');
                    }}
                  >
                    Отменить
                  </button>
                </div>
              </>
            ) : (
              <div className="settings-value" onClick={() => setIsEditingTelegram(true)}>
                {user?.telegram || 'Не указан'}
              </div>
            )}
          </div>

          
          <div className="settings-item">
            <label className="settings-label">Пароль</label>
            {isEditingPassword ? (
              <>
                <input
                  type="password"
                  className="settings-input"
                  value={passwordData.oldPassword}
                  onChange={(e) =>
                    setPasswordData({ ...passwordData, oldPassword: e.target.value })
                  }
                  placeholder="Старый пароль"
                />
                <input
                  type="password"
                  className="settings-input"
                  value={passwordData.newPassword}
                  onChange={(e) =>
                    setPasswordData({ ...passwordData, newPassword: e.target.value })
                  }
                  placeholder="Новый пароль"
                />
                <input
                  type="password"
                  className="settings-input"
                  value={passwordData.confirmPassword}
                  onChange={(e) =>
                    setPasswordData({ ...passwordData, confirmPassword: e.target.value })
                  }
                  placeholder="Подтвердите новый пароль"
                />
                <div className="settings-actions">
                  <button className="settings-save-btn" onClick={handleSavePassword}>
                    Сохранить изменения
                  </button>
                  <button
                    className="settings-cancel-btn"
                    onClick={() => {
                      setIsEditingPassword(false);
                      setPasswordData({ oldPassword: '', newPassword: '', confirmPassword: '' });
                    }}
                  >
                    Отменить
                  </button>
                </div>
              </>
            ) : (
              <div className="settings-value" onClick={() => setIsEditingPassword(true)}>
                ••••••••
              </div>
            )}
          </div>
        </div>

        
        <div className="settings-block">
          <h2 className="settings-block-title">Приложение</h2>
          <p className="settings-block-description">Выберите тему приложения</p>
          <label className="settings-label">Тема</label>
          <div className="settings-theme-buttons">
            <button
              className={`settings-theme-btn ${theme === 'dark' ? 'active' : ''}`}
              onClick={() => handleThemeChange('dark')}
            >
              <span className="theme-icon">🌙</span>
              <span>Тёмная</span>
            </button>
            <button
              className={`settings-theme-btn ${theme === 'light' ? 'active' : ''}`}
              onClick={() => handleThemeChange('light')}
            >
              <span className="theme-icon">☀️</span>
              <span>Светлая</span>
            </button>
          </div>
        </div>

        
        <div className="settings-block">
          <button className="settings-logout-btn" onClick={handleLogout}>
            Выйти из профиля
          </button>
        </div>
      </div>

      <BottomNavigation />

      {toast && <Toast message={toast.message} type={toast.type} onClose={() => setToast(null)} />}
    </div>
  );
}

export default Settings;
