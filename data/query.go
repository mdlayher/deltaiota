package data

const (
	// sqlSelectAllUsers is the SQL statement used to select all Users
	sqlSelectAllUsers = `
		SELECT * FROM users;
	`

	// sqlInsertUser is the SQL statement used to insert a new User
	sqlInsertUser = `
		INSERT INTO users (
			"username"
		) VALUES (?);
	`
)
