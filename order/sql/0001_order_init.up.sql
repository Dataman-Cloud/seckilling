create database if not exists seckilling
    default character set utf8
    default collate utf8_general_ci;

use seckilling;

create table if not exists event (
    id bigint unsigned primary key auto_increment,
    name varchar(128) not null,
    summary varchar(1024),
    client varchar(64), -- customer info
    type tinyint unsigned not null, -- 1 counter 2 instances
    qty int unsigned not null,
    balance int unsigned default 0, -- current ordered quantity
    ext text, -- json for extensions
    created timestamp,
    updated timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    ) character set utf8 collate utf8_general_ci;

alter table event add index name (name);
alter table event add index client (client);

create table if not exists prod_inst (
    id bigint unsigned primary key auto_increment,
    eid bigint not null,
    seq int unsigned not null,
    sku int unsigned not null,
    status tinyint unsigned default 0,
    created timestamp,
    updated timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    ) character set utf8 collate utf8_general_ci;
alter table prod_inst add index unique uni_eid_seq (eid, seq);
alter table prod_inst add unique (eid, seq);

create table if not exists order (
    id bigint unsigned primary key auto_increment,
    eid bigint unsigned not null,
    uid varchar(64) not null, -- user token
    status tinyint(1) not null,
    eid bigint not null,
    seq int unsigned not null,
    ext text, -- json for extension, eg. headers
    created timestamp,
    updated timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    ) character set utf8 collate utf8_general_ci;

ALTER TABLE ORDER ADD INDEX UID (UID); 

