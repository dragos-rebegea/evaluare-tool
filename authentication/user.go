package authentication

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Materie string

const (
	Matematica Materie = "Matematica"
	Biologie           = "Biologie"
	Fizica             = "Fizica"
)

type User struct {
	gorm.Model
	Nume     string `json:"nume"`
	Prenume  string `json:"prenume"`
	Username string `json:"username" gorm:"unique"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Type     string `json:"type"`
}

type Student struct {
	User
	Clasa string `json:"clasa"`
}

func NewStudent(nume, prenume, clasa, email, password string) *Student {
	return &Student{
		User: User{
			Nume:     nume,
			Prenume:  prenume,
			Username: nume + "_" + prenume,
			Email:    email,
			Password: password,
			Type:     "student",
		},
		Clasa: clasa,
	}
}

type Profesor struct {
	User
	IsAdmin bool   `json:"is_admin"`
	Materie string `json:"materie"`
}

func NewProfesor(nume, prenume, materie, email, password string) *Profesor {
	return &Profesor{
		User: User{
			Nume:     nume,
			Prenume:  prenume,
			Username: nume + "_" + prenume,
			Email:    email,
			Password: password,
			Type:     "profesor",
		},
		Materie: materie,
	}
}

func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

func (user *User) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}
