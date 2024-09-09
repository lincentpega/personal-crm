BEGIN;

CREATE TABLE IF NOT EXISTS public.persons (
    id SERIAL,
    first_name VARCHAR(256) NOT NULL,
    last_name VARCHAR(256),
    second_name VARCHAR(256),
    birth_date DATE,
    CONSTRAINT pk_persons PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public.contact_infos (
    person_id INT NOT NULL,
    method_name VARCHAR(256) NOT NULL,
    contact_data VARCHAR(256) NOT NULL,
    CONSTRAINT fk_contact_infos_persons
    FOREIGN KEY (person_id) REFERENCES persons (id)
);

CREATE TABLE IF NOT EXISTS public.job_infos (
    person_id INT NOT NULL,
    company VARCHAR(256) NOT NULL,
    job_position VARCHAR(256) NOT NULL,
    is_current BOOLEAN NOT NULL,
    CONSTRAINT fk_job_infos_persons
    FOREIGN KEY (person_id) REFERENCES persons (id)
);

CREATE TABLE IF NOT EXISTS public.person_settings (
    person_id INT NOT NULL,
    birthday_notify BOOLEAN NOT NULL,
    CONSTRAINT fk_person_settings_persons
    FOREIGN KEY (person_id) REFERENCES persons (id)
);

CREATE TYPE notification_type AS ENUM ('keep_in_touch');

CREATE TABLE IF NOT EXISTS public.notifications (
    person_id INT NOT NULL,
    type notification_type,
    notification_time TIMESTAMP,
    description TEXT,
    CONSTRAINT fk_notifications_persons
    FOREIGN KEY (person_id) REFERENCES persons (id)
);

COMMIT;
