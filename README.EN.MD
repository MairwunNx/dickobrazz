# Dickobrazz

[![AI Capable](https://img.shields.io/badge/AI-Capable-brightgreen?style=flat&logo=openai&logoColor=white)](https://github.com/mairwunnx/ximanager)

**Language**: [🇷🇺 Russian](README.MD) | 🇺🇸 English

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
- **Cock Dynamics** — detailed personal statistics
- **Cock Seasons** — season history and winners

### How it works

1. Open any chat in Telegram
2. Type `@dickobrazz_bot` 
3. Choose the option you're interested in from the list
4. Send the result to chat or view privately

## Compatibility

- **Telegram Bot API** — works in any chats
- **Database**: MongoDB 5.0+
- **Cache**: Redis 6.0+
- **Runtime**: Go 1.23+
- **Deployment**: Docker + Docker Compose

## Installation

### Quick Start with Docker Compose

1. Clone the repository:
```bash
git clone https://github.com/mairwunnx/dickobrazz
cd dickobrazz
```

2. Configure environment variables:
```bash
# Copy example configuration
cp .env.draft .env

# Edit the .env file, filling in necessary variables
```

3. Start services:
```bash
docker compose up -d
```

4. Check status:
```bash
docker compose logs dickbot
```

### Build from Source

Requirements:
- Go 1.23+
- Docker for containerization

```bash
# Build binary
go build -o dickbot program.go

# Or build Docker image
docker build -t dickobrazz .
```

## Tech Stack

- **Go 1.23** — main development language
- **Telegram Bot API (tgbotapi/v5)** — Telegram integration
- **MongoDB + Mongo Driver** — main database with aggregations
- **Redis + go-redis/cache** — caching with TinyLFU
- **Random.org API** — source of true randomness
- **ChaCha8** — cryptographically strong PRNG for fallback
- **Docker + Docker Compose** — containerization and orchestration
- **Structured Logging (slog)** — JSON logging with context

## Architecture

### Structure
```
dickbot (Go app) ←→ Redis (cache)
       ↓
   MongoDB (data)
```

### Key Components

- **application.go** — main update processing loop
- **queries.go** — inline query handlers
- **mongo.go** — MongoDB logic and aggregations
- **redis.go** — size caching
- **random.go** — double randomization system
- **database/pipelines.go** — complex aggregation queries
- **messages.go** — message templates with formatting

### Inline Queries System

Each query type is handled separately:
- `CockSize` — main size with caching
- `CockRace` — season ranking through aggregation
- `CockRuler` — daily top from Redis
- `CockDynamic` — complex statistics through MongoDB pipelines
- `CockSeason` — season history and winners

## Implementation Features

### Randomness System
- **Random.org** — HTTP API for quantum randomness
- **ChaCha8** — backup PRNG with fixed seed for determinism during testing

### Caching
- **Daily Caching** — user size is saved until midnight Moscow time
- **TinyLFU** — eviction algorithm considering frequency and recency of access

### Season System
- **Automatic Seasons** — every 3 months from the date of first cock
- **Respect System** — points for season placements using formula `3000 / place^1.2`
- **Top-3 Winners** — visual rewards for users

### Performance
- **Operation Tracing** — execution time of each operation is logged
- **Aggregation Pipelines** — complex analytics calculated at database level
- **Batching** — multiple operations combined into pipelines

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

![image](./media.jpg)

🇷🇺 **Made in Russia with love.** ❤️

**Dickobrazz** is about honest randomness and healthy competition. No cheating, just math and luck.

> 🫡 Made by Pavel Erokhin (Павел Ерохин), aka mairwunnx.
