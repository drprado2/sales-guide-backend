CREATE TABLE question
(
    id            uuid                 PRIMARY KEY,
    company_id    uuid        NOT NULL,
    treinament_id uuid        NOT NULL,
    enunciated    TEXT        NOT NULL,
    timeout       integer     NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL,
    updated_at    TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (company_id) REFERENCES company (id),
    FOREIGN KEY (treinament_id) REFERENCES treinament (id)
);