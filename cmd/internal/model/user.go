package model

import (
	"crypto/sha512"
	"encoding/base64"
)

type User struct {
	UserID   string `json:"userId,omitempty" db:"user_id"`
	Login    string `json:"login" db:"login"`
	Password string `json:"password" db:"password"`
}

func (u *User) PasswordHash() string {
	hash := []byte(u.Password)
	dig := sha512.Sum512(hash)
	for i := 1; i < 5000; i++ {
		dig = sha512.Sum512(append(dig[:], hash[:]...))
	}
	return base64.StdEncoding.EncodeToString(dig[:])
}
