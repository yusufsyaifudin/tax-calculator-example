package user

var (
	sqlInsertUser         = `INSERT INTO users(username, password) VALUES(?, ?) RETURNING *;`
	sqlFindUserByID       = `SELECT * FROM users WHERE id = ?;`
	sqlFindUserByUsername = `SELECT * FROM users WHERE username = ?;`
)
