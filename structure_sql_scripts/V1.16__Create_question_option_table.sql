CREATE TABLE question_option
(
    id          uuid                 DEFAULT uuid_generate_v4() PRIMARY KEY,
    company_id  uuid        NOT NULL,
    question_id uuid        NOT NULL,
    content     TEXT        NOT NULL,
    correct     BOOLEAN     NOT NULL DEFAULT false,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (company_id) REFERENCES company (id),
    FOREIGN KEY (question_id) REFERENCES question (id)
);

CREATE TRIGGER set_question_option_updated_at
    BEFORE UPDATE
    ON question_option
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();