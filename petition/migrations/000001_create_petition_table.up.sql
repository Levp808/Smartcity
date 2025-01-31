CREATE TYPE petition_status AS ENUM ('created', 'moderated', 'in progress', 'done');
CREATE TYPE department_service AS ENUM ('Коммунальная служба', 'Дорожная служба');

CREATE TABLE IF NOT EXISTS petitions (
    petition_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    location TEXT NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    status petition_status NOT NULL,
    department department_service,
    done_at TIMESTAMP DEFAULT NULL,
    report_id UUID,
    content_job TEXT DEFAULT NULL
);
