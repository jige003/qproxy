# qproxy
### 简介
> qproxy 是一个http迷你代理型服务器， 基于golang的goproxy开发， 主要的功能是根据一些简单的规则提取requests info包， 功能比较简单。

### 使用说明
>qproxy -l :9010 -v -d 
参数意思：
* -l 指明监听的端口， 默认为9010， 冒号不能少
* -v 主要是代理中request和response日志输出到stdout
* -d 打开调试， 查看requests info信息
* 默认为把满足要求的requests info 输出到qproxylog文件当中去
