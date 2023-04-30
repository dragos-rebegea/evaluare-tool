package core

type Class struct {
	Nume  string `json:"nume"`
	Elevi []struct {
		Nume    string `json:"nume"`
		Prenume string `json:"prenume"`
		Email   string `json:"email"`
		Exam    string `json:"exam"`
	} `json:"elevi"`
	ProfMate   uint `json:"prof_mate"`
	ProfFizica uint `json:"prof_fizica"`
	ProfBio    uint `json:"prof_bio"`
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
	Numar    uint     `json:"numar"`
	Variante []string `json:"variante"`
	Materie  string   `json:"materie"`
	Exam     string   `json:"exam"`
}

type Calificativ struct {
	Student   uint   `json:"student_id"`
	Profesor  uint   `json:"profesor_id"`
	Exam      string `json:"exam"`
	Exercitiu int    `json:"exercitiu"`
	Varianta  string `json:"varianta"`
}
