package affinitas

var productPersistorChannel = make(chan Perist_Product_Event, AffinitasQueueSize)

func EnqueuePersistProductEvent(persistProductEvent Perist_Product_Event) {
	productPersistorChannel <- persistProductEvent
}

func persistProductWorker(persistProductChannel <-chan Perist_Product_Event) {
	for persistEvent := range persistProductChannel {
		processPersistProductEvent(&persistEvent)
	}
}

func StartProductPersistWorkers(workerCount int) {
	for i := 0; i < workerCount; i++ {
		go persistProductWorker(productPersistorChannel)
	}
}

func processPersistProductEvent(persistEvent *Perist_Product_Event) {
	insertProductRecord(&persistEvent.Product)
	Log("Persisted Product: " + persistEvent.Product.Product_title)
}
