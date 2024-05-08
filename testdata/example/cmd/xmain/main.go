package main

import (
	"example/domain"
	"example/domain/adapter/bolt"
	"example/domain/adapter/mysql"
	"example/supporting"
	"time"
)

// #[markdown]
func main() {
	var repo domain.ZeitlogRepo
	if time.Now().UnixMilli()%10 == 0 {
		repo = mysql.Repo{}
	} else {
		repo = bolt.Repo{}
	}
	service := domain.NewZeiterfassung(repo)
	service.ZeitBuchen(supporting.User{}, domain.Mitarbeiter{}, 12)
}
