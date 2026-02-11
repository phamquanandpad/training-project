ALTER TABLE todos ADD COLUMN status_str VARCHAR(20) NOT NULL DEFAULT 'pending';

UPDATE todos
SET status_str =
    CASE status
        WHEN 0 THEN 'pending'
        WHEN 1 THEN 'in_progress'
        WHEN 2 THEN 'done'
    END;

ALTER TABLE todos DROP COLUMN status;

ALTER TABLE todos CHANGE status_str status VARCHAR(20) NOT NULL;

ALTER TABLE todos
ADD CONSTRAINT check_todos_status
CHECK (status IN ('pending','in_progress','done'));
