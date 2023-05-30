package affinitas

type Event_Metadata struct {
	Site     Site
	Crawl    Crawl
	Category string
}

type Initiate_Crawl_Event struct {
	Event_metadata Event_Metadata
	Domain_id      string
	Search_terms   []string
	Sites_to_crawl []string
}

type Crawl_Search_Page_Event struct {
	Event_metadata   Event_Metadata
	Search_page_link string
}

type Crawl_Product_Page_Event struct {
	Event_metadata   Event_Metadata
	Search_page_link string
	Product_url      string
}

type Crawl_Image_Event struct {
	Event_metadata    Event_Metadata
	Product           Product
	Primary_image_url string
}

type Perist_Product_Event struct {
	Event_metadata Event_Metadata
	Product        Product
}

type Perist_Image_Event struct {
	Event_metadata Event_Metadata
	Product_image  Product_Image
}

type Generate_Blog_Event struct {
	Event_metadata Event_Metadata
	Product        Product
	Product_image  Product_Image
}

type Upload_Blog_Event struct {
	Event_metadata Event_Metadata
	Blog           Blog
}
