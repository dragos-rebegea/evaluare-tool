package core

type Class struct {
	Nume  string `json:"nume"`
	Elevi []struct {
		Nume    string `json:"nume"`
		Prenume string `json:"prenume"`
		Email   string `json:"email"`
	} `json:"elevi"`
}
