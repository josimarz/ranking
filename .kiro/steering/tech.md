---
inclusion: always
---

# Steering Document – Technology Stack for **Ranking**

This document defines the official technology stack and infrastructure standards that must be used to implement and deploy the **Ranking** product.

It serves as a single source of truth for any AI agent, developer, or system involved in building this application.

The primary goals of this stack are:

- Simplicity and low operational overhead  
- Low cost at small scale  
- Fast iteration during development  
- Full alignment with the product’s “no friction” philosophy  

---

## General Instruction – Documentation via MCP Context7

Whenever implementing, configuring, or reasoning about any library, framework, SDK, or tool referenced in this document (or introduced later), **you must use the MCP `context7` server to retrieve up-to-date official documentation**.

Rules:

- Always query `context7` before:
  - Using APIs from a library or framework  
  - Writing configuration for tools such as AWS CDK, Gin, Next.js, Swagger, DynamoDB SDKs, or LocalStack  
  - Making assumptions about versions, defaults, or best practices  
- Treat information from `context7` as the **authoritative source of truth**.
- Do not rely on:
  - Training-time knowledge  
  - Memory of previous projects  
  - Outdated examples from the web  
- If there is any ambiguity or version mismatch, resolve it by re-querying `context7`.

All implementation decisions must be aligned with the **current official documentation** as returned by MCP `context7`.

---

## 1. Backend

### 1.1 Language and Framework

The backend must be implemented as a **REST API** using:

- **Go** — version **1.25.6**
- **Gin-Gonic** framework  
  https://gin-gonic.com/en/

Gin must be used for:

- HTTP routing  
- Request validation  
- Middleware  
- JSON serialization  

The API must be fully stateless. All persistence must be handled via DynamoDB and S3.

---

### 1.2 API Documentation (Swagger / OpenAPI)

The backend **must provide a complete and up-to-date Swagger (OpenAPI) documentation** for the REST API.

This documentation must be implemented using:

- **swaggo/gin-swagger**  
  https://github.com/swaggo/gin-swagger  

Rules:

- All endpoints must be documented using Swag annotations in the Go source code.
- Every route must define:
  - Summary and description  
  - Request parameters  
  - Request body (when applicable)  
  - Response schemas  
  - HTTP status codes  
- The Swagger UI must be exposed by the backend in all environments:
  - Local development
  - Production (protected by API Key)

An endpoint that is not documented in Swagger **must be considered incomplete**.

---

### 1.3 Runtime Environment

The backend must be deployed on **AWS Lambda** with the following configuration:

- Architecture: **ARM64**
- Operating System: **Amazon Linux 2023**
- Memory: **128 MB**

Each Lambda function must be:

- Small and focused  
- Cold-start optimized  
- Designed for fast execution  

The backend must be exposed via **Amazon API Gateway**.

---

### 1.4 API Gateway & Security

- The REST API must be published through **API Gateway**.
- Access must be protected using an **API Key** in production environments.
- In **local development mode**, the API Key must be **disabled**.

---

## 2. Data Storage

### 2.1 Primary Database – DynamoDB

The application must use **Amazon DynamoDB** as its only database for structured data:

- Rankings  
- Attributes  
- Items  
- User identifiers  
- User ratings  

There must be no relational database, no ORM, and no secondary structured storage layer.

---

### 2.2 Mandatory Data Modeling Strategy – Single Table Design

DynamoDB **must be modeled using the Single Table Design pattern**.

This is a **strict requirement**, not an implementation detail.

Rules:

- The application must use **one single DynamoDB table** for all entities.
- All access patterns must be defined **before** defining the table structure.
- Data must be modeled using:
  - Composite primary keys (PK + SK)
  - Entity types encoded in keys and/or attributes
  - Sparse indexes where needed
- Secondary indexes (GSI) may be used, but must follow the same access-pattern-first approach.

Designing multiple tables for separate entities (e.g., `rankings`, `items`, `users`, `ratings`) is **explicitly forbidden**.

If new features introduce new entities, they must be incorporated into the same table.

---

### 2.3 Object Storage – Images (Amazon S3)

All user-provided images for items **must be stored in Amazon S3**.

The backend is fully responsible for handling image ingestion, optimization, and storage.

Rules:

- When the user provides an **image URL**:
  - The backend must download the image from that URL.
  - The downloaded image must be processed and stored in S3.

- When the user uploads an **image file directly**:
  - The backend must receive the file.
  - The image must be processed and stored in S3.

- Before storing any image:
  - The backend **must optimize the image**.
  - The backend **must generate a thumbnail version**.
  - Both the optimized original and the thumbnail must be stored in S3.

- The frontend must **never** access S3 objects directly.
- To display an image, the backend must generate a **pre-signed URL** and return it to the frontend.

Objectives:

- Reduce bandwidth and storage costs  
- Ensure consistent image sizes and formats  
- Keep S3 buckets private  
- Prevent public access to raw objects  
- Maintain full control over asset delivery  

S3 is the **only** allowed storage for user-provided images.

---

### 2.4 Local Development

For local development and testing:

- Use **LocalStack** to simulate:
  - DynamoDB
  - S3
- The backend must be able to run fully offline using this setup.

---

## 3. Infrastructure as Code

All backend infrastructure must be provisioned using:

- **AWS CDK**
- Language: **TypeScript**

The CDK project must define:

- Lambda functions  
- API Gateway  
- DynamoDB table (single table)  
- S3 buckets for images  
- IAM roles and policies  
- Environment variables  

No manual configuration in the AWS Console is allowed for production resources.

---

## 4. Frontend

### 4.1 Framework

The frontend must be implemented using:

- **Next.js** — LTS version **16.1.3**  
  https://nextjs.org/

---

### 4.2 Hosting

The frontend must be hosted using **AWS Amplify**.

The frontend must consume:

- The production API Gateway endpoint  
- The local API endpoint during development  
- Pre-signed image URLs provided by the backend  

---

## 5. Environment Separation

The system must support at least two environments:

- **Local Development**
- **Production**

Rules:

- Local:
  - No API Key required
  - Uses LocalStack for DynamoDB and S3
  - Backend runs locally
  - Frontend runs via `next dev`
  - Swagger UI is publicly accessible

- Production:
  - API Key is mandatory
  - Uses real AWS services
  - Backend runs on Lambda
  - Frontend runs on Amplify
  - Swagger UI is protected by API Key

---

## 6. Non-Goals

The following are explicitly **out of scope** for this product:

- Authentication systems (OAuth, login, passwords, sessions)
- Relational databases
- Multiple DynamoDB tables
- Public S3 buckets
- Undocumented APIs
- Microservices architecture
- Message queues or event buses
- Kubernetes or container orchestration
- Multi-region deployments

---

This document defines **how the product must be built**, just as the Product Steering Document defines **what the product must do**.  
Any deviation from this stack — especially the use of multiple DynamoDB tables, public S3 access, or undocumented endpoints — must be treated as a product-level decision.