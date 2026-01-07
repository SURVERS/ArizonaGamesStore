package docs

type RegisterRequest struct {
	Nickname string `json:"nickname" example:"CoolGamer123" binding:"required"`
	Email    string `json:"email" example:"gamer@arizona.rp" binding:"required,email"`
	Password string `json:"password" example:"SuperSecret123!" binding:"required,min=8"`
	RecaptchaToken string `json:"recaptcha_token" example:"03AGdBq26..." binding:"required"`
}

type LoginRequest struct {
	Nickname string `json:"nickname" example:"CoolGamer123" binding:"required"`
	Password string `json:"password" example:"SuperSecret123!" binding:"required"`
	RecaptchaToken string `json:"recaptcha_token" example:"03AGdBq26..." binding:"required"`
}

type VerifyEmailRequest struct {
	Email string `json:"email" example:"gamer@arizona.rp" binding:"required,email"`
	Code  string `json:"code" example:"123456" binding:"required,len=6"`
}

type CreateAdRequest struct {
	Server      string  `form:"server" example:"ViceCity" binding:"required"`
	Title       string  `form:"title" example:"Продам крутой дом" binding:"required,max=25"`
	Description string  `form:"description" example:"Дом в центре города, 3 этажа, гараж на 2 машины" binding:"required,max=500"`
	Type        string  `form:"type" example:"Продать" binding:"required"`
	Currency    string  `form:"currency" example:"VC" binding:"required"`
	Price       float64 `form:"price" example:"5000000" binding:"required"`
	Category    string  `form:"category" example:"house" binding:"required"`
	Nickname    string  `form:"nickname" example:"CoolGamer123" binding:"required"`
	ImagePath   string  `form:"imagePath" example:"ads/house/123456_user.jpg" binding:"required"`
	RentalHoursLimit int `form:"rentalHoursLimit" example:"24"`
}

type AdResponse struct {
	ID              int     `json:"id" example:"42"`
	Server          string  `json:"server" example:"ViceCity"`
	Title           string  `json:"title" example:"Продам крутой дом"`
	Description     string  `json:"description" example:"Дом в центре города"`
	Type            string  `json:"type" example:"Продать"`
	Currency        string  `json:"currency" example:"VC"`
	Price           float64 `json:"price" example:"5000000"`
	Category        string  `json:"category" example:"house"`
	Nickname        string  `json:"nickname" example:"CoolGamer123"`
	Image           string  `json:"image" example:"http://localhost:8080/uploads/ads/house/123456_user.jpg"`
	Views           int     `json:"views" example:"420"`
	AuthorRating    float64 `json:"author_rating" example:"4.8"`
	CreatedAt       string  `json:"created_at" example:"2026-01-07T12:00:00Z"`
	RentalHoursLimit int    `json:"rental_hours_limit" example:"24"`
}

type CreateFeedbackRequest struct {
	AdOwnerNickname string `json:"ad_owner_nickname" example:"SellerNick" binding:"required"`
	Rating          int    `json:"rating" example:"5" binding:"required,min=1,max=5"`
	Comment         string `json:"comment" example:"Отличный продавец, всё быстро и честно!" binding:"required,max=500"`
}

type FeedbackResponse struct {
	ID              int    `json:"id" example:"1"`
	AdOwnerNickname string `json:"ad_owner_nickname" example:"SellerNick"`
	FromNickname    string `json:"from_nickname" example:"BuyerNick"`
	Rating          int    `json:"rating" example:"5"`
	Comment         string `json:"comment" example:"Отличный продавец!"`
	ConfirmFeedback bool   `json:"confirm_feedback" example:"true"`
	CreatedAt       string `json:"created_at" example:"2026-01-07T12:00:00Z"`
}

type CreateReportRequest struct {
	AdID        int    `json:"ad_id" example:"42" binding:"required"`
	Reason      string `json:"reason" example:"Мошенничество" binding:"required"`
	Description string `json:"description" example:"Продавец требует предоплату" binding:"max=500"`
}

type UserResponse struct {
	UserID                  int     `json:"user_id" example:"1"`
	Nickname                string  `json:"nickname" example:"CoolGamer123"`
	Email                   string  `json:"email" example:"gamer@arizona.rp"`
	Telegram                string  `json:"telegram" example:"@coolgamer"`
	Avatar                  string  `json:"avatar" example:"http://localhost:8080/uploads/avatars/user1.jpg"`
	BackgroundAvatarProfile string  `json:"background_avatar_profile" example:"http://localhost:8080/uploads/backgrounds/bg1.jpg"`
	Rating                  float64 `json:"rating" example:"4.8"`
	ReviewsCount            int64   `json:"reviews_count" example:"15"`
	UserRole                string  `json:"user_role" example:"user"`
	UserDescription         string  `json:"user_description" example:"Активный продавец на Arizona RP"`
	Theme                   string  `json:"theme" example:"dark"`
	LastSeenAt              string  `json:"last_seen_at" example:"2026-01-07T12:00:00Z"`
}

type UpdateProfileRequest struct {
	Nickname    string `json:"nickname" example:"NewNickname" binding:"max=20"`
	Email       string `json:"email" example:"new@arizona.rp" binding:"email"`
	Password    string `json:"password" example:"NewPassword123!" binding:"min=8"`
	Theme       string `json:"theme" example:"dark"`
	Description string `json:"description" example:"Новое описание профиля" binding:"max=500"`
	Telegram    string `json:"telegram" example:"@newusername" binding:"max=50"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"Что-то пошло не так"`
}

type SuccessResponse struct {
	Message string `json:"message" example:"Операция выполнена успешно!"`
}
