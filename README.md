# gold-mailtplsrv

Email template rendering service in Go. Takes an [MJML](https://mjml.io/) template plus JSON data and
returns responsive HTML ready to be handed to an SMTP/ESP sender.

This is a **template repository**: create your repository from it, drop in your own templates and static
assets, and you have a mailtplsrv for your project. Day-to-day work happens in `tmpl/` and `static/` — the Go
code rarely needs to change.

* MJML is rendered by [gomjml](https://github.com/preslavrachev/gomjml) — pure Go, no Node.js runtime required.
* Templates are Go [html/template](https://pkg.go.dev/html/template) files, so data binding, conditionals and
  loops work as usual.
* Templates and static assets are embedded into the binary for production and read from disk in development,
  so editing a template does not require a restart.

Working on the service itself rather than on templates? See [DEVELOPMENT.md](DEVELOPMENT.md).

## Quick start

```bash
make init    # copy Makefile.mk.dist -> Makefile.mk and cfg/local.toml.dist -> cfg/local.toml
make run     # start on http://localhost:8075
```

Open http://localhost:8075/template/ for the list of templates and
http://localhost:8075/template/test/welcome for the bundled example.

## Making a new service out of this repository

1. Create your repository from this one and clone it.
2. Replace `tmpl/welcome/` and the placeholder assets in `static/` with your own template and assets.
3. Set `RPC.CallerNamespace` in the config to your project name.
4. Point `deployments/*.vars.hcl` and `deployments/service.nomad.hcl` at your registry, domains and Sentry DSN.

No Go code needs to change.

The module is named `mailtplsrv` and the binary follows `NAME` in `Makefile.mk`. Keeping both as-is across
projects is intentional: templates, docker images and Nomad jobs stay comparable between services.

## Writing templates

A template is MJML markup with Go template actions in it: `html/template` runs first and produces MJML,
then `gomjml` turns that MJML into HTML. Two references worth keeping open:

* [text/template](https://pkg.go.dev/text/template) — the action syntax: `{{ .field }}`, `{{ if }}`,
  `{{ range }}`, `{{ with }}`, `{{ define }}`, pipelines.
* [html/template](https://pkg.go.dev/html/template) — the package actually used here. Same syntax, but values
  are context-aware escaped, which is what keeps user data from breaking the markup.

```
tmpl/
  _shared/                     shared partials and styles, never listed as templates
    partial-header.mjml
  welcome/                     one directory per email
    welcome.mjml               the email itself
    welcome-body.mjml          optional content fragment
  fixture/
    welcome.json               default test data for `welcome`
    test.json                  extra data set, see ?fixture=
static/
  logo.png                     placeholder wordmark, {{ static }}logo.png
  test.jpg                     placeholder banner, {{ static }}test.jpg
```

Rules the renderer applies:

| Rule | Meaning |
|------|---------|
| A template name is the **file name**, not the path | `tmpl/welcome/welcome.mjml` is rendered as `welcome` and included as `{{ template "welcome.mjml" . }}` |
| File names must be **globally unique** | Two `header.mjml` in different directories are a startup error, not a silent override |
| `_shared/` and `fixture/` | Parsed, but never listed as standalone templates |
| `partial-*` prefix | Parsed, but never listed — for includes |
| `*-body`, `*-styles` suffix | Parsed, but never listed — fragments of another template |

Everything else under `tmpl/` is a standalone email and shows up in `/template/` and in `render.list`.

Nesting is free-form: a flat `tmpl/welcome.mjml` works exactly as well as `tmpl/welcome/welcome.mjml`.
Use directories once a template grows its own partials, styles and images.

### Where to get MJML

[mjml.io/try-it-live](https://mjml.io/try-it-live) is a live editor with a side-by-side preview — the fastest
way to lay out an email and see the generated HTML. [mjml.io/templates](https://mjml.io/templates) has free
ready-made templates worth using as a starting point or just for ideas; the `welcome` example bundled here is
one of them.

The workflow is: build the layout in the editor, paste the MJML into `tmpl/<name>/<name>.mjml`, replace the
hardcoded strings with `{{ .field }}` and image URLs with `{{ static }}`. Then check the result at
`/template/test/<name>` — `gomjml` is a Go reimplementation, not the JS compiler, so complex layouts may need
adjusting (see [caveats](#gomjml-caveats)).

`.claude/commands/mjmldoc.md` is a full MJML component reference kept in the repo for offline lookups.

### Template functions

| Function | Example | Result |
|----------|---------|--------|
| `static` | `<mj-image src="{{ static }}test.jpg" />` | prefixes with `Renderer.PublicWebPath` |
| `add`, `sub` | `{{ add $index 1 }}` | integer arithmetic, handy inside `range` |
| `pluralize` | `{{ pluralize .count "товар" "товара" "товаров" }}` | Russian plural forms; accepts `int` and JSON `float64` |

Never hardcode absolute asset URLs — `{{ static }}` resolves to `Renderer.PublicWebPath`, which differs
between environments.

### Fixtures and preview

`GET /template/test/<name>` renders a template with test data:

* `tmpl/fixture/<name>.json` is used by default, if present;
* `?fixture=<other>` loads `tmpl/fixture/<other>.json` instead;
* no fixture at all renders the template with `nil` data.

```
http://localhost:8075/template/test/welcome
http://localhost:8075/template/test/welcome?fixture=test
```

Check both viewports before shipping: desktop (~648px) and mobile (<480px). When rendering fails in
development mode, the intermediate MJML is written to `debug.html` — that is usually enough to tell a Go
template error from an MJML one.

### gomjml caveats

* `mj-wrapper` hardcodes `max-width: 600px`. If `mj-body` is wider, override it via a CSS class in a
  non-inline `<mj-style>`, both on the class and on `.class > table`.
* `background-color` on `mj-wrapper` is overridden by the inline style gomjml generates — set the background
  through a CSS class instead of the attribute.
* CSS in `<mj-style inline="inline">` is inlined onto the outer element, while the MJML `padding` attribute
  applies to the inner cell. Setting padding in both places doubles it.
* Media queries only work in a non-inline `<mj-style>`; inline styles cannot be overridden by them.
* Columns stack on mobile (<480px). Wrap them in `mj-group` to keep them side by side — widths inside a group
  must be in `%`.
* `mj-attributes` cannot be nested: a shared `attributes.mjml` should contain only `mj-class` definitions,
  with the `<mj-attributes>` wrapper left in the template that includes it.

## API

Service communication goes through JSON-RPC:

| Endpoint | Description |
|----------|-------------|
| `POST /int/rpc/` | JSON-RPC 2.0: `render.renderTemplate`, `render.list` |
| `GET /int/rpc/doc/` | SMD Box — interactive API docs |
| `GET /int/rpc/api.go` | generated Go client |

```bash
curl -s localhost:8075/int/rpc/ -H 'Content-Type: application/json' -d '{
  "jsonrpc": "2.0", "id": 1, "method": "render.renderTemplate",
  "params": {"templateName": "welcome", "data": {"discount": 50, "promoCode": "SALEONSALE"}}
}'
```

A plain REST endpoint is available for non-Go clients:

```bash
curl -s -X POST localhost:8075/template/welcome -H 'Content-Type: application/json' -d '{"discount": 50}'
```

Service endpoints: `/status` (healthcheck), `/metrics` (Prometheus), `/debug/pprof/`, `/debug/metadata`,
`/mail/*` (static assets, path configurable), and `/` with the route list in devel mode.

## Configuration

`cfg/local.toml` (see `cfg/local.toml.dist`):

```toml
[Server]
Host      = "localhost"
Port      = 8075
IsDevel   = true          # exposes route list and verbose RPC errors

[Renderer]
PublicWebPath = "http://localhost:8075/mail/"   # value of {{ static }} — must be reachable by mail clients
WebPath       = "/mail/"                        # where this service serves static assets
IsEmbed       = false                           # false: read tmpl/ from disk on every render

[RPC]
CallerNamespace = ""                            # your project name, used in the generated Go client

[Sentry]
DSN         = ""
Environment = ""
```

`RPC.CallerNamespace` only affects the Go client served at `/int/rpc/api.go` — the caller name that consuming
services report to this one. Empty value gives `mailtplsrv`, `"myproject"` gives `myproject-mailtplsrv`.

Every field is also settable via env with the `MAILTPLSRV_` prefix (`MAILTPLSRV_CONFIG`, `MAILTPLSRV_VERBOSE`,
`MAILTPLSRV_JSON`).

`IsEmbed = false` is the development mode: templates are re-read from disk on every request, so a browser
refresh is enough to see a change. Production always runs with `IsEmbed = true`, which parses templates once
at startup from the binary.

> `PublicWebPath` must be an absolute, publicly reachable URL in production — mail clients load images from
> the outside, not from the service.

## Deployment

`deployments/Dockerfile` builds an alpine image around the static binary; `service.nomad.hcl` plus
`devel.vars.hcl` / `master.vars.hcl` describe the Nomad job. The config is delivered as a Nomad template into
`local/config.toml`, so changing it does not require a rebuild.

## License

MIT
