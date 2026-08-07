# Development

Notes for working on the service itself. For writing templates and calling the API see [README.md](README.md).

Requires Go 1.26. Dependencies are vendored (`GOFLAGS=-mod=vendor`), so builds are offline and reproducible.

## Commands

```bash
make init          # first-time setup: Makefile.mk + cfg/local.toml
make run           # run with dev logger and cfg/local.toml
make build         # CGO_ENABLED=0 build into ./mailtplsrv
make test          # go test ./...
make lint          # golangci-lint, must be clean before committing
make fmt           # golangci-lint fmt
make generate      # regenerate pkg/rpc/rpc_zenrpc.go
make mod           # go mod tidy && go mod vendor && git add vendor
make tools         # install golangci-lint of the pinned version
```

Run a single test with `go test ./pkg/renderer/ -run TestRendererRender`.

## Layout

```
cmd/mailtplsrv/     entry point: flags, config, sentry, graceful shutdown
pkg/app/            echo server, routes, middlewares, metrics
pkg/renderer/       template discovery, parsing and MJML rendering + REST handlers
pkg/rpc/            zenrpc JSON-RPC service (generated rpc_zenrpc.go)
embed.go            //go:embed of tmpl/ and static/
cfg/                config
deployments/        Dockerfile and Nomad job
```

`main.go` parses flags and TOML into `app.Config`, initialises the logger and Sentry, then hands over to
`app.New` / `app.Run`. `App.Run` registers metrics, template handlers, debug handlers, the RPC server and
service metadata, in that order, and starts echo.

## How the renderer works

`pkg/renderer/renderer.go` is the core of the service:

* `parseTemplates(fsys)` walks `tmpl/` in an `fs.FS` and parses every `.mjml` file. The template name is
  `filepath.Base(path)`, so file names must be unique across directories — duplicates return an error instead
  of silently overriding each other.
* `fsys()` returns `mailtplsrv.Templates` (embed) or `os.DirFS(".")` (live) depending on `Config.IsEmbed`.
  Both modes therefore go through exactly the same walk and parse code.
* Embed mode parses once in `New` and stores the result in `parsedTemplates`; live mode re-parses on every
  render, which is what makes edits visible without a restart.
* `List()` walks the same tree and hides `_shared/`, `fixture/`, `partial-*`, `*-body`, `*-styles`.
* Rendering is two-stage: `html/template` produces MJML, then `mjml.Render` produces HTML. On failure in live
  mode the intermediate MJML is dumped to `debug.html`.

Live mode resolves `tmpl/` and `static/` relative to the working directory — run the binary from the repo
root, and use `t.Chdir` in tests.

Template funcs are registered in `newTemplateWithFunc`; add project-specific helpers there.

## RPC

`pkg/rpc/server.go` wires the zenrpc server and its middleware chain (devel, headers, sentry, metrics, timing,
slog). `pkg/rpc/render.go` is the service itself — two methods, both thin wrappers over `renderer`.

After changing `render.go`, run `make generate`. `rpc_zenrpc.go` is generated and committed; never edit it by
hand. The Go client served at `/int/rpc/api.go` is generated at runtime by `rpcgen` from the SMD schema — its
package name comes from `appName`, its caller name from `RPC.CallerNamespace` in the config.

## Tests

`pkg/renderer/renderer_test.go` covers template discovery and rendering in **both** embed and live mode, plus
`pluralize`. If a template stops parsing or a fixture stops matching the template, `make test` fails — that is
the main safety net for template changes, so keep fixtures in sync.

Tests use `t.Chdir` to run from the repo root; no database and no network are needed, which is why CI runs
them directly.

## CI

`.github/workflows/go.yml` runs tests and build on Go 1.26; `.github/workflows/golangci-lint.yml` runs the
linter. Downstream projects usually replace these with their own GitLab CI pipeline — the Nomad job in
`deployments/` is there for that.

## Dependency pinning

`getsentry/sentry-go` is held at `v0.44.0`. From `v0.45.0` its echo integration moved to `echo/v5`, while
`vmkteam/appkit` is still on `echo/v4`, so newer versions fail to compile. Both `sentry-go` and
`sentry-go/echo` must be bumped together, and only once appkit migrates.

`go get -u ./...` is otherwise safe; after any update run `make mod`, `make test` and `make lint`.

## Deployment

`deployments/Dockerfile` copies the static binary into an alpine image and sets
`LABEL service_name="$CI_PROJECT_NAME"`, which is what the log pipeline uses to attribute logs.

`service.nomad.hcl` takes `image_version`, `commit_sha` and a `service` object; `devel.vars.hcl` and
`master.vars.hcl` supply that object, including the whole `config.toml` as a heredoc. The config is written
by a Nomad template into `local/config.toml` with `change_mode = "restart"`, so a config change restarts the
task without a rebuild.

Production always runs with `IsEmbed = true`.
