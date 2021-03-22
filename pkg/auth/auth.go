package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"postit/model"

	"github.com/dgrijalva/jwt-go"
)

const (
	publicKeyLocation = "/app/jwt/public.key"
)

func Validate(token string) (model.ShortenedUser, error) {
	public, err := ioutil.ReadFile(publicKeyLocation)
	var shorten model.ShortenedUser
	if err != nil {
		return shorten, err
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(public)
	if err != nil {
		return shorten, fmt.Errorf("validate: parse key: %w", err)
	}

	tok, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}

		return key, nil
	})

	if err != nil {
		return shorten, fmt.Errorf("validate: %w", err)
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return shorten, fmt.Errorf("validate: invalid")
	}

	jsonString, err := json.Marshal(claims["dat"])
	if err != nil {
		return shorten, err
	}

	err = json.Unmarshal(jsonString, &shorten)
	if err != nil {
		return shorten, err
	}

	return shorten, nil
}
