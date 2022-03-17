package domain

import "time"

type CryptoService interface {
	CreateToken(data map[string]interface{}, expires time.Duration) (string, error)
	VerifyToken(t string) (map[string]interface{}, error)

	HashPassword(password string) (string, error)
	VerifyPassword(hash, password string) bool
}
