/* deltaiota sqlite schema */
BEGIN TRANSACTION;
/* notifications */
CREATE TABLE "notifications" (
	"id"          INTEGER PRIMARY KEY AUTOINCREMENT
	, "user_id"   INTEGER NOT NULL
	, "timestamp" INTEGER NOT NULL
	, "read"      INTEGER NOT NULL
	, "text"         TEXT NOT NULL
	, "uri"          TEXT NOT NULL

	, FOREIGN KEY(user_id) REFERENCES users(id)
);
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
