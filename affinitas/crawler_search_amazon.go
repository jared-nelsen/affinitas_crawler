package affinitas

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	mapset "github.com/deckarep/golang-set/v2"
)

var amazonSearchCrawlChannel = make(chan Crawl_Search_Page_Event, AffinitasQueueSize)

func EnqueueAmazonSearchCrawlEvent(searchCrawlEvent Crawl_Search_Page_Event) {
	amazonSearchCrawlChannel <- searchCrawlEvent
}

func amazonSearchCrawlWorker(searchPageChannel <-chan Crawl_Search_Page_Event) {
	for crawlEvent := range searchPageChannel {
		processSearchPageLink(&crawlEvent)
	}
}

func StartAmazonSearchWorkers(workerCount int) {
	for i := 0; i < workerCount; i++ {
		go amazonSearchCrawlWorker(amazonSearchCrawlChannel)
	}
}

func processSearchPageLink(crawlEvent *Crawl_Search_Page_Event) {
	linksOnPage := getAllProductLinksOnPage(crawlEvent)
	enqueueProductLinks(crawlEvent, linksOnPage)
}

func enqueueProductLinks(crawlEvent *Crawl_Search_Page_Event, links []string) {
	Log("Crawled Amazon Search Link: " + crawlEvent.Search_page_link)
	for _, link := range links {
		productCrawlEvent := Crawl_Product_Page_Event{
			Event_metadata:   crawlEvent.Event_metadata,
			Search_page_link: crawlEvent.Search_page_link,
			Product_url:      link,
		}
		EnqueueAmazonProductCrawlEvent(productCrawlEvent)
	}
}

func getAllProductLinksOnPage(crawlEvent *Crawl_Search_Page_Event) []string {
	response := RequestPage(crawlEvent.Event_metadata.Crawl, crawlEvent.Search_page_link)
	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		Log("Error loading HTTP response body... Retrying!")
		EnqueueAmazonSearchCrawlEvent(*crawlEvent)
		return []string{}
	}

	switch amazonLocale := crawlEvent.Event_metadata.Crawl.Sites_to_crawl[0]; amazonLocale {
	case "www.amazon.co.uk":
		return getAllLinksOnAmazonUKPage(crawlEvent, document)
	case "www.amazon.ca":
		return getAllLinksOnAmazonCAPage(crawlEvent, document)
	default:
		return getAllLinksOnAmazonUSPage(crawlEvent, document)
	}
}

func getAllLinksOnAmazonUKPage(crawlEvent *Crawl_Search_Page_Event, document *goquery.Document) []string {
	linkSet := mapset.NewSet[string]()
	document.Find(".s-underline-link-text").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists && strings.Contains(href, "%2Fdp%2F") {
			startOfAsinIndex := strings.Index(href, "%2Fdp%2F") + 8
			endOfAsinIndex := startOfAsinIndex + 10
			asin := href[startOfAsinIndex:endOfAsinIndex]
			nonRelativeAmazonLink := "https://" + crawlEvent.Event_metadata.Crawl.Sites_to_crawl[0] + "/dp/" + asin
			linkSet.Add(nonRelativeAmazonLink)
		}
	})
	return linkSet.ToSlice()
}

func getAllLinksOnAmazonCAPage(crawlEvent *Crawl_Search_Page_Event, document *goquery.Document) []string {
	linkSet := mapset.NewSet[string]()
	document.Find(".s-underline-link-text").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists && strings.Contains(href, "%2Fdp%2F") {
			startOfAsinIndex := strings.Index(href, "%2Fdp%2F") + 8
			endOfAsinIndex := startOfAsinIndex + 10
			asin := href[startOfAsinIndex:endOfAsinIndex]
			nonRelativeAmazonLink := "https://" + crawlEvent.Event_metadata.Crawl.Sites_to_crawl[0] + "/dp/" + asin
			linkSet.Add(nonRelativeAmazonLink)
		}
	})
	return linkSet.ToSlice()
}

func getAllLinksOnAmazonUSPage(crawlEvent *Crawl_Search_Page_Event, document *goquery.Document) []string {
	linkSet := mapset.NewSet[string]()
	document.Find(".s-underline-link-text").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists && strings.Contains(href, "/dp/B0") {
			startOfAsinIdent := strings.Index(href, "/dp/B0")
			endOfAsin := startOfAsinIdent + 14
			relativeProductLink := href[:endOfAsin]
			nonRelativeAmazonLink := "https://" + crawlEvent.Event_metadata.Crawl.Sites_to_crawl[0] + relativeProductLink
			linkSet.Add(nonRelativeAmazonLink)
		}
	})
	return linkSet.ToSlice()
}
