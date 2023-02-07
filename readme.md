
一个伟大的作品

# hjc 2023.02.03更新
## 更新的结构

DOUYIN
- controller
    - MessageChat: 返回聊天记录

- repository
    - message: message对应的user1和user2的返回
    - model: message的model
    - video: 更改错误的字段

- service
    - MessageService: message service层

main.go 开放"/message/chat/"接口
