CREATE TABLE rooms
(
    session_id      VARCHAR(36) NOT NULL,
    user_id         VARCHAR(36) NOT NULL,
    role_id         VARCHAR(36) NOT NULL,
    is_active       BOOLEAN DEFAULT TRUE NOT NULL
);
