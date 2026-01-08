# Система JWT аутентификации

## Описание

Реализована двухуровневая система JWT токенов:
- **Access Token** - живет 3 минуты, используется для доступа к защищенным эндпоинтам
- **Refresh Token** - живет 30 дней, используется для автоматического обновления access токена

## Как это работает

1. При регистрации или логине пользователь получает оба токена в cookie
2. Access token автоматически обновляется через refresh token когда истекает
3. Пользователь не выкидывается из системы при истечении access token (3 минуты)
4. Только через 30 дней (истечение refresh token) требуется повторная авторизация

## API Endpoints

### Регистрация
```
POST /api/register
Content-Type: application/json

{
  "nickname": "username",
  "password": "password123"
}
```

Ответ: устанавливает `access_token` и `refresh_token` в cookie

### Логин
```
POST /api/login
Content-Type: application/json

{
  "nickname": "username",
  "password": "password123"
}
```

Ответ: устанавливает `access_token` и `refresh_token` в cookie

### Обновление токена
```
POST /api/refresh
```

Использует `refresh_token` из cookie для генерации нового `access_token`

### Выход
```
POST /api/logout
```

Удаляет токены из cookie и базы данных

## Использование middleware

Для защиты эндпоинтов используйте middleware `AuthRequired()`:

```go
import "arizonagamesstore/backend/middleware"

router.GET("/api/protected", middleware.AuthRequired(), handlers.ProtectedHandler)
```

Внутри обработчика доступ к данным пользователя:
```go
func ProtectedHandler(c *gin.Context) {
    userID := c.GetUint("user_id")
    nickname := c.GetString("nickname")

    c.JSON(200, gin.H{
        "user_id": userID,
        "nickname": nickname,
    })
}
```

## Переменные окружения

Добавьте в `.env`:
```
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_REFRESH_SECRET=your-super-secret-refresh-jwt-key-change-this-in-production
```

## Валидация данных

### Никнейм
- Длина: 3-20 символов
- Только английские буквы, цифры и подчёркивание (_)
- Должен начинаться с буквы
- Без пробелов

### Пароль
- Длина: 6-100 символов
- Разрешены: буквы, цифры и специальные символы (!@#$%^&*()_+={}[]|:";',./?~ и другие)
- Запрещены опасные символы: < > (HTML теги)
- Запрещены SQL/JS конструкции: --, /*, */, script, eval, и т.д.

## Безопасность

- Токены хранятся в HttpOnly cookie (защита от XSS)
- Refresh токены хранятся в базе данных для валидации
- При выходе refresh token удаляется из БД
- Используются разные секреты для access и refresh токенов
- Валидация защищает от SQL и XSS инъекций
