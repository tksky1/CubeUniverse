package main

import (
	"control-backend/cubeControl"
	"net/http"
)

// 测试用
func test() {
	req, err := http.NewRequest("GET", cubeControl.CephApiBase+"api/pool?stats=true", nil)
	if err != nil {
		panic(err)
	}
	resJson, err := cubeControl.SendHttpsForJson(req)
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
