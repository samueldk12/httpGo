package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()
	
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://ondetrabalho.com.br", nil)
	
	if err != nil {
		panic(err)
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil{
		panic(err)
	}

	data, err := io.ReadAll(resp.Body)
	
	if err != nil{
		panic(err)
	}
	
	fmt.Println(string(data))
}