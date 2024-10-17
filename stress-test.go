package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"
)

// Change here the number of requests to simulate each iteration
var requestsToSimulate []int = []int{32, 64, 96, 128, 160, 192, 224, 256, 288, 320, 352, 384, 416, 448, 480, 512, 544, 576, 608, 640, 672, 704, 736, 768, 800, 832, 864, 896, 928, 960, 992, 1000}

//var requestsToSimulate []int = []int{32, 64, 96, 128}

var requestList []string

type Result struct {
	Duration time.Duration
	Err      error
}

func getRequestList(path string) []string {
	fp, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error while opening : %s\n", err)
		os.Exit(1)
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	reqList := make([]string, 0)
	//read line by line
	for scanner.Scan() {
		reqList = append(reqList, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error while scanning file: %s\n", err)
	}

	return reqList

}

func timeResponse(reqStr string, f func(string) ([]byte, error), c chan Result) {
	before := time.Now()
	_, err := f(reqStr)
	after := time.Now()
	c <- Result{after.Sub(before), err}
}

func makeGetRequest(reqStr string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, reqStr, nil)
	if err != nil {
		fmt.Printf("Error creating request: %s\n", err)
		os.Exit(1)
	}
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making request: %s\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body")
		os.Exit(1)
	}
	if resp.StatusCode == http.StatusGatewayTimeout {
		return body, err
	}
	fmt.Println("request status: ", resp.Status)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		fmt.Printf("The following request failed: %s\nwith the message: %s\nand status: %s", req.URL.String(), string(body), string(resp.StatusCode))
		return body, errors.New("StatusError")
		// os.Exit(1)
	}
	return body, err
}

// The API path is looked as an environment variable in API_PATH
func main() {
	apiPath, ok := os.LookupEnv("API_PATH")
	if !ok {
		panic("no API_PATH found")
	}
	b4Total := time.Now()
	s := getRequestList("text_Thu Oct 10 11:20:19 2024.txt")
	// select a random sample of requests
	reqList := make([][]string, len(requestsToSimulate))
	for i, v := range requestsToSimulate {
		for j := 0; j < v; j++ {
			reqList[i] = append(reqList[i], s[rand.Intn(len(s))])
		}
		fmt.Println(len(reqList[i]))
	}
	// fmt.Println(reqList)
	speeds := make([]time.Duration, 0)
	resultsChan := make(chan Result)
	errorsList := make([]error, 0)
	for _, requests := range reqList {
		fmt.Printf("launching %d goroutines\n", len(requests))
		for _, reqTxt := range requests {
			testUrl := apiPath + reqTxt
			go timeResponse(testUrl, makeGetRequest, resultsChan)

		}
		// We want to wait until all requests are served before
		// going onto the next batch of tests
		for i := 0; i < len(requests); i++ {
			speed := <-resultsChan
			if speed.Err != nil {
				errorsList = append(errorsList, speed.Err)
			}
			speeds = append(speeds, speed.Duration)
		}
	}
	// We create a new file with the current time
	times, err := os.Create(fmt.Sprintf("test_%s.txt", time.Now().Format(time.Stamp)))
	if err != nil {
		fmt.Printf("Error while creating file: %s", err)
	}

	// Write the times onto said file
	fmt.Println("times: ", speeds)
	writer := bufio.NewWriter(times)
	writer.WriteString(fmt.Sprintf("%s", speeds))
	writer.WriteString(fmt.Sprintf("number of errors: %s", string(len(errorsList))))
	writer.Flush()
	totalTime := time.Now().Sub(b4Total)
	fmt.Println("total execution time:", totalTime)
	fmt.Println("total number of errors:", len(errorsList))

}
