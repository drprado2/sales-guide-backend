CREATE TABLE question_option
(
    id          uuid                 PRIMARY KEY,
    company_id  uuid        NOT NULL,
    question_id uuid        NOT NULL,
    content     TEXT        NOT NULL,
    correct     BOOLEAN     NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL,
    updated_at  TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (company_id) REFERENCES company (id),
    FOREIGN KEY (question_id) REFERENCES question (id)
);