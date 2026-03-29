-- name: GetCommentsBySlug :many
WITH RECURSIVE CommentHierarchy AS (
    SELECT id, name, content, slug, created_at, parent_id, 0 AS depth, printf('%020d', id) AS path
    FROM comments ch
    WHERE ch.parent_id IS NULL AND ch.slug = ?

    UNION ALL
    
    SELECT c.id, c.name, c.content, c.slug, c.created_at, c.parent_id, ch.depth + 1, ch.path || '.' || printf('%020d', c.id)
    FROM comments c
    JOIN CommentHierarchy ch ON c.parent_id = ch.id
) 
SELECT * FROM CommentHierarchy
ORDER BY path, created_at;
-- WITH RECURSIVE CommentHierarchy AS (
--     SELECT id, name, content, slug, created_at, parent_id, 0 AS depth, name as path
--     FROM comments ch
--     WHERE ch.parent_id IS NULL AND ch.slug = '/writing/2-types-of-work'
--     UNION ALL
    
--     SELECT c.id, c.name, c.content, c.slug, c.created_at, c.parent_id, ch.depth + 1, ch.path || ' > ' || c.name
--     FROM comments c
--     JOIN CommentHierarchy ch ON c.parent_id = ch.id
-- ) 
-- SELECT * FROM CommentHierarchy;


-- name: CreateComment :exec
INSERT INTO comments (name, content, slug, parent_id) VALUES (?, ?, ?, ?);
-- INSERT INTO comments (name, content, slug, parent_id) VALUES ('adem', 'yo waddup', '/writing/2-types-of-work', NULL);

-- name: DeleteComment :exec
DELETE FROM comments WHERE id = ?;

-- name: UpdateComment :exec
UPDATE comments SET content = ? WHERE id = ?;