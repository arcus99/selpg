# Selpg的go语言实现

由于老师已经给出了selpg的C语言实现，所以大部分代码是简单地将go语言翻译成C语言即可。相关的知识见连接[selpg的实现][1]。

需要改动的地方：
---
###1、process_args中关于命令行的实现
这里使用老师介绍的flag包对命令行进行读取和分析，同时参考了以下博客
[go学习过程中flag包的应用][2]
这一部分是对原C语言实现的代码改动最大的部分。
###2、关于go语言中读取的实现
[stin][3]以及网上一些零散的材料。学习之后可以看出来读取和C语言是十分相似的，所以直接对原有的代码进行改写即可。

##实验结果：

1、把“input_file”的第1页写至标准输出（也就是屏幕）

2、selpg 读取标准输入，标准输入已被shell／内核重定向为来自“inputfile”而不是显式命名的文件名参数

3、selpg 将第一页写至标准输出；标准输出被shell／内核重定向至“output_file”

4、“other_command”（这里为cat）的标准输出被shell／内核重定向至selpg的标准输入。将第一页写至 selpg 的标准输出（屏幕）。

结果截图在附件中上传。
  [1]: https://www.ibm.com/developerworks/cn/linux/shell/clutil/index.html
  [2]: https://studygolang.com/articles/5608
  [3]: https://stackoverflow.com/questions/29060922/reading-from-stdin-in-golang