-- +goose Up
-- +goose StatementBegin
INSERT INTO users (email, password_hash, name, preference, weight_unit, height_unit, weight, height, image_uri)
VALUES 
    ('user1@example.com', '$2a$10$u9EJsKFbT5HwnHT0dKDiCehO4GqE9ssypqQO/BvprR5pWRJkJ0XW2', 'User One', 'default', 'kg', 'cm', 70, 175, NULL),
    ('user2@example.com', '$2a$10$u9EJsKFbT5HwnHT0dKDiCehO4GqE9ssypqQO/BvprR5pWRJkJ0XW2', 'User Two', 'default', 'kg', 'cm', 65, 168, NULL),
    ('user3@example.com', '$2a$10$u9EJsKFbT5HwnHT0dKDiCehO4GqE9ssypqQO/BvprR5pWRJkJ0XW2', 'User Three', 'default', 'kg', 'cm', 80, 180, NULL);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM users WHERE email IN ('user1@example.com', 'user2@example.com', 'user3@example.com');
-- +goose StatementEnd
