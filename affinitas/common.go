package affinitas

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

var AffinitasQueueSize = 100

func RequestPage(crawlData Crawl, pageURL string) *http.Response {
	scrapeDoQueryVals := url.Values{}
	scrapeDoQueryVals.Add("token", "YOUR_TOKEN")
	scrapeDoQueryVals.Add("url", pageURL)
	scrapeDoQueryVals.Add("geoCode", selectGeoCode(crawlData))

	scrapeDoURI := "http://api.scrape.do?" + scrapeDoQueryVals.Encode()

	res, err := http.Get(scrapeDoURI)
	if err != nil {
		fmt.Println(err.Error())
	}
	return res
}

func selectGeoCode(crawlData Crawl) string {
	switch amazonLocale := crawlData.Sites_to_crawl[0]; amazonLocale {
	case "www.amazon.co.uk":
		return "gb"
	case "www.amazon.ca":
		return "us"
	default:
		return "us"
	}
}

func regularizeFileName(productTitle string) string {
	fileName := strings.ReplaceAll(productTitle, ",", "-")
	fileName = strings.Trim(fileName, " ")
	fileName = strings.ReplaceAll(fileName, "/", "-")
	fileName = strings.ReplaceAll(fileName, "\\", "-")
	fileName = strings.ReplaceAll(fileName, ".", "-")
	fileName = strings.ReplaceAll(fileName, " ", "-")
	return fileName
}
