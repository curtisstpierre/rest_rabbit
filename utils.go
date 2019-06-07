package main

import (
	"fmt"
	"log"

	"github.com/zond/gotomic"
)

// String this is a string
type String string

//Compare thread safe string to a string
func (s String) Compare(t gotomic.Thing) int {
	return compStrings(string(s), t.(string))
}

// failOnError if and error is passed in log and panic
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

// RandomString return a random string of length len
func RandomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(65 + r.Intn(25))
	}
	return string(bytes)
}

func compStrings(i, j string) int {
	l := len(i)
	if len(j) < l {
		l = len(j)
	}
	for ind := 0; ind < l; ind++ {
		if i[ind] < j[ind] {
			return -1
		} else if i[ind] > j[ind] {
			return 1
		}
	}
	if len(i) < len(j) {
		return -1
	} else if len(i) > len(j) {
		return 1
	}
	if i == j {
		return 0
	}
	panic(fmt.Errorf("wtf, %v and %v are not the same! ", i, j))
}
