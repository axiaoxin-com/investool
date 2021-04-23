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
- 按选择策略自动选择筛选出的结果并计算对应的持仓金额 TODO
- 将筛选结果导出为 JSON 文件
- 将筛选结果导出为 CSV 文件
- 将筛选结果导出为 EXCEL 文件，并按行业、价格、历史波动率分工作表
- 将筛选结果导出为股票代码图片便于东方财富 APP 上导入到自选列表 TODO
- 支持搜索给定股票名称或代码并对其进行评估 TODO
- 提供终端界面操作 TODO
- 导入自选通过接口实现自动化 TODO
- 获取 K 线图数据 TODO
- 按我的 X 均线交易法则筛选可以买入的短期股票 TODO
- 监控自选和持仓的股票，触发交易纪律则进行提醒 TODO
- 提供 WEB 界面操作 TODO

## 我的选股规则

根据各种指标筛选值得长期持有其股票进行投资的优质公司。（优质公司不代表当前股价在涨，长线存股）

### 1. 财务优质

- 最新 ROE 高于 8%
- ROE 平均值小于 20 时，至少 3 年内逐年递增
- EPS 至少 3 年内逐年递增
- 营业总收入至少 3 年内逐年递增
- 净利润至少 3 年内逐年递增
- 负债率低于 60%
- 配发股利

### 2. 估值优质

- 估值较低或中等
- 股价低于合理价格（`合理价格 = 历史市盈率中位数 * (EPS * (1 + 今年 Q1 营收增长比))`）

### 3. 低波动

- 选择行业分散，均衡配置
- 按历史波动率进行持仓占比： 0.5 以下:0.5 以上 = 6:4

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
