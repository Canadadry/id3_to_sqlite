package save

import (
	"app/model"
	"context"
	"flag"
	"fmt"

	id3v2 "github.com/bogem/id3v2/v2"
)

const Action = "save"

func Run(args []string) error {
	input := "db.sqlite"
	limit := 100

	fs := flag.NewFlagSet("arg", flag.ContinueOnError)
	fs.StringVar(&input, "i", input, "intput database")

	err := fs.Parse(args)
	if err != nil {
		fs.PrintDefaults()
		return err
	}

	ctx := context.Background()
	db, err := model.Open(ctx, input)
	if err != nil {
		return err
	}
	defer db.Close()

	offset := 0
	for {
		files, err := model.Fetch(db, ctx, limit, offset)
		if err != nil {
			return fmt.Errorf("while fecting file tag : %w", err)
		}
		if len(files) == 0 {
			break
		}
		offset = offset + limit

		err = processFiles(files)
		if err != nil {
			return fmt.Errorf("while processing files : %w", err)
		}
	}
	return nil
}

func processFiles(files []model.File) error {
	for _, f := range files {
		tag, err := id3v2.Open(f.Name, id3v2.Options{Parse: true})
		if err != nil {
			return fmt.Errorf("cannot open file : %w", err)
		}
		defer tag.Close()

		for name, value := range f.Fields {
			tag.AddTextFrame(tag.CommonID(name), tag.DefaultEncoding(), value)
		}
		if err = tag.Save(); err != nil {
			return fmt.Errorf("cannot save %s : %w", f.Name, err)
		}
	}
	return nil
}
