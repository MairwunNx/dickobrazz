# Dickobrazz

[![AI Capable](https://img.shields.io/badge/AI-Capable-brightgreen?style=flat&logo=openai&logoColor=white)](https://github.com/mairwunnx/dickobrazz)
[![Docker](https://img.shields.io/badge/Docker-Available-2496ED?style=flat&logo=docker&logoColor=white)](https://github.com/MairwunNx/dickobrazz/pkgs/container/dickobrazz)
[![GitHub Release](https://img.shields.io/github/v/release/mairwunnx/dickobrazz?style=flat&logo=github&color=blue)](https://github.com/mairwunnx/dickobrazz/releases)
[![Go Version](https://img.shields.io/badge/Go-1.26-00ADD8?style=flat&logo=go&logoColor=white)](https://go.dev/)

**Language**: [🇷🇺 Russian](README.MD) | 🇺🇸 English | [🇨🇳 中文](README.CN.MD)

🌶️ Dickobrazz bot, also known as dickobot, can estimate your "size" down to a centimeter. Do not worry if it is not huge today - it is just a joke.

> A modern and high-tech cockmeter: you ask, the bot returns a scientifically "validated" size and even jokingly maps it to a Russian region number. No ruler or microscope needed.

> Product inspired by South Park episode 1504 (#213), T.M.I.

## Features

### Cock Size

Every day the bot gives a unique value from 0 to 61 cm. The value resets at midnight MSK, includes seasonal emojis, and maps to a Russian region number.

### Leaderboards and Competitions

- **Cock Ruler** — top-13 best sizes for the current day
- **Cock Ladder** — all-time ranking by total accumulated size
- **Cock Race** — 3-month seasonal competition with accumulated sizes
- **Season System** — every 3 months a new season with winners and rewards
- **Cock-Respect™** — points for season wins and completed achievements

### Achievements

50+ unique achievements: from the first measurement to rare numeric combinations. Every achievement grants respect points and includes progress tracking.

### Advanced Analytics

- **CRI (Cock Rating Index)** — personal rating (T.M.I.-style metric from South Park)
- **Detailed statistics** — average size, median, records, and change dynamics
- **Luck coefficient**, **cock volatility**, **growth speed**
- **Dominance percentage** — user contribution share in overall system stats

## Usage

Dickobrazz works through **inline queries** in Telegram. Type `@dickobrazz_bot` in any chat and choose an option:

### Available inline commands

- **Cock Size** — get your size for today with region mapping
- **Cock Ruler** — top-13 players for the current day (daily ranking)
- **Cock Ladder** — all-time ranking by total accumulated size
- **Cock Race** — current season ranking (3 months)
- **Cock Dynamics** — detailed personal statistics and analytics
- **Cock Seasons** — season history and winners with navigation
- **Cock Achievements** — achievements and progress with pagination

### How it works

1. Open any Telegram chat
2. Type `@dickobrazz_bot`
3. Choose the needed option from the list
4. Send it to chat or view privately

## Run & Deployment 🚀

### Before you start 📦

If the repository is not cloned yet:
```bash
git clone https://github.com/mairwunnx/dickobrazz
cd dickobrazz
```

---

### Dev Containers (recommended for development) 🧰

**Requirements**:
- **1 GB RAM** for containers
- Running **Docker**
- **VS Code** or **Cursor**

**How to open**:
1. Open the project in VS Code/Cursor
2. Click **Reopen in Container**

**Environment variables**:
- File: `.devcontainer/devcontainer.env`
- Minimum required:
```bash
BOT_CSOT="your_internal_csot_token"
BOT_SERVER_BASE_URL="https://api.dickobrazz.com"
BOT_TG_TOKEN="your_telegram_bot_token"
BOT_TG_ENV="test"
```

**Run dev bot** *(uses `.devcontainer/devcontainer.env`)*:
```bash
air
```
`air` automatically restarts the app on source changes.
or
```bash
go run program.go
```

---

### Production via Docker Compose 🐳

**Environment variables**:
- File: `.env`
- Minimum required:
```bash
BOT_CSOT="your_internal_csot_token"
BOT_SERVER_BASE_URL="https://api.dickobrazz.com"
BOT_TG_TOKEN="your_telegram_bot_token"
BOT_TG_ENV="production"
```

**Start**:
```bash
docker compose up -d
```

---

### Ready image + your own services 🧩

Use the ready image from GitHub Container Registry:

```yaml
services:
  dickobrazz:
    image: ghcr.io/mairwunnx/dickobrazz:latest
    env_file: .env
```

---

### Build from source 🛠️

**Requirements**:
- Go `1.26+`

```bash
go build -o dickobrazz program.go
```

---

### Build with Docker 🧱

**Requirements**:
- Docker

```bash
docker build -t dickobrazz .
```

## Tech Stack

- **Go 1.26** — core development language
- **Telegram Bot API (tgbotapi/v5)** — Telegram integration
- **Resty** — HTTP client for backend interaction
- **uber/fx** — dependency injection and lifecycle management
- **Prometheus** — metrics and monitoring
- **go-i18n** — internationalization (6 languages)
- **Docker + Docker Compose** — containerization and orchestration
- **Structured Logging (slog)** — context-aware JSON logging

## AI Participation

AI was used to create parts of documentation and commits, and to optimize some MongoDB aggregation pipelines and statistics algorithms.

Also, as of September 23, 2025, AI hints (prompts/rules) were added to `.cursor/rules/`. They can help when implementing new features with LLM agents in the future.

## Links to Related Projects

[Xi Manager](https://github.com/mairwunnx/xi) — 🀄️ AI-powered Telegram bot styled as Xi's personal assistant. A personal assistant for the great leader, ready to answer questions from ordinary people.

[Louisepizdon](https://github.com/MairwunNx/louisepizdon) — 🥀 Louisepizdon, an AI-powered Telegram bot that is more honest than your grandmother. It evaluates you "fairly" and breaks down your outfit pricing from photos.

## From the "Same Author" Series

[Mo'Bosses](https://github.com/mairwunnx/mobosses) — 🏆 **Mo'Bosses** is the best RPG plugin that turns ordinary mobs into epic bosses with an **advanced player progression system**. Unlike other plugins, every fight matters here, and every level unlocks new possibilities. ⚔

[Mo'Joins](https://github.com/mairwunnx/mojoins) — 🎉 Custom join/leave effects: messages, sounds, particles, fireworks, and protection after join. All for PaperMC.

[Mo'Afks](https://github.com/mairwunnx/moafks) — 🛡️ Online-time pause is now possible. A PaperMC plugin that provides safe AFK mode: damage immunity, no collisions, mob ignore, automatic inactivity detection, and neat visual effects.

[McBuddy Server](https://github.com/mcbuddy-ai/mcbuddy-server) — 🛠️⚡ Backend for the MCBuddy AI assistant with OpenRouter integration and request handling

[McBuddy Telegram](https://github.com/mcbuddy-ai/mcbuddy-bot) — 🤖📱 Telegram bot to talk with MCBuddy outside the game

[McBuddy Spigot](https://github.com/mcbuddy-ai/mcbuddy-spigot) — 💬 Spigot plugin for MCBuddy integration — adds `/ask` command for AI assistant questions directly in Minecraft server chat. 🎮

---

![image](./media.png)

🇷🇺 **Made in Russia with love.** ❤️

**Dickobrazz** is about honest randomness and healthy competition. No cheating, only math and luck.

> 🫡 Made by Pavel Erokhin (Павел Ерохин), aka mairwunnx.
