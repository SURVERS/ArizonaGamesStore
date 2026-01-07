# Инструкция по применению миграций

## Миграция: Добавление полей настроек

Эта миграция добавляет поля для функционала настроек профиля.

### Что добавляется:
- `theme` - тема оформления (dark/light), по умолчанию 'dark'
- `last_settings_change` - время последнего изменения настроек (для cooldown 1 минута)
- `last_nickname_change` - время последнего изменения никнейма (для cooldown 1 неделя)
- `last_email_change` - время последнего изменения email (для cooldown 1 неделя)

### Применение миграции:

#### Вариант 1: Через командную строку PostgreSQL
```bash
psql -U ваш_пользователь -d название_базы -f migrations/add_settings_fields.sql
```

#### Вариант 2: Через DBeaver или другой GUI клиент
1. Откройте файл `migrations/add_settings_fields.sql`
2. Скопируйте содержимое
3. Выполните SQL в вашей базе данных

#### Вариант 3: Используя Go migrate (если установлен)
```bash
migrate -path migrations -database "postgresql://user:password@localhost:5432/dbname?sslmode=disable" up
```

### Проверка:
После применения миграции выполните:
```sql
\d accounts
```
или
```sql
SELECT column_name, data_type, column_default
FROM information_schema.columns
WHERE table_name = 'accounts'
AND column_name IN ('theme', 'last_settings_change', 'last_nickname_change', 'last_email_change');
```

Вы должны увидеть новые поля в таблице accounts.

### Откат миграции (при необходимости):
```sql
ALTER TABLE accounts DROP COLUMN IF EXISTS theme;
ALTER TABLE accounts DROP COLUMN IF EXISTS last_settings_change;
ALTER TABLE accounts DROP COLUMN IF EXISTS last_nickname_change;
ALTER TABLE accounts DROP COLUMN IF EXISTS last_email_change;
```
