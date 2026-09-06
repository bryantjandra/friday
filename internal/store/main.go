package store

/*
	database/sql is a generic interface, it defines what a database looks like and can do without knowing whether the underlying DB is
	SQLite, MySQL, or Postgres.

	The actual DB-specific code lives in a seperate package called a `driver` (modernc.org/sqlite). This driver is the one that knows the details
	of speaking SQLite. This 2-Layer design is used such that our code looks nearly identical whether we are using Postgre or SQLite.

	Whenever we import a package in Go, its `init` function runs automatically before the `main` function runs. The SQLite driver's init registers
	itself with database/sql. Essentially, it says "hey if anyone asks for a driver named SQLite, I'm the one who handles it.". The SQLite driver's
	use in this file is only for that `init` function and thus we use a "_" to ensure Go doesn't throw a compile error.
*/

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func Open(path string) (*sql.DB, error) {

	/* Step 1. Create the parent directories */

	/* Chops off the filename and gives the path for the parent (e.g. ~/.friday/friday.db --> ~/friday */
	dir := filepath.Dir(path)

	/*
		os.MkdirAll creates that parent directory and any missing parents of that directory.
		0700 is setting the permissions for that directory. 7 (read (4) + write (2) + execute (1), for the owner of the directory). For groups and everyone else,
		they're not allowed to do aything.
	*/
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}

	/* Step 2. Setup handle to the DB*/

	/* The Open() function is lazy. Thus, it defers the real connection to the DB whenever we actually run a query.*/
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	/* Step 3. Force a real connection to the DB now using Ping() */

	/* This is such that we know whether our connection to the DB is fine.*/
	if err := db.Ping(); err != nil {
		return nil, err
	}

	/* Step 4. Turn on WAL mode. */

	/*
		PRAGMA is a specific command for congifuring the DB itself.
		journal_mode=WAL switches on Write-Ahead logging so that reads and writes don't block each other --> useful for background daemons.
		db.Exec runs a SQL statement but it doesn't return rows, db.Query returns rows.
	*/
	if _, err := db.Exec("PRAGMA journal_mode=WAL;"); err != nil {
		return nil, err
	}

	return db, nil
}
