CREATE TABLE users
(
    id         VARCHAR(36) PRIMARY KEY,
    username       VARCHAR(50) NOT NULL,
    password    VARCHAR(150) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    is_active INT(11) NOT NULL,
    first_name  VARCHAR(50) NOT NULL,
    last_name   VARCHAR(50) NOT NULL
);
