package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	resp, err := http.Post("https://ondetrabalho.com.br/api/v1/status", "application/json", nil)
	
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
}