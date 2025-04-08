package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	from := flag.String("from", "", "source currency, required")
	to := flag.String("to", "EUR", "target currency")
	flag.Parse()

	value := flag.Arg(0)
	if value == "" {
		_, _ = fmt.Fprintln(os.Stderr, "missing amount to convert")
		flag.Usage()
		os.Exit(1)
	}

	fmt.Println(*from, *to, value)
}
