package main

import (
	"encoding/json"
	"fmt"

	"github.com/justone/pod-justone-archiver/babashka"
	"github.com/mholt/archiver/v3"
)

type opts struct {
	OverwriteExisting *bool `json:"overwrite-existing"`
	CompressionLevel  *int  `json:"compression-level"`
}

var trueValue = true
var defaultOpts = opts{
	OverwriteExisting: &trueValue,
	CompressionLevel:  nil,
}

func setOverwrite(a interface{}, overwriteExisting *bool) {
	if overwriteExisting != nil {
		switch v := a.(type) {
		case *archiver.FileCompressor:
			v.OverwriteExisting = *overwriteExisting
		case *archiver.Rar:
			v.OverwriteExisting = *overwriteExisting
		case *archiver.Tar:
			v.OverwriteExisting = *overwriteExisting
		case *archiver.TarBz2:
			v.OverwriteExisting = *overwriteExisting
		case *archiver.TarGz:
			v.OverwriteExisting = *overwriteExisting
		case *archiver.TarLz4:
			v.OverwriteExisting = *overwriteExisting
		case *archiver.TarSz:
			v.OverwriteExisting = *overwriteExisting
		case *archiver.TarXz:
			v.OverwriteExisting = *overwriteExisting
		case *archiver.Zip:
			v.OverwriteExisting = *overwriteExisting
		}
	}
}

func setCompressionLevel(a interface{}, compressionLevel *int) {
	if compressionLevel != nil {
		switch v := a.(type) {
		case *archiver.Bz2:
			v.CompressionLevel = *compressionLevel
		case *archiver.Gz:
			v.CompressionLevel = *compressionLevel
		case *archiver.Lz4:
			v.CompressionLevel = *compressionLevel
		case *archiver.TarBz2:
			v.CompressionLevel = *compressionLevel
		case *archiver.TarGz:
			v.CompressionLevel = *compressionLevel
		case *archiver.TarLz4:
			v.CompressionLevel = *compressionLevel
		case *archiver.Zip:
			v.CompressionLevel = *compressionLevel
		}
	}
}

func processMessage(message *babashka.Message) (interface{}, error) {
	switch message.Op {
	case "describe":
		return &babashka.DescribeResponse{
			Format: "json",
			Namespaces: []babashka.Namespace{
				{
					Name: "pod.archiver",
					Vars: []babashka.Var{
						{
							Name: "archive",
							Meta: "{:name \"archive\", :doc \"Create an archive, combining all the sources into the destination. The type\\n  of archive is determined by the destination extension.\", :arglists ([sources destination opts])}",
						},
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

			o := defaultOpts
			if len(args) == 3 {
				if err := json.Unmarshal([]byte(args[2]), &o); err != nil {
					return nil, err
				}
			}

			a, err := archiver.ByExtension(dest)
			if err != nil {
				return nil, err
			}

			setOverwrite(a, o.OverwriteExisting)
			setCompressionLevel(a, o.CompressionLevel)

			return true, a.(archiver.Archiver).Archive(sources, dest)
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

			a, err := archiver.ByExtension(src)
			if err != nil {
				return nil, err
			}
			setOverwrite(a, o.OverwriteExisting)

			return true, a.(archiver.Unarchiver).Unarchive(src, dest)
		case "pod.archiver/extract":
			args := []json.RawMessage{}
			if err := json.Unmarshal([]byte(message.Args), &args); err != nil {
				return nil, err
			}

			var src string
			if err := json.Unmarshal([]byte(args[0]), &src); err != nil {
				return nil, err
			}

			var file string
			if err := json.Unmarshal([]byte(args[1]), &file); err != nil {
				return nil, err
			}

			var dest string
			if err := json.Unmarshal([]byte(args[2]), &dest); err != nil {
				return nil, err
			}

			o := defaultOpts
			if len(args) == 4 {
				if err := json.Unmarshal([]byte(args[3]), &o); err != nil {
					return nil, err
				}
			}

			a, err := archiver.ByExtension(src)
			if err != nil {
				return nil, err
			}
			setOverwrite(a, o.OverwriteExisting)

			return true, a.(archiver.Extractor).Extract(src, file, dest)
		case "pod.archiver/compress-file":
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

			a, err := archiver.ByExtension(dest)
			if err != nil {
				return nil, err
			}

			c := &archiver.FileCompressor{
				Compressor: a.(archiver.Compressor),
			}

			setCompressionLevel(a, o.CompressionLevel)
			setOverwrite(c, o.OverwriteExisting)

			return true, c.CompressFile(src, dest)
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
			c := &archiver.FileCompressor{
				Compressor:   a.(archiver.Compressor),
				Decompressor: a.(archiver.Decompressor),
			}

			setOverwrite(c, o.OverwriteExisting)

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

		res, err := processMessage(message)
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
