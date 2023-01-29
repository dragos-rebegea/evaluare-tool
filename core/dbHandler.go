package core

import (
	"math/rand"
	"sync"

	"github.com/dragos-rebegea/evaluare-tool/authentication"
	"gorm.io/gorm"
)

// DatabaseHandler is the interface that wraps the basic database operations

type DatabaseHandler struct {
	mutex    sync.RWMutex
	database *gorm.DB
}

// NewDatabaseHandler returns a new DatabaseHandler instance
func NewDatabaseHandler(connectionstring string) (*DatabaseHandler, error) {
	db, err := authentication.Connect(connectionstring)
	if err != nil {
		return nil, err
	}
	err = authentication.Migrate(db)
	if err != nil {
		return nil, err
	}
	return &DatabaseHandler{database: db}, nil
}

// GetStudentByEmail returns a student by email
func (db *DatabaseHandler) GetStudentByEmail(email string) (*authentication.Student, error) {
	var student authentication.Student
	record := db.database.Where("email = ?", email).First(&student)
	if record.Error != nil {
		return nil, record.Error
	}
	return &student, nil
}

// GetProfesorByEmail returns a profesor by email
func (db *DatabaseHandler) GetProfesorByEmail(email string) (*authentication.Profesor, error) {
	var profesor authentication.Profesor
	record := db.database.Where("email = ?", email).First(&profesor)
	if record.Error != nil {
		return nil, record.Error
	}
	return &profesor, nil
}

// GetStudentsByClass returns all the users from a class
func (db *DatabaseHandler) GetStudentsByClass(clasa string) ([]authentication.Student, error) {
	var students []authentication.Student
	record := db.database.Table("students").Where("clasa = ?", clasa).Select("nume,prenume,clasa").Scan(&students)
	if record.Error != nil {
		return nil, record.Error
	}
	return students, nil
}

// CreateProfesor creates a new profesor
func (db *DatabaseHandler) CreateProfesor(profesor *authentication.Profesor) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	profesor.Type = "profesor"
	record := db.database.Create(profesor)
	if record.Error != nil {
		return record.Error
	}
	return nil
}

// CreateClass creates a new class
func (db *DatabaseHandler) CreateClass(class *Class) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	for _, student := range class.Elevi {
		// generate random password
		password := GenerateRandomString(10)
		studentDb := authentication.NewStudent(student.Nume, student.Prenume, class.Nume, student.Email, password)
		err := db.CreateStudent(studentDb)
		if err != nil {
			continue
		}
	}
	return nil
}

// CreateStudent creates a new student
func (db *DatabaseHandler) CreateStudent(student *authentication.Student) error {
	//db.mutex.RLock()
	//defer db.mutex.RUnlock()

	record := db.database.Create(student)
	if record.Error != nil {
		return record.Error
	}
	return nil
}

// IsAdmin returns true if the user is an admin
func (db *DatabaseHandler) IsAdmin(email string) (bool, error) {
	user, err := db.GetProfesorByEmail(email)
	if err != nil {
		return false, err
	}
	return user.IsAdmin, nil
}

// IsProfesor returns true if the user is a profesor
func (db *DatabaseHandler) IsProfesor(email string) (bool, error) {
	user, err := db.GetProfesorByEmail(email)
	if err != nil {
		return false, err
	}
	return user.Type == "profesor", nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (db *DatabaseHandler) IsInterfaceNil() bool {
	return db == nil
}

func GenerateRandomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^&*1234567890")

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}
