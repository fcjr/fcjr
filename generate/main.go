package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"text/template"
	"time"

	_ "embed"

	"github.com/mmcdole/gofeed"
)

//go:embed README.md.tmpl
var readmeTemplate string

type ReadmeData struct {
	LatestBlogPost *BlogPost
}

type BlogPost struct {
	Title string
	URL   string
}

func main() {

	if len(os.Args) < 2 || os.Args[1] == "" {
		fmt.Fprintln(os.Stderr, "No output path provided")
		os.Exit(1)
	}
	var outPath = os.Args[1]

	outFile, err := os.OpenFile(outPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	latestBlogPost, err := getLatestBlogPost()
	if err != nil {
		panic(err)
	}

	readmeData := ReadmeData{
		LatestBlogPost: latestBlogPost,
	}

	err = renderReadmeToWriter(&readmeData, outFile)
	if err != nil {
		panic(err)
	}
}

func renderReadmeToWriter(readmeData *ReadmeData, writer io.Writer) error {
	template, err := template.New("").Parse(readmeTemplate)
	if err != nil {
		return err
	}

	return template.ExecuteTemplate(writer, "", readmeData)
}

func getLatestBlogPost() (*BlogPost, error) {
	fp := gofeed.NewParser()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	feed, err := fp.ParseURLWithContext("https://frankchiarulli.com/blog/rss.xml", ctx)
	if err != nil {
		return nil, fmt.Errorf("error parsing RSS feed: %w", err)
	}

	if len(feed.Items) == 0 {
		return nil, fmt.Errorf("no blog posts found")
	}

	latest := feed.Items[0]
	return &BlogPost{
		Title: latest.Title,
		URL:   latest.Link,
	}, nil
}
