CREATE TABLE roles
(
    id         VARCHAR(36) PRIMARY KEY,
    name       VARCHAR(50) NOT NULL,
    constraint roles_name_uindex
        unique (name)
);
