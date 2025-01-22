-- +goose Up
INSERT INTO clients (id, name, cnpj, role, created_at, updated_at) VALUES ('78bc62fa-949b-485c-b8fd-a786417c8bd0', 'ContaFacil', '00000000000000', 'ADMIN', NOW(), NOW());
