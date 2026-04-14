# mmw-notifications — Event Consumer Module

This Go module is part of the [mmw](https://github.com/pivaldi/mmw) project that demonstrates the implementation of the [Go Modular Monolith White Paper](https://github.com/pivaldi/go-modular-monolith-white-paper).

This module cannot run independently. It is wired into [the platform runner](https://github.com/pivaldi/mmw) inside `cmd/mmw/main.go` alongside the `auth` and `todo` modules.

---

## Overview

`mmw-notifications` is a **pure event consumer module**. It has no HTTP server, no domain layer, no repository, and no outbox. Its sole responsibility is to subscribe to Watermill topics produced by other modules and react to the events it receives.

In the current POC it:

- **Logs every received event** to the structured logger (`log/slog`)
- **Forwards events to RocketChat** (optional — activated by `WithNotifer: true` and the environment variables below)

Because notifications only consumes events from the shared in-process Watermill `GoChannel`, it is **intentionally excluded from Mode 3** (distributed processes). In a real distributed deployment it would subscribe to a shared broker (e.g. RabbitMQ) instead.
