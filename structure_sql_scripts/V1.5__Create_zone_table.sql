CREATE TABLE zone
(
    id                   uuid                   PRIMARY KEY NOT NULL,
    company_id           uuid NOT NULL,
    name                 VARCHAR(250)  NOT NULL,
    total_sellers         integer       NOT NULL,
    total_proves_sent      integer       NOT NULL,
    total_treinament_dones integer       NOT NULL,
    description          VARCHAR(2000) NULL,
    created_at           TIMESTAMPTZ   NOT NULL,
    updated_at           TIMESTAMPTZ   NOT NULL,
    FOREIGN KEY (company_id) REFERENCES company (id)
);