create table if not exists question
(
    id           varchar(255) comment "id" primary key,
    title        varchar(512)                       null comment "标题",
    content      text                               null comment "内容",
    tags         varchar(1024)                      null comment "标签列表json数组",
    answer       text                               null comment "题目答案",
    submit_num   int      default 0                 not null comment "题目提交数",
    accept_num   int      default 0                 not null comment "题目通过数",
    judge_case   text                               null comment "判题用例json数组",
    judge_config text                               null comment "判题配置json对象",
    thum_num     int      default 0                 not null comment "点赞数",
    user_id      varchar(256)                       not null comment "创建用户id",
    create_time  datetime default CURRENT_TIMESTAMP not null comment "创建时间",
    update_time  datetime default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment "更新时间",
    is_delete    tinyint  default 0                 not null comment "是否删除",
    index idx_userId (user_id)
) comment "题目" collate = utf8mb4_unicode_ci;