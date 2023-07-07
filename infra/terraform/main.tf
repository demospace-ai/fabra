terraform {
  required_providers {
    google = {
      source = "hashicorp/google"
    }
  }

  backend "gcs" {
    bucket = "fabra-344902-tfstate"
    prefix = "terraform/state"
  }
}

provider "google-beta" {
  project = "fabra-344902"
}

provider "google" {
  project = "fabra-344902"
}

resource "google_container_registry" "registry" {
  location = "US"
}

resource "google_compute_network" "vpc" {
  name                    = "fabra-vpc"
  routing_mode            = "GLOBAL"
  auto_create_subnetworks = true
}

resource "google_compute_subnetwork" "subnet" {
  name          = "fabra-vpc"
  region        = "us-west1"
  network       = google_compute_network.vpc.id
  ip_cidr_range = "10.138.0.0/20"
}

# Setup IP block for VPC
resource "google_compute_global_address" "private_ip_block" {
  name          = "private-ip-block"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  ip_version    = "IPV4"
  prefix_length = 20
  network       = google_compute_network.vpc.self_link
}

# Connection that allows services in our private VPC to access underlying Google Cloud Services VPC
resource "google_service_networking_connection" "private_vpc_connection" {
  network                 = google_compute_network.vpc.self_link
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_block.name]
}

resource "google_sql_database" "main_database" {
  name     = "fabra-db"
  instance = google_sql_database_instance.main_instance.name
}

resource "google_sql_database_instance" "main_instance" {
  name             = "fabra-database-instance"
  region           = "us-west1"
  database_version = "POSTGRES_11"
  settings {
    availability_type = "REGIONAL"
    tier              = "db-custom-1-3840"

    ip_configuration {
      ipv4_enabled                                  = true
      private_network                               = google_compute_network.vpc.self_link
      enable_private_path_for_google_cloud_services = true
    }
  }

  deletion_protection = "true"

  depends_on = [google_service_networking_connection.private_vpc_connection]
}

data "google_secret_manager_secret_version" "db_password" {
  secret = "fabra-db-password"
}

resource "google_sql_user" "db_user" {
  name     = "db_user"
  instance = google_sql_database_instance.main_instance.name
  password = data.google_secret_manager_secret_version.db_password.secret_data
}

resource "google_cloudbuild_worker_pool" "builder_pool" {
  name     = "fabra-pool"
  location = "us-west1"
  worker_config {
    no_external_ip = false
    disk_size_gb   = 100
    machine_type   = "e2-medium"
  }
  network_config {
    peered_network = google_compute_network.vpc.id
  }
}

resource "google_cloudbuild_trigger" "terraform-build-trigger" {
  name = "terraform-trigger"

  included_files = ["infra/terraform/**"]

  github {
    name  = "fabra"
    owner = "fabra-io"

    push {
      branch       = "main"
      invert_regex = false
    }
  }

  filename = "infra/cloudbuild/terraform.yaml"
}

resource "google_cloudbuild_trigger" "backend-build-trigger" {
  name = "backend-trigger"

  included_files = ["backend/server/**"]
  ignored_files  = ["backend/server/migrations/**"]

  github {
    name  = "fabra"
    owner = "fabra-io"

    push {
      branch       = "main"
      invert_regex = false
    }
  }

  filename = "infra/cloudbuild/backend.yaml"
}

resource "google_cloud_run_service" "fabra" {
  name     = "fabra"
  location = "us-west1"

  template {
    spec {
      service_account_name = google_service_account.fabra-backend.email
      containers {
        image = "gcr.io/fabra-344902/fabra"
        env {
          name  = "DB_USER"
          value = google_sql_user.db_user.name
        }
        env {
          name  = "DB_NAME"
          value = google_sql_database.main_database.name
        }
        env {
          name  = "DB_HOST"
          value = google_sql_database_instance.main_instance.private_ip_address
        }
        env {
          name  = "DB_PORT"
          value = "5432"
        }
        env {
          name  = "IS_PROD"
          value = "true"
        }
      }
    }

    metadata {
      annotations = {
        # Limit scale up to prevent any cost blow outs!
        "autoscaling.knative.dev/maxScale" = 5
        # Use the VPC Connector
        "run.googleapis.com/vpc-access-connector" = google_vpc_access_connector.connector.id
        # all egress from the service should go through the VPC Connector
        "run.googleapis.com/vpc-access-egress" = "all-traffic"
        "run.googleapis.com/client-name"       = "cloud-console"
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }

  lifecycle {
    ignore_changes = [
      template.0.metadata.0.annotations,
      template.0.metadata.0.labels,
      metadata.0.annotations,
      template.0.spec.0.containers.0.image,
    ]
  }

  autogenerate_revision_name = true
}

resource "google_vpc_access_connector" "connector" {
  name          = "vpcconn"
  region        = "us-west1"
  ip_cidr_range = "10.8.0.0/28"
  network       = google_compute_network.vpc.name
}

resource "google_cloud_run_service_iam_member" "all_users_member" {
  location = google_cloud_run_service.fabra.location
  project  = google_cloud_run_service.fabra.project
  service  = google_cloud_run_service.fabra.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}

resource "google_cloudbuild_trigger" "database-migration-trigger" {
  name = "database-migration-trigger"

  included_files = ["backend/server/migrations/**"]

  github {
    name  = "fabra"
    owner = "fabra-io"

    push {
      branch       = "main"
      invert_regex = false
    }
  }

  filename = "infra/cloudbuild/database-migration.yaml"
}

resource "google_compute_backend_service" "default" {
  security_policy                 = google_compute_security_policy.fabra-security-policy.id
  affinity_cookie_ttl_sec         = 0
  connection_draining_timeout_sec = 300
  enable_cdn                      = false
  load_balancing_scheme           = "EXTERNAL"
  name                            = "fabra-lb-backend-default"
  port_name                       = "http"
  protocol                        = "HTTP"
  session_affinity                = "NONE"
  timeout_sec                     = 30

  backend {
    balancing_mode               = "UTILIZATION"
    capacity_scaler              = 1
    group                        = google_compute_region_network_endpoint_group.fabra_neg.id
    max_connections              = 0
    max_connections_per_endpoint = 0
    max_connections_per_instance = 0
    max_rate                     = 0
    max_rate_per_endpoint        = 0
    max_rate_per_instance        = 0
    max_utilization              = 0
  }
}

resource "google_compute_global_address" "default" {
  address_type  = "EXTERNAL"
  name          = "fabra-lb-address"
  prefix_length = 0
}

resource "google_compute_global_forwarding_rule" "http" {
  ip_address            = google_compute_global_address.default.id
  ip_protocol           = "TCP"
  load_balancing_scheme = "EXTERNAL"
  name                  = "fabra-lb"
  port_range            = "80"
  target                = google_compute_target_http_proxy.default.id
}

resource "google_compute_global_forwarding_rule" "https" {
  ip_address            = google_compute_global_address.default.id
  ip_protocol           = "TCP"
  load_balancing_scheme = "EXTERNAL"
  name                  = "fabra-lb-https"
  port_range            = "443"
  target                = google_compute_target_https_proxy.default.id
}

locals {
  managed_domains = tolist(["app.fabra.io", "connect.fabra.io", "api.fabra.io"])
}

resource "random_id" "cert-name" {
  byte_length = 4
  prefix      = "issue6147-cert-"

  keepers = {
    domains = join(",", local.managed_domains)
  }
}

resource "google_compute_managed_ssl_certificate" "cert" {
  name = random_id.cert-name.hex
  type = "MANAGED"

  lifecycle {
    create_before_destroy = true
  }

  managed {
    domains = local.managed_domains
  }
}

resource "google_compute_url_map" "default" {
  name            = "fabra-lb-url-map"
  default_service = google_compute_backend_bucket.frontend_backend.id
  host_rule {
    hosts = [
      "app.fabra.io",
    ]
    path_matcher = "fabra-lb-path-matcher"
  }

  host_rule {
    hosts = [
      "connect.fabra.io",
    ]
    path_matcher = "fabra-connect-path-matcher"
  }

  host_rule {
    hosts = [
      "api.fabra.io",
    ]
    path_matcher = "fabra-api-path-matcher"
  }

  path_matcher {
    name            = "fabra-lb-path-matcher"
    default_service = google_compute_backend_bucket.frontend_backend.id
  }

  path_matcher {
    name            = "fabra-api-path-matcher"
    default_service = google_compute_backend_service.default.id
  }

  path_matcher {
    name            = "fabra-connect-path-matcher"
    default_service = google_compute_backend_bucket.connect_backend.id
  }
}

resource "google_compute_url_map" "https_redirect" {
  name = "fabra-lb-https-redirect"

  default_url_redirect {
    https_redirect         = true
    redirect_response_code = "MOVED_PERMANENTLY_DEFAULT"
    strip_query            = false
  }
}

resource "google_compute_target_http_proxy" "default" {
  name       = "fabra-lb-http-proxy"
  proxy_bind = false
  url_map    = google_compute_url_map.https_redirect.id
}

resource "google_compute_target_https_proxy" "default" {
  name          = "fabra-lb-https-proxy"
  proxy_bind    = false
  quic_override = "NONE"
  ssl_certificates = [
    google_compute_managed_ssl_certificate.cert.id,
  ]
  url_map = google_compute_url_map.default.id
}

resource "google_compute_region_network_endpoint_group" "fabra_neg" {
  provider              = google
  name                  = "fabra-neg"
  network_endpoint_type = "SERVERLESS"
  region                = "us-west1"
  cloud_run {
    service = google_cloud_run_service.fabra.name
  }
}

resource "google_storage_bucket" "fabra_frontend_bucket" {
  name          = "fabra-frontend-bucket"
  location      = "US"
  storage_class = "STANDARD"

  uniform_bucket_level_access = true

  website {
    main_page_suffix = "index.html"
    not_found_page   = "index.html"
  }

  cors {
    max_age_seconds = 3600
    method = [
      "GET",
    ]
    origin = [
      "*",
    ]
    response_header = [
      "Content-Type",
    ]
  }
}

resource "google_storage_bucket_iam_member" "public_frontend_read_access" {
  bucket = google_storage_bucket.fabra_frontend_bucket.name
  role   = "roles/storage.objectViewer"
  member = "allUsers"
}

resource "google_compute_backend_bucket" "frontend_backend" {
  name        = "frontend-backend-bucket"
  description = "Static react web app for Fabra"
  bucket_name = google_storage_bucket.fabra_frontend_bucket.name
  enable_cdn  = true
}

resource "google_cloudbuild_trigger" "frontend-build-trigger" {
  name = "frontend-trigger"

  included_files = ["frontend/**"]

  github {
    name  = "fabra"
    owner = "fabra-io"

    push {
      branch       = "main"
      invert_regex = false
    }
  }

  filename = "infra/cloudbuild/frontend.yaml"
}

resource "google_storage_bucket" "fabra_connect_bucket" {
  name          = "fabra-connect-bucket"
  location      = "US"
  storage_class = "STANDARD"

  uniform_bucket_level_access = true

  website {
    main_page_suffix = "connect.html"
    not_found_page   = "connect.html"
  }

  cors {
    max_age_seconds = 3600
    method = [
      "GET",
    ]
    origin = [
      "*",
    ]
    response_header = [
      "Content-Type",
    ]
  }
}

resource "google_storage_bucket_iam_member" "public_connect_read_access" {
  bucket = google_storage_bucket.fabra_connect_bucket.name
  role   = "roles/storage.objectViewer"
  member = "allUsers"
}

resource "google_compute_backend_bucket" "connect_backend" {
  name        = "connect-backend-bucket"
  description = "Static react web app for Fabra Connect"
  bucket_name = google_storage_bucket.fabra_connect_bucket.name
  enable_cdn  = true
}

# TODO: figure out how to only trigger this on connect changes
resource "google_cloudbuild_trigger" "connect-build-trigger" {
  name = "connect-trigger"

  included_files = ["frontend/**"]

  github {
    name  = "fabra"
    owner = "fabra-io"

    push {
      branch       = "main"
      invert_regex = false
    }
  }

  filename = "infra/cloudbuild/connect.yaml"
}

resource "google_kms_key_ring" "data-connection-keyring" {
  name     = "data-connection-keyring"
  location = "global"
}

resource "google_kms_crypto_key" "data-connection-key" {
  name            = "data-connection-key"
  key_ring        = google_kms_key_ring.data-connection-keyring.id
  rotation_period = "100000s"

  lifecycle {
    prevent_destroy = true
  }
}

resource "google_kms_key_ring_iam_binding" "data-connection-key-ring-binding" {
  key_ring_id = google_kms_key_ring.data-connection-keyring.id
  role        = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  members = [
    "serviceAccount:fabra-sync@fabra-344902.iam.gserviceaccount.com",
    "serviceAccount:fabra-backend@fabra-344902.iam.gserviceaccount.com"
  ]
}

resource "google_kms_key_ring" "api-key-keyring" {
  name     = "api-key-keyring"
  location = "global"
}

resource "google_kms_crypto_key" "api-key-key" {
  name            = "api-key-key"
  key_ring        = google_kms_key_ring.api-key-keyring.id
  rotation_period = "100000s"

  lifecycle {
    prevent_destroy = true
  }
}

resource "google_kms_key_ring_iam_binding" "api-key-key-ring-binding" {
  key_ring_id = google_kms_key_ring.api-key-keyring.id
  role        = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  members = [
    "serviceAccount:fabra-sync@fabra-344902.iam.gserviceaccount.com",
    "serviceAccount:fabra-backend@fabra-344902.iam.gserviceaccount.com"
  ]
}

resource "google_compute_router" "fabra-ip-router" {
  name    = "fabra-ip-router"
  network = google_compute_network.vpc.name
  region  = "us-west1"
}

resource "google_compute_address" "egress-ip-address" {
  name   = "egress-static-ip"
  region = "us-west1"
}

resource "google_compute_router_nat" "fabra-nat" {
  name   = "fabra-static-nat"
  router = google_compute_router.fabra-ip-router.name
  region = "us-west1"

  min_ports_per_vm       = 64
  nat_ip_allocate_option = "MANUAL_ONLY"
  nat_ips = [
    google_compute_address.egress-ip-address.self_link,
  ]

  source_subnetwork_ip_ranges_to_nat = "ALL_SUBNETWORKS_ALL_IP_RANGES"
}

resource "google_container_cluster" "fabra-sync-cluster" {
  name     = "fabra-sync-gke-cluster"
  location = "us-west1"

  network    = google_compute_network.vpc.name
  subnetwork = google_compute_subnetwork.subnet.name

  ip_allocation_policy {}

  enable_autopilot = true

  private_cluster_config {
    enable_private_nodes = true
  }

  cluster_autoscaling {
    auto_provisioning_defaults {
      service_account = google_service_account.fabra-sync.email
    }
  }
}

# TODO: figure out how to only trigger this when dependencies change
resource "google_cloudbuild_trigger" "worker-build-trigger" {
  name = "worker-trigger"

  included_files = ["backend/sync/**"]

  github {
    name  = "fabra"
    owner = "fabra-io"

    push {
      branch       = "main"
      invert_regex = false
    }
  }

  filename = "infra/cloudbuild/worker.yaml"
}

resource "google_service_account" "fabra-backend" {
  account_id   = "fabra-backend"
  display_name = "Fabra Backend Service"
}

resource "google_service_account" "fabra-sync" {
  account_id   = "fabra-sync"
  display_name = "Fabra Sync Service"
}

resource "google_service_account_iam_binding" "fabra-worker-gke-binding" {
  service_account_id = google_service_account.fabra-sync.id
  role               = "roles/iam.workloadIdentityUser"

  members = [
    "serviceAccount:fabra-344902.svc.id.goog[default/default]",
  ]
}

resource "google_kms_key_ring" "webhook-verification-key-keyring" {
  name     = "webhook-verification-key-keyring"
  location = "global"
}

resource "google_kms_crypto_key" "webhook-verification-key-key" {
  name            = "webhook-verification-key-key"
  key_ring        = google_kms_key_ring.webhook-verification-key-keyring.id
  rotation_period = "100000s"

  lifecycle {
    prevent_destroy = true
  }
}

resource "google_kms_key_ring_iam_binding" "webhook-verification-key-ring-binding" {
  key_ring_id = google_kms_key_ring.webhook-verification-key-keyring.id
  role        = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  members = [
    "serviceAccount:fabra-sync@fabra-344902.iam.gserviceaccount.com",
    "serviceAccount:fabra-backend@fabra-344902.iam.gserviceaccount.com"
  ]
}

resource "google_kms_key_ring" "end-customer-api-key-keyring" {
  name     = "end-customer-api-key-keyring"
  location = "global"
}

resource "google_kms_crypto_key" "end-customer-api-key-key" {
  name            = "end-customer-api-key-key"
  key_ring        = google_kms_key_ring.end-customer-api-key-keyring.id
  rotation_period = "100000s"

  lifecycle {
    prevent_destroy = true
  }
}

resource "google_kms_key_ring_iam_binding" "end-customer-api-key-ring-binding" {
  key_ring_id = google_kms_key_ring.end-customer-api-key-keyring.id
  role        = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  members = [
    "serviceAccount:fabra-sync@fabra-344902.iam.gserviceaccount.com",
    "serviceAccount:fabra-backend@fabra-344902.iam.gserviceaccount.com"
  ]
}

resource "google_kms_key_ring" "jwt-signing-key-keyring" {
  name     = "jwt-signing-key-keyring"
  location = "global"
}

resource "google_kms_crypto_key" "jwt-signing-key-key" {
  name     = "jwt-signing-key-key"
  key_ring = google_kms_key_ring.jwt-signing-key-keyring.id
  purpose  = "MAC"

  version_template {
    algorithm = "HMAC_SHA256"
  }

  lifecycle {
    prevent_destroy = true
  }
}

resource "google_kms_key_ring_iam_binding" "jwt-signing-key-ring-binding" {
  key_ring_id = google_kms_key_ring.jwt-signing-key-keyring.id
  role        = "roles/cloudkms.signerVerifier"
  members = [
    "serviceAccount:fabra-sync@fabra-344902.iam.gserviceaccount.com",
    "serviceAccount:fabra-backend@fabra-344902.iam.gserviceaccount.com"
  ]
}

resource "google_compute_security_policy" "fabra-security-policy" {
  name = "fabra-cloud-armor"

  adaptive_protection_config {
    layer_7_ddos_defense_config {
      rule_visibility = "STANDARD"
      enable          = true
    }
  }

  rule {
    action   = "throttle"
    preview  = false
    priority = 10

    match {
      versioned_expr = "SRC_IPS_V1"

      config {
        src_ip_ranges = [
          "*",
        ]
      }
    }

    rate_limit_options {
      ban_duration_sec = 0
      conform_action   = "allow"
      enforce_on_key   = "ALL"
      exceed_action    = "deny(429)"

      rate_limit_threshold {
        count        = 200
        interval_sec = 60
      }
    }
  }

  rule {
    action   = "allow"
    priority = "2147483647"
    match {
      versioned_expr = "SRC_IPS_V1"
      config {
        src_ip_ranges = ["*"]
      }
    }
    description = "Default rule, higher priority overrides it"
  }

  timeouts {}
}
