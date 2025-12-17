-- +goose Up
-- +goose StatementBegin
-- Расширение для генерации UUID
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Таблица с "сырым" паспортом (JSON как есть)
CREATE TABLE passports_raw (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    snapshot_hash TEXT NOT NULL,

    passport JSONB NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Таблица с "вытащенными" (сериализованными) данными
CREATE TABLE passports_serialized (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    passport_id UUID NOT NULL,

    cluster_id TEXT NOT NULL,
    namespace TEXT NOT NULL,
    resource_type TEXT NOT NULL,
    resource_name TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT fk_passport
        FOREIGN KEY (passport_id)
        REFERENCES passports_raw(id)
        ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS passports_serialized;
DROP TABLE IF EXISTS passports_raw;
-- +goose StatementEnd