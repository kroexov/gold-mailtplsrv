# CLAUDE.md

Guidance for Claude Code (claude.ai/code) when working in this repository.

`mailtplsrv` renders MJML email templates with Go `html/template` data binding and returns HTML over
JSON-RPC. It is a template repository — most work is adding or editing templates under `tmpl/`, not changing
Go code.

## Where to look

Read the file that matches the task; do not read all of them up front.

| Task | File |
|------|------|
| Add or edit a template, fixture, static asset; debug rendering; call the API; change config | `README.md` |
| Change Go code, tests, CI, deployment; update dependencies | `DEVELOPMENT.md` |
| Look up an MJML component, attribute or head tag | `.claude/commands/mjmldoc.md` |

## Rules

* Comments and identifiers in English; user-facing template copy in whatever language the project uses.
* `pkg/rpc/rpc_zenrpc.go` is generated — run `make generate` after editing `pkg/rpc/render.go`, never edit it
  by hand.
* Template file names must be unique across all of `tmpl/`; a duplicate is a startup error.
* Images go through `{{ static }}` — never hardcode an absolute asset URL.
* A new template needs a fixture in `tmpl/fixture/<name>.json`, otherwise `/template/test/<name>` renders
  with `nil` data and the tests cover nothing.
* `make lint` and `make test` must be clean before committing.
* Keep `vendor/` in sync via `make mod`; `getsentry/sentry-go` is pinned — see DEVELOPMENT.md before bumping it.
