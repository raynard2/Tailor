package user

import (
	"Mlops/config"
	"Mlops/model"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

var SignedKey = config.GetHmacSignKey()

func GenerateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"email":     user.Email,
		"user_id":   user.ID,
		"is_Admin":  user.IsAdmin,
		"issued_at": time.Now(),
		"expire_at": time.Now().Add(time.Minute * 72).Unix(),
	}

	raw := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	token, err := raw.SignedString(SignedKey)

	if err != nil {
		return "", err
	}

	return token, err
}
func CreateCookie(user *models.User) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = "mlops_cookie"
	cookie.Value = string(user.ID)
	cookie.Expires = time.Now().Add(30 * time.Minute)

	return cookie
}
