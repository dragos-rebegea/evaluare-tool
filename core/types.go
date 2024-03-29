package core

type Class struct {
	Nume  string `json:"nume"`
	Elevi []struct {
		Nume        string `json:"nume"`
		Prenume     string `json:"prenume"`
		Email       string `json:"email"`
		ExamStiinta string `json:"exam_stiinta"`
		ExamLimba   string `json:"exam_limba"`
	} `json:"elevi"`
	ProfMate    string `json:"prof_mate"`
	ProfFizica  string `json:"prof_fizica"`
	ProfBio     string `json:"prof_bio"`
	ProfRomana  string `json:"prof_romana"`
	ProfEngleza string `json:"prof_engleza"`
}

type AbsentStatus struct {
	Id     uint `json:"id"`
	Absent bool `json:"absent"`
}

type Exam struct {
	Nume      string      `json:"nume"`
	Exercitii []Exercitiu `json:"exercitii"`
}

type Exercitiu struct {
	Numar    string   `json:"numar"`
	Variante []string `json:"variante"`
	Materie  string   `json:"materie"`
	Exam     string   `json:"exam"`
}

type Calificativ struct {
	Student   uint   `json:"student_id"`
	Profesor  uint   `json:"profesor_id"`
	Exam      string `json:"exam"`
	Exercitiu string `json:"exercitiu"`
	Varianta  string `json:"varianta"`
}
