package mysql

import (
	"example/domain"
)

type Repo struct {
}

func (r Repo) Each(yield func(domain.Zeitlog) bool) {
	//TODO implement me
	panic("implement me")
}

func (r Repo) FindById(id int) (domain.Zeitlog, error) {
	//TODO implement me
	panic("implement me")
}

func (r Repo) Save(z domain.Zeitlog) (domain.Zeitlog, error) {
	//TODO implement me
	panic("implement me")
}
