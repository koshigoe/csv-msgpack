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
	return Encode(os.Stdin, os.Stdout)
}

func decodeAction (c *cli.Context) error {
	return Decode(os.Stdin, os.Stdout)
}

func main() {
	app := cli.NewApp()

	app.Commands = []cli.Command{
		{
			Name:    "encode",
			Aliases: []string{"e"},
			Usage:   "output MessagePack",
			Action:  encodeAction,
		},
		{
			Name:    "decode",
			Aliases: []string{"d"},
			Usage:   "output CSV",
			Action:  decodeAction,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
