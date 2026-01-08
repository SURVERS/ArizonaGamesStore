import React, { createContext, useState, useContext, useEffect, useCallback } from 'react';

const AuthContext = createContext(null);

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const PENDING_VERIFICATION_KEY = 'pendingVerificationEmail';
  const getInitialVerificationEmail = () => {
    if (typeof window === 'undefined') {
      return '';
    }
    return sessionStorage.getItem(PENDING_VERIFICATION_KEY) || '';
  };
  const [pendingVerificationEmail, setPendingVerificationEmail] = useState(getInitialVerificationEmail);

  const persistVerificationEmail = useCallback((email) => {
    if (typeof window !== 'undefined') {
      sessionStorage.setItem(PENDING_VERIFICATION_KEY, email);
    }
    setPendingVerificationEmail(email);
  }, []);

  const clearVerificationFlow = useCallback(() => {
    if (typeof window !== 'undefined') {
      sessionStorage.removeItem(PENDING_VERIFICATION_KEY);
    }
    setPendingVerificationEmail('');
  }, []);

  const checkAuth = async () => {
    try {
      const response = await fetch('http://localhost:8080/api/me', {
        credentials: 'include',
      });

      if (response.ok) {
        const data = await response.json();
        clearVerificationFlow();
        setUser(data);
      } else {
        setUser(null);
      }
    } catch (error) {
      setUser(null);
    } finally {
      setLoading(false);
    }
  };

  const login = async (nickname, password, recaptchaToken) => {

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

    const response = await fetch('http://localhost:8080/api/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify({ nickname, password, recaptcha_token: recaptchaToken, client_ip: clientIP }),
    });

    if (response.ok) {
      const data = await response.json();
      clearVerificationFlow();

      await checkAuth();
      return { success: true, data };
    } else {
      const errorData = await response.json();
      return { success: false, error: errorData.error || 'Ошибка входа' };
    }
  };

  const register = async (nickname, password, email, recaptchaToken) => {
    const response = await fetch('http://localhost:8080/api/register', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify({ nickname, password, email, recaptcha_token: recaptchaToken }),
    });

    if (response.ok) {
      const data = await response.json();
      if (!data.requires_verify) {
        clearVerificationFlow();

        await checkAuth();
      }
      if (data.requires_verify) {
        persistVerificationEmail(email);
      }
      return { success: true, data };
    } else {
      const errorData = await response.json();
      return { success: false, error: errorData.error || 'Ошибка регистрации' };
    }
  };

  const logout = async () => {
    try {
      await fetch('http://localhost:8080/api/logout', {
        method: 'POST',
        credentials: 'include',
      });
    } catch (error) {
      console.error('Logout error:', error);
    } finally {
      clearVerificationFlow();
      setUser(null);
    }
  };

  useEffect(() => {
    checkAuth();
  }, []);


  useEffect(() => {
    if (user?.theme) {
      if (user.theme === 'light') {
        document.body.classList.add('light-theme');
        document.body.classList.remove('dark-theme');
      } else {
        document.body.classList.add('dark-theme');
        document.body.classList.remove('light-theme');
      }
    } else {

      document.body.classList.add('dark-theme');
      document.body.classList.remove('light-theme');
    }
  }, [user?.theme]);

  return (
    <AuthContext.Provider
      value={{
        user,
        loading,
        login,
        register,
        logout,
        checkAuth,
        pendingVerificationEmail,
        clearVerificationFlow,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within AuthProvider');
  }
  return context;
};
