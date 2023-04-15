package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"
)

const url = "http://192.168.177.201:30402/osspgdl"

func Test1A(t *testing.T) {
	// Test 1A: put big obj
	byteData, err := os.ReadFile("./testCase/testCase1/1GB")
	if err != nil {
		panic(err)
	}
	stringdata := string(byteData)
	fmt.Println("load complete")
	time1 := time.Now()
	fmt.Println(0)
	PostHTTP("bigfile", stringdata, 1, 0)
	costTime := time.Now().Sub(time1).Milliseconds()
	fmt.Println(costTime)
	// Test 1B: get big obj
	// GetHTTP("bigfile", 1, 0)
}

func Test2A(t *testing.T) {
	// Test 2A: put 1 small obj
	byteData, err := os.ReadFile("./testCase/testCase3/file1.txt")
	if err != nil {
		panic(err)
	}
	stringdata := string(byteData)
	fmt.Println("load complete")
	time1 := time.Now()
	fmt.Println(0)
	PostHTTP("smallfile", stringdata, 1, 0)
	costTime := time.Now().Sub(time1).Milliseconds()
	fmt.Println(costTime)
}

func Test2B(t *testing.T) {
	for i := 1; i <= 1000; i++ {
		file := "./testCase3/file" + strconv.Itoa(i) + ".txt"
		byteData, err := os.ReadFile(file)
		if err != nil {
			panic(err)
		}
		stringdata := string(byteData)
		key := "smallfile" + strconv.Itoa(i)
		fmt.Println(key)
		time1 := time.Now()
		PostHTTP(key, stringdata, 1, 0)
		costTime := time.Now().Sub(time1).Milliseconds()
		fmt.Println(strconv.Itoa(int(costTime)) + "ms")
	}
}

func PostHTTP(key, value string, block, index int) {
	params := map[string]string{
		"namespace": "cubeuniverse",
		"name":      "testbucketclaim",
		"key":       key,
		"X-action":  "put",
		"value":     value,
		"block":     (string)(block),
		"index":     (string)(index),
	}
	reqBody, _ := json.Marshal(params)
	_, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		panic(err)
	}
}

func GetHTTP(key string, block, index int) []byte {
	params := map[string]string{
		"namespace": "cubeuniverse",
		"name":      "testbucketclaim",
		"key":       key,
		"X-action":  "get",
		"block":     (string)(block),
		"index":     (string)(index),
	}
	reqBody, _ := json.Marshal(params)
	_, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		panic(err)
	}
	return reqBody
}

func ListHTTP() {
	params := map[string]string{
		"namespace": "cubeuniverse",
		"name":      "testbucketclaim",
		"X-action":  "list",
	}
	reqBody, _ := json.Marshal(params)
	_, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		panic(err)
	}
}

// func main() {

// // Test 2A: put small objs

// // Test 2B: get small obj

// // Test C1: test AI module

// // Test C2: list all

// }
