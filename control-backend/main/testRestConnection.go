package main

import (
	"net/http"
)

func test() {
	req, err := http.NewRequest("GET", cephApiBase+"api/pool", nil)
	if err != nil {
		panic(err)
	}
	resJson, err := SendHttpsForJson(req, true)
	if err != nil {
		panic(err)
	}

	out, err := resJson.Encode()
	if err != nil {
		panic(err)
	}
	println(string(out))

	select {}
}
