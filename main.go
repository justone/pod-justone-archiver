package main

import (
	"fmt"

	"github.com/babashka/pod-babashka-fswatcher/babashka"
	"github.com/mholt/archiver/v3"
)

func ProcessMessage(message *babashka.Message) (interface{}, error) {
	fmt.Println("Here")
	return nil, nil
}

func unarchive() {
	archiver.Unarchive("test.tar.gz", "test")
}

func main() {
	fmt.Println("Hello, World!")
}

func foo() {
	for {
		message, err := babashka.ReadMessage()
		if err != nil {
			babashka.WriteErrorResponse(message, err)
			continue
		}

		res, err := ProcessMessage(message)
		if err != nil {
			babashka.WriteErrorResponse(message, err)
			continue
		}

		describeRes, ok := res.(*babashka.DescribeResponse)
		if ok {
			babashka.WriteDescribeResponse(describeRes)
			continue
		}
		babashka.WriteInvokeResponse(message, res)
	}
}
