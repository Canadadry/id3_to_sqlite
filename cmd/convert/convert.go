package convert

import (
	"app/pkg/runner"
	"app/pkg/walk"
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"
)

const Action = "convert"
const tagHeaderSize = 10

var (
	id3Identifier = []byte("ID3")
)

func Run(args []string) error {
	entry := "."
	ext := ".mp3"
	command := "eyeD3 --to-v2.3 \"%s\""
	dryRun := false

	fs := flag.NewFlagSet("arg", flag.ContinueOnError)
	fs.StringVar(&entry, "i", entry, "input folder to convert")
	fs.StringVar(&ext, "e", ext, "filter file by extension separated by coma")
	fs.StringVar(&command, "cmd", command, "change conversion command")
	fs.BoolVar(&dryRun, "dry-run", dryRun, "dry run, print command line")

	err := fs.Parse(args)
	if err != nil {
		fs.PrintDefaults()
		return err
	}
	extensions := strings.Split(ext, ";")

	converFunc := runner.Run
	if dryRun {
		converFunc = func(command string) error {
			fmt.Println(command)
			return nil
		}
	}

	return walk.Walk(entry, func(path string) error {
		found := false
		for _, e := range extensions {
			if strings.HasSuffix(path, e) {
				found = true
			}
		}
		if !found {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("cannot open file %s : %w", path, err)
		}
		defer f.Close()
		start := [tagHeaderSize]byte{}
		f.Read(start[:])

		if !isID3Tag(start) {
			return nil
		}
		if isID3TagVersionSupported(start) {
			return nil
		}

		return converFunc(fmt.Sprintf(command, path))
	})
}
func isID3Tag(data [tagHeaderSize]byte) bool {
	return bytes.Equal(data[0:3], id3Identifier)
}

func isID3TagVersionSupported(data [tagHeaderSize]byte) bool {
	return data[3] >= 3
}
