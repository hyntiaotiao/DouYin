
一个伟大的作品

## 项目结构优化
一、命名规范
1. 将同一类型的接口合并在一个go文件中，如Login和Register接口合并在User.go中
2. 将原本的like改为favorite
3. service层命名与controller层统一
4. common中的对象结构体添加VO后缀，与repository的model区分
5. 局部变量统一采用小驼峰法，全局变量全字母大写，函数/方法、结构体、结构体字段统一采用大驼峰法

二、添加容错
1. 在控制台中输出错误信息

三、
