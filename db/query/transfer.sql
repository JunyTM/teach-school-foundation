-- name: CreateTranfer :one
INSERT INTO transfers (
    from_account_id,
    to_account_id,
    amount,
    currency
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: DeleteTranfer :exec
DELETE FROM transfers WHERE id = $1;

-- name: UpdateTranfer :one
UPDATE transfers SET amount = $2 WHERE id = $1 RETURNING *;

-- name: GetTranfer :one
SELECT * FROM transfers WHERE id = $1 LIMIT 1;

-- name: ListTranfer :many
SELECT * FROM transfers ORDER BY id LIMIT $1 OFFSET $2;
