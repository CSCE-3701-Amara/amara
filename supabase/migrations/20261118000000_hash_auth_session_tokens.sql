-- +goose Up
-- Sessions are looked up by the SHA-256 hex digest of the bearer token, so a
-- copy of this table cannot be replayed as cookies. Convert existing plaintext
-- tokens in place so signed-in users stay signed in. Issued tokens are 43-char
-- base64url strings and can never look like a 64-char hex digest, which keeps
-- this statement safe to run more than once.
UPDATE auth_sessions
SET token = encode(sha256(convert_to(token, 'UTF8')), 'hex')
WHERE token !~ '^[0-9a-f]{64}$';

-- +goose Down
-- Manual rollback reference only; Supabase executes the entire migration file.
/*
Digests cannot be reversed; roll back by deleting sessions so users sign in again:
DELETE FROM auth_sessions;
*/
