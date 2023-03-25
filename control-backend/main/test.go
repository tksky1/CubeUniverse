package main

import (
	"control-backend/cubecontrol"
	"net/http"
)

// 测试用
func test() {
	req, err := http.NewRequest("GET", cubecontrol.CephApiBase+"api/pool?stats=true", nil)
	if err != nil {
		panic(err)
	}
	resJson, err := cubecontrol.SendHttpsForJson(req)
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
