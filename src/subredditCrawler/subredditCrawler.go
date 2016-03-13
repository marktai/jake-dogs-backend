package subredditCrawler

import (
	"email"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type subredditResponse struct {
	Kind string
	Data data
}

type data struct {
	After    string
	Before   string
	Modhash  string
	Children []Post
}

type Post struct {
	Kind string
	Data PostData
}

type PostData struct {
	Title       string
	Url         string
	Score       int
	Over_18     bool
	Created_utc float64
	Id          string
	Permalink   string
}

func GetPosts(subreddit string) (posts []Post, err error) {
	url := fmt.Sprintf("https://www.reddit.com/r/%s/new.json?sort=new", subreddit)

	res, err := http.Get(url)
	if err != nil {
		return
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var response subredditResponse
	err = decoder.Decode(&response)
	if err != nil {
		return
	}

	posts = response.Data.Children
	return
}

func GetMatchingPost(posts []Post, exp string) *Post {
	for _, post := range posts {
		if strings.Contains(strings.ToLower(post.Data.Title), strings.ToLower(exp)) {
			return &post
		}
	}
	return nil
}

func CheckAndEmail(subreddit, exp string) {
	posts, err := GetPosts(subreddit)
	if err != nil {
		log.Println(err)
		return
	}
	if post := GetMatchingPost(posts, exp); post != nil {
		var mail email.Email
		mail.Subject = post.Data.Title + " " + time.Now().String()
		mail.Recipient = "taifighterm@gmail.com"
		mail.Body = fmt.Sprintf("This post matches %s: \n%s\n\nThe reddit link is here: \n https://www.reddit.com%s", exp, post.Data.Url, post.Data.Permalink)

		email.SendMail("mail.marktai.com:25", mail)
		log.Println(fmt.Sprintf("Sent email about %s", post.Data.Title))
	} else {
		log.Println(fmt.Sprintf("No matching post for %s", exp))
	}
}

func Run(subreddit string, exp string, wait time.Duration, killChan chan bool) {

	log.Println(fmt.Sprintf("Scanning /r/%s for %s every %s", subreddit, exp, wait.String()))

	ticker := time.NewTicker(wait)

	CheckAndEmail(subreddit, exp)

	for {
		select {
		case <-ticker.C:
			CheckAndEmail(subreddit, exp)
		case <-killChan:
			log.Println("Killing subredditCrawler")
			ticker.Stop()
			break
		}
	}
}
