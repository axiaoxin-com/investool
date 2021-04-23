# x-stock

[![Open Source Love](https://badges.frapsoft.com/os/v1/open-source.svg?v=103)](https://github.com/axiaoxin-com/x-stock/)
[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org)
[![visitors](https://visitor-badge.glitch.me/badge?page_id=axiaoxin-com.x-stock)](https://github.com/axiaoxin-com/x-stock)
[![GitHub release](https://img.shields.io/github/release/axiaoxin-com/x-stock.svg)](https://gitHub.com/axiaoxin-com/x-stock/releases/)
[![Github all releases](https://img.shields.io/github/downloads/axiaoxin-com/x-stock/total.svg)](https://gitHub.com/axiaoxin-com/x-stock/releases/)

使用 Golang 实现的股票机器人，股票数据来源于东方财富网和亿牛网。

## 功能

- 按指定条件的默认值自动筛选可以长期持有其股票的优质公司
- 按指定条件的自定义值自动筛选可以长期持有其股票的优质公司
- 将筛选结果导出为 JSON 文件
- 将筛选结果导出为 CSV 文件
- 将筛选结果导出为 EXCEL 文件，并按行业分工作表
- 根据给定股票代码/股票名称判断是否为可以长期持有其股票的优质公司 TODO
- 按行业筛选可以长期持有其股票的优质公司 TODO
- 按历史波动率和行业计算持仓金额 TODO
- 提供 WEB 界面筛选可以长期持有其股票的优质公司 TODO
- 提供终端界面筛选可以长期持有其股票的优质公司 TODO
- 添加给定股票到东方财富自选列表 TODO
- 按我的 X 均线交易法则筛选可以买入的短期股票 TODO
- 监控自选和持仓的股票，触发交易纪律则进行提醒 TODO

## 我的选股规则

根据各种指标筛选值得长期持有其股票进行投资的优质公司。（优质公司不代表当前股价在涨，长线存股）

### 1. 财务优质

- 最新 ROE 高于 8%
- ROE 平均值小于 20 时，至少 3 年内逐年递增
- EPS 至少 3 年内逐年递增
- 营业总收入至少 3 年内逐年递增
- 净利润至少 3 年内逐年递增
- 负债率低于 60%

### 2. 估值优质

- 估值较低或中等
- 股价低于合理价格（合理价格 = 历史市盈率中位数 _ (EPS _ (1 + 今年 Q1 营收增长比))）

### 3. 低波动

- 行业分散，均衡配置
- 历史波动率在 1 以内（持仓占比： 0.1:0.1-0.5:0.5-1 = 3:3:4 ）

## 我的交易纪律

### 长线

TODO

### 短线

TODO

## 使用方法

[点击下载对应系统最新版本的可执行文件](https://github.com/axiaoxin-com/x-stock/releases/)，按下列使用说明操作，可执行文件名需替换为你下载的版本。

查看使用说明：

```
./x-stock_darwin_amd64 -h
Usage of ./x-stock_darwin_amd64:
  -f string
    	export filename (default "./docs/x-stock.20210424.xlsx")
  -l string
    	loglevel: debug|info|warn|error (default "info")
  -run string
    	processor: exportor|tview|webserver
```

### exportor

exportor 是数据导出器，不使用参数默认导出 EXCEL 文件。

支持导出以下类型的数据：

- 导出 JSON 文件：

```
./x-stock_darwin_amd64 -run exportor -f ./x-stock.20210424.json
```

- 导出 CSV 文件：

```
./x-stock_darwin_amd64 -run exportor -f ./x-stock.20210424.csv
```

- 导出 EXCEL 文件：

```
./x-stock_darwin_amd64 -run exportor -f ./x-stock.20210424.xlsx
```

## 欢迎 Star

[![Stargazers over time](https://starchart.cc/axiaoxin-com/x-stock.svg)](https://githuv.com/axiaoxin-com/x-stock)
