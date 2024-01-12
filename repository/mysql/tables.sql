CREATE DATABASE IF NOT EXISTS payhere COLLATE utf8_general_ci;

USE payhere;

CREATE TABLE users
(
    user_id      int UNSIGNED AUTO_INCREMENT
        PRIMARY KEY,
    user_name    varchar(50)                        NOT NULL,
    phone_number char(13)                           NOT NULL,
    password     varchar(255)                       NOT NULL,
    created_at   datetime DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT phone_number
        UNIQUE (phone_number)
);

CREATE TABLE item_categories
(
    category_id   int UNSIGNED AUTO_INCREMENT
        PRIMARY KEY,
    user_id       int UNSIGNED                       NOT NULL,
    category_name varchar(20)                        NOT NULL,
    created_at    datetime DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT uidx_user_id_category_name
        UNIQUE (user_id, category_name),
    CONSTRAINT item_categories_ibfk_1
        FOREIGN KEY (user_id) REFERENCES users (user_id)
            ON DELETE CASCADE
);

CREATE TABLE items
(
    item_id      int UNSIGNED AUTO_INCREMENT
        PRIMARY KEY,
    user_id      int UNSIGNED                       NOT NULL,
    category_id  int UNSIGNED                       NOT NULL,
    item_name    varchar(255)                       NOT NULL,
    price        int UNSIGNED                       NOT NULL,
    cost         int UNSIGNED                       NOT NULL,
    description  varchar(255)                       NOT NULL,
    barcode      varchar(255)                       NOT NULL,
    barcode_type varchar(20)                        NOT NULL,
    item_size    varchar(255)                       NOT NULL,
    created_at   datetime DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT uidx_user_id_item_name
        UNIQUE (user_id, item_name),
    CONSTRAINT items_ibfk_1
        FOREIGN KEY (user_id) REFERENCES users (user_id)
            ON DELETE CASCADE,
    CONSTRAINT items_ibfk_2
        FOREIGN KEY (category_id) REFERENCES item_categories (category_id)
            ON DELETE CASCADE
);

CREATE TABLE token_blacklist (
    token VARCHAR(500) NOT NULL PRIMARY KEY,
    expires_at datetime NOT NULL
);