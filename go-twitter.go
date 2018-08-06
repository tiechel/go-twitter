package main

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func main() {
	config := oauth1.NewConfig(
		"consumerKey",
		"consumerSecret",
	)
	token := oauth1.NewToken(
		"accessToken",
		"accessSecret",
	)

	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)

	params := &twitter.StreamUserParams{
		With:          "followings",
		StallWarnings: twitter.Bool(true),
	}

	stream, err := client.Streams.User(params)

	if err != nil {
		log.Fatal("Failed to start streaming.")
		return
	}

	for message := range stream.Messages {
		switch message := message.(type) {
		case *twitter.FriendsList:
			go funcFriendsList(message)
		case *twitter.Tweet:
			go printTweet(message)
		case *twitter.StatusDeletion:
			go funcStatusDeletion(message)
		case *twitter.Event:
			go funcEvent(message)
		default:
			log.Println("unknown type: ", reflect.TypeOf(message))
		}
	}

}

func funcFriendsList(friendsList *twitter.FriendsList) {
	// do something
	fmt.Println("FriendsList")
}

func printTweet(tweet *twitter.Tweet) {
	if tweet.RetweetedStatus != nil {
		tweet = tweet.RetweetedStatus
	}
	createdAt, err := tweet.CreatedAtTime()
	if err != nil {
		log.Println("tweet.CreatedAtTime() retruns error.")
	}

	// タイムゾーン対応
	createdAt = createdAt.In(time.FixedZone("Asia/Tokyo", 9*60*60))
	fmt.Printf("%s % 6dRT % 6dFav %s\n", createdAt.Format("2006-01-02 15:04:05"), tweet.RetweetCount, tweet.FavoriteCount, tweet.Text)
}

func funcEvent(event *twitter.Event) {
	// do something
	fmt.Println("Event")
}

func funcStatusDeletion(deletion *twitter.StatusDeletion) {
	// do something
	fmt.Println("StatusDeletion")
}
