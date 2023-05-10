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

type Varianta string

const (
	A Varianta = "A"
	B          = "B"
	C          = "C"
	D          = "D"
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

type Clasa struct {
	Nume        string `gorm:"primarykey" json:"nume"`
	ProfMate    uint   `gorm:"foreignkey" json:"prof_mate"`
	ProfFizica  uint   `gorm:"foreignkey" json:"prof_fizica"`
	ProfBio     uint   `gorm:"foreignkey" json:"prof_bio"`
	ProfRomana  uint   `gorm:"foreignkey" json:"prof_romana"`
	ProfEngleza uint   `gorm:"foreignkey" json:"prof_engleza"`
}

type Student struct {
	User
	Absent      bool   `json:"absent"`
	Clasa       string `gorm:"foreignkey" json:"clasa"`
	ExamStiinta string `gorm:"foreignkey" json:"exam_stiinta"`
	ExamLimba   string `gorm:"foreignkey" json:"exam_limba"`
}

type Calificativ struct {
	Student   uint   `gorm:"primarykey" json:"student_id"`
	Profesor  uint   `json:"profesor_id"`
	Exam      string `gorm:"primarykey" json:"exam"`
	Exercitiu int    `gorm:"primarykey" json:"exercitiu"`
	Varianta  string `json:"varianta"`
}

func NewStudent(nume, prenume, clasa, email, password, examStiinta, examLimba string) *Student {
	return &Student{
		User: User{
			Nume:     nume,
			Prenume:  prenume,
			Username: nume + "_" + prenume,
			Email:    email,
			Password: password,
			Type:     "student",
		},
		Absent:      false,
		Clasa:       clasa,
		ExamStiinta: examStiinta,
		ExamLimba:   examLimba,
	}
}

type Profesor struct {
	User
	IsAdmin bool   `json:"is_admin"`
	Materie string `json:"materie"`
}

type Exam struct {
	Nume string `gorm:"primarykey" json:"nume"`
}

type Exercitiu struct {
	Numar    uint   `gorm:"primarykey" json:"numar"`
	Variante string `json:"variante"`
	Materie  string `json:"materie"`
	Exam     string `gorm:"primarykey" json:"exam"`
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
