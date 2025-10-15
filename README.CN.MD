# Dickobrazz

[![AI Capable](https://img.shields.io/badge/AI-Capable-brightgreen?style=flat&logo=openai&logoColor=white)](https://github.com/mairwunnx/dickobrazz)
[![Docker](https://img.shields.io/badge/Docker-Available-2496ED?style=flat&logo=docker&logoColor=white)](https://github.com/MairwunNx/dickobrazz/pkgs/container/dickobrazz)
[![GitHub Release](https://img.shields.io/github/v/release/mairwunnx/dickobrazz?style=flat&logo=github&color=blue)](https://github.com/mairwunnx/dickobrazz/releases)
[![Go Version](https://img.shields.io/badge/Go-1.25-00ADD8?style=flat&logo=go&logoColor=white)](https://go.dev/)
[![MongoDB](https://img.shields.io/badge/MongoDB-8.0.15-47A248?style=flat&logo=mongodb&logoColor=white)](https://www.mongodb.com/)
[![Redis](https://img.shields.io/badge/Redis-8.2-DC382D?style=flat&logo=redis&logoColor=white)](https://redis.io/)

**语言**: 🇨🇳 中文 | [🇷🇺 Русский](README.MD) | [🇺🇸 English](README.EN.MD)

🌶️ Dickobrazz 机器人，也叫 dickobot，能够精确到厘米地测量你的"尺寸"，如果没能炫耀大尺寸请别介意，这只是个玩笑。

> 现代化高科技的测量器：你问，机器人就给出科学合理的尺寸，甚至以幽默的方式将你的尺寸与俄罗斯地区编号进行对比。不再需要尺子或显微镜！

> 产品灵感来源于《南方公园》第 1504 集（第 213 集）(T.M.I.)。

## 功能特性

### 核心功能

- **尺寸测量** — 生成0到61*厘米的尺寸，配有季节性表情符号。
- **每日限制** — 每天一次测量，莫斯科时间午夜重置。
- **季节性表情** — 根据季节变化的不同表情符号。
- **GSK** — 与全球系统集成。（全球系统，不仅限于 telegram）

> * - 这是一个变成功能的错误。原计划最高60，但结果是61 :)。

### 游戏化和竞赛

- **排行榜** — 当日前13名最佳尺寸。
- **竞赛** — 季节性比赛，累计尺寸总和。
- **季节系统** — 每3个月一个新赛季，产生获胜者。
- **尊重值™** — 参与赛季的积分系统。

### 高级分析

- **IRC（尺寸指数）** — 相对于所有用户的个人排名（类似《南方公园》中的 T.M.I）。
- **动态统计** — 增长统计、记录、优势百分比。
- **详细统计** — 中位数、平均值、尺寸分布。
- **资产** — 变化动态和日增长。

### 技术可靠性

- **双重随机化** — Random.org + ChaCha8后备方案，确保最大随机性。
- **高性能数据库** — MongoDB 聚合管道。
- **缓存系统** — Redis 配合 TinyLFU 快速响应。

## 使用方法

Dickobrazz 通过 Telegram 的 **内联查询** 工作。只需在任何聊天中输入 `@dickobrazz_bot` 并选择所需选项：

### 可用的内联命令

- **尺寸测量** — 获取今日尺寸
- **排行榜** — 当日玩家排行榜  
- **竞赛** — 当前赛季排名
- **全球排行榜** — 全球玩家累计尺寸排行榜
- **动态统计** — 详细个人统计
- **赛季历史** — 赛季历史和获胜者

### 工作原理

1. 在 Telegram 中打开任何聊天
2. 输入 `@dickobrazz_bot` 
3. 从列表中选择感兴趣的选项
4. 将结果发送到聊天或私下查看

## 部署

### Dev Containers（推荐用于开发）

开发的最简单方式 — 在 VS Code 或类似 IDE 中使用 Dev Containers：

1. 克隆仓库：
```bash
git clone https://github.com/mairwunnx/dickobrazz
cd dickobrazz
```

2. 只需配置两个必需的令牌（在 `.env.dev` 文件中）：
```bash
TELEGRAM_TOKEN="你的 telegram_bot_token"
RANDOMORG_TOKEN="你的 random_org_api_key"
```

3. 在 VS Code 中打开项目并选择 "Reopen in Container"

其余所有内容（MongoDB、Redis、依赖项）将自动配置。简单快速！

### Docker Compose（生产环境）

1. 克隆仓库：
```bash
git clone https://github.com/mairwunnx/dickobrazz
cd dickobrazz
```

2. 配置两个必需的令牌（在 `.env` 文件中）：
```bash
TELEGRAM_TOKEN="你的 telegram_bot_token"
RANDOMORG_TOKEN="你的 random_org_api_key"
```

> **注意**：如需要，请配置其余环境变量（Redis 和 MongoDB 密码）。

3. 启动服务：
```bash
docker compose up -d
```

### Docker Compose 使用现成镜像

如果想使用现成镜像并自行管理 MongoDB/Redis：

1. 在 `docker-compose.yml` 中声明 dickobrazz 服务：

```yaml
services:
  dickobrazz:
      image: ghcr.io/mairwunnx/dickobrazz:1.23.1
    env_file: .env
    environment:
      DD_DB_HOST: dickobrazz-mongo
      MONGO_INITDB_DATABASE: ${MONGO_INITDB_DATABASE}
      MONGO_INITDB_ROOT_PASSWORD: ${MONGO_INITDB_ROOT_PASSWORD}
      MONGO_INITDB_ROOT_USERNAME: ${MONGO_INITDB_ROOT_USERNAME}
      MONGODB_URI: ${MONGODB_URI}
      RANDOMORG_TOKEN: ${RANDOMORG_TOKEN}
      REDIS_ADDRESS: ${REDIS_ADDRESS}
      REDIS_PASSWORD: ${REDIS_PASSWORD}
      TELEGRAM_TOKEN: ${TELEGRAM_TOKEN}
      TZ: ${TZ}
```

> **注意**：需要传递 `.env` 文件中的所有环境变量或通过其他方便的方式。

2. 在 `docker-compose.yml` 中声明 MongoDB 和 Redis 服务。
> 重要：**MongoDB** 版本必须是 `8.2+`，**Redis** 版本必须是 `8.2+`。

3. 启动服务：
```bash
docker compose up -d
```

### 从源码构建

#### 手动直接构建方式

要求：
- Go 1.25+

```bash
go build -o dickobrazz program.go
```

#### 使用 Docker 构建

要求：
- Docker

```bash
docker build -t dickobrazz .
```

## 技术栈

- **Go 1.25** — 主要开发语言
- **Telegram Bot API (tgbotapi/v5)** — Telegram 集成
- **MongoDB + Mongo Driver** — 主数据库，支持聚合
- **Redis + go-redis/cache** — TinyLFU 缓存
- **Random.org API** — 真随机数源
- **ChaCha8** — 加密安全的 PRNG 后备方案
- **Docker + Docker Compose** — 容器化和编排
- **Structured Logging (slog)** — 带上下文的 JSON 日志

## AI参与

AI用于创建部分文档和提交，以及优化一些MongoDB聚合管道和改进统计算法。

截至2025年9月23日，在 `.cursor/rules/` 文件夹中添加了AI提示（prompts/rules），可能有助于将来使用LLM代理实现新功能。

## 相关项目链接

[Xi Manager](https://github.com/mairwunnx/xi) — 🀄️ AI 驱动的 Telegram 机器人，风格化为个人助理Xi。伟大领袖的个人助理，准备回答普通民众的问题。

[Louisepizdon](https://github.com/MairwunNx/louisepizdon) — 🥀 Louisepizdon，比你奶奶还诚实的 AI Telegram 机器人。会按价值评价你，从照片分析你的服装定价！

## 来自同一作者的系列

[Mo'Bosses](https://github.com/mairwunnx/mobosses) — 🏆 **Mo'Bosses** — 最佳 RPG 插件，将普通怪物变成史诗级boss，具有 **高级玩家进展系统**。与其他插件不同，这里每场战斗都有意义，每个等级都开启新可能！ ⚔

[Mo'Joins](https://github.com/mairwunnx/mojoins) — 🎉 自定义进入/退出：消息、声音、粒子、烟花和进入后保护。全部适用于 PaperMC。

[Mo'Afks](https://github.com/mairwunnx/moafks) — 🛡️ 在线暂停现在成为可能。PaperMC 插件，为玩家提供安全的AFK模式：伤害免疫、无碰撞、怪物忽略、自动检测不活跃和精美视觉效果。

[McBuddy Server](https://github.com/mcbuddy-ai/mcbuddy-server) — 🛠️⚡ MCBuddy AI助手后端，集成 OpenRouter 和请求处理

[McBuddy Telegram](https://github.com/mcbuddy-ai/mcbuddy-bot) — 🤖📱 在游戏外与 MCBuddy 交流的 Telegram 机器人

[McBuddy Spigot](https://github.com/mcbuddy-ai/mcbuddy-spigot) — 💬 MCBuddy集成的Spigot插件 — 添加 `/ask` 命令，直接在 Minecraft 服务器聊天中向AI助手提问！🎮

---

![image](./media.png)

🇷🇺 **俄罗斯制造，用心打造。** ❤️

**Dickobrazz** — 关于诚实随机和健康竞争。没有欺骗，只有数学和运气。

> 🫡 Made by Pavel Erokhin (Павел Ерохин), aka mairwunnx.
