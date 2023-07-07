# Fabra

## Introduction

This super basic setup should work for any Go project.

You only need to manually setup a few things:
- GCP Project
- Google Cloud Bucket
- Google Cloud Build
    - Special service account with access to create/modify/delete GCP resources

The Google Cloud Bucket is used as a backing store for Terraform, so must be manually setup.

Google Cloud Build is used for a various automatic actions triggered by pushes to the main Github branch:
- Run Terraform to build any new infrastructure
- Build Docker image for the Go code and deploy it to GCR
- Run database migrations

Once those two things are setup, everything else can be configured with Infrastructure-as-code in Terraform,
including:
- Google Cloud Repository bucket
- Cloud Run services
    - Includes setting any environment variables and Cloud SQL connections
- Cloud SQL instances
