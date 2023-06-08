package main

import (
	"fmt"
	"os"

	"github.com/preferit/evotalk"
)

func main() {
	page := evotalk.Presentation().Page()
	if err := page.SaveAs("index.html"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
