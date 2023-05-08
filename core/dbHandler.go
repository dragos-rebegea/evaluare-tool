package core

import (
	"errors"
	"math/rand"
	"strings"
	"sync"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/dragos-rebegea/evaluare-tool/authentication"
	"gorm.io/gorm"
)

// DatabaseHandler is the interface that wraps the basic database operations

var dbLogger = logger.GetOrCreate("dbHandler")

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

// GetStudentByID returns a student by id
func (db *DatabaseHandler) GetStudentByID(id uint) (*authentication.Student, error) {
	var student authentication.Student
	record := db.database.Where("id = ?", id).First(&student)
	if record.Error != nil {
		return nil, record.Error
	}
	return &student, nil
}

// GetExamByName returns a student by name
func (db *DatabaseHandler) GetExamByName(name string) (*authentication.Exam, error) {
	var exam authentication.Exam
	record := db.database.Where("name = ?", name).First(&exam)
	if record.Error != nil {
		return nil, record.Error
	}
	return &exam, nil
}

// GetExercitiuByName returns a student by name
func (db *DatabaseHandler) GetExercitiuByExamAndNumber(exam string, number int) (*authentication.Exercitiu, error) {
	var exercitiu authentication.Exercitiu
	record := db.database.Where("exam = ? AND numar = ?", exam, number).First(&exercitiu)
	if record.Error != nil {
		return nil, record.Error
	}
	return &exercitiu, nil
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
	record := db.database.
		Table("students").
		Where("clasa = ? AND deleted_at IS NULL", clasa).
		Select("id,nume,prenume,clasa,absent").
		Scan(&students)
	if record.Error != nil {
		return nil, record.Error
	}
	return students, nil
}

// GetAllClasses returns all the users from a class
func (db *DatabaseHandler) GetAllClasses(profEmail string) ([]string, error) {
	var classes []authentication.Clasa
	record := db.database.Table("clasas").Find(&classes)
	if record.Error != nil {
		return nil, record.Error
	}

	profesor, err := db.GetProfesorByEmail(profEmail)
	if err != nil {
		return nil, err
	}
	var classList []string
	for _, class := range classes {
		err = db.checkProfesor(profesor, &class)
		if err != nil {
			continue
		}
		classList = append(classList, class.Nume)
	}
	return classList, nil
}

// SetAbsent sets a student as absent
func (db *DatabaseHandler) SetAbsent(status *AbsentStatus) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	record := db.database.Table("students").Where("id = ?", status.Id).Update("absent", status.Absent)
	if record.Error != nil {
		return record.Error
	}
	return nil
}

// CreateProfesor creates a new profesor
func (db *DatabaseHandler) CreateProfesor(profesor *authentication.Profesor) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	profesor.Type = "profesor"
	if profesor.Password == "" {
		profesor.Password = GenerateRandomString(10)
	}
	password := profesor.Password
	if err := profesor.HashPassword(password); err != nil {
		return errors.New("error hashing password")
	}
	record := db.database.Create(profesor)
	if record.Error != nil {
		return record.Error
	}

	profesor.Password = password
	return nil
}

// CreateClass creates a new class
func (db *DatabaseHandler) CreateClass(class *Class) ([]*authentication.Student, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	clasa := authentication.Clasa{
		Nume:       class.Nume,
		ProfMate:   class.ProfMate,
		ProfBio:    class.ProfBio,
		ProfFizica: class.ProfFizica,
	}
	record := db.database.Create(&clasa)
	if record.Error != nil {
		dbLogger.Error(record.Error.Error())
	}

	students := make([]*authentication.Student, 0)
	for _, student := range class.Elevi {
		// generate random password
		password := GenerateRandomString(10)
		studentDb := authentication.NewStudent(student.Nume, student.Prenume, class.Nume, student.Email, password, student.Exam)
		if err := studentDb.HashPassword(password); err != nil {
			return nil, errors.New("error hashing password")
		}
		err := db.CreateStudent(studentDb)
		if err != nil {
			dbLogger.Error(err.Error())
			continue
		}
		studentDb.Password = password
		students = append(students, studentDb)
	}
	return students, nil
}

// CreateStudent creates a new student
func (db *DatabaseHandler) CreateStudent(student *authentication.Student) error {
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

func (db *DatabaseHandler) DeleteStudent(u *uint) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	record := db.database.Delete(&authentication.Student{}, *u)
	if record.Error != nil {
		return record.Error
	}
	return nil
}

func (db *DatabaseHandler) CreateExam(a *Exam) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	exam := authentication.Exam{Nume: a.Nume}
	record := db.database.Create(&exam)
	if record.Error != nil {
		dbLogger.Error(record.Error.Error())
	}

	for _, ex := range a.Exercitii {
		exercitiu := authentication.Exercitiu{
			Numar:    ex.Numar,
			Variante: strings.Join(ex.Variante, ";"),
			Materie:  ex.Materie,
			Exam:     exam.Nume,
		}
		record = db.database.Create(&exercitiu)
		if record.Error != nil {
			dbLogger.Error(record.Error.Error())
		}
	}
	return nil
}

func (db *DatabaseHandler) AddCalificativ(profEmail string, calificativ *Calificativ) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	err := db.checkCalificativ(profEmail, calificativ)
	if err != nil {
		return err
	}

	record := db.database.Create(&calificativ)
	if record.Error != nil {
		dbLogger.Error(record.Error.Error())
	}
	return nil
}

func (db *DatabaseHandler) GetCalificativByStudentAndExercitiu(id uint, exercitiu uint) (*Calificativ, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	var calificativ Calificativ
	record := db.database.Where("student = ? AND exercitiu = ?", id, exercitiu).First(&calificativ)
	if record.Error != nil {
		return nil, record.Error
	}
	return &calificativ, nil
}

func (db *DatabaseHandler) UpdateCalificativ(profEmail string, calificativ *Calificativ) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	var calificativDb Calificativ
	record := db.database.Where("student = ? AND exercitiu = ?", calificativ.Student, calificativ.Exercitiu).First(&calificativDb)
	if record.Error != nil {
		return record.Error
	}
	if record.RowsAffected == 0 {
		return errors.New("calificativ not found")
	}

	err := db.checkCalificativ(profEmail, calificativ)
	if err != nil {
		return err
	}
	record = db.database.Where("student = ? AND exercitiu = ?", calificativ.Student, calificativ.Exercitiu).Save(&calificativ)
	if record.Error != nil {
		return record.Error
	}
	return nil

}

func (db *DatabaseHandler) checkCalificativ(profEmail string, calificativ *Calificativ) error {
	student, err := db.GetStudentByID(calificativ.Student)
	if student == nil {
		return errors.New("student not found")
	}

	exercitiu, err := db.GetExercitiuByExamAndNumber(calificativ.Exam, calificativ.Exercitiu)
	if err != nil {
		return err
	}
	if exercitiu == nil {
		return errors.New("exercitiu not found")
	}

	class, err := db.GetClassByID(calificativ.Student)
	prof, err := db.GetProfesorByEmail(profEmail)
	if err != nil {
		return err
	}
	if prof == nil {
		return errors.New("profesor not found")
	}
	err = db.checkProfesor(prof, class)
	if err != nil {
		return err
	}
	calificativ.Profesor = prof.ID

	variante := strings.Split(exercitiu.Variante, ";")
	if !contains(variante, calificativ.Varianta) {
		return errors.New("varianta invalida")
	}

	return nil
}

func (db *DatabaseHandler) GetCalificative(email string, student string) ([]*Calificativ, error) {
	prof, err := db.GetProfesorByEmail(email)
	if err != nil {
		return make([]*Calificativ, 0), err
	}
	if prof == nil {
		return make([]*Calificativ, 0), errors.New("profesor not found")
	}

	if student == "" {
		return make([]*Calificativ, 0), nil
	}

	var calificative []*Calificativ
	record := db.database.
		Table("calificativs").
		Where("student = ? AND profesor = ? ", student, prof.ID).
		Scan(&calificative)
	if record.Error != nil {
		return nil, record.Error
	}
	return calificative, nil
}

func (db *DatabaseHandler) GetExercitiiForProfesorAndStudent(email string, studentId string) ([]*Exercitiu, error) {
	prof, err := db.GetProfesorByEmail(email)
	if err != nil {
		return make([]*Exercitiu, 0), err
	}
	if prof == nil {
		return make([]*Exercitiu, 0), errors.New("profesor not found")
	}

	if studentId == "" {
		return make([]*Exercitiu, 0), nil
	}

	var student authentication.Student
	record := db.database.Where("id = ?", studentId).First(&student)
	if record.Error != nil {
		return nil, record.Error
	}

	var class authentication.Clasa
	record = db.database.Where("nume = ?", student.Clasa).First(&class)
	if record.Error != nil {
		return nil, record.Error
	}

	err = db.checkProfesor(prof, &class)
	if err != nil {
		return make([]*Exercitiu, 0), err
	}
	var exercitii []authentication.Exercitiu
	record = db.database.Where("materie = ? AND exam = ?", prof.Materie, student.Exam).Find(&exercitii)
	if record.Error != nil {
		return nil, record.Error
	}

	exercitiiReturn := make([]*Exercitiu, 0)
	for _, exercitiu := range exercitii {
		variante := strings.Split(exercitiu.Variante, ";")
		exercitiuReturn := &Exercitiu{
			Numar:    exercitiu.Numar,
			Variante: variante,
			Materie:  exercitiu.Materie,
			Exam:     exercitiu.Exam,
		}
		exercitiiReturn = append(exercitiiReturn, exercitiuReturn)
	}
	return exercitiiReturn, nil
}

func (db *DatabaseHandler) GetClassByID(studentId uint) (*authentication.Clasa, error) {
	var student authentication.Student
	record := db.database.Where("id = ?", studentId).First(&student)
	if record.Error != nil {
		return nil, record.Error
	}

	var class authentication.Clasa
	record = db.database.Where("nume = ?", student.Clasa).First(&class)
	if record.Error != nil {
		return nil, record.Error
	}
	return &class, nil
}

func (db *DatabaseHandler) checkProfesor(prof *authentication.Profesor, class *authentication.Clasa) error {
	if prof.ID == class.ProfMate {
		return nil
	}
	if prof.ID == class.ProfBio {
		return nil
	}
	if prof.ID == class.ProfFizica {
		return nil
	}
	return errors.New("profesor invalid")
}

func GenerateRandomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^&*1234567890")

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
