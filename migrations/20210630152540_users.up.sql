CREATE TABLE users
(
    id         VARCHAR(36) PRIMARY KEY,
    username       VARCHAR(50) NOT NULL,
    email       VARCHAR(100) NOT NULL,
    password    VARCHAR(150) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    is_active BOOLEAN DEFAULT TRUE NOT NULL,
    first_name  VARCHAR(50) default 'Jhon' NOT NULL,
    last_name   VARCHAR(50) default 'Doe' NOT NULL,
    constraint users_username_uindex
        unique (username),
    constraint users_email_uindex
        unique (email)
);
