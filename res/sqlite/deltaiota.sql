/* deltaiota sqlite schema */
BEGIN TRANSACTION;
/* sessions */
CREATE TABLE "sessions" (
	"id"        INTEGER PRIMARY KEY AUTOINCREMENT
	, "user_id" INTEGER NOT NULL
	, "key"        TEXT NOT NULL
	, "expire"  INTEGER NOT NULL

	, FOREIGN KEY(user_id) REFERENCES users(id)
);
CREATE UNIQUE INDEX "sessions_unique_key" ON "sessions" ("key");
/* users */
CREATE TABLE "users" (
	"id"           INTEGER PRIMARY KEY AUTOINCREMENT
	, "username"      TEXT NOT NULL
	, "first_name"    TEXT NOT NULL
	, "last_name"     TEXT NOT NULL
	, "email"         TEXT NOT NULL
	, "phone"         TEXT NOT NULL
	, "password"      TEXT NOT NULL
);
CREATE UNIQUE INDEX "users_unique_username" ON "users" ("username");
CREATE UNIQUE INDEX "users_unique_email" ON "users" ("email");
CREATE UNIQUE INDEX "users_unique_password" ON "users" ("password");
COMMIT;
