# ayachan

[![codebeat badge](https://codebeat.co/badges/3482bd1e-45d7-4e83-af70-3f1ccb874fcd)](https://codebeat.co/projects/github-com-6qhtsk-ayachan-development)
[![Go Report Card](https://goreportcard.com/badge/github.com/6QHTSK/ayachan)](https://goreportcard.com/report/github.com/6QHTSK/ayachan)
![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/6QHTSK/ayachan)
![GitHub](https://img.shields.io/github/license/6QHTSK/ayachan)
![Libraries.io dependency status for GitHub repo](https://img.shields.io/librariesio/github/6QHTSK/ayachan)

一个异步获取Bestdori站点自制谱面，并提供高级搜索功能、谱师谱面分析（WIP）等功能

同时具有谱面特征提取，谱面难度计算模块（将在未来拆分）

## 安装

该项目使用go、mysql、meilisearch。

### 数据库初始化

1. （WIP mysql数据库启动）
2. （WIP meilisearch数据库index启动）

### 主程序

从release下载对应版本或下载源码

```bash
    git clone https://github.com/6QHTSK/ayachan
    go build && ./ayachan
```

首次运行会生成yaml配置文件，请配置好数据库、设定运行地址、远程Bestdori抓取API后再次运行。

## 使用方法

### log

该项目会在控制台上打印debug级别以上的log,同时生成logs/console.log, logs/warnings.log文件存储一般的log和warning级别以上的log

### ...(WIP)

