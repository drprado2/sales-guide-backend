CREATE TABLE company
(
    id                   uuid PRIMARY KEY NOT NULL,
    name                 VARCHAR(150) NOT NULL,
    document             VARCHAR(30)  NOT NULL,
    logo                 TEXT         NOT NULL,
    total_colaborators   integer      NOT NULL,
    primary_color        VARCHAR(30)  NOT NULL,
    primary_font_color   VARCHAR(30)  NOT NULL,
    secondary_color      VARCHAR(30)  NOT NULL,
    secondary_font_color VARCHAR(30)  NOT NULL,
    created_at           TIMESTAMPTZ  NOT NULL,
    updated_at           TIMESTAMPTZ  NOT NULL
);
