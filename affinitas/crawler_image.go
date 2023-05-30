package affinitas

import (
	"io/ioutil"
	"log"

	"github.com/google/uuid"
)

var imageCrawlChannel = make(chan Crawl_Image_Event, AffinitasQueueSize)

func EnqueueImageCrawlEvent(imageCrawlEvent Crawl_Image_Event) {
	imageCrawlChannel <- imageCrawlEvent
}

func imageCrawlWorker(imageCrawlChannel <-chan Crawl_Image_Event) {
	for crawlEvent := range imageCrawlChannel {
		processImageCrawlEvent(&crawlEvent)
	}
}

func StartImageCrawlerWorkers(workerCount int) {
	for i := 0; i < workerCount; i++ {
		go imageCrawlWorker(imageCrawlChannel)
	}
}

func processImageCrawlEvent(crawlEvent *Crawl_Image_Event) {
	imageResponse := RequestPage(crawlEvent.Event_metadata.Crawl, crawlEvent.Primary_image_url)

	imageBytes, err := ioutil.ReadAll(imageResponse.Body)
	defer imageResponse.Body.Close()
	if err != nil {
		log.Println("Unable to read image data from image response")
		EnqueueImageCrawlEvent(*crawlEvent)
	}

	imageData := Product_Image{
		ID:              uuid.New().String(),
		Product_id:      crawlEvent.Product.ID,
		Image_url:       crawlEvent.Primary_image_url,
		Image_file_name: regularizeFileName(crawlEvent.Product.Product_title) + ".jpg",
		Image_byte_arr:  imageBytes,
	}

	enqueueImageData(crawlEvent, imageData)
}

func enqueueImageData(crawlEvent *Crawl_Image_Event, imageData Product_Image) {
	Log("Crawled Image: " + imageData.Image_file_name)
	persistImageEvent := Perist_Image_Event{
		Event_metadata: crawlEvent.Event_metadata,
		Product_image:  imageData,
	}
	EnqueuePersistImageEvent(persistImageEvent)
}
