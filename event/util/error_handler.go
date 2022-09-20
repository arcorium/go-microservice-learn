package util

import (
	"log"
)

func LogError(err_ error) {
	if err_ != nil {
		log.Println(err_)
	}
}

func GetReturn[T any](val_ T, err_ error) T {
	if err_ != nil {
		log.Println(err_)
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
