package dump

import (
	"app/model"
	"app/pkg/walk"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	id3v2 "github.com/bogem/id3v2/v2"
)

const Action = "dump"

func Run(args []string) error {
	entry := "."
	output := "db.sqlite"
	columns := "Album/Movie/Show title;Artist;Title;Year"
	ext := ".mp3"

	fs := flag.NewFlagSet("arg", flag.ContinueOnError)
	fs.StringVar(&entry, "i", entry, "input folder to scrap")
	fs.StringVar(&output, "o", output, "output database")
	fs.StringVar(&columns, "c", columns, "data to scrap spearated by coma")
	fs.StringVar(&ext, "e", ext, "filter file by extension separated by coma")

	err := fs.Parse(args)
	if err != nil {
		fs.PrintDefaults()
		return err
	}

	ctx := context.Background()

	db, err := model.Open(ctx, output)
	if err != nil {
		return err
	}
	defer db.Close()

	err = model.Migrate(db, ctx, strings.Split(columns, ";"))
	if err != nil {
		return err
	}

	extensions := strings.Split(ext, ";")

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
		tag, err := id3v2.ParseReader(f, id3v2.Options{Parse: true})
		if err != nil {
			f.Seek(0, io.SeekStart)
			start := [10]byte{}
			f.Read(start[:])
			return fmt.Errorf("cannot parse file %s starting with %v : %w", path, start, err)
		}
		mf := model.File{
			Name:   path,
			Fields: map[string]string{},
		}

		for _, c := range strings.Split(columns, ";") {
			mf.Fields[c] = tag.GetTextFrame(tag.CommonID(c)).Text
		}

		err = model.Upsert(db, ctx, []model.File{mf})
		if err != nil {
			return fmt.Errorf("cannot upsert file %s : %w", path, err)
		}
		return nil
	})
}
