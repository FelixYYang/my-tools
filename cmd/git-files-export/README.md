## 导出当前工作目录下git项目的最新两个提交的差异文件
### 安装
```
$ go install github.com/FelixYYang/my-tools/cmd/git-files-export@latest
``` 

用法示例（windows）：导出当前提交点上的文件
~~~
$ git-files-export.exe
~~~
用法示例（windows）：导出指定提交点间的文件
~~~
$ git-files-export.exe HEAD~ HEAD
~~~