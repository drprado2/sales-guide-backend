CREATE TABLE app_user
(
    id                    uuid     PRIMARY KEY NOT NULL,
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
    created_at            TIMESTAMPTZ       NOT NULL,
    updated_at            TIMESTAMPTZ       NOT NULL,
    FOREIGN KEY (company_id) REFERENCES company (id)
);
