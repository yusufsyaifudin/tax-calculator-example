package tax

var (
	sqlInsertTax        = `INSERT INTO taxes(user_id, name, tax_code, price) VALUES(?, ?, ?, ?) RETURNING *;`
	sqlGetTaxesByUserId = `SELECT * FROM taxes WHERE user_id = ? ORDER BY id DESC;`
)
