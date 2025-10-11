BEGIN;

ALTER TABLE conf_info ADD COLUMN IF NOT EXISTS child_objects jsonb;

COMMIT;