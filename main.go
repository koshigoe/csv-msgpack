package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"encoding/csv"

	"github.com/urfave/cli"
	"github.com/vmihailenco/msgpack"
)

func Encode(r io.Reader, w io.Writer) error {
	reader := csv.NewReader(r)

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		b, err := msgpack.Marshal(row)
		if err != nil {
			return err
		}

		fmt.Fprintf(w, "%s", b)
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

func decodeAction (c *cli.Context) error {
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
					Name: "input, i",
					Usage: "Input file path. (default STDIN)",
				},
				cli.StringFlag{
					Name: "output, o",
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
					Name: "input, i",
					Usage: "Input file path. (default STDIN)",
				},
				cli.StringFlag{
					Name: "output, o",
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
