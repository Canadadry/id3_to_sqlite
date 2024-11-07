package main

import (
	"app/model"
	"app/pkg/walk"
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	id3v2 "github.com/bogem/id3v2/v2"
)

func main() {
	if err := Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, "failed", err)
		os.Exit(1)
	}
}

func Run(args []string) error {
	entry := "."
	output := "db.sqlite"
	columns := "Album/Movie/Show title;Artist;Title;Year"

	fs := flag.NewFlagSet("arg", flag.ContinueOnError)
	fs.StringVar(&entry, "i", entry, "input folder to scrap")
	fs.StringVar(&output, "o", output, "output database")
	fs.StringVar(&columns, "c", columns, "data to scrap")

	err := fs.Parse(args)
	if err != nil {
		fs.PrintDefaults()
		return err
	}
	ctx := context.Background()
	db, err := model.Open(ctx, output, strings.Split(columns, ";"))
	if err != nil {
		return err
	}
	defer db.Close()

	return walk.Walk(entry, func(path string) error {
		tag, err := id3v2.Open(path, id3v2.Options{Parse: true})
		if err != nil {
			return fmt.Errorf("cannot open file : %w", err)
		}
		defer tag.Close()
		f := model.File{
			Name:   path,
			Fields: map[string]string{},
		}

		for _, c := range strings.Split(columns, ";") {
			f.Fields[c] = tag.GetTextFrame(tag.CommonID(c)).Text
		}
		err = model.Upsert(db, ctx, []model.File{f})
		if err != nil {
			return fmt.Errorf("cannot upsert file %s : %w", path, err)
		}
		return nil
	})
}
