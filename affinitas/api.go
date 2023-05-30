package affinitas

import (
	"fmt"
)

var searchPageCount = 300
var amazonUSPrototypeURL = "https://www.amazon.com/s?k=%s&page=%d&rh=p_85%%3A2470955011&dc"
var amazonUKPrototypeURL = "https://www.amazon.co.uk/s?k=%s&page=%d&rh=p_85%%3A2470955011&dc"
var amazonCAPrototypeURL = "https://www.amazon.ca/s?k=%s&page=%d&rh=p_85%%3A2470955011&dc"

func InitiateCrawl(crawlMessage Initiate_Crawl_Event) {
	searchPageCrawlEvents := generateSearchPageCrawlEvents(crawlMessage)
	for _, searchPageCrawlEvent := range searchPageCrawlEvents {
		EnqueueAmazonSearchCrawlEvent(searchPageCrawlEvent)
	}
}

func generateSearchPageCrawlEvents(crawlMessage Initiate_Crawl_Event) []Crawl_Search_Page_Event {
	var crawlEvents []Crawl_Search_Page_Event
	for _, searchTerm := range crawlMessage.Search_terms {
		for searchPageIndex := 1; searchPageIndex <= searchPageCount; searchPageIndex++ {
			formattedLink := formSearchLink(crawlMessage, searchTerm, searchPageIndex)
			crawlMessage := Crawl_Search_Page_Event{
				Event_metadata:   crawlMessage.Event_metadata,
				Search_page_link: formattedLink,
			}
			crawlEvents = append(crawlEvents, crawlMessage)
		}
	}
	return crawlEvents
}

func formSearchLink(crawlMessage Initiate_Crawl_Event, searchTerm string, searchPageIndex int) string {
	switch siteToCrawl := crawlMessage.Event_metadata.Crawl.Sites_to_crawl[0]; siteToCrawl {
	case "www.amazon.co.uk":
		return fmt.Sprintf(amazonUKPrototypeURL, searchTerm, searchPageIndex)
	case "www.amazon.ca":
		return fmt.Sprintf(amazonCAPrototypeURL, searchTerm, searchPageIndex)
	default:
		return fmt.Sprintf(amazonUSPrototypeURL, searchTerm, searchPageIndex)
	}
}
