package main

import (
	"example/domain"
	"example/domain/adapter/bolt"
	"example/domain/adapter/mysql"
	"example/supporting"
	"time"
)

// #[markdown "out":"README.md"]
func main() {
	var repo domain.ZeitlogRepo
	if time.Now().UnixMilli()%10 == 0 {
		repo = mysql.Repo{}
	} else {
		repo = bolt.Repo{}
	}
	service := domain.NewZeiterfassung(repo)
	service.ZeitBuchen(supporting.User{}, domain.Mitarbeiter{}, 12)
	domain.Aufstehen(NagoAuditor{})
}

type NagoAuditor struct {
}

func (n NagoAuditor) Audit(s string) error {
	panic("implement me")
}
