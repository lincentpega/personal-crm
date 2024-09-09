BEGIN;

DROP TABLE IF EXISTS job_infos;
DROP TABLE IF EXISTS contact_infos;
DROP TABLE IF EXISTS persons;
DROP TABLE IF EXISTS person_settings;
DROP TYPE IF EXISTS notification_type;
DROP TABLE IF EXISTS keep_in_touch_notifications;

COMMIT;
