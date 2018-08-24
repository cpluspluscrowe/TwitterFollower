package main

import (
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"net/http"
	"os"
	"strings"
	"time"
)

func doEvery(d time.Duration, f func(string), arg string) {
	for _ = range time.Tick(d) {
		f(arg)
	}
}

func main() {
	doEvery(1507*time.Second, followUserRounds, "#Golang")
	doEvery(2126*time.Second, followUserRounds, "#Scala")
	doEvery(3125*time.Second, followUserRounds, "#C++")
	doEvery(4124*time.Second, followUserRounds, "#Java")
	doEvery(5123*time.Second, followUserRounds, "#PhD")
}

func followUserRounds(search4 string) {
	followUsers(search4)
}

func isUserMyFriend(userId int64) bool {
	client := getTwitterClient()
	relationship, _, err := client.Friendships.Show(&twitter.FriendshipShowParams{TargetID: userId})
	if err != nil {
		panic(err)
	}
	return relationship.Source.Following
}

func TweetUsersFirstTweet(userId int64) {
	client := getTwitterClient()
	tweets, _, err := client.Timelines.UserTimeline(&twitter.UserTimelineParams{UserID: userId})
	if err != nil {
		panic(err)
	}
	firstTweetId := tweets[0].ID
	client.Statuses.Retweet(firstTweetId, &twitter.StatusRetweetParams{})
}

func followUsers(toSearchFor string) {
	client := getTwitterClient()
	search, _, err := client.Search.Tweets(&twitter.SearchTweetParams{
		Query: toSearchFor,
		Count: 10,
	})
	if err != nil {
		fmt.Println(err)
		time.Sleep(300 * time.Second)
		return
	}

	var users []UserEntity
	for _, element := range search.Statuses {
		if element.Lang == "en" {
			user := UserEntity{
				ScreenName: element.User.ScreenName,
				UserID:     element.User.ID,
			}
			if !isUserMyFriend(user.UserID) {
				follow(user.ScreenName, client)
				fmt.Println("Followed User! ", user, isUserMyFriend(user.UserID))
				users = append(users, user)
				TweetUsersFirstTweet(user.UserID)
				return
			} else {
				fmt.Println("You are already following user: ", user)
			}
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
		fmt.Println(err)
		time.Sleep(300 * time.Second)
		return
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
	search, _, err := twitterClient.Search.Tweets(&twitter.SearchTweetParams{
		Query: look4,
	})
	if err != nil {
		panic(err)
	}
	for _, tweet := range search.Statuses {
		tweet_id := tweet.ID
		text := tweet.Text
		if strings.Contains(text, look4) {
			twitterClient.Statuses.Retweet(tweet_id, &twitter.StatusRetweetParams{})
		}
		break
	}
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
