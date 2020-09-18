package main

import "fmt"

func main(){
	message := make(chan string)

	go func(){
		message <- "testing"
		message <- "123"
	}()

	fmt.Println(<-message)
	fmt.Println(<-message)
}
