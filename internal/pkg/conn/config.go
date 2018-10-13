package conn

import (
	"github.com/yusufsyaifudin/tax-calculator-example/pkg/db"
)

var connection db.SQL

// SetDBConnection will set the main database connection. This connection will be use as main DB.
func SetDBConnection(conn db.SQL) {
	connection = conn
}

// GetDBConnection get the main database connection set by SetDBConnection.
func GetDBConnection() db.SQL {
	return connection
}
