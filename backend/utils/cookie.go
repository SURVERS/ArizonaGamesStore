package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetAuthCookie(c *gin.Context, name, value string, maxAge int) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   maxAge,
		Path:     "/",
		Domain:   "",
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
		HttpOnly: true,
	}
	http.SetCookie(c.Writer, cookie)
}
