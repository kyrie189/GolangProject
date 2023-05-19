package common

import (
	"ginProjectGRPC/model"
	"time"
	"github.com/dgrijalva/jwt-go"
)

var Secret = []byte("a_secret_crect") //HS的方式

type Claims struct {
	UserId uint
	jwt.StandardClaims
}

// 获得token
func ReleaseToken(user model.User) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)

	claims := &Claims{
		UserId: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "weilifeng",
			Subject:   "user token",
		},
	}
	//token 是*jwt.Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(Secret) //使用密钥生成token

}

// 解析token
func ParseToken(tokenString string)(*jwt.Token,*Claims,error){
	claims := &Claims{}
	token,err := jwt.ParseWithClaims(tokenString,claims, 
		func(token *jwt.Token) (interface{}, error) {
		return Secret,nil
	})
	return token,claims,err
}

