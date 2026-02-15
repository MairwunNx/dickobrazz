# Dickobrazz

[![AI Capable](https://img.shields.io/badge/AI-Capable-brightgreen?style=flat&logo=openai&logoColor=white)](https://github.com/mairwunnx/dickobrazz)
[![Docker](https://img.shields.io/badge/Docker-Available-2496ED?style=flat&logo=docker&logoColor=white)](https://github.com/MairwunNx/dickobrazz/pkgs/container/dickobrazz)
[![GitHub Release](https://img.shields.io/github/v/release/mairwunnx/dickobrazz?style=flat&logo=github&color=blue)](https://github.com/mairwunnx/dickobrazz/releases)
[![Go Version](https://img.shields.io/badge/Go-1.26-00ADD8?style=flat&logo=go&logoColor=white)](https://go.dev/)

**语言**: 🇨🇳 中文 | [🇷🇺 Русский](README.MD) | [🇺🇸 English](README.EN.MD)

🌶️ Dickobrazz 机器人（也叫 dickobot）可以精确到厘米给出你的“尺寸”。如果今天数值不大，也别介意，这只是个玩笑。

> 现代化高科技测量器：你提问，机器人会给出“科学依据”的尺寸，还会用玩笑方式映射到俄罗斯地区编号。不再需要尺子或显微镜。

> 产品灵感来源于《南方公园》第 1504 集（第 213 集）T.M.I.

## 功能特性

### 尺寸测量

机器人每天给出一个 0 到 61 cm 的唯一数值。数值在莫斯科时间午夜重置，附带季节表情，并映射到俄罗斯地区编号。

### 排行与竞赛

- **尺寸排行榜** — 当日前 13 名最佳尺寸
- **永久天梯** — 按历史累计尺寸的全时排名
- **赛季竞速** — 3 个月赛季累计尺寸竞赛
- **赛季系统** — 每 3 个月开启新赛季并发放奖励
- **Cock-Respect™** — 赛季胜利与成就完成可获得积分

### 成就系统

50+ 个独特成就：从首次测量到稀有数字组合。每个成就都会奖励 respect，并显示进度。

### 高级分析

- **CRI（Cock Rating Index）** — 个人评分（类似 South Park 的 T.M.I.）
- **详细统计** — 平均值、中位数、纪录与变化趋势
- **运气系数**、**波动率**、**增长速度**
- **主导占比** — 用户对整体系统统计的贡献比例

## 使用方法

Dickobrazz 通过 Telegram **inline queries（内联查询）** 工作。只需在任意聊天输入 `@dickobrazz_bot` 并选择对应选项：

### 可用内联命令

- **尺寸测量** — 获取今日尺寸并进行地区映射
- **尺寸排行榜** — 当日前 13 名玩家（每日排名）
- **永久天梯** — 历史累计尺寸全时排名
- **赛季竞速** — 当前赛季排名（3 个月）
- **尺寸动态** — 详细个人统计与分析
- **赛季历史** — 历史赛季与冠军导航
- **成就系统** — 成就与进度分页

### 工作方式

1. 在 Telegram 打开任意聊天
2. 输入 `@dickobrazz_bot`
3. 从列表中选择需要的选项
4. 发送到聊天或私下查看

## 运行与部署 🚀

### 开始前 📦

如果仓库尚未克隆：
```bash
git clone https://github.com/mairwunnx/dickobrazz
cd dickobrazz
```

---

### Dev Containers（推荐开发） 🧰

**要求**:
- 容器至少 **1 GB RAM**
- 已运行 **Docker**
- **VS Code** 或 **Cursor**

**打开方式**:
1. 在 VS Code/Cursor 中打开项目
2. 点击 **Reopen in Container**

**环境变量**:
- 文件: `.devcontainer/devcontainer.env`
- 最低必需:
```bash
BOT_CSOT="your_internal_csot_token"
BOT_SERVER_BASE_URL="https://api.dickobrazz.com"
BOT_TG_TOKEN="your_telegram_bot_token"
BOT_TG_ENV="test"
```

**运行开发机器人**（使用 `.devcontainer/devcontainer.env`）:
```bash
air
```
`air` 会在源码变更时自动重启应用。
或
```bash
go run program.go
```

---

### 通过 Docker Compose 运行生产环境 🐳

**环境变量**:
- 文件: `.env`
- 最低必需:
```bash
BOT_CSOT="your_internal_csot_token"
BOT_SERVER_BASE_URL="https://api.dickobrazz.com"
BOT_TG_TOKEN="your_telegram_bot_token"
BOT_TG_ENV="production"
```

**启动**:
```bash
docker compose up -d
```

---

### 现成镜像 + 自有服务 🧩

可直接使用 GitHub Container Registry 的现成镜像：

```yaml
services:
  dickobrazz:
    image: ghcr.io/mairwunnx/dickobrazz:latest
    env_file: .env
```

---

### 从源码构建 🛠️

**要求**:
- Go `1.26+`

```bash
go build -o dickobrazz program.go
```

---

### 通过 Docker 构建 🧱

**要求**:
- Docker

```bash
docker build -t dickobrazz .
```

## 技术栈

- **Go 1.26** — 核心开发语言
- **Telegram Bot API (tgbotapi/v5)** — Telegram 集成
- **Resty** — 与后端交互的 HTTP 客户端
- **uber/fx** — 依赖注入与生命周期管理
- **Prometheus** — 指标与监控
- **go-i18n** — 国际化（6 种语言）
- **Docker + Docker Compose** — 容器化与编排
- **Structured Logging (slog)** — 带上下文的 JSON 日志

## AI 参与

AI 用于创建部分文档和提交，也用于优化部分 MongoDB 聚合流水线和统计相关算法。

另外，截至 2025 年 9 月 23 日，`.cursor/rules/` 目录中已添加 AI 提示（prompts/rules），可用于未来通过 LLM agents 实现新功能。

## 相关项目链接

[Xi Manager](https://github.com/mairwunnx/xi) — 🀄️ AI 驱动的 Telegram 机器人，风格为 Xi 的私人助理。伟大领袖的私人助手，随时回答普通民众的问题。

[Louisepizdon](https://github.com/MairwunNx/louisepizdon) — 🥀 Louisepizdon，一个比你奶奶还诚实的 AI Telegram 机器人。会“认真”评价你，并从照片分析你的穿搭定价。

## “同一作者”系列

[Mo'Bosses](https://github.com/mairwunnx/mobosses) — 🏆 **Mo'Bosses** 是顶级 RPG 插件，可将普通怪物变成史诗 Boss，并提供**高级玩家成长系统**。与其他插件不同，这里每场战斗都很重要，每次升级都能解锁新可能。⚔

[Mo'Joins](https://github.com/mairwunnx/mojoins) — 🎉 自定义进出服效果：消息、音效、粒子、烟花与进服后保护。适用于 PaperMC。

[Mo'Afks](https://github.com/mairwunnx/moafks) — 🛡️ 在线时间“暂停”成为可能。适用于 PaperMC 的 AFK 安全模式：免伤、无碰撞、怪物忽略、自动挂机检测与精致视觉效果。

[McBuddy Server](https://github.com/mcbuddy-ai/mcbuddy-server) — 🛠️⚡ MCBuddy AI 助手后端，集成 OpenRouter 与请求处理

[McBuddy Telegram](https://github.com/mcbuddy-ai/mcbuddy-bot) — 🤖📱 在游戏外与 MCBuddy 对话的 Telegram 机器人

[McBuddy Spigot](https://github.com/mcbuddy-ai/mcbuddy-spigot) — 💬 MCBuddy 集成用 Spigot 插件 —— 添加 `/ask` 命令，可直接在 Minecraft 服务器聊天中向 AI 助手提问。🎮

---

![image](./media.png)

🇷🇺 **Made in Russia with love.** ❤️

**Dickobrazz** 代表真实随机与健康竞争。没有作弊，只有数学和运气。

> 🫡 Made by Pavel Erokhin (Павел Ерохин), aka mairwunnx.
