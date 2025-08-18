package main

import (
	"context"
	"encoding/json"
	"fmt"
	"kenja2/bin/quicklog"
	"kenja2/documents"
	"kenja2/emood"
	"kenja2/mongodb"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type result struct {
	Query   string   `json:"query"`
	Results []string `json:"results"`
}

func run(a *mongodb.AtlasClient) error {
	ctx := context.Background()
	res := []result{}

	for k, v := range emood.EmoodTextMap {
		for _, emd := range v {
			c, err := a.UnsafeVectorSearch(
				ctx,
				documents.ITEM_TYPE_ANIME,
				documents.RATING_ALL_AGES,
				documents.VFieldName("text_vector"),
				emd.Tensor.Binary(),
			)
			if err != nil {
				return err
			}

			buf := []documents.FlatDocument{}
			if err := c.All(ctx, &buf); err != nil {
				return err
			}

			r := make([]string, len(buf))
			for i, d := range buf {
				r[i] = d.Name
			}

			q := fmt.Sprintf("%s : %s", k, emd.Src)

			res = append(res, result{Query: q, Results: r})

			time.Sleep(time.Millisecond * 100)
		}
	}

	for _, emd := range emood.EmoodImgList {
		c, err := a.UnsafeVectorSearch(
			ctx,
			documents.ITEM_TYPE_ANIME,
			documents.RATING_ALL_AGES,
			documents.VFieldName("image_vector"),
			emd.Tensor.Binary(),
		)
		if err != nil {
			return err
		}

		buf := []documents.FlatDocument{}
		if err := c.All(ctx, &buf); err != nil {
			return err
		}

		r := make([]string, len(buf))
		for i, d := range buf {
			r[i] = d.Name
		}

		res = append(res, result{Query: emd.Src, Results: r})

		time.Sleep(time.Millisecond * 100)
	}

	b, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile("result.json", b, 0664); err != nil {
		return err
	}

	return nil
}

func main() {
	l := quicklog.NewQuickLogger()

	err := godotenv.Load("../../.env")
	if err != nil {
		l.Fatal(err)
	}

	uri := os.Getenv("SEARCHENGINE_URI")
	if len(uri) == 0 {
		l.Fatal("env for search engine uri is not set")
	}

	atlas, err := mongodb.DirectConnect(uri, 3, l)
	if err != nil {
		l.Fatal(err)
	}

	if err := run(atlas); err != nil {
		l.Fatal(err)
	}
}
