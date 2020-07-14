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


func GenerateToken (user *model.User) (*jwt.Token, string,error) {
	claims := jwt.MapClaims{
		"email": user.Email,
		"user_id": user.ID,
		"issued_at": time.Now(),
		"expire_at": time.Now().Add(time.Minute * 72).Unix(),
	}
	rawtoken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	token, err := rawtoken.SignedString(SignedKey)
	if err != nil {
		log.Println("error generating token")
		return  rawtoken,"",err
	}
	log.Println(rawtoken)
	return rawtoken,token,err
}
func CreateCookie (user *model.User) *http.Cookie {
	cookie := http.Cookie{
		Name: "mlops",
		Value: "mlops_id",
		Expires: time.Now().Add(30 * time.Minute),
	}
	return &cookie
}
