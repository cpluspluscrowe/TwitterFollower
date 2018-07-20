package main

import (
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"net/http"
	"time"
)

func doEvery(d time.Duration, f func()) {
	for _ = range time.Tick(d) {
		f()
	}
}

func main() {
	doEvery(500000*time.Millisecond, followUserRounds)
}

func followUserRounds() {
	//	followUsers("Python")
	followUsers("#Golang")
	//	followUsers("C++")
	//	followUsers("Kotlin")
	//	followUsers("Scala")
	//	followUsers("Java")
}

func followUsers(toSearchFor string) {
	client := getTwitterClient()
	search, _, err := client.Search.Tweets(&twitter.SearchTweetParams{
		Query: toSearchFor,
		Count: 20,
	})
	if err != nil {
		panic(err)
	}

	var users []UserEntity
	for _, element := range search.Statuses {
		if element.Lang == "en" {
			user := UserEntity{
				ScreenName: element.User.ScreenName,
				UserID:     element.User.ID,
			}
			follow(user.ScreenName, client)
			fmt.Println(user)
			users = append(users, user)
		}
	}
	fmt.Println()
}

func addrTrue() *bool {
	boolean := true
	return &boolean
}

func follow(user string, client *twitter.Client) {
	//preventReachingLimit()
	_, _, err := client.Friendships.Create(&twitter.FriendshipCreateParams{
		ScreenName: user,
		Follow:     addrTrue(),
	})
	if err != nil {
		panic(err)
	}
}

type UserEntity struct {
	ScreenName        string `json:"screenName"`
	UserID            int64  `json:"userID"`
	FollowedTimestamp int64  `json:"followedTimestamp"`
}

func Tweet(tweetText string) {
	httpClient := getClient()
	client := twitter.NewClient(httpClient)

	_, _, err := client.Statuses.Update(tweetText, nil)
	if err != nil {
		panic(err)
	}
}

func FakeTweet(tweetText string) {
	fmt.Println(tweetText)
}

func GetTweets(highlight string) {
	httpClient := getClient()
	client := twitter.NewClient(httpClient)

	search, resp, err := client.Search.Tweets(&twitter.SearchTweetParams{
		Query: highlight,
	})
	fmt.Println(search, resp, err)
}

func Search(look4 string) {
	twitterClient := getTwitterClient()
	search, resp, err := twitterClient.Search.Tweets(&twitter.SearchTweetParams{
		Query: "Golang",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(search, resp, err)
}

func getClient() *http.Client {
	consumerKey := os.Getenv("consumerKey")
	consumerSecret := os.Getenv("consumerSecret")
	accessToken := os.Getenv("accessToken")
	accessSecret := os.Getenv("accessSecret")
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	return httpClient
}

func getTwitterClient() *twitter.Client {
	httpClient := getClient()
	client := twitter.NewClient(httpClient)
	return client
}
