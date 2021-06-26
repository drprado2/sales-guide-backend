CREATE TABLE company
(
    id                   uuid                  DEFAULT uuid_generate_v4() PRIMARY KEY,
    name                 VARCHAR(150) NOT NULL,
    logo                 TEXT         NOT NULL,
    total_colaborators   integer      NOT NULL DEFAULT 0,
    primary_color        VARCHAR(30)  NOT NULL,
    primary_font_color   VARCHAR(30)  NOT NULL,
    secondary_color      VARCHAR(30)  NOT NULL,
    secondary_font_color VARCHAR(30)  NOT NULL,
    created_at           TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TRIGGER set_company_updated_at
    BEFORE UPDATE
    ON company
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();