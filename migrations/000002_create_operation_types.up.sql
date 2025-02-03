CREATE TABLE operation_types (
    id SERIAL PRIMARY KEY,
    description VARCHAR(50) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);


INSERT INTO operation_types (id, description) VALUES
(1, 'COMPRA A VISTA'),
(2, 'COMPRA PARCELADA'),
(3, 'SAQUE'),
(4, 'PAGAMENTO');
