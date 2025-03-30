-- +goose Up
-- +goose StatementBegin
CREATE TABLE expenses
(
    id SERIAL PRIMARY KEY,
    amount INTEGER NOT NULL,
    subject TEXT NOT NULL,
    receiver TEXT DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX expenses_id_idx ON expenses (id);
CREATE INDEX expenses_subject_idx ON expenses USING gin (to_tsvector('russian', subject));
CREATE INDEX expenses_receiver_idx ON expenses USING gin (to_tsvector('russian', receiver));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE expenses;
-- +goose StatementEnd
