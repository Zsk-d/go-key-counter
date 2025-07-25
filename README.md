# key-counter
新一代的键盘按键次数记录, go语言实现, 5分钟保存一次

输出结果
![Image](https://raw.githubusercontent.com/Zsk-d/Keyboard-counter/master/img/result.png)

## counter/main.go
记录器

- autoSaveJSON() 中调整记录频率

## export/main.go
导出工具

- KeyCodeMap 变量配置按键码与表格模板单元格位置
## export_temp.xlsx
结果导出表格模板

## build-*.bat
- build-console.bat 编译控制台版本
- build-headless.bat 编译无控制台版本
- build-export.bat 编辑导出程序
