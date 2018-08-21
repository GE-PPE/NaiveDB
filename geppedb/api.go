package naivedb

import "os"

/* NaiveDB: the database API */
type NaiveDB struct {
	Directory   string
	filePointer *os.File
}

// New NaiveDB
func New(Directory string) NaiveDB {
	return NaiveDB{Directory, nil}
}

func (db NaiveDB) Connect() {
	fp, err := os.Create(findDatabaseFile(db.Directory))
	if err != nil {
		panic(err)
	}
	db.filePointer = fp
}

func findDatabaseFile(directory string) string {
	return ""
}
