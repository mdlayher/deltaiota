package data

const (
	// sqlSelectAllUsers is the SQL statement used to select all Users
	sqlSelectAllUsers = `
		SELECT * FROM users;
	`
	// sqlSelectUserByID is the SQL statement used to select a single user by ID
	sqlSelectUserByID = `
		SELECT * FROM users WHERE id = ?;
	`

	// sqlInsertUser is the SQL statement used to insert a new User
	sqlInsertUser = `
		INSERT INTO users (
			"username"
		) VALUES (?);
	`
)
