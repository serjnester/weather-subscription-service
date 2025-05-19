-- name: CreateSubscription :one
INSERT INTO subscriptions (email, city, frequency, token)
VALUES ($1, $2, $3, $4)
    RETURNING *;

-- name: GetSubscriptionByToken :one
SELECT * FROM subscriptions
WHERE token = $1;

-- name: ConfirmSubscription :exec
UPDATE subscriptions
SET confirmed = true
WHERE token = $1;

-- name: Unsubscribe :exec
DELETE FROM subscriptions
WHERE token = $1;

-- name: IsAlreadySubscribed :one
SELECT COUNT(*) FROM subscriptions
WHERE email = $1 AND city = $2;