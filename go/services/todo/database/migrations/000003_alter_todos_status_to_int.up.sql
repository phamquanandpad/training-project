ALTER TABLE todos ADD COLUMN status_int TINYINT UNSIGNED NOT NULL DEFAULT 0;

UPDATE todos
SET status_int =
    CASE status
        WHEN 'pending' THEN 0
        WHEN 'in_progress' THEN 1
        WHEN 'done' THEN 2
    END;

ALTER TABLE todos DROP COLUMN status;

ALTER TABLE todos CHANGE status_int status TINYINT UNSIGNED NOT NULL;

ALTER TABLE todos
ADD CONSTRAINT check_todos_status CHECK (status IN (0,1,2));
