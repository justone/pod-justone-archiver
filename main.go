package main

import (
	"encoding/json"
	"fmt"

	"github.com/babashka/pod-babashka-fswatcher/babashka"
	"github.com/mholt/archiver/v3"
)

func ProcessMessage(message *babashka.Message) (interface{}, error) {
	switch message.Op {
	case "describe":
		return &babashka.DescribeResponse{
			Format: "json",
			Namespaces: []babashka.Namespace{
				{
					Name: "pod.archiver",
					Vars: []babashka.Var{
						{Name: "archive"},
						{Name: "unarchive"},
						{Name: "extract"},
						{Name: "compress-file"},
						{Name: "decompress-file"},
					},
				},
			},
		}, nil
	case "invoke":
		switch message.Var {
		case "pod.archiver/archive":
			args := []json.RawMessage{}
			if err := json.Unmarshal([]byte(message.Args), &args); err != nil {
				return nil, err
			}

			var sources []string
			if err := json.Unmarshal([]byte(args[0]), &sources); err != nil {
				return nil, err
			}
			var dest string
			if err := json.Unmarshal([]byte(args[1]), &dest); err != nil {
				return nil, err
			}

			return true, archiver.Archive(sources, dest)
		case "pod.archiver/unarchive":
			var args []string
			if err := json.Unmarshal([]byte(message.Args), &args); err != nil {
				return nil, err
			}

			return true, archiver.Unarchive(args[0], args[1])
		case "pod.archiver/extract":
			var args []string
			if err := json.Unmarshal([]byte(message.Args), &args); err != nil {
				return nil, err
			}

			return true, archiver.Extract(args[0], args[1], args[2])
		case "pod.archiver/compress-file":
			var args []string
			if err := json.Unmarshal([]byte(message.Args), &args); err != nil {
				return nil, err
			}

			return true, archiver.CompressFile(args[0], args[1])
		case "pod.archiver/decompress-file":
			var args []string
			if err := json.Unmarshal([]byte(message.Args), &args); err != nil {
				return nil, err
			}

			return true, archiver.DecompressFile(args[0], args[1])
		default:
			return nil, fmt.Errorf("Unknown var %s", message.Var)
		}
	default:
		return nil, fmt.Errorf("Unknown op %s", message.Op)
	}
}

func main() {
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
