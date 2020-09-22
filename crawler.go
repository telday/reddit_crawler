package main

import (
    "fmt"
    "sync"
    "regexp"
    "strings"

    "github.com/turnage/graw/reddit"
    "github.com/turnage/graw/streams"
    "gorm.io/gorm"
    "gorm.io/driver/sqlite"
)

type Subreddit struct {
	Name string `gorm:"primaryKey"`
	DisplayName string
}

const MAX_COMMENTS float64 = 1e50
var reddit_re, _ = regexp.Compile("r/[a-zA-Z0-9]+")
var db, err = gorm.Open(sqlite.Open("subreddits.db"), &gorm.Config{
	SkipDefaultTransaction: true,
})

func commentHasSubreddit(comment *reddit.Comment) []string{
    slice := reddit_re.FindAllString(comment.Body, -1)
	var subreddits []string
    for i := 0; i < len(slice); i++ {
        if strings.ToLower(comment.Subreddit) != strings.ToLower(slice[i][2:]){
			subreddits = append(subreddits, slice[i][2:])
        }
    }
    return subreddits
}

func insertIntoDatabase(name string){
	var dbSubs Subreddit
	transx := db.Begin()
	nameLower := strings.ToLower(name)
	res := transx.Where("name = ?", nameLower).First(&dbSubs)
	if res.RowsAffected > 0 {
		transx.Rollback()
	}else{
		subreddit := Subreddit{Name: nameLower, DisplayName: name}
		err := transx.Create(&subreddit).Error
		if err != nil {
			fmt.Printf("Error adding subreddit %s to database\n", name)
			transx.Rollback()
		}else{
			fmt.Printf("Adding %s to database\n", name)
			transx.Commit()
		}
	}
}

func crawlSubreddit(bot *reddit.Bot, subreddit string, wg *sync.WaitGroup) int {
    defer wg.Done()
    fmt.Println("Starting bot for subreddit: ", subreddit)
    kill, stream_errors := make(chan bool), make(chan error)

    comments, err := streams.SubredditComments(*bot, kill, stream_errors, subreddit)

    if err != nil {
        fmt.Println("Error getting comments stream")
    }
    var counter float64 = 0
    for ;counter < MAX_COMMENTS; counter++ {
        select {
        case s_error := <-stream_errors:
            fmt.Printf("Stream Error: %s\n", s_error)
        case comment := <-comments:
			subreddits := commentHasSubreddit(comment)

			for i := 0; i < len(subreddits); i++ {
				insertIntoDatabase(subreddits[i])
			}
        }
    }
    return 0
}

func main(){
    var wg sync.WaitGroup
    bot, err := reddit.NewBotFromAgentFile("crawler.agent", 0)

    if err != nil {
        fmt.Println("There was an error getting a bot handle: ", err)
        return
    }

    go crawlSubreddit(&bot, "AskReddit", &wg)
    wg.Add(1)

    wg.Wait()
}
