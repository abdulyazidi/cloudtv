package main

import (
	"fmt"
	"log"

	"github.com/golang-jwt/jwt/v5"
)

var (
	key []byte
	t   *jwt.Token
	s   string
)

type CustomClaims struct {
	Foo   int `json:"foo"`
	LMFAO int `json:"lmfao"`
	jwt.RegisteredClaims
}

func main() {

	key = []byte("XD")
	t = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "playground-server",
		"sub": "john cena",
		"foo": 1000,
		"XD":  "005",
	})
	s, err := t.SignedString(key)

	if err != nil {
		log.Println(err)
	}
	parsed, err := jwt.ParseWithClaims(s, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %s", t.Header["alg"])
		}
		return key, nil
	}, jwt.WithIssuer("playground-server"))

	if err != nil {
		log.Fatalf("token rejected: %v", err)
	}

	claims, ok := parsed.Claims.(*CustomClaims)
	if !ok || !parsed.Valid {
		log.Fatal("token failed signature or standard-claim validation")
	}

	// fmt.Println(parsed.Raw)
	fmt.Println(claims.Foo)
	// fmt.Println("JWT TOKEN:", s)
}
