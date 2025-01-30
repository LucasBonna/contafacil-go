-- +goose Up
-- +goose StatementBegin
ALTER TABLE gnre_emissions
    ALTER COLUMN guia_amount TYPE double precision;
-- +goose StatementEnd

