-- name: GetEmail :one
SELECT c.email FROM candidates c WHERE c.email = $1
UNION
SELECT i.email FROM interviewers i WHERE i.email = $1
UNION
SELECT co.email FROM companies co WHERE co.email = $1;

-- name: CreateCandidate :one
INSERT INTO candidates (id, name, email, password_hash, salt)
VALUES (gen_random_uuid(), $1, $2, $3, $4)
RETURNING *;