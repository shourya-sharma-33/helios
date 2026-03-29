package utils

import "fmt"

func SendEmail(to string, token string) {

	link := "http://localhost:8080/verify/" + token

	// dev mode
	fmt.Println("VERIFY LINK:", link)
}
