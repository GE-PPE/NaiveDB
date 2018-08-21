package naivedb

import "os"

/* NaiveDB: the database API */
type NaiveDB struct {
	Directory   string
	filePointer *os.File
}

// New NaiveDB
func New(Directory string) NaiveDB {
	db := NaiveDB{Directory, nil}
	fp, err := os.Create(findDatabaseFile(db.Directory))
	if err != nil {
		panic(err)
	}
	db.filePointer = fp
	return db
}

func findDatabaseFile(directory string) string {
	return ""
}
