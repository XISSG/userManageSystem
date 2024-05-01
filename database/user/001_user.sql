-- +goose Up
create table user
(
    id           bigint auto_increment comment 'id'
        primary key,
    user_name     varchar(256)                         null comment '用户昵称',
    user_account  varchar(256)                         null comment '用户账号',
    avatar_url    varchar(1024)                        null comment '用户头像',
    user_password varchar(512)                         not null comment '用户密码',
    create_time   datetime   default CURRENT_TIMESTAMP null comment '创建时间',
    update_time   datetime   default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP comment '更新时间',
    is_delete     tinyint(1) default 0                 not null comment '是否删除',
    user_role     int        default 0                 not null comment '0为普通用户, 1为会员用户'
)
    comment '用户';
-- +goose Down
drop table  if exists user;
