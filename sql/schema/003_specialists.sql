-- +goose Up 
CREATE TABLE specialists (
    specialist_id BIGSERIAL PRIMARY KEY,
    user_id INT UNIQUE NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    specialization VARCHAR(255),
    license_number VARCHAR(50)
);
CREATE INDEX idx_specialists_user_id ON specialists(user_id);
CREATE INDEX idx_specialists_specialization ON specialists(specialization);
-- +goose Down
DROP TABLE specialists;
