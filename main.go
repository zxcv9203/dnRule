package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	token, err := getToken()

	if err != nil {
		log.Fatalf("토큰을 가져오는데 실패했습니다 [ %v ]", err)
	}
	fmt.Println(token)
}

func getToken() (string, error) {
	token := os.Getenv("INPUT_TOKEN")
	if token == "" {
		return "", fmt.Errorf("GITHUB_TOKEN이 존재하지 않습니다.")
	}
	return token, nil
}
