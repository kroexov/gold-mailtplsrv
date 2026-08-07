service = {
  instances = 1
  cpu       = 300
  memory    = 256

  config = <<EOF
[Server]
Host      = "0.0.0.0"
Port      = 8075
IsDevel   = true

[Renderer]
PublicWebPath = "https://mailtplsrv.example.dev/mail/"
WebPath       = "/mail/"
IsEmbed       = true

[RPC]
CallerNamespace = ""

[Sentry]
DSN         = ""
Environment = "devel"
EOF

  env = {
    MAILTPLSRV_VERBOSE = "true"
    MAILTPLSRV_JSON    = "true"
    MAILTPLSRV_CONFIG  = "config.toml"
  }
}
