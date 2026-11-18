package main

import (
	"fmt"
	"orderflow/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		fmt.Println(err)
		return
	}
}
