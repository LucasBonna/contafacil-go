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

CREATE TABLE files (
    id UUID PRIMARY KEY,
    file_name VARCHAR(255) NOT NULL,
    extension VARCHAR(50) NOT NULL,
    content_type VARCHAR(100) NOT NULL,
    file_path VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE users (
    id UUID PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    api_key VARCHAR(255) UNIQUE NOT NULL,
    role VARCHAR(50),
    client_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    CONSTRAINT fk_client FOREIGN KEY (client_id) REFERENCES clients(id) ON DELETE CASCADE,
    CONSTRAINT unique_client_username UNIQUE (client_id, username)
);

CREATE TABLE clients (
    id UUID PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    cnpj VARCHAR(20) UNIQUE NOT NULL,
    role VARCHAR(50),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

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
    guia_amount double precision NOT NULL,
    numero_recibo VARCHAR(255),
    chave_nota VARCHAR(255) NOT NULL,
    cod_barras_guia VARCHAR(255),
    num_nota VARCHAR(255) NOT NULL,
    destinatario VARCHAR(255) NOT NULL
);
