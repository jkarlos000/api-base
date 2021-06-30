CREATE TABLE nurses
(
    id         VARCHAR(36) PRIMARY KEY,
    userid       VARCHAR(36) NOT NULL,
    is_working INT(11) NOT NULL,
    latitude    DECIMAL(10,16) NOT NULL ,
    longitude   DECIMAL (10,16) NOT NULL
);