---
inclusion: auto
name: docker
description: Guide and best practices for Docker e Docker Compose. Apply this skill when creating Dockerfile, docker-compose-yml, or running docker commands from terminal.
---

# Best Practices for Writing `docker-compose.yml` Files

This guide outlines best practices for creating and maintaining `docker-compose.yml` files using **Docker Compose V2**. Following these recommendations helps ensure reliability, security, scalability, and long-term maintainability of containerized applications.

---

## 1. Do Not Specify the Compose File Version

- Do **not** declare the `version` field.
- The Compose specification is now **versionless**, and Docker automatically uses the latest supported schema.
- Explicit version declarations are deprecated and no longer required.

---

## 2. Keep Services Simple and Single-Responsibility

- Define **one service per responsibility** (e.g., application, database, cache).
- Avoid running multiple unrelated processes in the same container.
- This improves observability, scalability, and fault isolation.

---

## 3. Use Explicit and Immutable Image Tags

- Never use the `latest` tag.
- Pin images to a **specific version** (e.g., `postgres:16.2`).
- Prefer immutable digests (`@sha256:...`) for production when possible.
- Regularly update pinned versions to receive security patches.

---

## 4. Externalize Configuration and Secrets

- Use `.env` files for non-sensitive, environment-specific configuration.
- Never commit credentials or secrets to version control.
- For production:
  - Prefer **Docker secrets**
  - Or integrate with a **dedicated secret manager** (Vault, AWS Secrets Manager, etc.)

---

## 5. Prefer Named Volumes for Persistent Data

- Use named volumes for databases and other persistent state.
- Named volumes are easier to manage, migrate, back up, and inspect.
- Avoid bind mounts in production unless strictly necessary.
- Limit host path mounts to development environments.

---

## 6. Define Explicit Networks

- Create custom networks to control service communication.
- Separate internal and external traffic when possible.
- Attach services only to the networks they actually need.
- This improves security and reduces accidental coupling.

---

## 7. Add Healthchecks and Dependency Conditions

- Define `healthcheck` for critical services.
- Use `depends_on` with `condition: service_healthy` where supported.
- This reduces startup race conditions and improves reliability.
- Healthchecks also help with monitoring and automated recovery.

---

## 8. Ensure Portability and Reproducibility

- Avoid absolute, machine-specific paths.
- Use relative paths and environment variables instead.
- Do not rely on local system configuration or tooling.
- Aim for a setup that works consistently across machines and CI environments.

---

## 9. Use Multiple Compose Files for Different Environments

- Keep a base `docker-compose.yml` with shared configuration.
- Add environment-specific files, such as:
  - `docker-compose.override.yml` (local development)
  - `docker-compose.prod.yml` (production)
- This avoids duplication and keeps concerns clearly separated.

---

## 10. Limit Resource Usage

- Define resource constraints to prevent containers from exhausting host resources:
  - CPU
  - Memory
- This is especially important in shared or production environments.
- Resource limits improve system stability and predictability.

---

## 11. Improve Security Posture

- Avoid running containers as `root` whenever possible.
- Expose only the ports that are strictly necessary.
- Prefer internal networking over published ports.
- Read-only filesystems and dropped capabilities are recommended for hardened setups.

---

## 12. Keep Files Organized and Well Documented

- Group related services together.
- Use consistent key ordering and indentation.
- Add comments explaining **why** decisions were made, not just **what** they do.
- Clear documentation reduces onboarding and maintenance costs.

---

## 13. Prefer `docker compose` Over `docker-compose`

- Always use the modern CLI:
  ```bash
  docker compose up
  ```

* The legacy `docker-compose` command is deprecated.

---

## 14. Review and Maintain Regularly

- Periodically:

  - Update image versions
  - Rotate secrets
  - Review exposed ports and volumes

- Test changes in staging environments before production rollout.

---

## Summary

When writing `docker-compose.yml` files:

- Avoid deprecated features and legacy syntax
- Be explicit with image versions and configurations
- Externalize configuration and protect secrets
- Design for portability, security, and multiple environments
- Use healthchecks, networks, and resource limits
- Keep files clean, documented, and easy to evolve

Adopting these practices results in Compose files that are **predictable, secure, maintainable, and production-ready**.