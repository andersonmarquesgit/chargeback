CREATE TABLE IF NOT EXISTS batch_files (
    file_id TEXT PRIMARY KEY,
    file_url TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    record_count INT NOT NULL,
    status TEXT NOT NULL, -- ready | sent | failed
    sent_at TIMESTAMP,
    retry_count INT DEFAULT 0,
    last_attempt_at TIMESTAMP
);