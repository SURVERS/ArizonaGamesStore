package docs

// Регистрация нового пользователя
// @Summary Регистрация
// @Description Создает нового пользователя. Сразу после регистрации на почту придет код подтверждения. Пароль должен быть минимум 8 символов, иначе хакеры взломают за 5 минут
// @Tags Аутентификация
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Данные для регистрации"
// @Success 201 {object} SuccessResponse "Аккаунт создан! Проверь почту и введи код подтверждения"
// @Failure 400 {object} ErrorResponse "Не хватает данных или формат неправильный"
// @Failure 409 {object} ErrorResponse "Такой ник или email уже занят, придумай другой"
// @Failure 429 {object} ErrorResponse "Слишком много попыток, подожди немного"
// @Failure 500 {object} ErrorResponse "Что-то сломалось на сервере, пишите в поддержку"
// @Router /register [post]
func RegisterEndpoint() {}

// Вход в систему
// @Summary Авторизация
// @Description Авторизует пользователя и возвращает JWT токены. Access токен живет 15 минут, refresh - 7 дней. Токены приходят в куках (http-only), так что в localStorage их не тырить не получится
// @Tags Аутентификация
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Логин и пароль"
// @Success 200 {object} SuccessResponse "Добро пожаловать! Токены уже в твоих куках"
// @Failure 400 {object} ErrorResponse "Не хватает логина или пароля"
// @Failure 401 {object} ErrorResponse "Неправильный ник или пароль, попробуй еще раз"
// @Failure 403 {object} ErrorResponse "Сначала подтверди email, проверь почту"
// @Failure 429 {object} ErrorResponse "Слишком много попыток входа, отдохни"
// @Failure 500 {object} ErrorResponse "Сервер приказал долго жить"
// @Router /login [post]
func LoginEndpoint() {}

// Обновление access токена
// @Summary Обновить токен
// @Description Обновляет access токен используя refresh токен из куки. Используй когда получишь 401 ошибку
// @Tags Аутентификация
// @Produce json
// @Success 200 {object} SuccessResponse "Токен обновлен, можешь продолжать работать"
// @Failure 401 {object} ErrorResponse "Refresh токен протух или невалидный, залогинься заново"
// @Failure 500 {object} ErrorResponse "Че-то сломалось"
// @Router /refresh [post]
func RefreshTokenEndpoint() {}

// Выход из системы
// @Summary Выход
// @Description Удаляет все токены из куки. После этого придется логиниться заново
// @Tags Аутентификация
// @Produce json
// @Success 200 {object} SuccessResponse "Токены удалены, пока!"
// @Router /logout [post]
func LogoutEndpoint() {}

// Подтверждение email
// @Summary Подтвердить email
// @Description Подтверждает email введя 6-значный код из письма. Код действителен 10 минут
// @Tags Аутентификация
// @Accept json
// @Produce json
// @Param request body VerifyEmailRequest true "Email и код подтверждения"
// @Success 200 {object} SuccessResponse "Email подтвержден! Теперь можешь пользоваться сайтом"
// @Failure 400 {object} ErrorResponse "Неправильный код или формат"
// @Failure 404 {object} ErrorResponse "Такого email не существует"
// @Failure 429 {object} ErrorResponse "Слишком много попыток, подожди"
// @Failure 500 {object} ErrorResponse "Проблемы на сервере"
// @Router /verify-email [post]
func VerifyEmailEndpoint() {}

// Переотправка кода подтверждения
// @Summary Переотправить код
// @Description Отправляет новый код подтверждения на email. Полезно если код протух или потерялся
// @Tags Аутентификация
// @Accept json
// @Produce json
// @Param request body map[string]string true "Email для переотправки" example(email="gamer@arizona.rp")
// @Success 200 {object} SuccessResponse "Новый код отправлен на почту!"
// @Failure 400 {object} ErrorResponse "Не указан email"
// @Failure 404 {object} ErrorResponse "Аккаунт с таким email не найден"
// @Failure 429 {object} ErrorResponse "Подожди немного перед следующей отправкой"
// @Failure 500 {object} ErrorResponse "Ошибка отправки письма"
// @Router /resend-code [post]
func ResendCodeEndpoint() {}

// Создание объявления
// @Summary Создать объявление
// @Description Создает новое объявление с картинкой. Картинка загружается на AWS S3. После создания объявление автоматически удалится через 48 часов (можно продлить). Между созданиями объявлений нужно ждать 60 секунд
// @Tags Объявления
// @Accept multipart/form-data
// @Produce json
// @Param server formData string true "Сервер (ViceCity, Phoenix, и т.д.)"
// @Param title formData string true "Название (макс. 25 символов)"
// @Param description formData string true "Описание (макс. 500 символов)"
// @Param type formData string true "Тип (Продать/Купить/Сдать в аренду)"
// @Param currency formData string true "Валюта (VC/$/BTC/EURO/Договорная)"
// @Param price formData number true "Цена"
// @Param category formData string true "Категория (house/business/vehicle/security/accs/others)"
// @Param nickname formData string true "Никнейм автора"
// @Param imagePath formData string true "Путь для сохранения картинки"
// @Param image formData file true "Изображение (макс. 10MB, разрешение 300x200 - 1920x1080)"
// @Param rentalHoursLimit formData int false "Лимит часов аренды (1-180)"
// @Success 200 {object} SuccessResponse "Объявление создано! ID: 42"
// @Failure 400 {object} ErrorResponse "Не хватает данных или картинка кривая"
// @Failure 413 {object} ErrorResponse "Картинка слишком большая (макс. 10MB)"
// @Failure 429 {object} ErrorResponse "Подожди 60 секунд перед созданием нового объявления"
// @Failure 500 {object} ErrorResponse "Ошибка загрузки на S3 или БД"
// @Router /createnewads [post]
func CreateAdEndpoint() {}

// Получение объявлений
// @Summary Список объявлений
// @Description Возвращает список объявлений с фильтрацией и сортировкой. По умолчанию возвращает 20 штук, можно подгружать дальше через offset
// @Tags Объявления
// @Produce json
// @Param category query string true "Категория" Enums(house, business, vehicle, security, accs, others)
// @Param server query string false "Фильтр по серверу"
// @Param limit query int false "Сколько объявлений вернуть (по умолчанию 20)"
// @Param offset query int false "Сколько пропустить для пагинации (по умолчанию 0)"
// @Param sort query string false "Сортировка" Enums(date_desc, date_asc, price_desc, price_asc, views_desc)
// @Param type query string false "Фильтр по типу" Enums(Продать, Купить, Сдать в аренду)
// @Param currency query string false "Фильтр по валюте" Enums(VC, $, BTC, EURO, Договорная)
// @Param price_min query number false "Минимальная цена"
// @Param price_max query number false "Максимальная цена"
// @Success 200 {object} map[string]interface{} "Список объявлений"
// @Failure 400 {object} ErrorResponse "Не указана категория"
// @Failure 500 {object} ErrorResponse "Ошибка БД"
// @Router /ads [get]
func GetAdsEndpoint() {}

// Случайные объявления
// @Summary Случайные объявления
// @Description Возвращает 8 случайных объявлений для главной страницы. Каждый раз разные!
// @Tags Объявления
// @Produce json
// @Success 200 {array} AdResponse "Массив случайных объявлений"
// @Failure 500 {object} ErrorResponse "Что-то пошло не так"
// @Router /ads/random [get]
func GetRandomAdsEndpoint() {}

// Объявления пользователя
// @Summary Объявления по нику
// @Description Возвращает все объявления конкретного пользователя
// @Tags Объявления
// @Produce json
// @Param nickname path string true "Никнейм пользователя"
// @Success 200 {array} AdResponse "Список объявлений пользователя"
// @Failure 500 {object} ErrorResponse "Ошибка БД"
// @Router /listings/user/{nickname} [get]
func GetAdsByNicknameEndpoint() {}

// Количество объявлений
// @Summary Счетчик объявлений
// @Description Возвращает количество объявлений в категории. Нужно для отображения "Всего объявлений: 420"
// @Tags Объявления
// @Produce json
// @Param CategoryName query string true "Название категории"
// @Success 200 {object} map[string]int "Количество объявлений"
// @Failure 500 {object} ErrorResponse "Ошибка подсчета"
// @Router /getadcount [get]
func GetAdCountEndpoint() {}

// Увеличить просмотры
// @Summary Записать просмотр
// @Description Увеличивает счетчик просмотров объявления. Вызывай когда пользователь открывает карточку объявления
// @Tags Объявления
// @Produce json
// @Param id path int true "ID объявления"
// @Success 200 {object} SuccessResponse "Просмотр засчитан!"
// @Failure 404 {object} ErrorResponse "Объявление не найдено"
// @Failure 500 {object} ErrorResponse "Ошибка обновления"
// @Router /ads/{id}/view [post]
func IncrementAdViewsEndpoint() {}

// Обновить объявление
// @Summary Обновить объявление
// @Description Обновляет данные объявления. Доступно только автору объявления
// @Tags Объявления
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID объявления"
// @Param request body CreateAdRequest true "Новые данные"
// @Success 200 {object} SuccessResponse "Объявление обновлено!"
// @Failure 400 {object} ErrorResponse "Некорректные данные"
// @Failure 401 {object} ErrorResponse "Не авторизован"
// @Failure 403 {object} ErrorResponse "Это не твое объявление!"
// @Failure 404 {object} ErrorResponse "Объявление не найдено"
// @Failure 500 {object} ErrorResponse "Ошибка обновления"
// @Router /ads/{id} [put]
func UpdateAdEndpoint() {}

// Удалить объявление
// @Summary Удалить объявление
// @Description Удаляет объявление. Доступно только автору
// @Tags Объявления
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID объявления"
// @Success 200 {object} SuccessResponse "Объявление удалено!"
// @Failure 401 {object} ErrorResponse "Не авторизован"
// @Failure 403 {object} ErrorResponse "Это не твое объявление!"
// @Failure 404 {object} ErrorResponse "Объявление не найдено"
// @Failure 500 {object} ErrorResponse "Ошибка удаления"
// @Router /ads/{id} [delete]
func DeleteAdEndpoint() {}

// Создать жалобу
// @Summary Пожаловаться на объявление
// @Description Отправляет жалобу на объявление. Причины: Мошенничество, Спам, Порнография, и т.д. Жалобы проверяются модераторами
// @Tags Жалобы
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateReportRequest true "Данные жалобы"
// @Success 200 {object} SuccessResponse "Жалоба отправлена! Мы её рассмотрим"
// @Failure 400 {object} ErrorResponse "Не указана причина или описание слишком короткое"
// @Failure 401 {object} ErrorResponse "Залогинься чтобы жаловаться"
// @Failure 404 {object} ErrorResponse "Объявление не найдено"
// @Failure 500 {object} ErrorResponse "Ошибка создания жалобы"
// @Router /reports [post]
func CreateReportEndpoint() {}

// Создать отзыв
// @Summary Оставить отзыв
// @Description Оставляет отзыв о продавце. Рейтинг от 1 до 5 звезд. Отзыв нужно подтвердить продавцу, чтобы он отобразился
// @Tags Отзывы
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateFeedbackRequest true "Данные отзыва"
// @Success 200 {object} SuccessResponse "Отзыв отправлен! Ждем подтверждения продавца"
// @Failure 400 {object} ErrorResponse "Некорректный рейтинг или комментарий"
// @Failure 401 {object} ErrorResponse "Нужна авторизация"
// @Failure 409 {object} ErrorResponse "Ты уже оставлял отзыв этому продавцу"
// @Failure 500 {object} ErrorResponse "Ошибка сохранения"
// @Router /feedback [post]
func CreateFeedbackEndpoint() {}

// Получить отзывы
// @Summary Отзывы продавца
// @Description Возвращает все отзывы о продавце
// @Tags Отзывы
// @Produce json
// @Param nickname path string true "Никнейм продавца"
// @Success 200 {array} FeedbackResponse "Список отзывов"
// @Failure 500 {object} ErrorResponse "Ошибка загрузки"
// @Router /feedback/{nickname} [get]
func GetFeedbacksEndpoint() {}

// Подтвердить отзыв
// @Summary Подтвердить отзыв
// @Description Подтверждает отзыв. Только продавец может подтвердить отзыв о себе. После подтверждения рейтинг пересчитывается
// @Tags Отзывы
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID отзыва"
// @Success 200 {object} SuccessResponse "Отзыв подтвержден! Твой рейтинг обновлен"
// @Failure 401 {object} ErrorResponse "Не авторизован"
// @Failure 403 {object} ErrorResponse "Ты не можешь подтвердить этот отзыв"
// @Failure 404 {object} ErrorResponse "Отзыв не найден"
// @Failure 500 {object} ErrorResponse "Ошибка подтверждения"
// @Router /feedback/{id}/confirm [put]
func ConfirmFeedbackEndpoint() {}

// Добавить просмотренное объявление
// @Summary Добавить в историю
// @Description Добавляет объявление в историю просмотров пользователя
// @Tags Просмотренное
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body map[string]int true "ID объявления" example(ad_id=42)
// @Success 200 {object} SuccessResponse "Добавлено в историю"
// @Failure 400 {object} ErrorResponse "Не указан ID"
// @Failure 401 {object} ErrorResponse "Нужна авторизация"
// @Failure 500 {object} ErrorResponse "Ошибка сохранения"
// @Router /viewed-ads [post]
func AddViewedAdEndpoint() {}

// Получить просмотренные объявления
// @Summary История просмотров
// @Description Возвращает историю просмотренных объявлений пользователя
// @Tags Просмотренное
// @Security BearerAuth
// @Produce json
// @Success 200 {array} AdResponse "Список просмотренных объявлений"
// @Failure 401 {object} ErrorResponse "Не авторизован"
// @Failure 500 {object} ErrorResponse "Ошибка загрузки"
// @Router /viewed-ads [get]
func GetViewedAdsEndpoint() {}

// Обновить фон профиля
// @Summary Изменить фон профиля
// @Description Загружает новый фон для профиля на S3. Макс. размер 10MB
// @Tags Профиль
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param background formData file true "Изображение фона"
// @Success 200 {object} map[string]string "URL нового фона"
// @Failure 400 {object} ErrorResponse "Файл не загружен или слишком большой"
// @Failure 401 {object} ErrorResponse "Не авторизован"
// @Failure 500 {object} ErrorResponse "Ошибка загрузки на S3"
// @Router /profile/update-background [post]
func UpdateBackgroundEndpoint() {}

// Удалить фон профиля
// @Summary Удалить фон
// @Description Удаляет фон профиля
// @Tags Профиль
// @Security BearerAuth
// @Produce json
// @Success 200 {object} SuccessResponse "Фон удален!"
// @Failure 401 {object} ErrorResponse "Не авторизован"
// @Failure 500 {object} ErrorResponse "Ошибка удаления"
// @Router /profile/delete-background [delete]
func DeleteBackgroundEndpoint() {}

// Обновить аватар
// @Summary Изменить аватар
// @Description Загружает новый аватар на S3. Макс. размер 5MB
// @Tags Профиль
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param avatar formData file true "Изображение аватара"
// @Success 200 {object} map[string]string "URL нового аватара"
// @Failure 400 {object} ErrorResponse "Файл не загружен или слишком большой"
// @Failure 401 {object} ErrorResponse "Не авторизован"
// @Failure 500 {object} ErrorResponse "Ошибка загрузки"
// @Router /profile/update-avatar [post]
func UpdateAvatarEndpoint() {}

// Обновить никнейм
// @Summary Изменить никнейм
// @Description Изменяет никнейм пользователя. Макс. 20 символов
// @Tags Профиль
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body map[string]string true "Новый никнейм" example(nickname="NewNick123")
// @Success 200 {object} SuccessResponse "Никнейм обновлен!"
// @Failure 400 {object} ErrorResponse "Никнейм слишком длинный или пустой"
// @Failure 401 {object} ErrorResponse "Не авторизован"
// @Failure 409 {object} ErrorResponse "Такой никнейм уже занят"
// @Failure 500 {object} ErrorResponse "Ошибка обновления"
// @Router /profile/update-nickname [put]
func UpdateNicknameEndpoint() {}

// Обновить email
// @Summary Изменить email
// @Description Изменяет email пользователя. После изменения нужно будет заново подтвердить email
// @Tags Профиль
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body map[string]string true "Новый email" example(email="new@arizona.rp")
// @Success 200 {object} SuccessResponse "Email обновлен! Проверь почту для подтверждения"
// @Failure 400 {object} ErrorResponse "Некорректный email"
// @Failure 401 {object} ErrorResponse "Не авторизован"
// @Failure 409 {object} ErrorResponse "Такой email уже используется"
// @Failure 500 {object} ErrorResponse "Ошибка обновления"
// @Router /profile/update-email [put]
func UpdateEmailEndpoint() {}

// Обновить пароль
// @Summary Изменить пароль
// @Description Изменяет пароль пользователя. Нужно ввести старый пароль для подтверждения
// @Tags Профиль
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body map[string]string true "Старый и новый пароли" example(old_password="OldPass123!" new_password="NewPass123!")
// @Success 200 {object} SuccessResponse "Пароль обновлен!"
// @Failure 400 {object} ErrorResponse "Новый пароль слишком короткий (мин. 8 символов)"
// @Failure 401 {object} ErrorResponse "Старый пароль неверный"
// @Failure 500 {object} ErrorResponse "Ошибка обновления"
// @Router /profile/update-password [put]
func UpdatePasswordEndpoint() {}

// Обновить тему
// @Summary Изменить тему
// @Description Меняет тему оформления (светлая/темная)
// @Tags Профиль
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body map[string]string true "Тема" example(theme="dark")
// @Success 200 {object} SuccessResponse "Тема обновлена!"
// @Failure 401 {object} ErrorResponse "Не авторизован"
// @Failure 500 {object} ErrorResponse "Ошибка обновления"
// @Router /profile/update-theme [put]
func UpdateThemeEndpoint() {}

// Обновить описание
// @Summary Изменить описание
// @Description Обновляет описание профиля. Макс. 500 символов
// @Tags Профиль
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body map[string]string true "Описание" example(description="Активный продавец на Arizona RP")
// @Success 200 {object} SuccessResponse "Описание обновлено!"
// @Failure 400 {object} ErrorResponse "Описание слишком длинное"
// @Failure 401 {object} ErrorResponse "Не авторизован"
// @Failure 500 {object} ErrorResponse "Ошибка обновления"
// @Router /profile/update-description [put]
func UpdateDescriptionEndpoint() {}

// Обновить Telegram
// @Summary Изменить Telegram
// @Description Обновляет Telegram контакт. Макс. 50 символов
// @Tags Профиль
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body map[string]string true "Telegram" example(telegram="@coolgamer")
// @Success 200 {object} SuccessResponse "Telegram обновлен!"
// @Failure 400 {object} ErrorResponse "Telegram слишком длинный"
// @Failure 401 {object} ErrorResponse "Не авторизован"
// @Failure 500 {object} ErrorResponse "Ошибка обновления"
// @Router /profile/update-telegram [put]
func UpdateTelegramEndpoint() {}

// Получить данные текущего пользователя
// @Summary Мой профиль
// @Description Возвращает данные авторизованного пользователя
// @Tags Профиль
// @Security BearerAuth
// @Produce json
// @Success 200 {object} UserResponse "Данные пользователя"
// @Failure 401 {object} ErrorResponse "Не авторизован"
// @Failure 500 {object} ErrorResponse "Ошибка загрузки данных"
// @Router /me [get]
func GetMeEndpoint() {}
