CREATE TABLE employee_type
(
    id          uuid                   DEFAULT uuid_generate_v4() PRIMARY KEY,
    company_id  uuid          NOT NULL,
    name        VARCHAR(250)  NOT NULL,
    description VARCHAR(2000) NULL,
    created_at  TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    FOREIGN KEY (company_id) REFERENCES company (id)
);

CREATE TRIGGER set_employee_type_updated_at
    BEFORE UPDATE
    ON employee_type
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();