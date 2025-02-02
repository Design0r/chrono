-- name: CacheExists :one
SELECT EXISTS(
    SELECT 1 FROM api_cache
    WHERE year = ?
);

-- name: CreateCache :exec
INSERT INTO api_cache (year)
VALUES (?);

-- name: GetApiCacheYears :many
SELECT year FROM api_cache
GROUP BY year;
