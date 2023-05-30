package affinitas

var imagePersistorChannel = make(chan Perist_Image_Event, AffinitasQueueSize)

func EnqueuePersistImageEvent(persistImageEvent Perist_Image_Event) {
	imagePersistorChannel <- persistImageEvent
}

func persistImageWorker(persistImageChannel <-chan Perist_Image_Event) {
	for persistEvent := range persistImageChannel {
		processPersistImageEvent(&persistEvent)
	}
}

func StartImagePersistWorkers(workerCount int) {
	for i := 0; i < workerCount; i++ {
		go persistImageWorker(imagePersistorChannel)
	}
}

func processPersistImageEvent(persistEvent *Perist_Image_Event) {
	insertProductImage(&persistEvent.Product_image)
	Log("Persisted Image: " + persistEvent.Product_image.Image_file_name)
}
