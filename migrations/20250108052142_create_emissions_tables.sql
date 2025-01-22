-- +goose Up
CREATE TABLE emissions (
    id UUID PRIMARY KEY,
    emission_type VARCHAR(50) NOT NULL,
    client_id UUID NOT NULL,
    message TEXT,
    status VARCHAR(50),
    user_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE gnre_emissions (
    id UUID PRIMARY KEY REFERENCES emissions(id) ON DELETE CASCADE,
    xml UUID NOT NULL,
    pdf UUID,
    comprovante_pdf UUID,
    guia_amount NUMERIC(15, 2) NOT NULL,
    numero_recibo VARCHAR(255),
    chave_nota VARCHAR(255) NOT NULL,
    cod_barras_guia VARCHAR(255),
    num_nota VARCHAR(255) NOT NULL,
    destinatario VARCHAR(255) NOT NULL
);
