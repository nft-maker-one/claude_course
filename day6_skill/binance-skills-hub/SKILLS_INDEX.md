# Binance Skills Hub — Skills Index

本文档索引 `skills/` 目录下所有 skill，说明各 skill 的定位与所提供的功能。

---

## 目录

- [credential-guard（凭证安全守卫）](#credential-guard-凭证安全守卫)
- [binance（CEX 综合）](#binance-cex-综合)
- [fiat（法币支付）](#fiat-法币支付)
- [onchain-pay（链上支付）](#onchain-pay-链上支付)
- [p2p（P2P 交易）](#p2p-p2p-交易)
- [payment（扫码支付）](#payment-扫码支付)
- [square-post（广场发帖）](#square-post-广场发帖)
- [binance-agentic-wallet（Web3 钱包）](#binance-agentic-wallet-web3-钱包)
- [binance-tokenized-securities-info（代币化证券）](#binance-tokenized-securities-info-代币化证券)
- [crypto-market-rank（加密市场排行）](#crypto-market-rank-加密市场排行)
- [meme-rush（Meme 代币追踪）](#meme-rush-meme-代币追踪)
- [query-address-info（地址持仓查询）](#query-address-info-地址持仓查询)
- [query-token-audit（代币安全审计）](#query-token-audit-代币安全审计)
- [query-token-info（代币信息查询）](#query-token-info-代币信息查询)
- [trading-signal（智能资金信号）](#trading-signal-智能资金信号)

---

## credential-guard（凭证安全守卫）

**路径**: `skills/binance/credential-guard/`
**版本**: 1.0.0 | **作者**: binance-skills-hub

所有需要认证的 skill 的**前置安全层**，强制从 `.env` 读取凭证，杜绝密钥出现在对话中。

| 规则 | 说明 |
|------|------|
| 禁止从对话获取密钥 | 用户粘贴的 Key 一律拒绝，引导写入 `.env` |
| 统一从 `.env` 加载 | `source .env` 注入环境变量，不硬编码 |
| 输出自动脱敏 | 显示前 5 + 后 4 位，中间替换为 `...` |
| 禁止写入其他文件 | Key 只存 `.env`，不得出现在 SKILL.md / CLAUDE.md 等受版本控制的文件 |
| 验证 `.gitignore` | 每次 git 操作前确认 `.env` 已被排除 |

**覆盖的凭证**:
- `BINANCE_API_KEY` / `BINANCE_API_SECRET` → binance、p2p、fiat（SAPI）
- `X_SQUARE_OPENAPI_KEY` → square-post
- `ONCHAIN_PAY_PRIVATE_KEY_PATH` → onchain-pay（RSA 密钥路径）

**触发时机**: 任何其他 skill 发起认证请求时自动前置执行。

---

## binance（CEX 综合）

**路径**: `skills/binance/binance/`  
**版本**: 1.1.0 | **作者**: Binance  
**依赖**: `binance-cli`（npm 包 `@binance/binance-cli`）

通过 `binance-cli` 操作 Binance 中心化交易所的全部功能，涵盖现货、衍生品、钱包等几乎所有模块，**所有操作均需认证**。

| 子命令 | 功能 |
|--------|------|
| `spot` | 现货交易（下单、查询、撤单） |
| `futures-usds` | USDS-M 永续合约交易 |
| `futures-coin` | COIN-M 合约交易 |
| `derivatives-options` | 期权交易 |
| `derivatives-portfolio-margin` | 统一保证金账户 |
| `derivatives-portfolio-margin-pro` | Pro 统一保证金账户 |
| `margin-trading` | 杠杆交易（借贷、还款、爆仓查询） |
| `convert` | 一键兑换 |
| `simple-earn` | 活期/定期理财 |
| `staking` | 质押（ETH、SOL 等） |
| `dual-investment` | 双币投资产品 |
| `crypto-loan` | 加密贷款 |
| `vip-loan` | VIP 贷款 |
| `wallet` | 钱包（充提、资产划转、账单） |
| `fiat` | 法币出入金（CLI 版） |
| `pay` | Binance Pay 支付 |
| `c2c` | C2C 法币兑换 |
| `copy-trading` | 跟单交易 |
| `algo` | 算法交易（TWAP/VP 等） |
| `alpha` | Alpha 项目 |
| `mining` | 矿池 |
| `rebate` | 返佣 |
| `gift-card` | 礼品卡 |
| `sub-account` | 子账号管理 |

> **安全提示**: 执行任何生产下单操作前，须等待用户明确输入 `CONFIRM`。

---

## fiat（法币支付）

**路径**: `skills/binance/fiat/`  
**版本**: 1.1.0 | **作者**: Binance

通过 **公开 REST API** 查询 Binance 法币支付能力；订单历史等需要 API Key 的功能详见 `references/sapi-endpoints.md`。

| API | 功能 |
|-----|------|
| `get_capabilities` | 查询指定国家支持的法币、加密货币及业务类型（买/卖/充/提） |
| `get_buy_and_sell_payment_methods` | 获取买卖加密货币的可用支付方式及限额、报价 |
| `get_deposit_and_withdraw_payment_methods` | 获取法币充提的可用支付方式及限额 |
| `get_price` | 查询指定法币/加密货币对的参考汇率 |
| 认证端点（SAPI） | 查询法币订单历史、支付记录（需 API Key） |

**亮点**:
- 自动推断国家参数（从货币或对话上下文）
- 支持多语言 Binance 跳转链接（买/卖/充/提）
- 永不将 `US` 用作 country 参数

---

## onchain-pay（链上支付）

**路径**: `skills/binance/onchain-pay/`  
**版本**: 0.1.1 | **作者**: onchain-pay-team

帮助合作伙伴集成「法币→加密货币→链上转账」一体化支付流程，需 RSA SHA256 请求签名。

| API | 功能 |
|-----|------|
| `payment-method-list` | 查询指定法币/加密货币对的可用支付方式（信用卡、P2P、Apple Pay、Google Pay 等）及限额 |
| `trading-pairs` | 列出所有支持的法币和加密货币 |
| `estimated-quote` | 获取实时报价（汇率、手续费、预计到账加密货币数量） |
| `pre-order` | 创建买单并获取跳转至 Binance 支付流程的 URL |
| `order` | 查询订单状态与详情（处理中/完成/失败等） |
| `crypto-network` | 获取支持的区块链网络及提币手续费和限额 |
| `p2p/trading-pairs` | 列出 P2P 专属交易对 |

**典型场景**: 电商收款、法币→链上一步到位、跨链桥接前置流程、智能合约交互。

---

## p2p（P2P 交易）

**路径**: `skills/binance/p2p/`  
**版本**: 见 CHANGELOG | **作者**: Binance

支持 Binance P2P（C2C）市场的自然语言查询与操作。

| 能力 | 是否需要认证 |
|------|------------|
| 查询 P2P 广告列表（价格、限额、支付方式） | 否 |
| 比较不同支付方式的报价 | 否 |
| 查看自己的 P2P 订单历史与汇总 | 是 |
| 查询订单详情及时间线 | 是 |
| 查询申诉/投诉状态与历史 | 是 |
| 提交申诉补充证据 | 是 |
| 查看投诉处理流程（CS 备注、操作记录） | 是 |
| 取消申诉（不可逆） | 是 |
| 发布/更新/管理 P2P 广告（需商家权限） | 是 |
| 查看商家主页及其广告 | 是 |
| 查询支持的数字货币和法币 | 是 |

> **不支持**: 现货/合约价格、充提操作、发送聊天消息、发起新申诉（补充证据除外）。

---

## payment（扫码支付）

**路径**: `skills/binance/payment/`  
**版本**: 2.0.0 | **作者**: Binance  
**依赖**: `payment_skill.py`（本地 Python 脚本）

处理 Binance Pay QR 码支付及收款码生成。

| 功能 | 说明 |
|------|------|
| 解析 QR 码 | 从图片/剪贴板/链接文本提取支付数据（支持 C2C 和 PIX） |
| 发起支付 | 根据 QR 数据完成付款确认与提交 |
| 取消支付 | 取消进行中的支付订单 |
| 查询订单状态 | 查看支付结果 |
| 生成收款码 | 生成 QR 码或支付链接用于收取加密货币 |

**QR 处理优先级**: 直接视觉识别 → 图片路径解码 → 请求用户复制到剪贴板（禁止自动使用剪贴板）。

---

## square-post（广场发帖）

**路径**: `skills/binance/square-post/`  
**版本**: 1.2 | **作者**: binance-square

向 Binance Square（Binance 社交平台）发布交易观点内容。

| 功能 | 说明 |
|------|------|
| 发布纯文字帖子 | 通过 OpenAPI Key 向 Binance Square 发布文本内容 |

**触发词**: `post to square`、`square post`。

---

## binance-agentic-wallet（Web3 钱包）

**路径**: `skills/binance-web3/binance-agentic-wallet/`  
**版本**: 1.0.1 | **作者**: binance-web3-team  
**依赖**: `baw` CLI（npm 包 `@binance/agentic-wallet`）

驱动 `baw` CLI 管理 Binance Web3 钱包，支持链上全套操作。

| 操作类别 | 具体功能 |
|----------|----------|
| 账户管理 | 登录/登出、查看钱包状态、获取地址 |
| 资产查询 | 查看代币余额、交易历史、每日剩余限额、待处理交易锁 |
| 网络信息 | 列出支持的区块链网络 |
| 安全设置 | 查看/配置钱包安全参数（每日限额、滑点、MEV 保护） |
| 转账 | 向外部地址发送代币 |
| 市价交易 | DEX swap 市价兑换、获取报价（不执行）、查询市价订单 |
| 限价订单 | 设置目标价买入/卖出、查询/取消限价单 |

**安全规则**: 永不记录或展示 session token、API Key、私钥、助记词；每次状态变更操作前必须获得用户明确确认。

---

## binance-tokenized-securities-info（代币化证券）

**路径**: `skills/binance-web3/binance-tokenized-securities-info/`  
**版本**: 1.1 | **作者**: binance-web3-team

查询 Ondo 在 Binance Web3 上的代币化美股数据。支持 Ethereum（chainId=1）和 BSC（chainId=56）。

| API | 功能 |
|-----|------|
| Token Symbol List | 列出所有代币化股票（代码、链、合约地址） |
| RWA Meta | 获取公司元数据（CEO、行业、概念标签、证明报告） |
| Market Status | 查询 Ondo 市场整体开/收市状态 |
| Asset Market Status | 查询单个资产交易状态（含停牌原因：财报、分红、拆股、合并） |
| RWA Dynamic V2 | 实时数据（链上价格、持有者数、市值、美股 P/E、52 周区间、下单限额） |
| Token K-Line | K 线/蜡烛图 OHLC 数据（技术分析） |

> **重要**: 每个代币代表 `multiplier` 股而非 1 股，参考价 = `tokenPrice ÷ sharesMultiplier`。

---

## crypto-market-rank（加密市场排行）

**路径**: `skills/binance-web3/crypto-market-rank/`  
**版本**: 2.1 | **作者**: binance-web3-team

加密货币市场排行榜与领袖榜，支持 BSC、Base、Solana。

| API | 功能 |
|-----|------|
| Social Hype Leaderboard | 社交热度排行（情绪分析、舆情摘要） |
| Unified Token Rank | 多类型代币排行（热门/搜索榜/Alpha/代币化股票，可按链过滤） |
| Smart Money Inflow Rank | 聪明钱净流入代币排行（发现机构/大户在买什么） |
| Meme Rank | Pulse 发射台 Meme 代币排行（最可能突破的 Meme） |
| Address PnL Rank | 顶级交易者 PnL 排行榜（胜率、盈亏、持仓） |

**排行类型** (`rankType`): 10=热门, 11=搜索榜, 20=Alpha, 40=代币化股票。

---

## meme-rush（Meme 代币追踪）

**路径**: `skills/binance-web3/meme-rush/`  
**版本**: 1.1 | **作者**: binance-web3-team

实时跟踪 Meme 代币发射台（Pump.fun、Four.meme 等）及 AI 热点叙事，支持 BSC 和 Solana。

| API | 功能 |
|-----|------|
| Meme Rush Rank List | 按阶段（新发布/即将迁移/已迁移 DEX）列出 Meme 代币 |
| Topic Rush | AI 生成市场热点话题及关联代币（按净流入排序） |

**高级过滤**: 开发者卖出占比、Top10 持仓%、狙击/内部人/捆绑者持仓%、绑定曲线进度、流动性、交易量。  
**支持协议**: Pump.fun、Moonit、Four.meme、Raydium、Orca、Meteora 等 18 个平台。

---

## query-address-info（地址持仓查询）

**路径**: `skills/binance-web3/query-address-info/`  
**版本**: 1.1 | **作者**: binance-web3-team

查询任意链上钱包地址的代币持仓（**无需认证**，公开 API）。

| 功能 | 说明 |
|------|------|
| 持仓列表 | 列出指定地址在指定链上持有的所有代币 |
| 实时价格 | 每个代币的当前价格 |
| 24h 涨跌幅 | 每个代币的 24 小时价格变化百分比 |
| 持仓数量 | 各代币的持有数量 |

**支持链**: BSC（56）、Base（8453），可通过 `chainId` 参数切换。

---

## query-token-audit（代币安全审计）

**路径**: `skills/binance-web3/query-token-audit/`  
**版本**: 1.4 | **作者**: binance-web3-team

交易前安全检查，检测貔貅盘（Honeypot）、Rug Pull 及恶意合约（**无需认证**）。

| 检测维度 | 说明 |
|----------|------|
| 合约风险 | 危险所有权函数、隐藏后门 |
| 交易风险 | 貔貅检测（只能买不能卖）、异常买卖税 |
| 诈骗检测 | 假代币识别、Rug Pull 风险标记 |
| 综合评分 | 风险等级与详细说明 |

**支持链**: BSC（56）、Base（8453）、Solana（CT_501）、Ethereum（1）。

---

## query-token-info（代币信息查询）

**路径**: `skills/binance-web3/query-token-info/`  
**版本**: 1.1 | **作者**: binance-web3-team

通过关键词、合约地址或链搜索代币，获取元数据与实时行情（**无需认证**）。

| API | 功能 |
|-----|------|
| Token Search | 按名称、符号或合约地址搜索代币 |
| Token Metadata | 静态信息（名称、符号、Logo、社交链接、创建者地址） |
| Token Dynamic Data | 实时市场数据（价格、成交量、持有者数、流动性、市值） |
| Token K-Line | K 线蜡烛图 OHLCV 数据（技术分析） |

**支持链**: BSC（56）、Base（8453）、Solana（CT_501）。

---

## trading-signal（智能资金信号）

**路径**: `skills/binance-web3/trading-signal/`  
**版本**: 1.1 | **作者**: binance-web3-team

订阅和获取链上聪明钱交易信号，辅助投资决策（**无需认证**）。

| 信号字段 | 说明 |
|----------|------|
| 信号类型 | 买入/卖出信号 |
| 触发价格 | 聪明钱建仓时的价格 |
| 当前价格 | 实时价格（可与触发价对比） |
| 最大涨幅 | 信号触发后的最大盈利空间 |
| 退出率 | 聪明钱的已退出比例 |
| 代币标签 | Pumpfun、DEX Paid 等标签 |

**支持链**: BSC（56）、Solana（CT_501）。

---

## 快速对照表

| Skill | 类别 | 是否需要认证 | 核心能力 |
|-------|------|------------|----------|
| `credential-guard` | 安全 | — | 强制 .env 加载，对话脱敏，前置于所有认证 skill |
| `binance` | CEX | 是（binance-cli） | 现货/合约/钱包/理财等全套 CEX 操作 |
| `fiat` | CEX | 部分 | 法币支付方式查询与汇率 |
| `onchain-pay` | CEX+链上 | 是（RSA 签名） | 法币→加密→链上转账一体化 |
| `p2p` | CEX | 部分 | P2P 广告查询与订单管理 |
| `payment` | CEX | 是（脚本） | QR 码扫码支付与收款码生成 |
| `square-post` | CEX | 是（OpenAPI Key） | Binance Square 发帖 |
| `binance-agentic-wallet` | Web3 | 是（baw CLI） | Web3 钱包全套链上操作 |
| `binance-tokenized-securities-info` | Web3 | 否 | 代币化美股实时数据与技术分析 |
| `crypto-market-rank` | Web3 | 否 | 代币排行榜与聪明钱流向 |
| `meme-rush` | Web3 | 否 | Meme 代币实时追踪与热点叙事 |
| `query-address-info` | Web3 | 否 | 任意地址链上持仓查询 |
| `query-token-audit` | Web3 | 否 | 代币合约安全审计 |
| `query-token-info` | Web3 | 否 | 代币搜索、元数据与实时行情 |
| `trading-signal` | Web3 | 否 | 聪明钱链上交易信号 |
