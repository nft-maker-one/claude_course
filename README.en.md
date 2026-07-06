# Claude Course

A practical tutorial repository for system-level development with Claude Code.

## Follow Us

This tutorial is continuously updated. Follow our channel for the latest content — likes, bookmarks, and shares are welcome!

<table>
  <tr>
    <td align="center"><b>Bilibili · 狂飙web3</b></td>
    <td align="center"><b>WeChat Official Account · 狂飙web3</b></td>
  </tr>
  <tr>
    <td><img src="bilibili.png" width="280"/></td>
    <td><img src="wx_account.png" width="280"/></td>
  </tr>
</table>

## Project Overview

This repository is a comprehensive practical tutorial series for Claude Code, covering a complete knowledge system from installation and configuration to advanced extensions. Organized in a modular structure, it includes 13 thematic chapters, each focusing on a distinct core functionality domain.

## Directory Structure

| Chapter | Topic | Core Content |
|---------|-------|--------------|
| day1 | Installation & Interface | Claude CLI installation, VS Code extension, permission modes, model switching |
| day2 | Long-term Memory | JSONL persistence, Git anchors, Auto Memory, Rewind and Fork |
| day3 | Short-term Memory | Context window management, parentUuid linked list, Compact mechanism |
| day4 | Permission Management | Permission tiers, dontAsk mode, workspace boundaries |
| day5 | Subagent | Context physical isolation, concurrent Fan-out/Fan-in, permission delegation |
| day6 | Skills | Skill system development, Binance Skills Hub, MCP integration |
| day7 | MCP Protocol | Model Context Protocol, scope isolation, transport layer implementation |
| day8 | Hooks | Interception mechanisms, PreToolUse/PostToolUse/Stop, DAG state machine |
| day9 | Channels | Remote control architecture, Telegram integration |
| day10 | Scheduled Tasks | Local polling, cloud scheduling, MCP mounting |
| day11 | Caching Engineering | Prompt caching, context optimization |
| day12 | Task Scheduling | Goal mode, Workflows mode, DAG orchestration |
| day13 | Code Review | Security mode, Plugin development, Marketplace distribution |

## Main Content

### Core Tutorial Documents

- **day1_install/install_intro.md** - Deep dive into downloading, installing, and interface tools of Claude Code
- **day2_longterm_memory/longterm_memory.md** - Git-level memory engine and state forking
- **day3_shortterm_memory/context_window_memory.md** - Short-term memory linked list reconstruction and Compact mechanism
- **day4_permission/permission_claude.md** - Permission management and dontAsk paradigm
- **day5_subagent/subagent.md** - Subagent context physical isolation mechanism
- **day6_skill/skill_usage.md** - From beginner to advanced Skills system development
- **day7_mcp/model_context_protocol.md** - MCP protocol reshaping Agent perception neural networks
- **day8_hook/hook.md** - Using Hooks as harness reins to control highly autonomous AI
- **day9_channel/channel.md** - Channels remote control architecture and underlying implementation
- **day10_schedule/schedule.md** - Scheduled tasks and cloud scheduling mechanisms
- **day11_cache/cache.md** - Prompt caching mechanisms and optimization strategies
- **day12_task_schedule/task_schedule.md** - Goal and Workflow task orchestration
- **day13_code_review/plugin_marketplace_guide.md** - Deep analysis of Plugins and Marketplace

### Practical Projects

#### Blockchain Tutorial
- `day4_permission/blockchain_tutorial/` - React + Vite blockchain visualization tutorial

#### Binance Skills Hub
- `day6_skill/binance-skills-hub/` - Complete Binance trading skill suite
  - `binance` - Spot trading
  - `fiat` - Fiat payment
  - `p2p` - P2P trading
  - `payment` - QR code payment
  - `onchain-pay` - On-chain payment
  - `binance-web3` - Web3 wallet and on-chain operations
  - `crypto-market-rank` - Cryptocurrency market ranking
  - `meme-rush` - Meme token tracking
  - `trading-signal` - Smart money signals

#### MCP Services
- `day7_mcp/image_mcp/server.py` - Image processing MCP service
- `day7_mcp/math_mcp/server.py` - Mathematical computation MCP service

#### Utility Tools
- `day5_subagent/agents/` - Various Subagent implementations
  - `crypto_price.py` - Cryptocurrency price query
  - `reddit_search.py` - Reddit search
  - `twitter_search.py` - Twitter search
  - `token-analyze.md` - Token deep analysis

- `day11_cache/` - Caching and monitoring
  - `smart_money_dashboard.html` - Solana smart money monitoring dashboard
  - `supabase/functions/` - Cloud functions

- `day12_task_schedule/` - Task scheduling tools
  - `todo-scanner.js` - TODO scanner
  - `file_subfix.py` - File extension scanner

## Quick Start

### Prerequisites

- Node.js 18+
- Python 3.8+
- Go 1.20+

### Install Claude Code

Refer to `day1_install/install_intro.md` for detailed installation instructions.

### Run Examples

```bash
# Clone the repository
git clone https://gitee.com/jerry_luo_03/claude_course.git
cd claude_course

# View blockchain tutorial
cd day4_permission/blockchain_tutorial
npm install
npm run dev

# Run MCP service
cd day7_mcp/image_mcp
python server.py
```

## Learning Path Recommendations

1. **Beginner Stage** (day1-day4)
   - Master installation, configuration, permission management, and core concepts

2. **Intermediate Stage** (day5-day9)
   - Learn Subagent, Skills, MCP, Hooks, and Channels

3. **Advanced Stage** (day10-day13)
   - Master scheduled tasks, caching optimization, task orchestration, and code review

## License

MIT License — see the LICENSE file in each subproject directory.