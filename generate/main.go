package main

import (
	"fmt"
	"io"
	"os"
	"text/template"
	"context"
	"time"

	"github.com/mmcdole/gofeed"
	_ "embed"
)

//go:embed README.md.tmpl
var readmeTemplate string

type ReadmeData struct{
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

	readmeData := ReadmeData{
		LatestBlogPost: getLatestBlogPost(),
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

func getLatestBlogPost() *BlogPost {
	fp := gofeed.NewParser()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	feed, err := fp.ParseURLWithContext("https://frankchiarulli.com/blog/rss.xml", ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to fetch blog RSS feed: %v\n", err)
		return nil
	}
	
	if len(feed.Items) == 0 {
		fmt.Fprintf(os.Stderr, "Warning: No blog posts found in RSS feed\n")
		return nil
	}
	
	latest := feed.Items[0]
	return &BlogPost{
		Title: latest.Title,
		URL:   latest.Link,
	}
}
