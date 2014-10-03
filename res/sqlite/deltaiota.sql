/* deltaiota sqlite schema */
PRAGMA foreign_keys = OFF;
BEGIN TRANSACTION;
/* users */
CREATE TABLE "users" (
	"id"         INTEGER PRIMARY KEY AUTOINCREMENT,
	"username"   TEXT,
	"first_name" TEXT,
	"last_name"  TEXT,
	"email"      TEXT,
	"phone"      TEXT,
	"password"   TEXT,
	"salt"       TEXT
);
CREATE UNIQUE INDEX "users_unique_username" ON "users" ("username");
CREATE UNIQUE INDEX "users_unique_email" ON "users" ("email");
CREATE UNIQUE INDEX "users_unique_password" ON "users" ("password");
CREATE UNIQUE INDEX "users_unique_salt" ON "users" ("salt");
COMMIT;
