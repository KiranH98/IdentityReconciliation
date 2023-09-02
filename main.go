package main

import (
	"identityreconciliation/core"
	"identityreconciliation/repository"
	"identityreconciliation/service"
)

func main() {
	storage := &repository.Repository{}
	service := &service.Service{}

	newCoreService := core.NewCoreService(storage, service)
	newCoreService.CallRun()
}
