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
	Audit(string) error
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

type Auditor interface {
	Audit(string) error
}

// hello
// #[@Usecase]
// #[go.permission.audit]
func Aufstehen(audit Auditor) error {
	if err := audit.Audit("de.worldiety.aufstehen"); err != nil {
		return err
	}

	return nil
}

// Aufstehen Zeiterfassungsmethode
// #[@Usecase "Aufstehen in der Zeiterfassung"]
// #[go.permission.audit]
func (z *Zeiterfassung) Aufstehen(audit Auditor) error {
	if err := audit.Audit("de.worldiety.zeiterfassung.aufstehen"); err != nil {
		return err
	}

	return nil
}

type Other struct {
}

// Aufstehen Zeiterfassungsmethode
// #[@Usecase "Aufstehen woanders"]
// #[go.permission.audit]
func (z *Other) Aufstehen(audit Auditor) error {
	if err := audit.Audit("de.worldiety.woanders.aufstehen"); err != nil {
		return err
	}

	return nil
}

type Human interface {
	// Human Aufstehen doc
	Aufstehen(audit Auditor) error
}

// Cooles Zeitbuchen ist angesagt.
// #[@Usecase "Zeiten loggen"]
// #[go.permission.audit]
func (z *Zeiterfassung) ZeitBuchen(user User, mitarbeiter Mitarbeiter, dauer time.Duration) (int, error) {
	if err := user.Audit("de.worldiety.aufstehen2"); err != nil {
		return 0, err
	}
	z.repo.Save(Zeitlog{
		Dauer: dauer,
		Text:  "gearbeitet",
	})
	//fmt.Println("zeit gebucht", user.Audit)
	return 0, nil
}

// #[@Usecase "Beschwerde einreichen"]
func (z *Zeiterfassung) BeschwerdeEinreichen(msg string) {

}
