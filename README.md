# x-stock

[![Open Source Love](https://badges.frapsoft.com/os/v1/open-source.svg?v=103)](https://github.com/axiaoxin-com/x-stock/)
[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org)
[![visitors](https://visitor-badge.glitch.me/badge?page_id=axiaoxin-com.x-stock)](https://github.com/axiaoxin-com/x-stock)
[![GitHub release](https://img.shields.io/github/release/axiaoxin-com/x-stock.svg)](https://gitHub.com/axiaoxin-com/x-stock/releases/)
[![Github all releases](https://img.shields.io/github/downloads/axiaoxin-com/x-stock/total.svg)](https://gitHub.com/axiaoxin-com/x-stock/releases/)

使用 Golang 实现的股票机器人，股票数据来源于东方财富网、亿牛网、腾讯证券。

该程序不构成任何投资建议，程序只是我个人辅助获取数据的工具，具体分析仍然需要自己判断。

## 功能

- 按指定条件的默认值自动筛选可以长期持有其股票的优质公司
- 按指定条件的自定义值自动筛选可以长期持有其股票的优质公司
- 实现股票检测器，封装基本面检测方法
- 将筛选结果导出为 JSON 文件
- 将筛选结果导出为 CSV 文件
- 将筛选结果导出为 EXCEL 文件，并按行业、价格、历史波动率分工作表
- 将筛选结果导出为股票代码图片便于东方财富 APP 上导入到自选列表
- 支持关键词搜索股票并对其进行评估
- 检测器支持对银行股按不同规则进行检测
- 支持净利率和毛利率稳定性判断
- 增加获取东方财富智能诊股中综合评价和价值评估信息
- TODO:
- 筛选参数和检测参数支持命令行自定义
- 按选择策略自动选择筛选并导出结果
- 提供终端界面操作
- 实现支持完整的东方财富选股器条件
- 导入自选通过接口实现自动化
- 获取 K 线图数据
- 检测器实现技术面检测方法
- 按我的 X 均线交易法则筛选可以买入的短期股票
- 监控自选和持仓的股票，触发交易纪律则进行提醒
- 提供 WEB 界面操作

## 我的选股规则

根据各种指标筛选值得长期持有其股票进行投资的优质公司。（优质公司不代表当前股价在涨，长线存股）

### 1. 财务优质

- 最新 ROE 高于 8%
- ROE 平均值小于 20 时，至少 3 年内逐年递增
- EPS 至少 3 年内逐年递增
- 营业总收入至少 3 年内逐年递增
- 净利润至少 3 年内逐年递增
- 负债率低于 60% （金融股除外）
- 配发股利
- 营收本业比高于 80%，经营内容具备基础、垄断、必要性服务等特征
- 净利率、毛利率稳定（标准差不高于 1 ）

### 2. 估值优质

- 估值较低或中等
- 股价低于合理价格（`合理价格 = 历史市盈率中位数 * (EPS * (1 + 今年 Q1 营收增长比))`）

### 3. 配置均衡

- 选择行业分散配置
- 历史波动率在 1 以内，进行持仓占比： 0.5 以下:0.5 以上 = 6:4

## 银行股选择标准

银行股由于盈利模式不同，因此选择标准略有不同。

1. ROE 高于 8%
2. ROA 高于 0.5%
3. 资本充足率高于 8%
4. 不良贷款率低于 3%
5. 不良贷款拨备覆盖率高于 100%

## 低价股选择标准

低价股定义：价格范围 10 - 30 元

1. 市净率不低于 1%
2. ROE 不低于 8%
3. 配发股利

## 我的交易纪律

### 长线

TODO

### 短线

TODO

## 使用方法

[点击下载对应系统最新版本解压得到可执行文件](https://github.com/axiaoxin-com/x-stock/releases/)

查看使用说明：

```
./x-stock -h
Usage of ./x-stock:
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
./x-stock -run exportor -f ./x-stock.20210424.json
```

- 导出 CSV 文件：

```
./x-stock -run exportor -f ./x-stock.20210424.csv
```

- 导出 EXCEL 文件：

```
./x-stock -run exportor -f ./x-stock.20210424.xlsx
```

- 导出 PNG 图片：

```
./x-stock -run exportor -f ./x-stock.20210424.png
```

- 导出全部支持的类型：

```
./x-stock -run exportor -l info -f docs/x-stock.20210426.all
```

### checker

给定关键词/股票代码搜索股票进行评估检测

命令：

```
./x-stock_darwin_amd64 -run checker -l error -k 比亚迪

+--------------------+--------------------------------+
|  未通过检测的指标  |              原因              |
+--------------------+--------------------------------+
| 净资产收益率 (ROE) | 最新 ROE:7.43 低于:8           |
+--------------------+--------------------------------+
| ROE 逐年递增       | 3 年内未逐年递增:[7.43 2.62    |
|                    | 4.96]                          |
+--------------------+--------------------------------+
| EPS 逐年递增       | 3 年内未逐年递增:[1.47 0.5     |
|                    | 0.93]                          |
+--------------------+--------------------------------+
| 营收逐年递增       | 3                              |
|                    | 年内未逐年递增:[1.56597691e+11 |
|                    | 1.27738523e+11 1.30054707e+11] |
+--------------------+--------------------------------+
| 净利润逐年递增     | 3 年内未逐年递增:[4.234267e+09 |
|                    | 1.61445e+09 2.780194e+09]      |
+--------------------+--------------------------------+
| 负债率             | 负债率:67.936140               |
|                    | 高于:60.000000                 |
+--------------------+--------------------------------+
|  比亚迪 002594 SZ  |             FAILED             |
+--------------------+--------------------------------+

```

## 欢迎 Star

[![Stargazers over time](https://starchart.cc/axiaoxin-com/x-stock.svg)](https://githuv.com/axiaoxin-com/x-stock)
