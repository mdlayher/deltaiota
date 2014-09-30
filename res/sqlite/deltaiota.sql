/* deltaiota sqlite schema */
PRAGMA foreign_keys = OFF;
BEGIN TRANSACTION;
/* users */
CREATE TABLE "users" (
	"id"       INTEGER PRIMARY KEY AUTOINCREMENT,
	"username" TEXT
);
CREATE UNIQUE INDEX "users_unique_username" ON "users" ("username");
COMMIT;
