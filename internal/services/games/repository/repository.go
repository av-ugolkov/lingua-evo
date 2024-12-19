package repository

import "github.com/av-ugolkov/lingua-evo/internal/db/transactor"

type Repo struct {
	tr *transactor.Transactor
}

func New(tr *transactor.Transactor) *Repo {
	return &Repo{
		tr: tr,
	}
}
