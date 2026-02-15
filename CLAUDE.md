# CLAUDE.md

## Project Overview

Telegram Bot client for Dickobrazz (dickobrazz.com) — a Telegram bot with daily size, leaderboards, seasons, achievements, and analytics.

Modern and high-tech cockmeter: you ask, and the bot provides a scientifically justified size and even jokingly compares your unit size with Russian region numbers. No ruler or microscope needed anymore!

Product inspired by episode 1504 (#213) of South Park series (T.M.I.).

## Architecture

This project follows **Feature-Action Architecture (FAA)** — a TypeScript/Bun/Golang adaptation of Feature-Sliced Design for backends. See `Feature-Action-Architecture/MANIFEST.md` for the full spec and `Feature-Action-Architecture/examples/ts-bun.md` or `golang-gin.md` with uber/fx. AI notes: `Feature-Action-Architecture/AI.md`.

### Layer Hierarchy (downward imports only)

```
App → Features → Entities → Shared
```

- **App** (`src/app/`): HTTP server, router, route factory, DI container, request context via AsyncLocalStorage.
- **Features** (`src/features/`): Business use cases. Each has `api/handler.ts`, `*.action.ts`, optionally `db/`, `lib/`, `types.ts`, `index.ts`.
- **Entities** (`src/entities/`): Domain objects. Each has `model.ts` (Mongoose schema), `dal.ts` (CRUD), optionally `cache.ts`, `lib/`, `types.ts`.
- **Shared** (`src/shared/`): Infrastructure (MongoDB, Redis, config, logging, metrics, random) and utilities (datetime, encoding, sync, API primitives).


### Features

#### Core Functionality

- **Cock Size** — generates sizes from 0 to 61 cm with corresponding seasonal emojis.
- **Cock Ruler** — top-13 best sizes for the current day (daily ranking).
- **Cock Ladder** — top-13 eternal ranking by total accumulated size over all time (global ranking).
- **Cock Race** — top-13 seasonal competition (3 months) with size summation.
- **Season System** — every 3 months a new season with winners and rewards.
- **Cock-Respect™** — point system for season victories and achievement completion.
- **Cock Achievements** — achievement system with respect rewards.
- **Seasonal Emojis** — different emojis depending on the season.
- **Regional Mapping** — compares size with Russian region numbers.

### Usage

Dickobrazz works through **inline queries** in Telegram. Simply type `@dickobrazz_bot` in any chat and select the desired option:

#### Available inline commands

- **Cock Size** — get your size for today with regional mapping
- **Cock Ruler** — top-13 players for the current day (daily ranking)
- **Cock Ladder** — top-13 eternal ranking by total accumulated size over all time (global ranking)
- **Cock Race** — top-13 seasonal competition (3 months) with size summation
- **Cock Dynamics** — detailed global and personal statistics and analytics
- **Cock Seasons** — season history and winners with navigation
- **Cock Achievements** — achievements and progress with pagination

**Stack:** Go, Telegram Bot API (tgbotapi/v5), Resty (resty.dev/v3), YAML, prom-client, Docker, Docker-Compose, devcontainers, GitHub Actions, Air.

## Commands

```bash
go fmt ./...
go build -o out/dickobrazz program.go
```

After task done, run: `go fmt ./...`

## Core Rules

- Target `Go 1.26`.
- Use modern Go features where appropriate:
  - Generics for type safety.
  - `context` for request lifecycle and cancellation.
  - Structured logging.
  - Embedding for code reuse.
  - Range-over-function (Go 1.25+).
- Package names: short, lowercase, no underscores or plurals.
- Group related functionality logically; avoid large "god packages".
- Keep `main` package minimal — delegate logic to internal packages.
- Avoid panics in normal control flow.
- Use **structured logging** consistently.
- Always include relevant context (IDs, keys, parameters).
- Apply proper levels: Debug, Info, Warn, Error.
- Always use `context.Context` for timeouts and cancellation.
- Manage goroutine lifecycles explicitly; avoid leaks.
- Prefer channels or message passing to shared mutable state.
- Use `sync` primitives only when truly required.
- Prefer **composition over inheritance**.
- Use **interfaces** to abstract behavior, not just for testing.
- Keep functions small and cohesive.
- Return early to reduce nesting.
- Avoid global state; prefer dependency injection or constructors.
- Clean up connections and resources (`defer close(...)`).

## Communication & Commits

- Communicate with user preferably in Russian; thinking in English is fine.
- Name commits shortly in Russian language.
- Use emoji before commit summary that describes the changes.
- After finishing the task, provide a summary with checkmarks ✅ emoji tasks you've completed.