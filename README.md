在线online judge系统go语言后端项目，在用户中心的基础上添加题目提交查询，判题等功能模块。
该项目实现了用户管理，题目管理，题目提交，题目判题等。
使用gin作为web开发框架，集成了gorm，docker，swagger， viper，session的redis存储，mysql存储后端数据，由于学识有限，尚未完全了解如何解决redis缓存数据同步的问题，因此，在反复修改后，最终移除了redis缓存部分，待后续逐渐改进。