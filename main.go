package main

import (
	"bytes"
	"encoding/json"
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
			if jsonOutput, err := json.Marshal(output); err != nil {
				fmt.Println("ERROR: ", err)
			} else {
				fmt.Println(string(jsonOutput))
			}
		}
	}
}
