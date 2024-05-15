create table if not exists question_submit
(
    id          varchar(256)                                                   comment "id" primary key,
    language    varchar(128)                                                   not null comment "编程语言",
    code        text                                                           not null comment "用户代码",
    judge_info  text                                                           null comment "判题信息json对象",
    status      int      default 0                                             not null comment "判题状态（0-待判题,1-判题中,2-成功,3-失败)",
    question_id varchar(256)                                                   not null comment "判题id",
    user_id     varchar(256)                                                   not null comment "创建用户id",
    create_time datetime default CURRENT_TIMESTAMP                             not null comment "创建时间",
    update_time datetime default CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP not null comment "更新时间",
    is_delete   tinyint  default 0                                             not null comment "是否删除",
    index idx_question_id (question_id),
    index idx_user_id (user_id)
) comment "题目提交";