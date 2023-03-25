package main

import "control-backend/cubeControl"

func test2() {
	json, err := cubeControl.ListObjectBucketClaim()
	if err != nil {
		panic(err)
	}
	out, err := json.Encode()
	if err != nil {
		panic(err)
	}
	println(string(out))
	select {}
}

// 测试用
func test() {
	test2()
	/*
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
	*/
}
