job "svc-sample" {
  type = "service"

  datacenters = [
    "ap-southeast-1a",
    "ap-southeast-1b",
    "ap-southeast-1c"
  ]

   constraint {
    attribute = "${attr.cpu.arch}"
    value     = "arm64"
  }

  group "api" {
    count = 2

    update {
      max_parallel     = 1
      min_healthy_time = "30s"
      healthy_deadline = "180s"
      stagger          = "60s"
    }

    network {
      port "http" {}
      port "grpc" {}

      dns {
        servers = [
          "172.17.0.1"
        ]
      }
    }

    spread {
      attribute = "${attr.platform.aws.placement.availability-zone}"
    }

    task "api" {
      driver = "docker"

      env {
        SENTRY_ENV     = "staging"
        SENTRY_DSN     = "https://e18c931088a04b2dab6906dbf428e327@o445007.ingest.sentry.io/6115639"
        SENTRY_RELEASE = "${CI_PIPELINE_IID}"

        POSTGRES_HOST = "rds.svc.staging.jungleegames.io"
        POSTGRES_DB   = "jungleegames_sample"
      }

      config {
        image        = "${CI_REGISTRY_IMAGE}/${DEPLOY_PROJECT}:${CI_PIPELINE_IID}"
        network_mode = "host"

        ports = [
          "http",
          "grpc"
        ]
      }

      service {
        name = "svc-sample-grpc"
        port = "grpc"

        check {
          type     = "grpc"
          port     = "grpc"
          interval = "5s"
          timeout  = "2s"
        }
      }

      service {
        name = "svc-sample-http"
        port = "http"

        check {
          type     = "tcp"
          port     = "http"
          interval = "5s"
          timeout  = "2s"
        }
      }

      resources {
        cpu    = 100
        memory = 126
      }
    }
  }
}
