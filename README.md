# StartMeUp – Kubernetes Infra Demo for a SaaS starter kit in Go
[![Edka](https://img.shields.io/badge/Try%20It%20On%20A%20Free%20Cluster-4c1?style=for-the-badge)](https://console.edka.io/signup)
[![Kubernetes](https://img.shields.io/badge/K3S%20Kubernetes-1.32-blue?logo=kubernetes&style=for-the-badge)](https://k3s.io)
[![Hetzner](https://img.shields.io/badge/Hetzner%20Cloud-red?logo=hetzner&style=for-the-badge)](https://www.hetzner.com/cloud)
[![GitOps](https://img.shields.io/badge/GitOps-Flux-blue?logo=flux&style=for-the-badge)](https://fluxcd.io)
[![OctoDNS](https://img.shields.io/badge/GitDNS-purple?logo=octodns&style=for-the-badge)](dns)
[![Preview Environments](https://img.shields.io/badge/Preview%20Envs-blue?style=for-the-badge)](charts/preview/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg?style=for-the-badge)](LICENSE)

This repository is a modified fork of [Pagoda](https://github.com/mikestefanello/pagoda), a full-stack web development starter kit and admin panel in Go. 

Here we are using it as a reference stack that showcases how the **[Edka](https://edka.io)** platform provisions and manages the infrastructure with GitOps, creates preview environments for every pull request and deploys code to production. Use it as a fast-start template or simply explore how everything fits together.

---

## Table of Contents
1. [Key Advantages](#key-advantages)
2. [What's in the Box?](#whats-in-the-box)
3. [Prerequisites](#prerequisites)
4. [Status & Contributing](#status--contributing)
5. [Credits](#credits)
6. [License](#license)

---

## Key Advantages
- **Cloud freedom, PaaS simplicity** – Provision a hardened Kubernetes cluster on Hetzner in minutes, manage it through a simple UI and Git.
- **Cost Efficiency** — Hetzner VMs and storage are often 80–90 % cheaper than comparable AWS, Azure, or GCP offerings.
- **Self-Service Deployments** — Push to Git, run automated tests, get a preview environment per pull request, merge and ship.
- **Extensible Add-ons** — Pick pre-configured add-ons for databases, ingress, observability, CI/CD, and more.
- **Open Standards** — Built on CNCF projects like k3s, Flux, cert-manager, CloudNativePG, External Secrets.

> 📖 Read the full story in our blog post: **[Join the Edka Beta (1 free cluster)](https://edka.io/blog/join-our-beta/)**.

---

## Demo Video

[![Watch the demo video](https://img.youtube.com/vi/-ybRo6qhnJ0/maxresdefault.jpg)](https://www.youtube.com/watch?v=-ybRo6qhnJ0)


## What's in the Box?
- **Web, API & background workers** → [`charts/startmeup`](charts/startmeup)
- **Database** – PostgreSQL with PITR & S3 backups → [`clusters/resources/startmeup/postgres`](clusters/resources/startmeup/postgres)
- **Background jobs** – [River](https://github.com/riverqueue/river) powered by Postgres
- **Versioned Schema management** – [Ent](https://entgo.io/)
- **CI/CD** – GitHub Actions examples for tests, preview environments, DNS management, and production deployments → [`.github`](.github/workflows/)
- **GitOps** – Declarative infra with Flux → [`clusters/resources/startmeup`](clusters/resources/startmeup)
- **Secrets** – [Doppler](https://www.doppler.com/) + External Secrets → [`clusters/resources/startmeup/secrets`](clusters/resources/startmeup/secrets)
- **Ingress & HTTPS** – NGINX + cert-manager & Let's Encrypt
- **DNS** – Git-stored records synced to Cloudflare → [`dns`](dns)
- **Preview Environments** – One per pull request → [`charts/preview`](charts/preview)
---

## Prerequisites
To try this demo, you'll need:

- **Hetzner Cloud** project and API token
- **Domain name** you control (for HTTPS & DNS)
- **Secrets provider** (Doppler, AWS Secrets Manager, 1Password, Vault, etc.)
- **Container registry** (GitHub Container Registry, Docker Hub, AWS ECR, Google Artifact Registry)
- **S3-compatible storage** for Postgres backups (AWS S3, Wasabi, MinIO...)
- **GitHub repository** containing your application code & manifests

---

## Status & Contributing
> **Status:** Demo / Beta. Not production-supported.
>
> Please open an issue if you spot a problem — PRs are welcome!

---

## Credits
- Original project: **[Pagoda](https://github.com/mikestefanello/pagoda)** by Mike Stefanello — thanks for the fantastic groundwork.

---

## License
This project is licensed under the [MIT License](LICENSE).
