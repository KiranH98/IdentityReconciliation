package service

import (
	"identityreconciliation/repository"
	"log"
	"os"
)

type Service struct {
	storage *repository.Repository
	log     *log.Logger
}

func NewService(storage *repository.Repository) *Service {
	return &Service{
		storage: storage,
		log:     log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile),
	}
}
