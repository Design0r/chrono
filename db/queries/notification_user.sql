-- name: CreateNotificationUser :exec
INSERT INTO notification_user (notification_id, user_id)
VALUES (?, ?);


-- name: GetUserNotifications :many
SELECT n.* from notification_user nu
JOIN notifications n ON nu.notification_id = n.id
WHERE nu.user_id = ?
AND n.viewed_at IS NULL;


-- name: ClearAllUserNotifications :exec
UPDATE notifications
SET viewed_at = CURRENT_TIMESTAMP
WHERE id IN (
    SELECT n.id
    FROM notifications n
    JOIN notification_user nu ON n.id = nu.notification_id
    WHERE nu.user_id = ?
    AND n.viewed_at IS NULL
);
