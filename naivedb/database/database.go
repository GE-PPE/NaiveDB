package database

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	emptyString = ""
)

type NaiveDB struct {
	database map[string]string
	filePath string
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func encodeDatabase(database map[string]string) ([]byte, error) {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(database)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func decodeDatabase(encodedDatabase []byte) (map[string]string, error) {
	var decodedMap map[string]string
	buffer := bytes.NewBuffer(encodedDatabase)
	decoder := gob.NewDecoder(buffer)

	err := decoder.Decode(&decodedMap)
	if err != nil {
		return nil, err
	}
	return decodedMap, nil
}

func New(filePath string) *NaiveDB {
	if filePath == emptyString {
		return &NaiveDB{make(map[string]string), filePath}
	}
	if fileExists(filePath) {
		dataFile, err := os.Open(filePath)
		if err != nil {
			panic(err)
		}
		encodedDatabase, err := ioutil.ReadAll(dataFile)
		if err != nil {
			panic(err)
		}

		var database map[string]string
		database, err = decodeDatabase(encodedDatabase)
		if err != nil {
			panic(err)
		}
		return &NaiveDB{database, filePath}
	}
	_, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	return &NaiveDB{make(map[string]string), filePath}
}

func (db *NaiveDB) Get(key string) (string, bool) {
	return db.database[key], true
}

func (db *NaiveDB) Set(key, value string) bool {
	if db.filePath == emptyString {
		db.database[key] = value
		return true
	}
	oldValue := db.database[key]

	dataFile, err := os.Create(db.filePath)
	defer dataFile.Close()

	if err != nil {
		return false
	}

	db.database[key] = value

	encodedDatabase, err := encodeDatabase(db.database)
	if err != nil {
		db.database[key] = oldValue
		return false
	}

	_, err = dataFile.Write(encodedDatabase)
	if err != nil {
		fmt.Print(err.Error())
		db.database[key] = oldValue
		return false
	}

	return true
}

func (db *NaiveDB) Delete(key string) bool {
	if db.filePath == emptyString {
		delete(db.database, key)
		return true
	}
	oldValue := db.database[key]

	dataFile, err := os.Create(db.filePath)
	defer dataFile.Close()

	if err != nil {
		return false
	}

	delete(db.database, key)

	encodedDatabase, err := encodeDatabase(db.database)
	if err != nil {
		db.database[key] = oldValue
		return false
	}

	_, err = dataFile.Write(encodedDatabase)
	if err != nil {
		db.database[key] = oldValue
		return false
	}

	return true
}
