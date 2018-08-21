package test

import (
	"testing"

	"../naivedb"
)

func TestApiInterfaceEagerness(T *testing.T) {
	test := naivedb.New("test.txt")
	if test.Directory != "test.txt" {
		T.FailNow()
	}
}
