drop database if exists httpapi_transponder;
create database httpapi_transponder CHARACTER SET utf8 COLLATE utf8_general_ci;
use httpapi_transponder;

drop table if exists user;
create table user
(
    id            int auto_increment primary key,
    user_id       varchar(48)                        not null,
    user_name     varchar(48)                        not null,
    user_password varchar(128)                       not null,
    user_type     int      default 0                 null,
    remark        mediumtext,
    create_user   int      default 0                 null,
    create_time   datetime default CURRENT_TIMESTAMP null,
    update_user   int      default 0                 null,
    update_time   datetime default CURRENT_TIMESTAMP null,
    constraint user_user_id_uindex
        unique (user_id),
    constraint user_user_name_uindex
        unique (user_name)
);
INSERT INTO `user` (user_id, user_name, user_password, user_type, remark, create_user, create_time, update_user, update_time) VALUES ('u0001', 'admin', '21232f297a57a5a743894a0e4a801fc3', 1, null, 0, '2021-08-17 03:18:49', 0, '2021-08-17 03:18:49');

drop table if exists user_token;
create table user_token
(
    token_id     int auto_increment     primary key,
    user_id      varchar(48)            not null,
    access_token varchar(256)           not null,
    fresh_token  varchar(256)           null,
    create_user  int      default 0     null,
    create_time  datetime default now() null,
    update_user  int      default 0     null,
    update_time  datetime default now() null
);

drop table if exists user_aksk;
create table user_aksk
(
    aksk_id     int auto_increment     primary key,
    user_id     varchar(48)            not null,
    access_key  varchar(64)            not null,
    secure_key  varchar(64)            null,
    create_user int      default 0     null,
    create_time datetime default now() null,
    update_user int      default 0     null,
    update_time datetime default now() null
);
