package subredditCrawler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

const USER_AGENT = "script:corgi.server:v0.1"

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
	Is_self     bool
}

func GetSubredditPosts(subreddit, querystring string) ([]Post, error) {
	url := fmt.Sprintf("https://www.reddit.com/r/%s/top/.json?%s", subreddit, querystring)
	log.Printf("Crawling %s", url)
	first := true

	var resp *http.Response
	var err error
	for { //reddit rate limits
		if first {
			first = false
		}

		// Create a request and add the proper headers.
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("User-Agent", USER_AGENT)

		resp, err = http.DefaultClient.Do(req)
		if resp != nil {
			defer resp.Body.Close()
		}
		if err != nil {
			return nil, err
		}
		if resp.StatusCode == 200 {
			break
		} else if resp.StatusCode == 429 {
			log.Println("received 429")
			time.Sleep(120 * time.Second)
			continue
		} else {
			return nil, errors.New(fmt.Sprintf("Error of %s", resp.Status))
		}
	}

	decoder := json.NewDecoder(resp.Body)
	var response subredditResponse
	err = decoder.Decode(&response)
	if err != nil {
		return nil, err
	}

	posts := response.Data.Children
	if len(posts) == 0 {
		log.Println("Retrying")
		return GetSubredditPosts(subreddit, querystring)
	}
	return posts, nil
}
