CREATE DATABASE IF NOT EXISTS payhere COLLATE utf8_general_ci;

# USE payhere;

CREATE TABLE users
(
    user_id      BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    phone_number CHAR(13)                           NOT NULL,
    password     VARCHAR(72)                        NOT NULL,
    created_at   DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT phone_number
        UNIQUE (phone_number)
);

CREATE TABLE items
(
    item_id     BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id     BIGINT UNSIGNED                    NOT NULL,
    category    VARCHAR(100)                       NOT NULL,
    item_name   VARCHAR(100)                       NOT NULL,
    price       INT UNSIGNED                       NOT NULL,
    cost        INT UNSIGNED                       NOT NULL,
    description TEXT                               NOT NULL,
    barcode     VARCHAR(100)                       NOT NULL,
    item_size   ENUM ('small', 'large')            NOT NULL,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    expiry_at   DATETIME                           NOT NULL,
    CONSTRAINT uidx_user_id_item_name
        UNIQUE (user_id, item_name),
    CONSTRAINT items_ibfk_1
        FOREIGN KEY (user_id) REFERENCES users (user_id)
            ON DELETE CASCADE
);

CREATE TABLE token_blacklist
(
    token      VARCHAR(500) NOT NULL PRIMARY KEY,
    expires_at datetime     NOT NULL
);