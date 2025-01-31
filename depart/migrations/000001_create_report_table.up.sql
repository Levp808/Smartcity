CREATE TABLE IF NOT EXISTS reports (
    petition_id UUID NOT NULL PRIMARY KEY,
    location TEXT NOT NULL,
    description_petition TEXT NOT NULL,
    report_id UUID DEFAULT gen_random_uuid(),
    content_job TEXT DEFAULT NULL,
    done_at TIMESTAMP DEFAULT NULL
);
