package entity

import "golang.org/x/crypto/bcrypt"

type Password struct {
	Hashed string
}

func NewPassword(password string) (Password, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return Password{}, err
	}

	return Password{Hashed: string(hashed)}, nil
}

func (p Password) Compare(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p.Hashed), []byte(password))
	return err == nil
}
