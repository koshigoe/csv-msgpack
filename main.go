package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/urfave/cli"
	"github.com/vmihailenco/msgpack"
)

type CsvResult struct {
	Error error
	Row   []string
}

type MsgpackResult struct {
	Error error
	Row   []byte
}

type WriteResult struct {
	Error error
}

func csvStreamOwner(done <-chan interface{}, r io.Reader) <-chan CsvResult {
	csvStream := make(chan CsvResult)

	go func(r io.Reader) {
		defer close(csvStream)

		reader := csv.NewReader(r)
		for {
			row, err := reader.Read()
			if err == io.EOF {
				return
			}
			result := CsvResult{Error: err, Row: row}
			select {
			case <-done:
				return
			case csvStream <- result:
			}

		}
	}(r)

	return csvStream
}

func msgpackStreamOwner(done <-chan interface{}, csvStream <-chan CsvResult) <-chan MsgpackResult {
	msgpackStream := make(chan MsgpackResult)

	go func(csvStream <-chan CsvResult) {
		defer close(msgpackStream)

		var result MsgpackResult

		for csvResult := range csvStream {
			if csvResult.Error != nil {
				result = MsgpackResult{Error: csvResult.Error, Row: nil}
				msgpackStream <- result
				break
			}

			b, err := msgpack.Marshal(csvResult.Row)
			result = MsgpackResult{Error: err, Row: b}
			select {
			case <-done:
				return
			case msgpackStream <- result:
			}
		}
	}(csvStream)

	return msgpackStream
}

func writeStreamOwner(done <-chan interface{}, msgpackStream <-chan MsgpackResult, w io.Writer) <-chan WriteResult {
	writeStream := make(chan WriteResult)

	go func(w io.Writer) {
		defer close(writeStream)

		var result WriteResult
		for msgpackResult := range msgpackStream {
			if msgpackResult.Error != nil {
				result = WriteResult{Error: msgpackResult.Error}
				writeStream <- result
				break
			}

			select {
			case <-done:
				return
			default:
				fmt.Fprintf(w, "%s", msgpackResult.Row)
			}
		}

		result = WriteResult{Error: nil}
		writeStream <- result
	}(w)

	return writeStream
}

func Encode(r io.Reader, w io.Writer) error {
	done := make(chan interface{})
	defer close(done)

	csvStream := csvStreamOwner(done, r)
	msgpackStream := msgpackStreamOwner(done, csvStream)
	writeStream := writeStreamOwner(done, msgpackStream, w)

	for result := range writeStream {
		if result.Error != nil {
			return result.Error
		}
		select {
		case <-done:
			break
		default:
		}
	}

	return nil
}

func Decode(r io.Reader, w io.Writer) error {
	decoder := msgpack.NewDecoder(r)
	writer := csv.NewWriter(w)

	var row []string
	for {
		// TODO: How do I accept variable-length array?
		err := decoder.Decode(&row)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		err = writer.Write(row)
		if err != nil {
			return err
		}
	}
	writer.Flush()

	return nil
}

func encodeAction(c *cli.Context) error {
	var r io.ReadCloser
	var w io.WriteCloser
	var err error

	inputPath := c.String("input")
	if len(inputPath) > 0 {
		r, err = os.Open(inputPath)
		if err != nil {
			return err
		}
		defer r.Close()
	} else {
		r = os.Stdin
	}

	outputPath := c.String("output")
	if len(outputPath) > 0 {
		w, err = os.Create(outputPath)
		if err != nil {
			return err
		}
		defer w.Close()
	} else {
		w = os.Stdout
	}

	return Encode(r, w)
}

func decodeAction(c *cli.Context) error {
	var r io.ReadCloser
	var w io.WriteCloser
	var err error

	inputPath := c.String("input")
	if len(inputPath) > 0 {
		r, err = os.Open(inputPath)
		if err != nil {
			return err
		}
		defer r.Close()
	} else {
		r = os.Stdin
	}

	outputPath := c.String("output")
	if len(outputPath) > 0 {
		w, err = os.Create(outputPath)
		if err != nil {
			return err
		}
		defer w.Close()
	} else {
		w = os.Stdout
	}

	return Decode(r, w)
}

func main() {
	app := cli.NewApp()

	app.Commands = []cli.Command{
		{
			Name:    "encode",
			Aliases: []string{"e"},
			Usage:   "output MessagePack",
			Action:  encodeAction,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "input, i",
					Usage: "Input file path. (default STDIN)",
				},
				cli.StringFlag{
					Name:  "output, o",
					Usage: "Output file path. (default STDOUT)",
				},
			},
		},
		{
			Name:    "decode",
			Aliases: []string{"d"},
			Usage:   "output CSV",
			Action:  decodeAction,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "input, i",
					Usage: "Input file path. (default STDIN)",
				},
				cli.StringFlag{
					Name:  "output, o",
					Usage: "Output file path. (default STDOUT)",
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
