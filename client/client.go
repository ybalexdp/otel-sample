package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	for {
		resp, err := http.Get("http://app:8080")
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Response status:", resp.Status)
			resp.Body.Close()
		}
		time.Sleep(10 * time.Second)
	}
}
