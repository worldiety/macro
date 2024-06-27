// Dieses Projekt ist ein Beispielprojekt und zeigt die Verwendung verschiedener Annotationen.
// #[@Project "Beispielprojekt"]
package main

import (
	"example/domain"
	"example/domain/adapter/bolt"
	"example/domain/adapter/mysql"
	"example/supporting"
	"time"
)

// #[markdown "out":"README.md", "omitSecurityChapter":false]
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
	var dh domain.Human
	dh.Aufstehen(nil)
	service.Aufstehen(nil)
}

type NagoAuditor struct {
}

func (n NagoAuditor) Audit(s string) error {
	panic("implement me")
}
