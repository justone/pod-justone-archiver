package main

import (
	"encoding/json"
	"fmt"

	"github.com/babashka/pod-babashka-fswatcher/babashka"
	"github.com/mholt/archiver/v3"
)

type opts struct {
	OverwriteExisting bool `json:"overwrite-existing"`
}

var defaultOpts = opts{true}

func setOverwrite(a interface{}, overwriteExisting bool) {
	switch v := a.(type) {
	case *archiver.TarGz:
		v.OverwriteExisting = overwriteExisting
	case *archiver.Zip:
		v.OverwriteExisting = overwriteExisting
	case *archiver.FileCompressor:
		v.OverwriteExisting = overwriteExisting
	}
}

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
			args := []json.RawMessage{}
			if err := json.Unmarshal([]byte(message.Args), &args); err != nil {
				return nil, err
			}

			var src string
			if err := json.Unmarshal([]byte(args[0]), &src); err != nil {
				return nil, err
			}

			var dest string
			if err := json.Unmarshal([]byte(args[1]), &dest); err != nil {
				return nil, err
			}

			o := defaultOpts
			if len(args) == 3 {
				if err := json.Unmarshal([]byte(args[2]), &o); err != nil {
					return nil, err
				}
			}
			// fmt.Fprintf(os.Stderr, "%+v\n", o)

			a, err := archiver.ByExtension(src)
			if err != nil {
				return nil, err
			}
			// fmt.Fprintf(os.Stderr, "%+v\n", a.(*archiver.TarGz).OverwriteExisting)
			// a.(*archiver.TarGz).OverwriteExisting = true
			// fmt.Fprintf(os.Stderr, "%+v\n", a.(*archiver.TarGz).OverwriteExisting)
			// fmt.Fprintf(os.Stderr, "%+v\n", a.(*archiver.TarGz))
			setOverwrite(a, o.OverwriteExisting)
			// fmt.Fprintf(os.Stderr, "%+v\n", a.(*archiver.TarGz).OverwriteExisting)

			return true, a.(archiver.Unarchiver).Unarchive(src, dest)
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
			args := []json.RawMessage{}
			if err := json.Unmarshal([]byte(message.Args), &args); err != nil {
				return nil, err
			}

			var src string
			if err := json.Unmarshal([]byte(args[0]), &src); err != nil {
				return nil, err
			}

			var dest string
			if err := json.Unmarshal([]byte(args[1]), &dest); err != nil {
				return nil, err
			}

			o := defaultOpts
			if len(args) == 3 {
				if err := json.Unmarshal([]byte(args[2]), &o); err != nil {
					return nil, err
				}
			}

			a, err := archiver.ByExtension(src)
			if err != nil {
				return nil, err
			}
			// fmt.Fprintf(os.Stderr, "%+v\n", a.OverwriteExisting)
			c := &archiver.FileCompressor{Compressor: a.(archiver.Compressor), Decompressor: a.(archiver.Decompressor)}

			setOverwrite(c, o.OverwriteExisting)
			// fmt.Fprintf(os.Stderr, "%+v\n", o.OverwriteExisting)
			// fmt.Fprintf(os.Stderr, "%+v\n", reflect.TypeOf(o))
			// fmt.Fprintf(os.Stderr, "%+v\n", c.OverwriteExisting)
			// fmt.Fprintf(os.Stderr, "%+v\n", reflect.TypeOf(c))

			return true, c.DecompressFile(src, dest)
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
