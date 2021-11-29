package main

import (
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
	req.Header.Add("clienttag", "clientvalue1")

	if err != nil {
		logger.Fatal("http get request failed: ", err)
		panic(err)
	}
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	logger.Println("status code: ", resp.StatusCode, err)
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
