package main

import "fmt"

func greet(language, name string) (string, error) {
	if language == "Spanish" {
		return fmt.Sprintf("Hola, %s", name), nil
	}
	return fmt.Sprintf("Hello, %s", name), nil
}

func main() {
	fmt.Println("Hello, World")
	greetText, err := greet("English", "tryu")
	if err == nil {
		fmt.Println(greetText)
	}
}
