package main

import (
	"database/sql"
	"identityreconciliation/core"
	"identityreconciliation/repository"
	"identityreconciliation/service"
)

func main() {
	storage := repository.NewRepository(&sql.DB{})
	service := service.NewService(storage)
	newCoreService := core.NewCoreService(storage, service)

	newCoreService.CallRun()
}
