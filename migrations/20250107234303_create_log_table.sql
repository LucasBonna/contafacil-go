-- +goose Up
CREATE TABLE access_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ip VARCHAR(255) NOT NULL,
    method VARCHAR(10) NOT NULL,
    endpoint VARCHAR(255),
    request_body TEXT,
    request_headers TEXT,
    request_query TEXT,
    request_params TEXT,
    response_body TEXT,
    response_headers TEXT,
    response_time VARCHAR(50),
    status_code INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
