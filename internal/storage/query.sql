-- name: GetUserByEmail :one
SELECT id, email, password_hash, created_at, updated_at
FROM users
WHERE email = $1;

-- name: UserExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE email = $1);

-- name: CreateUser :exec
INSERT INTO users (id, email, password_hash, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5);

-- name: GetProxy :one
SELECT id, listen_url, mode, condition, path_key, created_at, updated_at
FROM proxies
WHERE id = $1;

-- name: GetProxies :many
SELECT id, listen_url, mode, condition, tags, path_key
FROM proxies
ORDER BY created_at DESC;

-- name: UpdateProxyCondition :exec
UPDATE proxies
SET condition  = $1,
    updated_at = $2
WHERE id = $3;

-- name: UpdateProxyTags :exec
UPDATE proxies
SET tags       = $1,
    updated_at = NOW()
WHERE id = $2;

-- name: GetAllTags :many
SELECT DISTINCT UNNEST(tags)::text as tags
FROM proxies
WHERE tags IS NOT NULL
ORDER BY 1;

-- name: GetProxyTags :one
SELECT tags
FROM proxies
WHERE id = $1;

-- name: GetProxiesByTags :many
SELECT DISTINCT p.id,
                p.listen_url,
                p.mode,
                p.condition,
                p.tags,
                p.path_key,
                p.created_at,
                p.updated_at
FROM proxies p
WHERE tags @> $1
ORDER BY p.created_at DESC;

-- name: GetTargetsByProxyID :many
SELECT id, url, weight, is_active
FROM targets
WHERE proxy_id = $1;

-- name: CreateProxy :exec
INSERT INTO proxies (id, listen_url, mode, path_key, condition, tags, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: CreateTarget :exec
INSERT INTO targets (id, proxy_id, url, weight, is_active)
VALUES ($1, $2, $3, $4, $5);

-- name: DeleteTargetByProxyID :exec
DELETE
FROM targets
WHERE proxy_id = $1;

-- name: CreateProxyChange :exec
INSERT INTO proxy_changes (id, proxy_id, change_type, previous_state, new_state, created_at, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: GetProxyChangesByProxyID :many
SELECT id, proxy_id, change_type, previous_state, new_state, created_at, created_by
FROM proxy_changes
WHERE proxy_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CreateVisit :exec
INSERT INTO visits (id, proxy_id, target_id, user_id, rid, rrid, ruid, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: GetStats :one
SELECT COALESCE(SUM(request_count), 0)::int as requests,
       COALESCE(SUM(error_count), 0)::int   as errors
FROM proxy_stats
WHERE timestamp BETWEEN to_timestamp(@from_time::text, 'YYYY-MM-DD HH24:MI:SS.MS')
          AND to_timestamp(@to_time::text, 'YYYY-MM-DD HH24:MI:SS.MS');

-- name: GetUniqueUsersCount :one
SELECT COUNT(DISTINCT users)
FROM proxy_stats,
     LATERAL jsonb_array_elements_text(unique_users) AS users
WHERE timestamp BETWEEN to_timestamp(@from_time::text, 'YYYY-MM-DD HH24:MI:SS.MS')
          AND to_timestamp(@to_time::text, 'YYYY-MM-DD HH24:MI:SS.MS');

-- name: GetTargetStats :many
SELECT target_id,
       timestamp,
       request_count as requests,
       error_count   as errors,
       jsonb_array_length(unique_users) as users_count
FROM proxy_stats
WHERE proxy_id = $1
  AND timestamp BETWEEN to_timestamp(@from_time::text, 'YYYY-MM-DD HH24:MI:SS.MS')
    AND to_timestamp(@to_time::text, 'YYYY-MM-DD HH24:MI:SS.MS');
