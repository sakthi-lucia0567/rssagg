package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Language    string    `xml:"language"`
		Item        []RSSItem `xml:"item"`
	}
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func urlToFeed(url string) (RSSFeed, error) {

	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := httpClient.Get(url)

	// fmt.Println("there is an error here", resp)

	if err != nil {
		return RSSFeed{}, err
	}
	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return RSSFeed{}, err
	}
	// fmt.Println("the data from the feed", dat)
	rssFeed := RSSFeed{}
	err = xml.Unmarshal(dat, &rssFeed)
	fmt.Println("the error here", err)
	if err != nil {
		return RSSFeed{}, err
	}
	// fmt.Println("this is my rss feed", rssFeed)
	return rssFeed, nil
}
