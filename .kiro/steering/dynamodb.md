---
inclusion: auto
name: dynamodb
description: Guide and best practices for modeling DynamoDB or interacting with it. Apply this skill when modeling the database or writing backend code that access the DynamoDB.
---

# DynamoDB Single Table Design – Steering Document

This document defines the guiding principles and best practices for modeling and interacting with Amazon DynamoDB using the **Single Table Design** approach. It serves as a reference for architects and developers to ensure consistency, scalability, predictable performance, and cost efficiency across all services.

This document is **mandatory** for all teams using DynamoDB.

---

# 🎯 Goals

- Achieve **horizontal scalability** and **predictable performance**
- Minimize the number of tables by using **one table per domain**
- Model data based on **access patterns**, not entities
- Ensure all reads use **Query**, never Scan
- Optimize for **low latency and low cost**
- Keep the design **evolvable** over time

---

# 🧠 Core Principles

## 1. Start With Access Patterns

The schema MUST be derived from application read requirements.

Before creating or modifying a table, teams must document:

- What data is read?
- How frequently?
- With what filters?
- Sorted by what attribute?
- What is the expected cardinality?

If an access pattern cannot be served with `Query`, the schema must be redesigned.

---

## 2. One Table per Domain

All related entities must live in a single table.

Cross-entity relationships are modeled via key design — not joins.

---

## 3. Composite Primary Key

Each table must use:

- `PK` (Partition Key)
- `SK` (Sort Key)

The sort key enables:
- Hierarchical modeling
- Entity grouping
- Efficient range queries
- Time-based ordering
- begins_with patterns

---

## 4. Keys Encode Meaning

Keys must be structured and descriptive:

Examples:
- `USER#<id>`
- `ORDER#<id>`
- `ITEM#<id>`
- `CATEGORY#<name>`

Never use opaque UUIDs without semantic prefixes.

---

## 5. Denormalize Aggressively

Duplicate data to optimize reads.

DynamoDB is optimized for read efficiency — not normalization.

---

## 6. Sparse GSIs Only

Global Secondary Indexes (GSIs) must:

- Support a clearly defined access pattern
- Contain only items that need to be indexed
- Not replicate the full table

---

# 🚨 Application Interaction Rules (MANDATORY)

This section defines how applications MUST interact with DynamoDB.

---

## 🔍 1. ALWAYS Prefer `Query` Over `Scan`

Applications MUST use `Query` whenever possible.

`Scan` is considered a design failure except for:

- One-time migrations
- Backfills
- Admin/debug scripts (never production paths)

### Why?

| Operation | Performance | Cost | Scalability |
|------------|------------|-------|-------------|
| Query | O(1) per partition | Low | Predictable |
| Scan | O(N) entire table | Expensive | Dangerous |

`Scan`:
- Reads every item in the table
- Consumes large RCUs
- Causes latency spikes
- Can create hot partitions
- Breaks cost predictability

If a feature requires `Scan`, STOP and redesign the schema.

---

## ✅ 2. All Reads Must Be Key-Based

Valid read patterns:

- `PK = value`
- `PK = value AND SK = value`
- `PK = value AND begins_with(SK, prefix)`
- GSI queries using key conditions

Invalid pattern:

- Filtering without key condition
- Using `FilterExpression` to reduce large result sets

### Rule

`FilterExpression` must NEVER be used to compensate for poor key design.

Filtering should only reduce already-small query results.

---

## 🧾 3. Use `Query` With Explicit KeyConditionExpression

Applications must always:

- Provide `KeyConditionExpression`
- Avoid broad partition reads unless intentionally grouped
- Limit result sets with `Limit`
- Use pagination correctly (`LastEvaluatedKey`)

Never rely on unbounded reads in high-traffic APIs.

---

## 📉 4. Avoid Full Partition Reads When Not Necessary

If a partition can grow large (e.g., user with millions of events):

- Use time-based bucketing
- Use entity sharding
- Introduce secondary access patterns

---

## 🧮 5. Understand RCU / WCU Cost Model

Developers must understand:

- `Query` cost depends on item size and consistency
- Strongly consistent reads double RCU cost
- Large items increase cost
- ProjectionExpression reduces network payload (but not RCU consumption)

Design decisions must consider cost impact.

---

## 🧵 6. Use ProjectionExpression When Appropriate

To reduce response size and network overhead:

- Use `ProjectionExpression`
- Avoid returning full items when unnecessary

---

## 🔐 7. Prefer Eventually Consistent Reads

Use strongly consistent reads only when strictly required.

Eventual consistency:
- Lower cost
- Higher scalability
- Better latency

---

## 🛑 8. No Relational Thinking

Applications must NOT:

- Perform multiple round trips to simulate joins
- Fetch large datasets and filter client-side
- Expect ad-hoc querying flexibility

DynamoDB requires **predefined access patterns**.

---

# 🗂 Example Domain: E-Commerce

## Access Patterns

| Use Case | Query Pattern |
|-----------|--------------|
| Get user profile | `PK = USER#<id>, SK = PROFILE` |
| List user orders | `PK = USER#<id>, SK begins_with ORDER#` |
| Get order details | `PK = ORDER#<id>, SK begins_with ITEM#` |
| List products by category | `GSI1PK = CATEGORY#<name>` |
| Get all items in an order | `PK = ORDER#<id>` |

All of the above are supported using `Query`.

None require `Scan`.

---

# 🧱 Table Schema

| Attribute | Purpose |
|-----------|----------|
| `PK` | Partition Key |
| `SK` | Sort Key |
| `GSI1PK` | GSI #1 Partition Key |
| `GSI1SK` | GSI #1 Sort Key |
| `Type` | Logical entity type |
| `SchemaVersion` | Version control for evolution |

---

# ⚠️ Anti-Patterns

- ❌ Creating one table per entity
- ❌ Using random UUID partition keys without structure
- ❌ Designing schema before defining access patterns
- ❌ Using `Scan` in production
- ❌ Heavy use of `FilterExpression`
- ❌ Client-side filtering
- ❌ Treating DynamoDB like a relational database

---

# 🔄 Evolution Strategy

- Add new item types without rewriting old ones
- Add new GSIs to support new access patterns
- Version item shapes (`SchemaVersion`)
- Prefer additive, backward-compatible changes

If a new feature requires `Scan`, the table design must be re-evaluated.

---

# 📋 Mandatory Design Review Checklist

Before approving a DynamoDB model:

Access Patterns
- [ ] All access patterns documented
- [ ] Each access pattern maps to a `Query`
- [ ] No production `Scan` required

Key Design
- [ ] PK/SK encode entity + relationship
- [ ] No hot partition risks
- [ ] High-cardinality partition keys

Cost & Performance
- [ ] Expected item size reviewed
- [ ] RCU/WCU estimation documented
- [ ] Pagination strategy defined
- [ ] Strong consistency justified (if used)

Indexes
- [ ] GSIs support explicit access patterns
- [ ] GSIs are sparse
- [ ] No unused indexes

---

# 📚 References

- Amazon DynamoDB Developer Guide – Data Modeling
- Amazon DynamoDB Best Practices for Designing and Using Partition Keys
- The DynamoDB Book – Alex DeBrie
- AWS re:Invent – Advanced Data Modeling with DynamoDB

---

# 🚨 Final Rule

If you cannot answer the question:

> “What exact Query will this feature execute?”

The data model is not ready.

This steering document must be followed by all teams designing new DynamoDB tables or evolving existing ones. Deviations require architectural review and explicit approval.