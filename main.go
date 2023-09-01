package main

import (
	"identityreconciliation/api"
	"identityreconciliation/core"
	db "identityreconciliation/database"
)

func main() {
	db := &db.DataBase{}
	apiHandler := &api.API{}

	newCoreService := core.NewCoreService(db, apiHandler)
	newCoreService.CallRun()
}
