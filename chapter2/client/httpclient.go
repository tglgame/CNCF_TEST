package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "http_client:", log.Ldate|log.Ltime|log.Lshortfile)

	url := "http://localhost/hello"
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	logger.Println("add 2 values to request header")
	req.Header.Add("ctag1", "value1")
	req.Header.Add("ctag2", "value2")
	if err != nil {
		logger.Fatal("http get request failed: ", err)
		panic(err)
	}
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	logger.Println("status code: ", resp.StatusCode, err)
	logger.Println("------------read body---------------")
	rb, _ := ioutil.ReadAll(resp.Body)
	for k, v := range rb {
		logger.Println(k, ":", v)
	}
	logger.Println("------------read header---------------")
	for k, v := range resp.Header {
		logger.Println(k, ":", v)
	}

	logger.Println("---------------check health-------------------")
	url = "http://localhost/healthz"
	resp, err = http.Get(url)
	if err != nil {
		logger.Fatal("http check health failed: ", err)
		panic(err)
	}
	logger.Println(resp.StatusCode)
}
