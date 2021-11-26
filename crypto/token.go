package crypto

import (
	"time"

	dq "doublequote"
	"github.com/golang-jwt/jwt"
)

func (s *Service) CreateToken(data map[string]interface{}, expires time.Duration) (string, error) {
	data["exp"] = s.now().Add(expires).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(data))
	tokenString, err := token.SignedString(s.key)

	return tokenString, err
}

func (s *Service) VerifyToken(t string) (map[string]interface{}, error) {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, dq.Errorf(dq.EINVALID, "Unexpected signing method: %v.", token.Header["alg"])
		}

		validExp := token.Claims.(jwt.MapClaims).VerifyExpiresAt(s.now().Unix(), true)
		if !validExp {
			return nil, dq.Errorf(dq.EINVALID, "Token expired.")
		}

		return s.key, nil
	})

	if err != nil {
		inner := err.(*jwt.ValidationError).Inner
		if inner != nil {
			return nil, inner
		}
		return nil, err
	}

	if !token.Valid {
		return nil, dq.Errorf(dq.EINVALID, "invalid token")
	}

	return token.Claims.(jwt.MapClaims), nil
}
