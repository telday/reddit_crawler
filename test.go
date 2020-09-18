package main

import (
	"fmt"

	"github.com/turnage/graw/reddit"
	"github.com/turnage/graw/streams"
)

func main(){
	bot, err := reddit.NewBotFromAgentFile("crawler.agent", 0)

	if err != nil {
		fmt.Println("There was an error getting a bot handle: ", err)
	}

	kill, errs := make(chan bool), make(chan error)

	posts, stream_errors := streams.SubredditComments(bot, kill, errs, "askreddit")

	for {
		select {
		case tmp_err := <-stream_errors:
			fmt.Println(tmp_err)
		case comment := <-posts:
			fmt.Println(comment)
		}
	}

	/*
	harvest, err := bot.Listing("/r/AskReddit", "")
	if err != nil {
		fmt.Println("There was an error getting the listings: ", err)
	}

	for _, post := range harvest.Posts[:1] {
		fmt.Printf("[%s] posted [%s] name: {%s}\n", post.Author, post.Title, post.Name)
		comment_harvest, err := bot.Listing("/r/AskReddit/comments", post.Name[3:])
		if err != nil {
			fmt.Println("Error getting post comments")
		}else{
			for _, comment := range comment_harvest.Comments {
				fmt.Printf("\tComment %s", comment.Body)
			}
			fmt.Println("Number of comments: ", len(comment_harvest.Comments))
		}
	}
	*/
}