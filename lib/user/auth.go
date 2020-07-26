package user

import (
	"Mlops/config"
	"Mlops/model"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"time"
)

var SignedKey = config.GetHmacSignKey()

func GenerateToken(user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"email":     user.Email,
		"user_id":   user.ID,
		"is_Admin":	user.IsAdmin,
		"issued_at": time.Now(),
		"expire_at": time.Now().Add(time.Minute * 72).Unix(),
	}
	rawtoken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	token, err := rawtoken.SignedString(SignedKey)
	if err != nil {
		log.Println("error generating token")
		return "", err
	}

	return token, err
}
func CreateCookie(user *model.User) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = "mlops_cookie"
	cookie.Value = string(user.ID)
	cookie.Expires = time.Now().Add(30 * time.Minute)

	return cookie
}
