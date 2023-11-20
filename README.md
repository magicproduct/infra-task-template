# Magic AI Infrastructure Task

# Scenario

Data engineers use this piece of code to read local files (located in the `assets` folder), scrub all PII such as user
IDs from them and then upload them to a GCS bucket.

A `verification.json` is uploaded too, this file contains validation info about the uploaded files.

# Task Instructions

- Use GitHub and invite us to your private repository
- Use Terraform for all infrastructure you provision (GCS, GKE, ...)
- Use GCP for all cloud infrastructure, no service outside of GitHub, GCP and Terraform Cloud (if you wish) should be
  used
- Take great care of security and harden all infrastructure and code you submit
- Deploy the provided code to a GKE cluster and make sure it runs, fix bugs as necessary

# Goal

The application should run in the GKE cluster without errors, all files should be scrubbed and uploaded to the GCS
bucket.