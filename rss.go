package main

import (
	"encoding/xml"
	"net/http"
	"time"
)

type RssFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Language    string    `xml:"language"`
		Item        []RssItem `xml:"item"`
	} `xml:"channel"`
}

type RssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func parseXmlToRssFeed(url string) (RssFeed, error) {
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	res, err := httpClient.Get(url)

	if err != nil {
		return RssFeed{}, err
	}

	decoder := xml.NewDecoder(res.Body)
	defer res.Body.Close()
	rssFeed := RssFeed{}
	err = decoder.Decode(&rssFeed)

	if err != nil {
		return RssFeed{}, err
	}

	return rssFeed, nil
}
