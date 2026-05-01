---
inclusion: auto
name: infra
description: Guide and best practices for write IaC. Apply this skill when writing IaC.
---

````md
# AWS Infrastructure Guidelines (Claude Code Edition — DynamoDB + LocalStack)

## 1 — Scope & high-level rules

* All infrastructure **must** be provisioned on **AWS**.
* Infrastructure code **must** be written using **AWS CDK (TypeScript)**.
* Observability must include **CloudWatch** (logs & metrics) and **X-Ray** (tracing).
* AWS services to use:
  * **Lambda** — backend REST API.
  * **API Gateway** — HTTP proxy to Lambda.
  * **CloudWatch** — logging & metrics.
  * **X-Ray** — distributed tracing.
  * **Amazon DynamoDB** — primary database.
* Local development **must** use **LocalStack**.

> Rationale: DynamoDB removes connection management and scales automatically.  
> LocalStack provides a high-fidelity AWS emulation, allowing the same CDK stacks and SDK code paths to be used locally and in AWS.

---

## 2 — Development rules

* Use **LocalStack** for all local AWS emulation (DynamoDB, Lambda, API Gateway, CloudWatch).
* **Run backend locally**.
* **Run frontend locally**.
* Use AWS CLI/CDK profile **`josimar`** by default.
* All data-access code **must**:
  * Be environment-agnostic.
  * Use table names via environment variables.
  * Avoid hard-coded AWS regions or ARNs.
* No production code may rely on “local-only” branches beyond endpoint configuration.

---

## 3 — Environment-variable conventions

All scripts (test, synth, deploy, destroy, local-start) must accept:

* `env` — environment name (default: `dev`)
* `profile` — AWS CLI/CDK profile (default: `josimar`)

**Behavior:** CDK stacks are deployed the same way for all environments.

Common runtime variables for Lambdas:

* `ENV`
* `AWS_REGION`
* `DYNAMODB_TABLE_MAIN`
* `AWS_ENDPOINT_URL` (only for local dev with LocalStack)

Example local values:

```env
ENV=local
AWS_REGION=us-east-1
AWS_ENDPOINT_URL=http://localhost:4566
DYNAMODB_TABLE_MAIN=myapp-main-local
````

---

## 4 — CDK-specific rules & best practices

* Use **TypeScript** and CDK v2.
* Structure CDK code into modular stacks:

  * `NetworkStack` — VPC, subnets, security groups (only if needed).
  * `DatabaseStack` — DynamoDB tables.
  * `BackendStack` — Lambda functions, IAM roles, env vars.
  * `ApiStack` — API Gateway & Lambda integration.
  * `ObservabilityStack` — CloudWatch dashboards, X-Ray.
* Use CDK context or environment variables for `env` and `profile`.
* Use **AWS Solutions Constructs** and **CDK Nag**.
* Centralize IAM least-privilege policies.
* Prefer `table.grantReadWriteData(fn)` over custom wildcard IAM policies.
* Parameterize stack names and tags with environment.
* Include unit & integration tests for CDK constructs.

**DynamoDB standards:**

* Use **On-Demand (PAY_PER_REQUEST)** billing by default.
* Define:

  * Partition key (required).
  * Sort key when modeling time-series or multi-entity patterns.
* Use **single-table design** when applicable.
* Prefer **GSIs** over table scans.
* Enable **Point-in-Time Recovery (PITR)** for all environments.
* Enable **TTL** for ephemeral data.
* Use **DynamoDB Streams** only when required.

---

## 5 — Local AWS (LocalStack) recommendations

* Use `localstack/localstack` Docker image.
* Provide a `docker-compose.yml` with:

  * Port `4566`.
  * Persistent volume.
  * Enabled services: `dynamodb`, `lambda`, `apigateway`, `cloudwatch`, `iam`, `logs`.
* Backend configuration:

  * Use `AWS_ENDPOINT_URL=http://localhost:4566` in local mode.
  * Use the **same AWS SDK code** for local and cloud.
  * Never auto-create tables in AWS; CDK is the single source of truth.
* CDK must be deployable to LocalStack using the same stacks.

Example `docker-compose.yml` snippet:

```yaml
services:
  localstack:
    image: localstack/localstack:latest
    ports:
      - "4566:4566"
    environment:
      - SERVICES=dynamodb,lambda,apigateway,cloudwatch,iam,logs
      - DEFAULT_REGION=us-east-1
      - DEBUG=1
    volumes:
      - ./.localstack:/var/lib/localstack
      - /var/run/docker.sock:/var/run/docker.sock
```

---

## 6 — Scripts (requirements)

Scripts must:

* Accept `env` and `profile`, with defaults (`dev` and `josimar`).
* Validate `env` values (`local`, `dev`, `staging`, `prod`).
* Print summary of created resources.
* Be executable via **terminal** and via **`npm run`**.

**Required scripts:**

* `test` — run lint, backend & frontend tests, CDK synth/Nag.
* `synth` — generate CloudFormation templates.
* `deploy` — deploy stacks for the given environment.
* `destroy` — destroy stacks for the given environment.
* `local-start` — start LocalStack, deploy stacks to LocalStack, run backend & frontend locally.

`local-start` must:

1. Start LocalStack via Docker.
2. Wait for health.
3. Deploy all CDK stacks to LocalStack.
4. Export required env vars.
5. Start backend & frontend.

---

## 7 — Environment parity

* All environments deploy the full set of stacks:
  `DatabaseStack`, `BackendStack`, `ApiStack`, `ObservabilityStack`, `NetworkStack` (if used).
* No conditional skipping based on environment.
* Only **configuration values** may differ (e.g., table names, log retention).

---

## 8 — Observability

* Enable CloudWatch Logs for all Lambdas.
* Enable X-Ray tracing for Lambda and API Gateway.
* Create CloudWatch dashboards for:

  * Lambda errors & duration.
  * API Gateway latency & 4xx/5xx.
  * DynamoDB:

    * ThrottledRequests
    * ConsumedReadCapacityUnits
    * ConsumedWriteCapacityUnits
    * SystemErrors
* Use structured JSON logs with:

  * `requestId`
  * `env`
  * `service`
  * `operation`

---

## 9 — Security & compliance

* Integrate **CDK Nag** rules in CI.
* Use least-privilege IAM roles.
* Encrypt data at rest and in transit (default for DynamoDB).
* Use Secrets Manager for secrets.
* Tag resources by `Owner`, `Environment`, `Project`, `CostCenter`.
* Never allow public access to tables.

---

## 10 — CI/CD recommendations

* CI pipeline should:

  1. Run tests & CDK synth with CDK Nag.
  2. Run integration tests against LocalStack or an ephemeral AWS environment.
  3. Promote artifacts to staging and prod with manual approvals.
* Store CDK context and long-lived parameters securely.

---

## 11 — Developer onboarding checklist

1. Install Node.js LTS, Yarn/npm, AWS CLI, Docker.
2. Add AWS profile `josimar`.
3. Run `npm run local:start` to start LocalStack and deploy all stacks locally.
4. Use `npm run synth` to confirm CDK app synthesizes correctly.
5. Use **Claude Code** to generate scaffolding for CDK constructs, scripts, and boilerplate safely, ensuring compliance with these guidelines.

---

## 12 — Recommended project layout

.
├─ package.json
├─ docker-compose.yml
├─ infra/
│ ├─ bin/
│ └─ lib/
├─ scripts/
├─ backend/
│ ├─ src/
│ └─ jest.config.js
├─ frontend/
├─ docs/
│ └─ data-model.md

`docs/data-model.md` must describe:

* Table schemas.
* Primary keys and GSIs.
* Access patterns (read/write paths).

---

## 13 — Extra best practices

* Automate CDK checks and linting in pre-commit hooks.
* Document all DynamoDB access patterns before implementation.
* Avoid scans in production code.
* Use batch operations for bulk writes.
* Set explicit timeouts and memory sizing for Lambdas.
* Use **Claude Code** for generating boilerplate code for new stacks, Lambda functions, and deployment scripts, but **always review generated code** before committing.

---

## 14 — CDK inclusion example

```ts
const app = new cdk.App();
const envName = app.node.tryGetContext('env') ?? process.env.ENV ?? 'dev';

new NetworkStack(app, `NetworkStack-${envName}`, { env: cdkEnv });
new DatabaseStack(app, `DatabaseStack-${envName}`, { env: cdkEnv });
new BackendStack(app, `BackendStack-${envName}`, { env: cdkEnv });
new ApiStack(app, `ApiStack-${envName}`, { env: cdkEnv });
new ObservabilityStack(app, `ObservabilityStack-${envName}`, { env: cdkEnv });
```

---

## 15 — Claude Code usage

* Use **Claude Code** as a productivity assistant for:

  * Generating CDK scaffolding and boilerplate.
  * Suggesting scripts that conform to `env` and `profile` conventions.
  * Following AWS CDK best practices and security guidelines.
* Always verify generated code against internal standards and **CDK Nag rules**.