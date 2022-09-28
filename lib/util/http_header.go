package util

import (
	"net/http"
)

func SetDefaultHeader(writer_ http.ResponseWriter) {
	//fmt.Printf("Address of writer_ in util.SetDefaultHeader : %p\n", writer_)
	head := writer_.Header()
	head.Set("Content-Type", "application/json")
	head.Set("Accept-Charset", "UTF-8")
}
