# stockcmd

## 介绍

一个简单的工具，可以在命令行查看股票的一些信息，上图：

![image](https://github.com/upwell/stockcmd/blob/master/images/sc1.png)

历史数据是从[baostock][1]抓取的，实时数据是从sina的接口拿的。

## 安装

- 下载编译好的binary，从release中下载binary，放到`/usr/local/bin/`下面即可。

暂时只编译了mac版的binary🤣

- 自行编译，需要有go的环境
```bash
> cd src
> go build
``` 
拷贝编译好的`stockcmd`到`/usr/local/bin/`下面

## 使用

只在mac系统上做了测试，以常规的一个使用路径举例

- 创建一个group，名称是hold
```bash
> stockcmd group create hold
```

- 加入股票立讯精密到group hold中，lx是立讯的拼音首字母
```bash
> stockcmd group add hold lx
```
会有个提示列表出来，上下键选择，按`回车`确认加入，`ctrl+c`取消；

![image](https://github.com/upwell/stockcmd/blob/master/images/sc2.png)

需要加入更多股票到分组的话，反复执行上面的命令即可。

- 显示这个分组的股票信息
```bash
> stockcmd show hold
```
加入后第一次显示会需要一点时间去抓取历史数据，之后就很快了。

- 快捷命令显示分组
有个快捷的script可以显示分组，需要把scripts/sc拷贝到`/usr/local/bin/`下面
```bash
> sc hold
or
# 不输入分组名称，默认会使用hold分组
> sc 
```

- 从分组中删除某只股票
```bash
> stockcmd group remove hold
```
上下键选择后，按`回车`确认删除，`ctrl+c`取消；


[1]: http://baostock.com