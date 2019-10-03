package main

import (
	"log"

	"github.com/go-shiori/warc"
)

func main() {
	// Define variables
	url := "https://apnews.com/6e151296fb194f85ba69a8babd972e4b"
	userAgent := "Shiori/2.0.0 (+https://github.com/go-shiori/shiori)"

	// Ceate archival request
	req := warc.ArchivalRequest{
		URL:        url,
		UserAgent:  userAgent,
		LogEnabled: true,
	}

	// Start archival
	err := warc.NewArchive(req, "ap-news")
	if err != nil {
		log.Fatalln(err)
	}
}
