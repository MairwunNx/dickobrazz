# Dickobrazz

[![AI Capable](https://img.shields.io/badge/AI-Capable-brightgreen?style=flat&logo=openai&logoColor=white)](https://github.com/mairwunnx/dickobrazz)
[![Build Status](https://img.shields.io/github/actions/workflow/status/mairwunnx/dickobrazz/build.yml?branch=master&style=flat&logo=githubactions&logoColor=white)](https://github.com/mairwunnx/dickobrazz/actions/workflows/build.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/mairwunnx/dickobrazz)](https://goreportcard.com/report/github.com/mairwunnx/dickobrazz)
[![Runs on Telegram](https://img.shields.io/badge/Runs_on-Telegram-26A5E4?style=flat&logo=telegram&logoColor=white)](https://t.me/dickobrazz_bot)
[![Docker](https://img.shields.io/badge/Docker-Available-2496ED?style=flat&logo=docker&logoColor=white)](https://github.com/MairwunNx/dickobrazz/pkgs/container/dickobrazz)
[![GitHub Release](https://img.shields.io/github/v/release/mairwunnx/dickobrazz?style=flat&logo=github&color=blue)](https://github.com/mairwunnx/dickobrazz/releases)
[![Go Version](https://img.shields.io/badge/Go-1.26-00ADD8?style=flat&logo=go&logoColor=white)](https://go.dev/)

**Language**: [🇷🇺 Russian](README.MD) | 🇺🇸 English | [🇨🇳 中文](README.CN.MD)

🌶️ Dickobrazz Telegram bot can estimate your "size" down to a centimeter. Do not worry if it is not huge today - this is just a joke.

> A modern and high-tech cockmeter: you ask, and the bot gives a scientifically "validated" size, then jokingly maps it to a Russian region number. No ruler or microscope needed anymore.

> [!NOTE]
> Product inspired by South Park episode 1504 (#213), T.M.I.

## Features

### Cock Size

Every day the bot gives a unique value from 0 to 61 cm. The value resets at midnight MSK, includes seasonal emojis, and maps to a Russian region number.

### Leaderboards and Competitions

- **Cock Ruler** - top-13 best sizes for the current day
- **Cock Ladder** - all-time ranking by total accumulated size
- **Cock Race** - 3-month seasonal competition with accumulated sizes
- **Season System** - every 3 months a new season with winners and rewards
- **Cock-Respect™** - points for season wins and completed achievements

### Achievements

50+ unique achievements: from the first measurement to rare numeric combinations. Every achievement grants respect points and includes progress tracking.

### Advanced Analytics

- **CRI (Cock Rating Index)** - personal rating (similar to T.M.I. from South Park)
- **Detailed statistics** - average size, median, records, and change dynamics
- **Luck coefficient**, **cock volatility**, **growth speed**
- **Dominance percentage** - user contribution share in overall system stats

## Usage

Dickobrazz works through **inline queries** in Telegram. Just type `@dickobrazz_bot` in any chat and choose an option:

### Available Inline Commands

- **Cock Size** - get your size for today with region mapping
- **Cock Ruler** - top-13 players for the current day (daily ranking)
- **Cock Ladder** - all-time ranking by total accumulated size
- **Cock Race** - current season ranking (3 months)
- **Cock Dynamics** - detailed personal statistics and analytics
- **Cock Seasons** - season history and winners with navigation
- **Cock Achievements** - achievements and progress with pagination

### How It Works

1. Open any Telegram chat
2. Type `@dickobrazz_bot`
3. Choose the option you need from the list
4. Send the result to the chat or view it privately

## Run and Deployment 🚀

### Before You Start 📦

If the repository has not been cloned yet:
```bash
git clone https://github.com/mairwunnx/dickobrazz.git
```

---

### Dev Containers (recommended for development and local dev run) 🧰

**Requirements**:
- Running **Docker**
- At least **1 GB RAM** for the development container (VS Code/Cursor/Goland)
- A tool that supports dev containers, or just a CLI utility to run a dev container

> [!IMPORTANT]
> Before running the dev container, set environment variables in `.devcontainer/devcontainer.env`.

**Environment Variables**:
- File: `.devcontainer/devcontainer.env`
- Minimum required:
```bash
BOT_CSOT="your_internal_csot_token" # Any trusted string received from backend.
BOT_SERVER_BASE_URL="https://api.dickobrazz.com" # Backend base URL.
BOT_TG_TOKEN="your_telegram_bot_token" # Telegram bot token.
BOT_TG_ENV="test" # Telegram bot environment (test/production).
```

Or edit the config file `./config.yaml`.

**How to Open** _(VS Code/Cursor/Goland)_:
1. Open the project in VS Code/Cursor/Goland
2. Find and run **Reopen in Container**

Your dev-container-enabled tool should automatically build and start the development container.

**Run Bot in Dev Container**:

> [!NOTE]
> To run the bot for development inside a dev container, use standard Go tooling or air with hot reload. _(all required tooling is already bundled in the container)_

| Tooling | Command | Description |
|---------|---------|---------|
| air | `air` | Run bot with hot reload |
| go | `go run program.go` | Regular bot run |

---

### Via Docker Compose 🐳

**Environment Variables**:
- File: `.devcontainer/devcontainer.env`
- Minimum required:
```bash
BOT_CSOT="your_internal_csot_token" # Any trusted string received from backend.
BOT_SERVER_BASE_URL="https://api.dickobrazz.com" # Backend base URL.
BOT_TG_TOKEN="your_telegram_bot_token" # Telegram bot token.
BOT_TG_ENV="test" # Telegram bot environment (test/production).
```

**Run**:
```bash
docker compose up -d
```

---

### Use Ready Image 🧩 (production)

Use the ready image from GitHub Container Registry:

```yaml
services:
  dickobrazz:
    image: ghcr.io/mairwunnx/dickobrazz:latest
    env_file: .env
```

> [!IMPORTANT]
> If you need to modify configuration file values, rebuild the image with new config, or mount a local config into container (the first approach is recommended).

---

## Stack

- **Go 1.26** - primary development language
- **Telegram Bot API (tgbotapi/v5)** - Telegram integration
- **Resty** - HTTP client for backend interaction
- **uber/fx** - dependency injection and lifecycle management
- **Prometheus** - metrics and monitoring
- **go-i18n** - internationalization (6 languages)
- **Docker + Docker Compose** - containerization and orchestration
- **Structured Logging (slog)** - JSON logging with context

## Links to Projects in the Same Ecosystem

[dickobrazz-server](https://github.com/MairwunNx/dickobrazz-server) - 🌶️ Backend server for @dickobrazz_bot: user data handling, personal data aggregation, and statistics generation.

## Links to Related Projects

[Emperor Xi](https://github.com/mairwunnx/xi) - 🀄️ AI-powered Telegram bot styled as Emperor Xi. Personal assistant of a great leader, ready to answer questions from ordinary people.

---

<img src="./media.png" alt="Russian power" width="500">

🇷🇺 **Made in Russia with love.** ❤️

**Dickobrazz** is about honest randomness and healthy competition. No cheating, only math and luck.

> 🫡 Made by Pavel Erokhin (Павел Ерохин), aka mairwunnx.
