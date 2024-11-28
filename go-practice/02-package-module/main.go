package main

import (
	"fmt"

	"github.com/tryuuu/udemy-k8s-operator/go/package-module/mypackage"
)

func main() {
	name := mypackage.GetName()
	fmt.Printf("Hello, %s\n", name)
}
