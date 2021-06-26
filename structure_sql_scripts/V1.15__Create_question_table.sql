CREATE TABLE question
(
    id            uuid                 DEFAULT uuid_generate_v4() PRIMARY KEY,
    company_id    uuid        NOT NULL,
    treinament_id uuid        NOT NULL,
    enunciated    TEXT        NOT NULL,
    timeout       integer     NOT NULL DEFAULT 10,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (company_id) REFERENCES company (id),
    FOREIGN KEY (treinament_id) REFERENCES treinament (id)
);

CREATE TRIGGER set_question_updated_at
    BEFORE UPDATE
    ON question
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();