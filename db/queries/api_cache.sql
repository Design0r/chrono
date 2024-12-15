-- name: CacheExists :one
SELECT EXISTS(
    SELECT 1 FROM api_cache
    WHERE year = ?
);

-- name: CreateCache :exec
INSERT INTO api_cache (year)
VALUES (?);
