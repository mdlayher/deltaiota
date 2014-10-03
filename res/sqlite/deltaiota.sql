/* deltaiota sqlite schema */
PRAGMA foreign_keys = OFF;
BEGIN TRANSACTION;
/* users */
CREATE TABLE "users" (
	"id"         INTEGER PRIMARY KEY AUTOINCREMENT,
	"username"   TEXT NOT NULL,
	"first_name" TEXT NOT NULL,
	"last_name"  TEXT NOT NULL,
	"email"      TEXT NOT NULL,
	"phone"      TEXT NOT NULL,
	"password"   TEXT NOT NULL,
	"salt"       TEXT NOT NULL
);
CREATE UNIQUE INDEX "users_unique_username" ON "users" ("username");
CREATE UNIQUE INDEX "users_unique_email" ON "users" ("email");
CREATE UNIQUE INDEX "users_unique_password" ON "users" ("password");
CREATE UNIQUE INDEX "users_unique_salt" ON "users" ("salt");
COMMIT;
