CREATE TABLE seller
(
    id                     uuid                  DEFAULT uuid_generate_v4() PRIMARY KEY,
    company_id             uuid         NOT NULL,
    employee_type_id       uuid         NOT NULL,
    zone_id                uuid         NOT NULL,
    name                   VARCHAR(250) NOT NULL,
    email                  VARCHAR(250) NOT NULL,
    document               VARCHAR(30)  NOT NULL,
    phone                  VARCHAR(15)  NOT NULL,
    birth_date             DATE         NULL,
    password               VARCHAR(80)  NOT NULL,
    enable                 BOOLEAN      NOT NULL,
    total_proves_sent      integer      NOT NULL DEFAULT 0,
    total_treinaments_done integer      NOT NULL DEFAULT 0,
    last_acess             DATE         NULL,
    avatar_image           TEXT         NOT NULL,
    created_at             TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at             TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    FOREIGN KEY (company_id) REFERENCES company (id),
    FOREIGN KEY (employee_type_id) REFERENCES employee_type (id),
    FOREIGN KEY (zone_id) REFERENCES zone (id)
);

CREATE TRIGGER set_seller_updated_at
    BEFORE UPDATE
    ON seller
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();