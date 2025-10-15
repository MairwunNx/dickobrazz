# Dickobrazz

[![AI Capable](https://img.shields.io/badge/AI-Capable-brightgreen?style=flat&logo=openai&logoColor=white)](https://github.com/mairwunnx/dickobrazz)
[![Docker](https://img.shields.io/badge/Docker-Available-2496ED?style=flat&logo=docker&logoColor=white)](https://github.com/MairwunNx/dickobrazz/pkgs/container/dickobrazz)
[![GitHub Release](https://img.shields.io/github/v/release/mairwunnx/dickobrazz?style=flat&logo=github&color=blue)](https://github.com/mairwunnx/dickobrazz/releases)
[![Go Version](https://img.shields.io/badge/Go-1.25-00ADD8?style=flat&logo=go&logoColor=white)](https://go.dev/)
[![MongoDB](https://img.shields.io/badge/MongoDB-8.0.15-47A248?style=flat&logo=mongodb&logoColor=white)](https://www.mongodb.com/)
[![Redis](https://img.shields.io/badge/Redis-8.2-DC382D?style=flat&logo=redis&logoColor=white)](https://redis.io/)

**Language**: [🇷🇺 Russian](README.MD) | 🇺🇸 English | [🇨🇳 中文](README.CN.MD)

🌶️ Dickobrazz bot, also known as dickobot, can accurately measure your unit, tap, or tackle down to the centimeter. Don't be offended if you can't boast a big device - it's just a joke.

> Modern and high-tech cockmeter: you ask, and the bot provides a scientifically justified size and even jokingly compares your unit size with Russian region numbers. No ruler or microscope needed anymore!

> Product inspired by episode 1504 (#213) of South Park series (T.M.I.).

## Features

### Core Functionality

- **Cock Size** — generates sizes from 0 to 61* cm with corresponding seasonal emojis.
- **Daily Limits** — one size per day, resets at midnight Moscow time.
- **Seasonal Emojis** — different emojis depending on the season.
- **GSC** — integration with Global System of Cocks. (Global cock system extending beyond telegram)

> * - a bug that became a feature. It was planned up to 60 inclusive, but turned out to be 61 :).

### Gamification and Competitions

- **Cock Ruler** — top-13 best sizes for the current day.
- **Cock Race** — seasonal competition with size summation.
- **Season System** — every 3 months a new season with winners.
- **Cock-Respect™** — point system for season participation.

### Advanced Analytics

- **CRI (Cock Rating Index)** — personal rating relative to all users (Analog of T.M.I from South Park).
- **Cock Dynamics** — growth statistics, records, dominance percentage.
- **Detailed Statistics** — median, averages, size distribution.
- **Cock Assets** — change dynamics and daily growth.

### Technical Reliability

- **Double Randomization** — Random.org + ChaCha8 fallback for maximum randomness.
- **High-Performance Database** — MongoDB with aggregation pipelines.
- **Caching** — Redis with TinyLFU for fast responses.

## Usage

Dickobrazz works through **inline queries** in Telegram. Simply type `@dickobrazz_bot` in any chat and select the desired option:

### Available inline commands

- **Cock Size** — get your size for today
- **Cock Ruler** — top players for the day  
- **Cock Race** — current season ranking
- **Cock Ladder** — global ranking of participants by total size of cock over all time
- **Cock Dynamics** — detailed personal statistics
- **Cock Seasons** — season history and winners

### How it works

1. Open any chat in Telegram
2. Type `@dickobrazz_bot` 
3. Choose the option you're interested in from the list
4. Send the result to chat or view privately

## Deployment

### Dev Containers (recommended for development)

The simplest way to run for development — using Dev Containers in VS Code or similar IDEs:

1. Clone the repository:
```bash
git clone https://github.com/mairwunnx/dickobrazz
cd dickobrazz
```

2. Configure only two required tokens (in `.env.dev` file):
```bash
TELEGRAM_TOKEN="your_telegram_bot_token"
RANDOMORG_TOKEN="your_random_org_api_key"
```

3. Open the project in VS Code and select "Reopen in Container"

Everything else (MongoDB, Redis, dependencies) will be configured automatically. Simple and fast!

### Docker Compose (for production)

1. Clone the repository:
```bash
git clone https://github.com/mairwunnx/dickobrazz
cd dickobrazz
```

2. Configure two required tokens (in `.env` file):
```bash
TELEGRAM_TOKEN="your_telegram_bot_token"
RANDOMORG_TOKEN="your_random_org_api_key"
```

> **Note**: If needed, configure other environment variables (passwords for Redis and MongoDB).

3. Start services:
```bash
docker compose up -d
```

### Docker Compose with Ready Image

If you want to use a ready image and manage MongoDB/Redis separately:

1. Declare dickobrazz service in `docker-compose.yml`:

```yaml
services:
  dickobrazz:
    image: ghcr.io/mairwunnx/dickobrazz:1.23.4
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

> **Note**: You'll need to pass all environment variables from `.env` file or any other convenient way.

2. Declare MongoDB and Redis services in `docker-compose.yml`. 
> Important: **MongoDB** version should be `8.2+`, and **Redis** should be `8.2+`.

3. Start services:
```bash
docker compose up -d
```

### Build from Source

#### Manual and Direct Build

Requirements:
- Go 1.25+

```bash
go build -o dickobrazz program.go
```

#### Docker Build

Requirements:
- Docker

```bash
docker build -t dickobrazz .
```

## Tech Stack

- **Go 1.25** — main development language
- **Telegram Bot API (tgbotapi/v5)** — Telegram integration
- **MongoDB + Mongo Driver** — main database with aggregations
- **Redis + go-redis/cache** — caching with TinyLFU
- **Random.org API** — source of true randomness
- **ChaCha8** — cryptographically strong PRNG for fallback
- **Docker + Docker Compose** — containerization and orchestration
- **Structured Logging (slog)** — JSON logging with context

## AI Participation

AI was used to create part of the documentation and commits, as well as to optimize some MongoDB aggregation pipelines and improve statistics algorithms.

Also, as of September 23, 2025, AI hints (prompts/rules) have been added to the `.cursor/rules/` folder, which can help when implementing new features using LLM agents in the future.

## Links to Related Projects

[Xi Manager](https://github.com/mairwunnx/xi) — 🀄️ AI-powered Telegram bot styled as Xi's personal assistant. Personal assistant to the great leader, ready to answer questions from ordinary people.

[Louisepizdon](https://github.com/MairwunNx/louisepizdon) — 🥀 Louisepizdon, AI-powered Telegram bot that's more honest than your grandmother. Will evaluate you properly, analyzing the pricing of your clothes from photos!

## From the "Same Author" Series

[Mo'Bosses](https://github.com/mairwunnx/mobosses) — 🏆 **Mo'Bosses** is the best RPG plugin that turns ordinary mobs into epic bosses with **advanced player progression system**. Unlike other plugins, here every fight matters, and every level unlocks new possibilities! ⚔

[Mo'Joins](https://github.com/mairwunnx/mojoins) — 🎉 Custom join/leave messages: messages, sounds, particles, fireworks and post-join protection. All for PaperMC.

[Mo'Afks](https://github.com/mairwunnx/moafks) — 🛡️ Pause in online time is now possible. Plugin for PaperMC that gives players safe AFK mode: damage immunity, no collisions, ignored by mobs, auto-inactivity detection and neat visual effects.

[McBuddy Server](https://github.com/mcbuddy-ai/mcbuddy-server) — 🛠️⚡ Backend for MCBuddy AI assistant with OpenRouter integration and request processing

[McBuddy Telegram](https://github.com/mcbuddy-ai/mcbuddy-bot) — 🤖📱 Telegram bot for communicating with MCBuddy outside the game

[McBuddy Spigot](https://github.com/mcbuddy-ai/mcbuddy-spigot) — 💬 Spigot plugin for MCBuddy integration — adds `/ask` command for AI assistant questions directly in Minecraft server chat! 🎮

---

![image](./media.png)

🇷🇺 **Made in Russia with love.** ❤️

**Dickobrazz** is about honest randomness and healthy competition. No cheating, just math and luck.

> 🫡 Made by Pavel Erokhin (Павел Ерохин), aka mairwunnx.
