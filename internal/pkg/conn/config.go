package conn

import (
	"github.com/yusufsyaifudin/tax-calculator-example/pkg/db"
)

var connection db.SQL

func SetDBConnection(conn db.SQL) {
	connection = conn
}

func GetDBConnection() db.SQL {
	return connection
}
