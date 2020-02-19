package env

import (
	"github.com/payfazz/secureenv"
)

var (
	Addr               = ":8080"
	PostgresConnection = ""
)

func init() {
	secureenv.String(&Addr, "Addr")
	secureenv.String(&PostgresConnection, "PostgresConnection")
}
