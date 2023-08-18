# 极简版抖音后端练习项目

## 技术栈
主体框架：Go、Gorm、Gin
数据库：MySQL
对象存储：七牛云OSS
数字音频处理：ffmpeg
token认证：jwt
配置信息管理：viper
加密：bcrypt

## 概述：
本项目为一个简单的抖音后端项目，采用MVC架构，按照项目方案说明接口分为Comment(评论)、
Favorite(喜欢)、Feed(视频流)、Publish(视频发布)、Relation(关系)、User(用户)、Message(消息)七大模块，
各模块分成controller、service、dao三层进行数据的传输与操作。


