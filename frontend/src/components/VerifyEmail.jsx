import React, { useState, useEffect } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { useRateLimit } from '../hooks/useRateLimit';
import SnowEffect from './SnowEffect';
import '../styles/Auth.css';

function VerifyEmail() {
  const [code, setCode] = useState('');
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();
  const location = useLocation();
  const { pendingVerificationEmail, clearVerificationFlow, checkAuth } = useAuth();
  const email = location.state?.email || pendingVerificationEmail || '';
  const { executeWithRateLimit: executeVerify, isOnCooldown: isVerifyCooldown, remainingTime: verifyRemaining } = useRateLimit(2);
  const { executeWithRateLimit: executeResend, isOnCooldown: isResendCooldown, remainingTime: resendRemaining } = useRateLimit(5);


  useEffect(() => {
    document.title = 'Arz Store | Подтверждение email';
  }, []);

  const handleVerify = async () => {
    if (!code) {
      setError('Введите код подтверждения.');
      return;
    }

    if (code.length !== 6) {
      setError('Код должен состоять из 6 цифр.');
      return;
    }

    if (isVerifyCooldown) {
      setError(`Подождите ${verifyRemaining} сек. перед следующей попыткой`);
      return;
    }

    setLoading(true);
    setError('');
    setSuccess('');

    const result = await executeVerify(async () => {
      try {

        let clientIP = '';
        try {
          const ipResponse = await fetch('https://api.seeip.org/jsonip?');
          if (ipResponse.ok) {
            const ipData = await ipResponse.json();
            clientIP = ipData.ip || '';
          }
        } catch (ipError) {
          console.warn('Не удалось получить IP адрес:', ipError);

        }

        const response = await fetch('http://localhost:8080/api/verify-email', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          credentials: 'include',
          body: JSON.stringify({ email, code, client_ip: clientIP }),
        });

        if (response.ok) {
          setSuccess('Email успешно подтвержден!');
          clearVerificationFlow();
          await checkAuth();
          setTimeout(() => {
            navigate('/feed');
          }, 2000);
        } else {
          const errorData = await response.json();
          setError(errorData.error || 'Не удалось подтвердить код.');
        }
        return { success: response.ok };
      } catch (err) {
        setError('Не удалось проверить код. Попробуйте снова.');
        return { success: false };
      } finally {
        setLoading(false);
      }
    });

    if (result.rateLimited) {
      setError(`Подождите ${verifyRemaining} сек. перед следующей попыткой`);
    }
  };

  const handleResend = async () => {
    if (isResendCooldown) {
      setError(`Подождите ${resendRemaining} сек. перед повторной отправкой`);
      return;
    }

    setLoading(true);
    setError('');
    setSuccess('');

    const result = await executeResend(async () => {
      try {
        const response = await fetch('http://localhost:8080/api/resend-code', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          credentials: 'include',
          body: JSON.stringify({ email }),
        });

        if (response.ok) {
          setSuccess('Новый код отправлен на почту.');
        } else {
          const errorData = await response.json();
          setError(errorData.error || 'Не удалось отправить код повторно.');
        }
        return { success: response.ok };
      } catch (err) {
        setError('Не удалось отправить запрос повторно.');
        return { success: false };
      } finally {
        setLoading(false);
      }
    });

    if (result.rateLimited) {
      setError(`Подождите ${resendRemaining} сек. перед повторной отправкой`);
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

        <h1 className="auth-title">Подтвердите Email</h1>

        <p className="auth-description">
          Введите 6-значный код, отправленный на {email}. Проверьте папку спам, если письмо не пришло сразу.
        </p>

        {error && <div className="error-message">{error}</div>}
        {success && <div className="success-message">{success}</div>}

        <div className="input-group">
          <div className="input-wrapper">
            <input
              type="text"
              placeholder="Код подтверждения"
              className="auth-input"
              value={code}
              onChange={(e) => setCode(e.target.value.replace(/\D/g, '').slice(0, 6))}
              maxLength={6}
              onKeyPress={(e) => e.key === 'Enter' && handleVerify()}
            />
          </div>
        </div>

        <button
          className="btn-primary"
          onClick={handleVerify}
          disabled={loading || isVerifyCooldown}
        >
          {loading ? 'Проверка...' : isVerifyCooldown ? `ПОДОЖДИТЕ (${verifyRemaining}с)` : 'Подтвердить'}
        </button>

        <button
          className="btn-secondary"
          onClick={handleResend}
          disabled={loading || isResendCooldown}
        >
          {isResendCooldown ? `ПОДОЖДИТЕ (${resendRemaining}с)` : 'Отправить код повторно'}
        </button>

        <p className="footer-text">
          Не получили письмо? Проверьте, что email указан верно, и попробуйте запросить код еще раз. Также можно проверить папку спам.
        </p>
      </div>
    </div>
  );
}

export default VerifyEmail;
