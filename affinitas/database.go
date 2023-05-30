package affinitas

import (
	"log"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var dbConnectionString = "YOUR_DATA_DB_URI"
var db *gorm.DB

func InitDB() {
	var err error
	db, err = gorm.Open(postgres.Open(dbConnectionString), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Panic(err)
	}
}

func RetrievePreviouslyCrawledProductLinks() []string {
	var products []Product
	db.Find(&products)
	links := []string{}
	for _, product := range products {
		links = append(links, product.Product_url)
	}
	return links
}

func getProdDBConnection() *gorm.DB {
	production_db_dsn := "YOUR_PROD_DB_URI"
	production_db, err := gorm.Open(postgres.Open(production_db_dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalln(err)
	}
	return production_db
}

func closeProdDBConn(db *gorm.DB) {
	DB, _ := db.DB()
	DB.Close()
}

func InsertCrawlRecord(crawl *Crawl) {
	db.Create(&crawl)
}

func InsertCrawlEvent(crawlEvent *Initiate_Crawl_Event) {
	event := Crawl_Event{
		ID:           uuid.New().String(),
		Crawl_id:     crawlEvent.Event_metadata.Crawl.ID,
		Category:     crawlEvent.Event_metadata.Category,
		Search_terms: crawlEvent.Search_terms,
	}
	db.Create(&event)
}

func RetrieveSiteList() []Site {
	db := getProdDBConnection()
	var sites []Site
	db.Find(&sites)
	closeProdDBConn(db)
	return sites
}

func insertProductRecord(productData *Product) error {
	result := db.Create(&productData)
	return result.Error
}

func insertProductImage(productImage *Product_Image) error {
	result := db.Create(&productImage)
	return result.Error
}
