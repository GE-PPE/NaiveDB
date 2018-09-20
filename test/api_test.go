package test

import (
	"testing"

	"../naivedb"
)

func TestApiInterfaceEagerness(T *testing.T) {
	test := naivedb.New("test")
	if test.Directory != "test" {
		T.FailNow()
	}
}

func TestWriteDatabase(T *testing.T) {
	test := naivedb.New("db")
	test.WriteString("Test writing into database")
}
