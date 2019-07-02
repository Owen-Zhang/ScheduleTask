package system

import (
	"github.com/dgrijalva/jwt-go"
	"time"
	"fmt"
	"crypto/sha256"
	"encoding/hex"
)

/*
  Encrypt 加密
*/
func Encrypt(content, key string ) (t string, err error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["content"] = content
	claims["exp"] = time.Now().Add(time.Minute * 10).Unix()

	t, err = token.SignedString([]byte(key))
	return
}

/*
  Decrypt 解密
*/

func Decrypt(tokenStr, key string) (interface{}, bool) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(key), nil
	})

	fmt.Println(token)

	if claims, ok := token.Claims.(jwt.MapClaims); ok && claims.Valid() != nil {
		return claims, ok
	} else {
		fmt.Println(claims["content"])

		fmt.Println("not ok")
		fmt.Println(err)
		return "", false
	}
}

func CryptoSHA256(key string) string {
	s := sha256.New()
	s.Write([]byte(key))
	bs := s.Sum(nil)
	return hex.EncodeToString(bs)
}
