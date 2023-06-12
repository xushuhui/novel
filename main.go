package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/carlmjohnson/requests"
	"github.com/spf13/cast"
	"github.com/urfave/cli/v2"
)

const baseURL = "https://api.aixdzs.com"

type Book struct {
	ID            string `json:"_id"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	ShortIntro    string `json:"shortIntro"`
	Cover         string `json:"cover"`
	Category      string `json:"cat"`
	FollowerCount string `json:"followerCount"`
	Status        string `json:"zt"`
	UpdatedTime   string `json:"updated"`
	LastChapter   string `json:"lastchapter"`
}

type Chapter struct {
	Title     string `json:"title"`
	Link      string `json:"link"`
	Unreadble bool   `json:"unreadble"`
}

type ChapterList struct {
	ID           string    `json:"_id"`
	ChapterCount int       `json:"chaptercount"`
	UpdatedTime  string    `json:"chaptersUpdated"`
	Chapters     []Chapter `json:"chapters"`
}
type ChapterResp struct {
	MixToc ChapterList `json:"mixToc"`
}
type ContentResp struct {
	Chapter Content `json:"chapter"`
}
type Content struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type BookInfo struct {
	ID            string `json:"_id"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	Cover         string `json:"cover"`
	LongIntro     string `json:"longIntro"`
	Status        string `json:"zt"`
	Category      string `json:"cat"`
	WordCount     string `json:"wordCount"`
	UpdatedTime   string `json:"updated"`
	ChaptersCount string `json:"chaptersCount"`
	LastChapter   string `json:"lastChapter"`
	FollowerCount string `json:"followerCount"`
}

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:   "help",
				Usage:  "Show available commands",
				Action: help,
			},
			{
				Name:   "search",
				Usage:  "Search books by keyword",
				Action: search,
			},
			{
				Name:   "info",
				Usage:  "Get detailed information of a book",
				Action: info,
			},
			{
				Name:   "chapter",
				Usage:  "Get a list of chapters of a book",
				Action: chapter,
			},
			{
				Name:   "content",
				Usage:  "Get the content of a chapter",
				Action: content,
			},
			{
				Name:  "download",
				Usage: "Download the content of a book",

				Action: download,
			},
		},
	}
	app.Run(os.Args)
}

func help(c *cli.Context) error {
	commands := []string{"help", "search", "info", "chapter", "content", "download"}
	descriptions := []string{
		"Show available commands",
		"Search books by keyword",
		"Get detailed information of a book",
		"Get a list of chapters of a book",
		"Get the content of a chapter",
		"Download the content of a book",
	}

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	fmt.Fprintln(w, "Available commands:")
	for i := 0; i < len(commands); i++ {
		fmt.Fprintf(w, "\t%s\t%s\n", commands[i], descriptions[i])
	}

	return nil
}

func search(c *cli.Context) error {
	keyword := c.Args().Get(0)
	if keyword == "" {
		return cli.ShowCommandHelp(c, "search")
	}
	var body struct {
		Books []Book `json:"books"`
	}

	err := requests.URL(baseURL+"/book/search").Param("query", keyword).ToJSON(&body).Fetch(context.Background())

	if err != nil {
		return cli.Exit(fmt.Sprintf("Failed to get book info: %s", err), 1)
	}

	results := map[string][]Book{"books": body.Books}
	resultsJSON, _ := json.MarshalIndent(results, "", "  ")
	fmt.Println(string(resultsJSON))

	return nil
}
func getBookInfo(id string) (info *BookInfo, err error) {
	endpoint := fmt.Sprintf("%s/book/%s", baseURL, id)

	err = requests.URL(endpoint).ToJSON(info).Fetch(context.Background())

	if err != nil {
		return nil, cli.Exit(fmt.Sprintf("Failed to get book info: %s", err), 1)
	}

	return
}
func info(c *cli.Context) error {
	bookID := c.Args().Get(0)
	if bookID == "" {
		return cli.ShowCommandHelp(c, "info")
	}

	info, _ := getBookInfo(bookID)

	results := map[string]*BookInfo{"info": info}
	resultsJSON, _ := json.MarshalIndent(results, "", "  ")
	fmt.Println(string(resultsJSON))

	return nil
}

func chapter(c *cli.Context) error {
	bookID := c.Args().Get(0)
	if bookID == "" {
		return cli.ShowCommandHelp(c, "chapter")
	}

	chapters, _ := getChapterCount(bookID)
	results := map[string]ChapterList{"chapters": chapters.MixToc}
	resultsJSON, _ := json.MarshalIndent(results, "", "  ")
	fmt.Println(string(resultsJSON))

	return nil
}

func content(c *cli.Context) error {
	bookID := c.Args().Get(0)
	chapterIndex, err := strconv.Atoi(c.Args().Get(1))
	if bookID == "" || err != nil {
		return cli.ShowCommandHelp(c, "content")
	}

	endpoint := fmt.Sprintf("%s/chapter/%s/%d", baseURL, bookID, chapterIndex)

	var content Content
	err = requests.URL(endpoint).ToJSON(&content).Fetch(context.Background())

	if err != nil {
		return cli.Exit(fmt.Sprintf("Failed to get chapter content: %s", err), 1)
	}

	results := map[string]Content{"content": content}
	resultsJSON, _ := json.Marshal(results)
	fmt.Println(string(resultsJSON))

	return nil
}

func download(c *cli.Context) error {
	bookID := c.Args().Get(0)
	if bookID == "" {
		return cli.ShowCommandHelp(c, "download")
	}

	info, err := getBookInfo(bookID)

	if err != nil {
		return cli.Exit(fmt.Sprintf("Failed to get chapter content: %s", err), 1)
	}
	filename := "./" + info.Title + ".txt"

	chapterCount := cast.ToInt(info.ChaptersCount)

	contents := make([]Content, 0, chapterCount)
	for i := 0; i < chapterCount; i++ {
		content, err := getChapterContent(bookID, i+1)
		if err != nil {
			fmt.Printf("Failed to get content of chapter %d: %s\n", i+1, err)
			continue
		}

		contents = append(contents, content)
		if i%5 == 0 {
			time.Sleep(1 * time.Second)
		}
	}

	err = writeContentToTXT(filename, contents)
	if err != nil {
		return cli.Exit(fmt.Sprintf("Failed to write content to txt file: %s", err), 1)
	}

	fmt.Printf("Content of book %s has been downloaded and saved to %s.\n", bookID, filename)

	return nil
}

func getChapterCount(bookID string) (r *ChapterResp, err error) {
	endpoint := fmt.Sprintf("%s/content/%s", baseURL, bookID)

	err = requests.URL(endpoint).ToJSON(r).Fetch(context.Background())

	if err != nil {
		return nil, cli.Exit(fmt.Sprintf("Failed to get chapter list: %s", err), 1)
	}

	return
}

func getChapterContent(bookID string, chapterIndex int) (c Content, err error) {
	endpoint := fmt.Sprintf("%s/chapter/%s/%d", baseURL, bookID, chapterIndex)
	var r *ContentResp
	err = requests.URL(endpoint).ToJSON(r).Fetch(context.Background())

	if err != nil {
		return c, cli.Exit(fmt.Sprintf("Failed to get chapter list: %s", err), 1)
	}

	return r.Chapter, nil
}

func writeContentToTXT(filename string, contents []Content) error {
	var data []byte
	for _, content := range contents {
		data = append(data, []byte(content.Title)...)
		data = append(data, []byte(content.Body)...)
		data = append(data, []byte("\n")...)
	}

	err := ioutil.WriteFile(filename, data, 0777)
	if err != nil {
		return err
	}

	return nil
}
