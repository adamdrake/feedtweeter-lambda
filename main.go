package main

import (
	"errors"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mmcdole/gofeed"
)

type feed struct {
	url, hashtags string
}

var (
	consumerKey       = getenv("TWITTER_CONSUMER_KEY")
	consumerSecret    = getenv("TWITTER_CONSUMER_SECRET")
	accessToken       = getenv("TWITTER_ACCESS_TOKEN")
	accessTokenSecret = getenv("TWITTER_ACCESS_TOKEN_SECRET")
)

func getenv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		panic("missing required environment variable " + name)
	}
	return v
}

func randomItemFromFeed(feedURL string) (gofeed.Item, error) {
	// Reddit (for example) will return an HTTP 429 if no User-Agent string is provided...
	userAgent := "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:58.0) Gecko/20100101 Firefox/58.0"
	req, err := http.NewRequest("GET", feedURL, nil)
	if err != nil {
		return gofeed.Item{}, errors.New("could not construct http request")
	}
	req.Header.Add("User-Agent", userAgent)
	client := http.Client{}
	resp, err := client.Do(req)

	fp := gofeed.NewParser()
	feed, err := fp.Parse(resp.Body)
	resp.Body.Close() // Close response body in order to avoid leaks
	if err != nil {
		return gofeed.Item{}, errors.New("could not parse feed")

	}
	if feed == nil {
		return gofeed.Item{}, errors.New("gofeed returned nil feed")

	}
	if len(feed.Items) < 1 {
		return gofeed.Item{}, errors.New("empty feed returned from url")
	}

	idx := rand.Intn(len(feed.Items))
	return *feed.Items[idx], nil

}

func randomTweetFromFeed(api *anaconda.TwitterApi, feedURL, hashtags string) {
	item, err := randomItemFromFeed(feedURL)
	if err != nil {
		log.Fatal(err)
	}

	tweetText := item.Title + " " + item.Link + " " + hashtags

	tweet, err := api.PostTweet(tweetText, url.Values{})
	if err != nil {
		log.Fatal("Could not post tweet: ", tweet, " ", err)
	}
}

func tweeter() {
	rand.Seed(time.Now().UnixNano())
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	api := anaconda.NewTwitterApi(accessToken, accessTokenSecret)

	feeds := make([]feed, 0)
	feeds = append(feeds, feed{"https://www.reddit.com/user/adrake/m/data/.rss", "#data #bigdata #ai #ml"})
	feeds = append(feeds, feed{"http://www.datatau.com/rss", "#data #bigdata #ai #ml #datascience"})
	feeds = append(feeds, feed{"https://cryptocurrencynews.com/feed/", "#cryptocurrency #blockchain #btc #eth #xrp #xrb #ltc"})

	idx := rand.Intn(len(feeds))

	randomTweetFromFeed(api, feeds[idx].url, feeds[idx].hashtags)
}

func main() {
	lambda.Start(tweeter)
}
