package domain

import (
	"time"
)

type Zeitlog struct {
	Dauer time.Duration
	Text  string
}

type Mitarbeiter struct {
	ID   string
	Name string
}

type User interface {
	Audit()
}

type ZeitlogRepo interface {
	Each(yield func(Zeitlog) bool)
	FindById(id int) (Zeitlog, error)
	Save(z Zeitlog) (Zeitlog, error)
}

type Zeiterfassung struct {
	repo ZeitlogRepo
}

func NewZeiterfassung(repo ZeitlogRepo) *Zeiterfassung {
	return &Zeiterfassung{repo: repo}
}

func (z *Zeiterfassung) ZeitBuchen(user User, mitarbeiter Mitarbeiter, dauer time.Duration) error {
	user.Audit()
	z.repo.Save(Zeitlog{
		Dauer: dauer,
		Text:  "gearbeitet",
	})
	//fmt.Println("zeit gebucht", user.Audit)
	return nil
}
