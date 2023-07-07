# Fabra

## Introduction

To deploy this system in production, there are only a few manual steps:
- Setup a GCP Project
- Create Terraform Google Cloud Storage bucket
- Perform initial Terraform run

### Setup GCP Project
Create a new GCP project and take note of the name and project ID.

### Create Terraform GCP Cloud Storage bucket
This Google Cloud Bucket is used as a backing store for Terraform, so must be manually setup. Note the name you use
as you'll need to edit `infra/terraform/main.tf` to point to it.

### Initial Terraform run
Once those two things are setup, everything else can be configured with Infrastructure-as-code in Terraform. However,
you must run `terraform apply` once manually to setup the correct permissions for the Cloud Build service account to
run Terraform automatically in the future.

## Deploy to new region

1. Create a new GCP project.
1. Create a new Terraform file for the region by copying infra/terraform/main.tf into a new subdirectory
    1. Modify the project ID in the new Terraform file
    1. Modify the Cloud Storage bucket names to match the new region
    1. Run `terraform init` in the new subdirectory
1. Enable all the GCP APIs needed:
    1. Cloud Build
    1. Cloud Engine
    1. Cloud Run
    1. Cloud SQL
    1. IAM
    1. KMS
    1. DNS
    1. Secret Manager
    1. Serverless VPC Access
    1. Service Networking
    1. Artifact Registry
    1. Kubernetes Engine
1. Create a new DB password in the new projects Secret Manager with the name `fabra-db-password`
1. Create a new Terraform bucket in Cloud Storage and add it to the Terraform file
1. Connect the Github repository to the new GCP project
1. Copy OAuth secrets to the new project's Secret Manager and ensure the code references them correctly
1. Enable Cloud Build to deploy to Cloud Run:

        gcloud iam service-accounts add-iam-policy-binding \
          fabra-backend@fabra-prod.iam.gserviceaccount.com \
          --member="serviceAccount:fabra-prod@cloudbuild.gserviceaccount.com" \
          --role="roles/iam.serviceAccountUser"
1. Run `terraform apply`

### Other Notes
Google Cloud Build is used for a various automatic actions triggered by pushes to the main Github branch:
- Run Terraform to build any new infrastructure
- Build Docker image for the Go code and deploy it to GCR
- Run database migrations
