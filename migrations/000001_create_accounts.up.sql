CREATE TABLE accounts (
    id SERIAL PRIMARY KEY,
    document_number VARCHAR(20) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
