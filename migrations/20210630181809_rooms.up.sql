CREATE TABLE rooms
(
    session_id      VARCHAR(36) NOT NULL,
    user_id         VARCHAR(36) NOT NULL,
    role_id         VARCHAR(36) NOT NULL,
    is_active INT default 1 NOT NULL
);
