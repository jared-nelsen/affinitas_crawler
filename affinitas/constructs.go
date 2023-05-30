package affinitas

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Site struct {
	gorm.Model
	ID                   string
	Site_title           string
	Site_base_url        string
	Admin_user           string
	Admin_api_password   string
	Author_count         int
	Categories           pq.StringArray `gorm:"type:string[]"`
	Site_category_titles pq.StringArray `gorm:"type:string[]"`
	Affiliate_id         string
}

type Crawl struct {
	gorm.Model
	ID             string
	Site_id        string
	Website        string
	Sites_to_crawl pq.StringArray `gorm:"type:string[]"`
}

type Crawl_Event struct {
	gorm.Model
	ID           string
	Crawl_id     string
	Category     string
	Search_terms pq.StringArray `gorm:"type:string[]"`
}

type Product struct {
	gorm.Model
	ID                     string
	Crawl_id               string
	Product_title          string
	Origin_site            string
	Search_page_link       string
	Unique_identifier      string
	Product_url            string
	Url_with_affiliate_tag string
	Product_category       string
}

type Product_Image struct {
	gorm.Model
	ID              string
	Product_id      string
	Image_url       string
	Image_file_name string
	Image_byte_arr  []byte `gorm:"type:byte[]"` // Will this go into the database?
}

type Blog struct {
	gorm.Model
	ID            string
	Product_id    string
	Blog_title    string
	Blog_string   string
	Blog_url      string
	Uploaded      bool
	Uploaded_date time.Time
	Prompts       pq.StringArray `gorm:"type:string[]"`
}
