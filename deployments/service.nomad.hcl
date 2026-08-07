locals {
  docker_image = "registry.example.com/mailtplsrv:${var.image_version}"
}

variable "image_version" {
  type        = string
  description = "The docker image version used for task."
}

variable "commit_sha" {
  type        = string
  description = "The short SHA1 of commit."
  default     = "no_version"
}

variable "service" {
  type = object({ cpu = number, memory = number, instances = number, config = string, env = map(string) })
  description = "Configuration object for service."
}

job "mailtplsrv" {
  datacenters = ["dc1"]
  type = "service"

  update {
    max_parallel      = var.service.instances
    canary            = var.service.instances
    min_healthy_time  = "30s"
    healthy_deadline  = "2m"
    auto_revert       = true
    auto_promote      = true
    progress_deadline = "3m"
  }

  meta {
    commit_SHA           = var.commit_sha
    docker_image_version = var.image_version
  }

  group "backend" {
    count = var.service.instances
    restart {
      attempts = 3
      interval = "5m"
      delay    = "5s"
      mode     = "fail"
    }

    network {
      port "http" {
        to = 8075
      }
    }

    task "mailtplsrv" {
      template {
        change_mode = "restart"
        destination = "local/config.toml"
        data        = var.service.config
      }

      driver = "docker"
      config {
        image   = local.docker_image
        volumes = ["local/config.toml:/opt/mailtplsrv/config.toml"]
        ports   = ["http"]
        labels  = {
          release        = var.commit_sha
          loki_json_logs = "true"
        }
      }
      env = var.service.env

      resources {
        cpu    = var.service.cpu
        memory = var.service.memory
      }
      kill_timeout = "5s"
    }

    service {
      name = "${NOMAD_JOB_NAME}"
      port = "http"
      tags = [
        "metrics",
        "traefik.http.routers.${NOMAD_JOB_NAME}-A.rule=Headers(`X-Application`, `${NOMAD_JOB_NAME}`)",
        "traefik.http.routers.${NOMAD_JOB_NAME}-B.rule=Host(`${NOMAD_JOB_NAME}.traefik.service.consul`)",
      ]
      check {
        type     = "tcp"
        interval = "10s"
        timeout  = "2s"
      }
    }
  }
}
