package domain

import (
	"time"
)

// #[@Value]
type Zeitlog struct {
	Dauer time.Duration
	Text  string
}

// Mitarbeiter arbeitet bei seinem Arbeitgeber.
// #[@Entity]
type Mitarbeiter struct {
	ID   string
	Name string
}

// #[@Aggregate]
type User interface {
	Audit()
}

// ZeitlogRepo manages the Zeitlogs.
// #[@Repository "Zeitaufzeichnungen"]
type ZeitlogRepo interface {
	Each(yield func(Zeitlog) bool)
	FindById(id int) (Zeitlog, error)
	Save(z Zeitlog) (Zeitlog, error)
}

// #[@DomainService]
type Zeiterfassung struct {
	repo ZeitlogRepo
}

func NewZeiterfassung(repo ZeitlogRepo) *Zeiterfassung {
	return &Zeiterfassung{repo: repo}
}

// #[@Usecase]
func Aufstehen() {

}

// Cooles Zeitbuchen ist angesagt.
// #[@Usecase]
func (z *Zeiterfassung) ZeitBuchen(user User, mitarbeiter Mitarbeiter, dauer time.Duration) error {
	user.Audit()
	z.repo.Save(Zeitlog{
		Dauer: dauer,
		Text:  "gearbeitet",
	})
	//fmt.Println("zeit gebucht", user.Audit)
	return nil
}
