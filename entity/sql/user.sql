create table user
(
    id            varchar(256) primary key comment "用户id",
    user_name     varchar(256) comment "用户名",
    user_account  varchar(256)                                                   not null comment "用户账户不允许重复",
    avatar_url    varchar(1024)                                                  null comment "用户头像",
    user_password longtext                                                       not null comment "用户密码",
    create_time   datetime default CURRENT_TIMESTAMP                             null comment "创建时间",
    update_time   datetime default CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP null comment "更新时间",
    is_delete     tinyint  default 0                                             not null comment "是否删除,0为不删除，1为删除",
    user_role     varchar(64)                                                    null comment "用户类型，有user,admin,ban"
) comment '用户' collate = utf8mb4_unicode_ci;;