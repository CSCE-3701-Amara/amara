-- +goose Up
CREATE TABLE companies (
    id                  UUID    PRIMARY KEY,
    email               TEXT    NOT NULL,
    password_hash       BYTEA   NOT NULL,
    salt                BYTEA   NOT NULL,
    legal_name          TEXT    NOT NULL,
    display_name        TEXT,
    description         TEXT,
    website_link        TEXT,
    logo_url            TEXT,
    contact_phone       TEXT,
    contact_email       TEXT,
    linkedin_profile    TEXT,
    created_at          TIMESTAMP,
    updated_at          TIMESTAMP,
    deleted_at          TIMESTAMP,
    deactivated_at      TIMESTAMP,
    last_login_at       TIMESTAMP,
    banned              BOOLEAN,
    ban_reason          TEXT
);

-- +goose Down
DROP TABLE companies;