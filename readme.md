
一个伟大的作品

## 项目结构优化
一、命名规范
1. 将同一类型的接口合并在一个go文件中，如Login和Register接口合并在User.go中
2. 将原本的like改为favorite
3. service层命名与controller层统一
4. common中的对象结构体添加VO后缀，与repository的model区分

二、添加容错
在控制台中输出错误信息

三、内部改动
1. 原本在controller层中发送http请求获得评论用户数据效率过低，已将其更改为直接通过数据库操作获得

四、添加配置信息管理，可以直接通过修改applicaiton.yml来对项目配置信息进行更改
## 项目结构介绍
本项目采用MVC结构