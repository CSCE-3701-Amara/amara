-- +goose Up
CREATE TABLE candidates (
    id                  UUID        PRIMARY KEY,
    name                TEXT        NOT NULL,
    email               TEXT        NOT NULL,
    password_hash       BYTEA       NOT NULL,
    salt                BYTEA       NOT NULL,
    created_at          TIMESTAMP   NOT NULL,
    updated_at          TIMESTAMP,
    deleted_at          TIMESTAMP,
    deactivated_at      TIMESTAMP,
    last_login_at       TIMESTAMP,
    phone_number        TEXT,
    github_profile      TEXT,
    linkedin_profile    TEXT,
    portfolio_link      TEXT,
    profile_picture_url TEXT,
    gender              TEXT,
    birth_date          DATE,
    nationality         TEXT,
    address             TEXT,
    biography           TEXT,
    years_of_experience INTEGER,
    available           BOOLEAN,
    banned              BOOLEAN,
    ban_reason          TEXT
);

-- +goose Down
DROP TABLE candidates;