package subredditCrawler

import (
	"email"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/golang-lru"
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

var sentPosts, _ = lru.New(256)
var firstBatch = true

func GetPosts(subreddit string) ([]Post, error) {
	url := fmt.Sprintf("https://www.reddit.com/r/%s/.json?sort=top&t=day", subreddit)
	first := true

	var resp *http.Response
	var err error
	for { //reddit rate limits
		if first {
			first = false
		}

		resp, err = http.Get(url)
		if resp != nil {
			defer resp.Body.Close()
		}
		if err != nil {
			return nil, err
		}
		if resp.StatusCode == 200 {
			break
		} else if resp.StatusCode == 429 {
			time.Sleep(30 * time.Second)
			continue
		} else {
			return nil, errors.New(fmt.Sprintf("Error of %s", resp.Status))
		}
	}

	decoder := json.NewDecoder(resp.Body)
	/*body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}*/
	var response subredditResponse
	err = decoder.Decode(&response)
	//err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	posts := response.Data.Children
	if len(posts) == 0 {
		log.Println("Retrying")
		return GetPosts(subreddit)
	}
	return posts, nil
}

func GetMatchingPostsString(posts []Post, exp string) []Post {
	retPosts := make([]Post, 0)
	for _, post := range posts {
		if strings.Contains(strings.ToLower(post.Data.Title), strings.ToLower(exp)) {
			retPosts = append(retPosts, post)
		}
	}
	return retPosts
}

func GetMatchingPostsPoints(posts []Post, threshold int) []Post {
	retPosts := make([]Post, 0)
	for _, post := range posts {
		if post.Data.Score >= threshold {
			retPosts = append(retPosts, post)
		}
	}
	return retPosts
}

func CheckAndEmail(subreddit, exp, recipient string) {
	posts, err := GetPosts(subreddit)
	if err != nil {
		log.Println(err)
		return
	}

	//if matchPosts := GetMatchingPostsString(posts, exp); matchPosts != nil && len(matchPosts) != 0 {
	if matchPosts := GetMatchingPostsPoints(posts, 60); matchPosts != nil && len(matchPosts) != 0 {
		for _, post := range matchPosts {
			if seen := sentPosts.Contains(post.Data.Id); seen {
				continue
			}
			var mail email.Email
			mail.Subject = post.Data.Title + " " + time.Now().Format(time.ANSIC)
			mail.Recipient = recipient
			//mail.Body = fmt.Sprintf("This post matches %s: \n%s\n\nThe reddit link is here: \n https://www.reddit.com%s", exp, post.Data.Url, post.Data.Permalink)
			mail.Body = fmt.Sprintf("This post has %d points: \n%s\n\nThe reddit link is here: \n https://www.reddit.com%s", post.Data.Score, post.Data.Url, post.Data.Permalink)

			if !firstBatch {
				err = email.SendMail("www.marktai.com:25", mail)
				if err != nil {
					log.Println(err)
				}
				log.Println(fmt.Sprintf("Sent email about %s", post.Data.Title))
			}
			sentPosts.Add(post.Data.Id, struct{}{})
		}
		firstBatch = false
	} else {
		log.Println(fmt.Sprintf("No matching post for %s", exp))
	}
}

func Run(subreddit string, exp string, recipient string, wait time.Duration, killChan chan bool) {

	log.Println(fmt.Sprintf("Scanning /r/%s for %s every %s", subreddit, exp, wait.String()))

	ticker := time.NewTicker(wait)

	CheckAndEmail(subreddit, exp, recipient)

	for {
		select {
		case <-ticker.C:
			CheckAndEmail(subreddit, exp, recipient)
		case <-killChan:
			log.Println("Killing subredditCrawler")
			ticker.Stop()
			break
		}
	}
}
