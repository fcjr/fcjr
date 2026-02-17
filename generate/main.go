package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	RecentRepos    []Repo
}

type BlogPost struct {
	Title string
	URL   string
}

type Repo struct {
	Name        string
	URL         string
	Description string
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

	recentRepos, err := getRecentRepos("fcjr", 5)
	if err != nil {
		panic(err)
	}

	readmeData := ReadmeData{
		LatestBlogPost: latestBlogPost,
		RecentRepos:    recentRepos,
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

func getRecentRepos(username string, count int) ([]Repo, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/repos?sort=pushed&per_page=%d&type=owner", username, count*3)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching repos: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var apiRepos []struct {
		Name        string `json:"name"`
		HTMLURL     string `json:"html_url"`
		Description string `json:"description"`
		Fork        bool   `json:"fork"`
		Private     bool   `json:"private"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiRepos); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	ignore := map[string]bool{
		"homebrew-fcjr":       true,
		"frankchiarulli.com":  true,
	}

	var repos []Repo
	for _, r := range apiRepos {
		if r.Fork || r.Private || r.Description == "" || ignore[r.Name] {
			continue
		}
		repos = append(repos, Repo{
			Name:        r.Name,
			URL:         r.HTMLURL,
			Description: r.Description,
		})
		if len(repos) >= count {
			break
		}
	}
	return repos, nil
}
