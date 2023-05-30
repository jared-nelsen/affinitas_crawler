package affinitas

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/google/uuid"
)

var amazonProductCrawlChannel = make(chan Crawl_Product_Page_Event, AffinitasQueueSize)
var crawledProductLinksSet = mapset.NewSet[string]()

func EnqueueAmazonProductCrawlEvent(productCrawlEvent Crawl_Product_Page_Event) {
	amazonProductCrawlChannel <- productCrawlEvent
}

func amazonProductCrawlWorker(productPageChannel <-chan Crawl_Product_Page_Event) {
	for crawlEvent := range productPageChannel {
		processProductPage(&crawlEvent)
	}
}

func StartAmazonProductWorkers(workerCount int) {
	for i := 0; i < workerCount; i++ {
		go amazonProductCrawlWorker(amazonProductCrawlChannel)
	}
}

func LoadPreviouslyCrawledProducts() {
	previouslyCrawledProductLinks := RetrievePreviouslyCrawledProductLinks()
	for _, link := range previouslyCrawledProductLinks {
		crawledProductLinksSet.Add(link)
	}
}

func processProductPage(crawlEvent *Crawl_Product_Page_Event) {

	if crawledProductLinksSet.Contains(crawlEvent.Product_url) {
		Log("Already crawled product at: " + crawlEvent.Product_url + "... Skipping!")
		return
	}

	response := RequestPage(crawlEvent.Event_metadata.Crawl, crawlEvent.Product_url)
	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		Log("Error loading HTTP response body... Retrying!")
		EnqueueAmazonProductCrawlEvent(*crawlEvent)
		return
	}

	productData := Product{
		ID:                     uuid.New().String(),
		Crawl_id:               crawlEvent.Event_metadata.Crawl.ID,
		Product_title:          scrapeProductTitle(document),
		Origin_site:            crawlEvent.Event_metadata.Crawl.Sites_to_crawl[0],
		Search_page_link:       crawlEvent.Search_page_link,
		Unique_identifier:      parseASIN(crawlEvent.Product_url),
		Product_url:            crawlEvent.Product_url,
		Url_with_affiliate_tag: formURLWithAffiliateTag(crawlEvent.Event_metadata.Site, crawlEvent.Product_url),
		Product_category:       crawlEvent.Event_metadata.Category,
	}
	primaryImageURL := scrapePrimaryImageURL(document)

	crawledProductLinksSet.Add(crawlEvent.Product_url)

	if crawlDataIsInvalid(primaryImageURL, &productData) {
		Log("Invalid Product Data Crawled... Skipping Product!")
		return
	}

	enqueueProductCrawledEvent(crawlEvent, productData)
	enqueueImageCrawlEvent(crawlEvent, productData, primaryImageURL)
}

func crawlDataIsInvalid(primaryImageURL string, productData *Product) bool {
	if primaryImageURL == "" || productData.Product_title == "" || productData.Unique_identifier == "" || productData.Product_url == "" {
		return true
	}
	return false
}

func formURLWithAffiliateTag(site Site, productURL string) string {
	return productURL + "?tag=" + site.Affiliate_id
}

func enqueueProductCrawledEvent(crawlEvent *Crawl_Product_Page_Event, productData Product) {
	Log("Crawled Amazon Product: " + productData.Product_title)
	persistProductEvent := Perist_Product_Event{
		Event_metadata: crawlEvent.Event_metadata,
		Product:        productData,
	}
	EnqueuePersistProductEvent(persistProductEvent)
}

func enqueueImageCrawlEvent(crawlEvent *Crawl_Product_Page_Event, productData Product, imageURL string) {
	imageCrawlEvent := Crawl_Image_Event{
		Event_metadata:    crawlEvent.Event_metadata,
		Product:           productData,
		Primary_image_url: imageURL,
	}
	EnqueueImageCrawlEvent(imageCrawlEvent)
}

func scrapeProductTitle(document *goquery.Document) string {
	title := document.Find("#productTitle").First().Text()
	title = regularizeFileName(title)
	return title
}

func scrapePrimaryImageURL(document *goquery.Document) string {
	href, exists := document.Find("#landingImage").First().Attr("data-old-hires")
	if exists {
		return strings.Trim(href, " ")
	}
	return ""
}

func parseASIN(productURL string) string {
	indexOfAsinStart := strings.Index(productURL, "/dp/")
	indexOfAsinStart = indexOfAsinStart + 4
	return productURL[indexOfAsinStart:]
}
