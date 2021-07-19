CREATE TABLE employee_type
(
    id          uuid                   PRIMARY KEY NOT NULL,
    company_id  uuid          NOT NULL,
    name        VARCHAR(250)  NOT NULL,
    description VARCHAR(2000) NULL,
    created_at  TIMESTAMPTZ   NOT NULL,
    updated_at  TIMESTAMPTZ   NOT NULL,
    FOREIGN KEY (company_id) REFERENCES company (id)
);