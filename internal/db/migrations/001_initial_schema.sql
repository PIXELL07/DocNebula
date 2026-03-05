-- JOBS TABLE
-- One row per uploaded processing job

CREATE TABLE IF NOT EXISTS jobs (
    id TEXT PRIMARY KEY,
    status TEXT NOT NULL CHECK (
        status IN ('UPLOADED','RUNNING','COMPLETED','FAILED')
    ),
    retry_count INT NOT NULL DEFAULT 0,
    idempotency_key TEXT UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_jobs_status
ON jobs(status);


-- FILES TABLE
-- One row per document extracted from ZIP


CREATE TABLE IF NOT EXISTS files (
    id TEXT PRIMARY KEY,
    job_id TEXT NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    path TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'PENDING' CHECK (
        status IN ('PENDING','PROCESSING','DONE','FAILED')
    ),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    -- Prevent duplicate files in same job
    UNIQUE(job_id, path)
);

CREATE INDEX IF NOT EXISTS idx_files_job_id
ON files(job_id);

CREATE INDEX IF NOT EXISTS idx_files_status
ON files(status);


-- PAGES TABLE
-- Per-page OCR checkpointing and text storage


CREATE TABLE IF NOT EXISTS pages (
    id TEXT PRIMARY KEY,
    file_id TEXT NOT NULL REFERENCES files(id) ON DELETE CASCADE,
    page_num INT NOT NULL,
    text TEXT,
    done BOOLEAN NOT NULL DEFAULT FALSE,
    processing_time_ms INT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    -- Prevent duplicate pages
    UNIQUE(file_id, page_num)
);

CREATE INDEX IF NOT EXISTS idx_pages_file_id
ON pages(file_id);

CREATE INDEX IF NOT EXISTS idx_pages_done
ON pages(done);

CREATE INDEX IF NOT EXISTS idx_pages_file_page
ON pages(file_id, page_num);


CREATE INDEX IF NOT EXISTS idx_pages_text_search
ON pages USING GIN(to_tsvector('english', text));


-- AUTO UPDATE updated_at TRIGGER

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


-- FILES TRIGGER


DROP TRIGGER IF EXISTS trg_files_updated_at ON files;

CREATE TRIGGER trg_files_updated_at
BEFORE UPDATE ON files
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();



-- PAGES TRIGGER

DROP TRIGGER IF EXISTS trg_pages_updated_at ON pages;

CREATE TRIGGER trg_pages_updated_at
BEFORE UPDATE ON pages
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();


DROP TRIGGER IF EXISTS trg_jobs_updated_at ON jobs;

CREATE TRIGGER trg_jobs_updated_at
BEFORE UPDATE ON jobs
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();