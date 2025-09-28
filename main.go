/*
Copyright © 2025 Tolu Adesina <tolu.adesina...>
*/
package main

import (
	"fmt"

	configuration "github.com/tolubydesign/todo-go/app/config"
	"github.com/tolubydesign/todo-go/cmd"
)

func main() {
	_, err := configuration.BuildConfiguration()
	if err != nil {
		fmt.Println(".env config error", err.Error())
		panic(err)
	}

	cmd.Execute()
}
