package database

import (
	"os"
)

type NaiveDB struct {
	database map[string]string
	dataFile *os.File
}

func New(filePath string) *NaiveDB {
	var dataFile *os.File
	if filePath != "" {
		var err error
		dataFile, err = os.Create(filePath)
		if err != nil {
			panic(err)
		}
	} else {
		dataFile = nil
	}
	return &NaiveDB{make(map[string]string), dataFile}
}

func (db *NaiveDB) Set(key, value string) bool {
	db.database[key] = value
	return true
}

func (db *NaiveDB) Get(key string) (string, bool) {
	return db.database[key], true
}

func (db *NaiveDB) Delete(key string) bool {
	delete(db.database, key)
	return true
}
