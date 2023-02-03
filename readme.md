
一个伟大的作品

# hjc 2023.02.03更新
## 更新的结构

DOUYIN
- controller
    - PublishList 返回发布列表
    - Comment 返回评论列表、发布和删除评论
    - UserInfo 原有的request的json格式与给定的接口格式不一致

- repository
    - video 在原有的video.go上做了修改
    - comment 评论的数据库CRUD
    - like 插入和删除like的事务操作

- service
    - PublishService 发布列表的service层
    - CommentService 评论的service层
    - LikeService 注释了查询操作，只保留插入和删除操作

main.go 开放"/publish/list/"、"/comment/action/"、"/comment/list/"接口
