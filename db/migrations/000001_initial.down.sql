BEGIN;
DROP TABLE IF EXISTS notifications;
DROP TYPE IF EXISTS notification_type;
DROP TYPE IF EXISTS notification_status;
DROP TABLE IF EXISTS job_infos;
DROP TABLE IF EXISTS contact_infos;
DROP TABLE IF EXISTS person_settings;
DROP TABLE IF EXISTS persons;
COMMIT;