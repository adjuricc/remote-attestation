# Remote Attestation in Google Compute Engine (GCP)
## Overview
This Go project implements remote attestation in Google Compute Engine (GCE) for Confidential VMs leveraging AMD SEV-SNP and Intel TDX technologies. It consists of two important roles:
- Verifier:
  - Represents a user who wants to perform attestation on expected measurements, in order to verify and validate the attestation report
  - Uses `go-sev-guest` and `go-tdx-guest` libraries for verification and validation of the requested attestation report
- Attester:
  - Represents a GCE VMs with enabled Confidential Computing that generates launch endorsement and attestation report
  - Uses `gce-tcb-verifier` library for launch endorsement generation
  - Types of attester:
    - Attester with enabled AMD SEV-SNP technology- uses `go-sev-guest` library for attestation report generation
    - Attester with enabled Intel TDX technology- uses `go-tdx-guest` library for attestation report generation

## Features
- End-to-end remote attestation flow
- Support for AMD SEV-SNP and Intel TDX attestation
- Secure generation and validation of attestation reports and launch endorsements
- Written in Go
