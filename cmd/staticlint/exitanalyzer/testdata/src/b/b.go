package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Hello")
}

func test() {
	os.Exit(1) // want "использование os.Exit вне функции main недопустимо"
}
