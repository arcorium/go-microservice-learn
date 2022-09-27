package util

import (
	"log"
)

func LogError(err_ error) {
	if err_ != nil {
		log.Println(err_)
	}
}

func ExitOnError(err_ error) {
	if err_ != nil {
		log.Fatalln(err_)
	}
}

func PackReturn[T any](val_ T, err_ error) T {
	if err_ != nil {
		log.Println(err_)
	}
	return val_
}

func PackReturnExit[T any](val_ T, err_ error) T {
	if err_ != nil {
		log.Fatalln(err_)
	}
	return val_
}

func IsError(err_ error) bool {
	if err_ != nil {
		log.Println(err_)
		return true
	}
	return false
}
