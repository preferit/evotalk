package main

import (
	"fmt"
	"os"

	"github.schneider-electric.com/SESA712749/evotalk"
)

func main() {
	page := evotalk.Presentation().Page()
	if err := page.SaveAs("index.html"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
