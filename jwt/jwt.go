package jwt

import (
	"fmt"
	"log"

	"github.com/golang-jwt/jwt/v4"
)

type JWT struct {
	HMACSecret []byte // 秘钥
}

type MyClaims struct {
	User string `json:"user,omitempty"`
	jwt.RegisteredClaims
}

// https://pkg.go.dev/github.com/golang-jwt/jwt/v4@v4.4.3#MapClaims.VerifyAudience
func (j *JWT) Sign(claims *MyClaims) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(j.HMACSecret)
	return tokenString, err
}

func (j *JWT) Validate(tokenString string) bool {
	flag := false // 登录是否有效

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名的算法alg，这里采用HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return token, nil
	})

	log.Fatalln(" jwt.Parse:", err)

	fmt.Println(token.Valid)
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		fmt.Println(claims)
		// token验证成功
		flag = true
	} else {
		// token验证失败
		flag = false
	}

	return flag
}
