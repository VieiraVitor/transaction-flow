CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    account_id INT NOT NULL,
    operation_type_id INT NOT NULL,
    amount NUMERIC(15, 2) NOT NULL,
    event_date TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),

    FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE,
    FOREIGN KEY (operation_type_id) REFERENCES operation_types(id) ON DELETE RESTRICT
);
