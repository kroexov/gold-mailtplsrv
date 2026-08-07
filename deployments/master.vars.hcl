service = {
  instances = 2
  cpu       = 500
  memory    = 512

  config = <<EOF
[Server]
Host      = "0.0.0.0"
Port      = 8075
IsDevel   = false

[Renderer]
PublicWebPath = "https://mailtplsrv.example.com/mail/"
WebPath       = "/mail/"
IsEmbed       = true

[RPC]
CallerNamespace = ""

[Sentry]
DSN         = ""
Environment = "production"
EOF

  env = {
    MAILTPLSRV_VERBOSE = "true"
    MAILTPLSRV_JSON    = "true"
    MAILTPLSRV_CONFIG  = "config.toml"
  }
}
