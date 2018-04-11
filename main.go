package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/walesey/dicelang/parser"
)

func main() {
	if len(os.Args) > 1 {
		buf := bytes.NewReader([]byte(os.Args[1]))
		if output, err := parser.NewParser(buf).Execute(); err != nil {
			fmt.Println("ERROR: ", err)
		} else {
			fmt.Println(output)
		}
	}
}
