package controllers

import (
	"fmt"
	"net/http"
)

func GetBasicInfo(r *http.Request) string {
	headString := fmt.Sprintln(fmt.Sprintf("[%s]", r.Method), r.URL.Path, r.Header["Content-Type"])
	return headString
}
