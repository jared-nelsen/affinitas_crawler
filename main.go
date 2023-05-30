package main

import (
	"cabo_affinitas/affinitas"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

func main() {
	affinitas.InitDB()
	affinitas.LoadPreviouslyCrawledProducts()
	startWorkers()

	fmt.Println("Choose a site by its index to crawl for:")
	sites := affinitas.RetrieveSiteList()
	for i, site := range sites {
		fmt.Printf("%d: %s\n", i, site.Site_title)
	}
	fmt.Println()
	var choiceStr string
	fmt.Scanf("%s", &choiceStr)
	choiceIdx, _ := strconv.Atoi(choiceStr)
	chosenSite := sites[choiceIdx]

	crawl := affinitas.Crawl{
		ID:             uuid.New().String(),
		Site_id:        chosenSite.ID,
		Website:        chosenSite.Site_title,
		Sites_to_crawl: pq.StringArray{"www.amazon.com"},
	}
	affinitas.InsertCrawlRecord(&crawl)

	initCoatsCategoryCrawl(chosenSite, crawl)
}

func initCoatsCategoryCrawl(site affinitas.Site, crawl affinitas.Crawl) {
	initiateCrawlEvent := affinitas.Initiate_Crawl_Event{
		Event_metadata: affinitas.Event_Metadata{
			Site:     site,
			Crawl:    crawl,
			Category: "coats",
		},
		Search_terms: []string{
			"coats",
		},
	}
	affinitas.InsertCrawlEvent(&initiateCrawlEvent)
	affinitas.InitiateCrawl(initiateCrawlEvent)
}

func startWorkers() {
	affinitas.StartAmazonSearchWorkers(5)
	affinitas.StartAmazonProductWorkers(5)
	affinitas.StartImageCrawlerWorkers(4)
	affinitas.StartImagePersistWorkers(8)
	affinitas.StartProductPersistWorkers(8)
}
