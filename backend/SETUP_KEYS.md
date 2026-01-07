# Инструкция по получению API ключей

## 1. Google reCAPTCHA v3

### Получение ключей:

1. Перейдите на https://www.google.com/recaptcha/admin/create
2. Войдите в Google аккаунт
3. Заполните форму:
   - **Label (Название)**: `Arizona Games Store`
   - **reCAPTCHA type**: Выберите **reCAPTCHA v3**
   - **Domains**: Добавьте:
     - `localhost`
     - `127.0.0.1`
     - ваш домен (если есть)
4. Нажмите **Submit**
5. Скопируйте ключи:
   - **Site Key** (для frontend)
   - **Secret Key** (для backend)

### Куда вставить ключи:

#### Backend (`backend/.env`):
```env
RECAPTCHA_SECRET_KEY=ваш_secret_key_сюда
```

#### Frontend (создайте `.env` в папке `frontend`):
```env
VITE_RECAPTCHA_SITE_KEY=ваш_site_key_сюда
```

---

## 2. Email (SMTP)

Для отправки email нужен SMTP сервер. Рекомендуемые варианты:

### Вариант 1: Gmail (Бесплатно, но с лимитами)

1. Перейдите в настройки Google аккаунта: https://myaccount.google.com/
2. **Security** → **2-Step Verification** (включите, если выключено)
3. **Security** → **App passwords**
4. Создайте новый App Password:
   - **Select app**: Mail
   - **Select device**: Other (custom name)
   - **Name**: `Arizona Games Store`
5. Скопируйте сгенерированный пароль (16 символов без пробелов)

**Лимиты Gmail:**
- 500 писем в день
- 100-150 писем в час

#### Куда вставить (`backend/.env`):
```env
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=ваш_email@gmail.com
SMTP_PASSWORD=сгенерированный_app_password
SMTP_FROM=ваш_email@gmail.com
SMTP_FROM_NAME=Arizona Games Store
```

---

### Вариант 2: SendGrid (Рекомендуется для продакшена)

1. Зарегистрируйтесь на https://sendgrid.com/ (бесплатно до 100 писем/день)
2. Создайте API Key:
   - **Settings** → **API Keys** → **Create API Key**
   - **Name**: `Arizona Games Store`
   - **Permissions**: Full Access (или Mail Send только)
3. Скопируйте API Key

#### Куда вставить (`backend/.env`):
```env
SMTP_HOST=smtp.sendgrid.net
SMTP_PORT=587
SMTP_USER=apikey
SMTP_PASSWORD=ваш_sendgrid_api_key
SMTP_FROM=noreply@yourdomain.com
SMTP_FROM_NAME=Arizona Games Store
```

---

### Вариант 3: Mailgun (Альтернатива)

1. Зарегистрируйтесь на https://www.mailgun.com/ (бесплатно до 5000 писем/месяц)
2. Добавьте домен или используйте sandbox домен
3. **Sending** → **Domain Settings** → **SMTP credentials**
4. Создайте SMTP credentials

#### Куда вставить (`backend/.env`):
```env
SMTP_HOST=smtp.mailgun.org
SMTP_PORT=587
SMTP_USER=postmaster@ваш_sandbox_домен
SMTP_PASSWORD=ваш_smtp_password
SMTP_FROM=noreply@ваш_sandbox_домен
SMTP_FROM_NAME=Arizona Games Store
```

---

## Финальный `.env` файл (backend)

После получения всех ключей, ваш `backend/.env` должен выглядеть так:

```env
# Database
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=13012005
DB_NAME=arzgamesstore_db
DB_PORT=5432

# JWT
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_REFRESH_SECRET=your-super-secret-refresh-jwt-key-change-this-in-production

# reCAPTCHA
RECAPTCHA_SECRET_KEY=ваш_recaptcha_secret_key

# SMTP Email
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=ваш_email@gmail.com
SMTP_PASSWORD=ваш_app_password_или_api_key
SMTP_FROM=ваш_email@gmail.com
SMTP_FROM_NAME=Arizona Games Store
```

## Frontend `.env` файл

Создайте файл `frontend/.env`:

```env
VITE_RECAPTCHA_SITE_KEY=ваш_recaptcha_site_key
```

---

## Проверка настроек

### Backend (в консоли):
```bash
cd backend
go run main.go
```

Должно быть без ошибок про отсутствующие env переменные.

### Frontend:
```bash
cd frontend
npm run dev
```

reCAPTCHA badge должен появиться в правом нижнем углу страницы.

---

## Тестирование

1. Зарегистрируйте тестовый аккаунт с реальным email
2. Проверьте почту - должно прийти письмо с кодом подтверждения
3. Введите код для подтверждения email
4. reCAPTCHA должна работать автоматически (v3 работает в фоне)

---

## Troubleshooting

### Письма не приходят:
1. Проверьте папку Spam
2. Убедитесь, что SMTP настройки верны
3. Проверьте логи backend на ошибки отправки

### reCAPTCHA не работает:
1. Проверьте, что site key правильный
2. Убедитесь, что домен добавлен в reCAPTCHA консоли
3. Проверьте консоль браузера на ошибки

### Rate limit при тестировании:
- Используйте разные IP (например, мобильный интернет)
- Или временно увеличьте лимиты в `middleware/ratelimit.go`
