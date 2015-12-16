create database if not exists seckilling
    default character set utf8
    default collate utf8_general_ci;

use seckilling;

create table if not exists order (
    id bigint unsigned primary key auto_increment,
    uid varchar(64) not null,
    status tinyint(1) not null,
    json text,
    created timestamp,
    updated timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    ) character set utf8 collate utf8_general_ci;
