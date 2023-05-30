
CREATE TABLE crawls (
    ID uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    -- FK to prod DB: Site ID
    site_id uuid NOT NULL,
    -- Defaults
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
	deleted_at TIMESTAMP WITH TIME ZONE,
    -- Data
    website TEXT NOT NULL,
    sites_to_crawl TEXT []
);