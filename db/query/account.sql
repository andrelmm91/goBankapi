-- name: CreateAccount :one
INSERT INTO account (
  owner,
  balance,
  currency
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetAccount :one
SELECT * FROM account
WHERE id = $1 LIMIT 1;

-- name: ListAccounts :many
SELECT * FROM account
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateAccount :one
UPDATE account
SET balance = $2
WHERE id = $1
RETURNING *;


-- name: DeleteAccount :exec
DELETE FROM account
WHERE id = $1;