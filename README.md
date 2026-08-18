# Ubuntu Server — Infrastructure, Security & Monitoring

Academic project focused on designing a resilient Linux server environment with secure networking, remote administration and system monitoring.

## Overview

The project demonstrates how Ubuntu Server can be used as the central component of a small infrastructure while addressing availability, integrity, reliability and security.

## Objectives

- Deploy an Ubuntu Server environment.
- Design a basic network architecture for the server.
- Configure secure remote administration with SSH.
- Apply firewall and brute-force protection.
- Identify infrastructure risks and define mitigation measures.
- Explore server monitoring and observability concepts.

## Environment

- Ubuntu Server 24.04 LTS
- Ethernet / Wi-Fi networking
- SSH
- UFW
- Fail2Ban
- Docker-based monitoring concepts
- Prometheus
- Grafana
- Grafana Loki
- Node Exporter
- cAdvisor

## Network Architecture

```text
Internet
   │
Router
   │
Firewall (UFW)
   │
Switch
   │
Ubuntu Server
   │
   ├── SSH
   ├── Services
   └── Monitoring Agents
           │
           ├── Node Exporter ──> Prometheus
           ├── cAdvisor ───────> Prometheus
           └── Log Collector ──> Loki
                                  │
                              Grafana
```

## Security

### SSH

Remote administration is protected using SSH, with public-key authentication recommended instead of password-only access.

### UFW

UFW provides host-level firewall rules controlling which services and ports are exposed.

### Fail2Ban

Fail2Ban can monitor authentication logs and temporarily block IP addresses that repeatedly fail authentication.

## Risk Matrix

| Risk | Example | Impact | Mitigation |
|---|---|---|---|
| Hardware failure | Power loss | High | UPS and backups |
| Software failure | Incorrect permissions/configuration | High | Updates, least privilege and backups |
| Security incident | SSH brute-force attempt | High | UFW, Fail2Ban and SSH keys |
| Network failure | Connectivity loss | Medium/High | Redundancy and network monitoring |

## Basic Commands

```bash
sudo apt update
sudo apt upgrade -y

sudo ufw status verbose
sudo ufw allow OpenSSH
sudo ufw enable

sudo systemctl status ssh
sudo systemctl status fail2ban
```

## Monitoring Concept

The monitoring architecture separates metric collection, log collection and visualization. Prometheus stores time-series metrics, Loki stores logs, and Grafana provides dashboards and visualization.

## What I Learned

This project provided practical experience with Linux server administration, network architecture, security controls, risk analysis and the fundamentals of observability.

## Author

Lucas Vinhal Ferreira — Computer Engineering student at UniCEUB.

[GitHub](https://github.com/romunac) · [LinkedIn](https://www.linkedin.com/in/lucas-vinhal-24334a2b3)
