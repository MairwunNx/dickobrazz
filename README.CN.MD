# Dickobrazz

[![AI Capable](https://img.shields.io/badge/AI-Capable-brightgreen?style=flat&logo=openai&logoColor=white)](https://github.com/mairwunnx/dickobrazz)
[![Build Status](https://img.shields.io/github/actions/workflow/status/mairwunnx/dickobrazz/build.yml?branch=master&style=flat&logo=githubactions&logoColor=white)](https://github.com/mairwunnx/dickobrazz/actions/workflows/build.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/mairwunnx/dickobrazz)](https://goreportcard.com/report/github.com/mairwunnx/dickobrazz)
[![Runs on Telegram](https://img.shields.io/badge/Runs_on-Telegram-26A5E4?style=flat&logo=telegram&logoColor=white)](https://t.me/dickobrazz_bot)
[![Docker](https://img.shields.io/badge/Docker-Available-2496ED?style=flat&logo=docker&logoColor=white)](https://github.com/MairwunNx/dickobrazz/pkgs/container/dickobrazz)
[![GitHub Release](https://img.shields.io/github/v/release/mairwunnx/dickobrazz?style=flat&logo=github&color=blue)](https://github.com/mairwunnx/dickobrazz/releases)
[![Go Version](https://img.shields.io/badge/Go-1.26-00ADD8?style=flat&logo=go&logoColor=white)](https://go.dev/)

**语言**: 🇨🇳 中文 | [🇷🇺 Русский](README.MD) | [🇺🇸 English](README.EN.MD)

🌶️ Dickobrazz Telegram 机器人可以精确到厘米给出你的“尺寸”。如果今天数值不大，也别介意，这只是个玩笑。

> 现代化高科技测量器：你提问，机器人会给出“科学依据”的尺寸，还会用玩笑方式映射到俄罗斯地区编号。不再需要尺子或显微镜。

> [!NOTE]
> 产品灵感来源于《南方公园》第 1504 集（第 213 集）T.M.I.

## 功能特性

### 尺寸测量

机器人每天给出一个 0 到 61 cm 的唯一数值。数值在莫斯科时间午夜重置，附带季节表情，并映射到俄罗斯地区编号。

### 排行与竞赛

- **尺寸排行榜** - 当日前 13 名最佳尺寸
- **永久天梯** - 按历史累计尺寸的全时排名
- **赛季竞速** - 3 个月赛季累计尺寸竞赛
- **赛季系统** - 每 3 个月开启新赛季并发放奖励
- **Cock-Respect™** - 赛季胜利与成就完成可获得积分

### 成就系统

50+ 个独特成就：从首次测量到稀有数字组合。每个成就都会奖励 respect，并显示进度。

### 高级分析

- **CRI（Cock Rating Index）** - 个人评分（类似 South Park 的 T.M.I.）
- **详细统计** - 平均值、中位数、纪录与变化趋势
- **运气系数**、**波动率**、**增长速度**
- **主导占比** - 用户对整体系统统计的贡献比例

## 使用方法

Dickobrazz 通过 Telegram **inline queries（内联查询）** 工作。只需在任意聊天输入 `@dickobrazz_bot` 并选择对应选项：

### 可用内联命令

- **尺寸测量** - 获取今日尺寸并进行地区映射
- **尺寸排行榜** - 当日前 13 名玩家（每日排名）
- **永久天梯** - 历史累计尺寸全时排名
- **赛季竞速** - 当前赛季排名（3 个月）
- **尺寸动态** - 详细个人统计与分析
- **赛季历史** - 历史赛季与冠军导航
- **成就系统** - 成就与进度分页

### 工作方式

1. 在 Telegram 打开任意聊天
2. 输入 `@dickobrazz_bot`
3. 从列表中选择需要的选项
4. 发送到聊天或私下查看

## 运行与部署 🚀

### 开始前 📦

如果仓库尚未克隆：
```bash
git clone https://github.com/mairwunnx/dickobrazz.git
```

---

### Dev Containers（推荐用于开发和本地 dev 运行） 🧰

**要求**:
- 已运行 **Docker**
- 开发容器至少 **1 GB RAM**（VS Code/Cursor/Goland）
- 支持 dev containers 的工具，或可启动 dev container 的 CLI 工具

> [!IMPORTANT]
> 启动 dev container 前，请先在 `.devcontainer/devcontainer.env` 中设置环境变量。

**环境变量**:
- 文件: `.devcontainer/devcontainer.env`
- 最低必需:
```bash
BOT_CSOT="your_internal_csot_token" # 从后端获得的可信字符串。
BOT_SERVER_BASE_URL="https://api.dickobrazz.com" # 后端基础 URL。
BOT_TG_TOKEN="your_telegram_bot_token" # Telegram 机器人 Token。
BOT_TG_ENV="test" # Telegram 机器人环境（test/production）。
```

或者编辑配置文件 `./config.yaml`。

**打开方式** _(VS Code/Cursor/Goland)_:
1. 在 VS Code/Cursor/Goland 中打开项目
2. 找到并运行 **Reopen in Container**

支持 dev containers 的工具应自动构建并启动开发容器。

**在 Dev Container 中运行机器人**:

> [!NOTE]
> 在 dev container 中进行开发时，可使用标准 Go 工具链或支持热重载的 air。_(容器内已包含所需工具)_

| 工具 | 命令 | 说明 |
|---------|---------|---------|
| air | `air` | 热重载运行机器人 |
| go | `go run program.go` | 常规运行机器人 |

---

### 通过 Docker Compose 运行 🐳

**环境变量**:
- 文件: `.devcontainer/devcontainer.env`
- 最低必需:
```bash
BOT_CSOT="your_internal_csot_token" # 从后端获得的可信字符串。
BOT_SERVER_BASE_URL="https://api.dickobrazz.com" # 后端基础 URL。
BOT_TG_TOKEN="your_telegram_bot_token" # Telegram 机器人 Token。
BOT_TG_ENV="test" # Telegram 机器人环境（test/production）。
```

**启动**:
```bash
docker compose up -d
```

---

### 使用现成镜像 🧩（production）

可直接使用 GitHub Container Registry 的现成镜像：

```yaml
services:
  dickobrazz:
    image: ghcr.io/mairwunnx/dickobrazz:latest
    env_file: .env
```

> [!IMPORTANT]
> 如果需要修改配置文件，可使用新配置重新构建镜像，或将本地配置挂载到容器中（推荐第一种方式）。

---

## 技术栈

- **Go 1.26** - 核心开发语言
- **Telegram Bot API (tgbotapi/v5)** - Telegram 集成
- **Resty** - 与后端交互的 HTTP 客户端
- **uber/fx** - 依赖注入与生命周期管理
- **Prometheus** - 指标与监控
- **go-i18n** - 国际化（6 种语言）
- **Docker + Docker Compose** - 容器化与编排
- **Structured Logging (slog)** - 带上下文的 JSON 日志

## 同生态项目链接

[dickobrazz-server](https://github.com/MairwunNx/dickobrazz-server) - 🌶️ @dickobrazz_bot 的后端服务：用户数据处理、个人数据聚合、统计生成。

## 相关项目链接

[Emperor Xi](https://github.com/mairwunnx/xi) - 🀄️ AI 驱动的 Telegram 机器人，风格为 Xi 皇帝。伟大领袖的私人助手，随时回答普通民众的问题。

---

<img src="./media.png" alt="Russian power" width="500">

🇷🇺 **Made in Russia with love.** ❤️

**Dickobrazz** 代表真实随机与健康竞争。没有作弊，只有数学和运气。

> 🫡 Made by Pavel Erokhin (Павел Ерохин), aka mairwunnx.
