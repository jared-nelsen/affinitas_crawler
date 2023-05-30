
CREATE TABLE crawl_events (
    ID uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    -- FK: Domain
    crawl_id uuid,
    -- Defaults
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
	deleted_at TIMESTAMP WITH TIME ZONE,
    -- Data
    category TEXT NOT NULL,
    search_terms TEXT [] NOT NULL
);