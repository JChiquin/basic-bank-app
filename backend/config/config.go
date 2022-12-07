package config

import (
	"bank-service/src/libs/database"
	"io"

	"bank-service/src/libs/logger"
)

/*
slice of dependencies, io.Closes is an interface with method Close() error
all package that makes connections implements it
*/
var dependenciesToClose []io.Closer

/*
SetupCommonDependencies calls setup for each necessary dependency
and registers them on one slice to be closed later
*/
func SetupCommonDependencies() {
	logger.SetupLogger()
	database.SetupBankGormDB()
	dependenciesToClose = []io.Closer{}
}

/*
TearDownCommonDependencies iterates each dependency and calls Close method
*/
func TearDownCommonDependencies() {
	for _, dependecy := range dependenciesToClose {
		dependecy.Close()
	}
}
