package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetBasicInfo(r *http.Request) string {
	headString := fmt.Sprintln(fmt.Sprintf("[%s]", r.Method), r.URL.Path, r.Header["Content-Type"])
	return headString
}

func loadJson(body io.Reader, dest interface{}) {
	var err error
	js, err := io.ReadAll(body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(js, dest)
	if err != nil {
		panic(err)
	}
}
