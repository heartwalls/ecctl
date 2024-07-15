package main

import (
	"ecctl/internal/utils"
	"fmt"
	"os"
)

func main() {
	if err := utils.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
