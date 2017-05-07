package corgi

import (
	"github.com/hashicorp/golang-lru"
	"log"
	"math/rand"
	"strings"
	"subredditCrawler"
	"time"
)

const MAX_CACHE_SIZE = 256

var Images []string
var Gifs []string

var lists map[string][]string
var max_lengths map[string]int
var indexes map[string]int
var seen map[string]*lru.Cache

func Init(max_images, max_gifs int, timeout *time.Duration) {
	Images = make([]string, 0)
	Gifs = make([]string, 0)

	lists = map[string][]string{
		"image": Images,
		"gif":   Gifs,
	}

	max_lengths = map[string]int{
		"image": max_images,
		"gif":   max_gifs,
	}

	indexes = map[string]int{
		"image": 0,
		"gif":   0,
	}

	image_cache, err := lru.New(MAX_CACHE_SIZE)
	if err != nil {
		log.Println(err)
	}

	gif_cache, err := lru.New(MAX_CACHE_SIZE)
	if err != nil {
		log.Println(err)
	}

	seen = map[string]*lru.Cache{
		"image": image_cache,
		"gif":   gif_cache,
	}

	init_image_crawler(timeout)
	init_gif_crawler(timeout)
}

func AddToList(url string, image_type string) bool {

	if image_type == "gif" {
		if strings.Contains(url, "imgur") {
			if strings.HasSuffix(url, ".gifv") {
				url = strings.Replace(url, ".gifv", ".gif", 1)
			} else if !strings.HasSuffix(url, ".gif") {
				url = url + ".gif"
			}
		}
	}

	if image_type == "image" {
		if strings.Contains(url, "imgur") {
			if !(strings.HasSuffix(url, ".jpg") || strings.HasSuffix(url, ".png")) {
				url = url + ".png"
			}
		}
	}

	// if this corgi has been seen recently
	if seen[image_type].Contains(url) {
		log.Printf("Already seen %s", url)
		return false
	}

	log.Printf("Adding %s", url)

	seen[image_type].Add(url, struct{}{})

	
	max_length := max_lengths[image_type]

	if len(lists[image_type]) < max_length {
		lists[image_type] = append(lists[image_type], url)
	} else {
		index_to_evict := rand.Intn(len(lists[image_type]))
		lists[image_type][index_to_evict] = url
	}

	return true
}

func AddToImage(url string) bool {
	return AddToList(url, "image")
}

func AddToGif(url string) bool {
	return AddToList(url, "gif")
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func Get(image_type string) string {
	list := lists[image_type]
	max_length := max_lengths[image_type]
	index := indexes[image_type]

	ret_url := list[index]
	_, ok := seen[image_type].Get(ret_url)
	if !ok {
		_ = seen[image_type].Add(ret_url, struct{}{})
	}

	indexes[image_type] = (indexes[image_type] + 1) % min(max_length, len(list))

	return ret_url
}

func GetImage() string {
	return Get("image")
}

func GetGif() string {
	return Get("gif")
}

func crawl_subreddit(subreddit string, query_strings []string, image_type string) {

	for _, query_string := range query_strings {
		posts, err := subredditCrawler.GetSubredditPosts(subreddit, query_string)
		if err != nil {
			log.Printf("Error crawling %s?%s: %s", subreddit, query_string, err.Error())
			continue
		}

		for _, post := range posts {
			if !post.Data.Is_self && post.Data.Url != "" {
				AddToList(post.Data.Url, image_type)
			}
		}
	}
}

func init_crawler(subreddit string, query_strings []string, wait *time.Duration, image_type string) *time.Ticker {

	ticker := time.NewTicker(*wait)

	go func(subreddit string, query_strings []string, image_type string) {
		for {
			crawl_subreddit(subreddit, query_strings, image_type)
			<-ticker.C
		}
	}(subreddit, query_strings, image_type)

	return ticker
}

func init_image_crawler(timeout *time.Duration) {

	subreddit := "corgi"
	query_strings := []string{"sort=top&t=week", "sort=top&t=month", "sort=top&t=all"}

	_ = init_crawler(subreddit, query_strings, timeout, "image")
}

func init_gif_crawler(timeout *time.Duration) {
	subreddit := "CorgiGifs"
	query_strings := []string{"sort=top&t=month", "sort=top&t=all"}

	_ = init_crawler(subreddit, query_strings, timeout, "gif")
}
