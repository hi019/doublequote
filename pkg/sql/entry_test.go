package sql

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"doublequote/pkg/domain"
)

func TestAA(t *testing.T) {
	s := NewSQL("file:ent?mode=memory&cache=shared&_fk=1")
	s.Open()

	svc := NewEntryService(s)
	fsvc := NewFeedService(s)

	_, err := fsvc.CreateFeed(context.Background(), &domain.Feed{
		Name:   "Test",
		RssURL: "TEst",
		Domain: "Test",
	})
	if err != nil {
		log.Fatalln(err)
	}

	//ID = {int} 0
	//Title = {string} "How to install the iOS 15.4 public beta"
	//URL = {string} "https://www.theverge.com/22906798/ios-iphone-15-4-face-id-mask-test-beta"
	//Author = {string} "Cameron Faulkner"
	//ContentKey = {string} ""
	//Feed = {dq.Feed}
	//FeedID = {int} 2
	//CreatedAt = {time.Time} 0001-01-01 00:00:00 +0000
	//UpdatedAt = {time.Time} 0001-01-01 00:00:00 +0000

	_, err = svc.CreateEntry(context.Background(), domain.Entry{
		Title:      "How to install the iOS 15.4 public beta",
		URL:        "https://www.theverge.com/22906798/ios-iphone-15-4-face-id-mask-test-beta",
		Author:     "Cameron Faulkner",
		ContentKey: "",
		FeedID:     1,
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
	})
	fmt.Println("ERR: ", err)

	if err != nil {
		log.Fatalln(err)
	}
}
