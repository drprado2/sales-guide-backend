CREATE TABLE app_user
(
    id                    uuid    DEFAULT uuid_generate_v4() PRIMARY KEY,
    company_id            uuid              NOT NULL,
    name                  VARCHAR(250)      NOT NULL,
    email                 VARCHAR(250)      NOT NULL,
    phone                 VARCHAR(15)       NULL,
    birth_date            DATE              NULL,
    password              VARCHAR(80)       NOT NULL,
    avatar_image          TEXT              NOT NULL,
    record_creation_count INTEGER DEFAULT 0 NOT NULL,
    record_editing_count  INTEGER DEFAULT 0 NOT NULL,
    record_deletion_count INTEGER DEFAULT 0 NOT NULL,
    last_acess            DATE              NULL,
    created_at            TIMESTAMPTZ       NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ       NOT NULL DEFAULT NOW(),
    FOREIGN KEY (company_id) REFERENCES company (id)
);

CREATE TRIGGER set_user_updated_at
    BEFORE UPDATE
    ON app_user
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();