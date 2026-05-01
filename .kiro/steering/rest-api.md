---
inclusion: auto
name: rest-api
description: Best practices for writing Rest API. Apply this skill when writing code in the backend.
---

# Best Practices for Writing REST APIs

Creating a well-designed REST API is key to building scalable, maintainable, and reliable web services. This document outlines the core best practices to follow when developing RESTful APIs in 2025.

## Use Clear and Consistent Resource Naming
- Use **nouns** for endpoint paths, not verbs. For example, use `/users` instead of `/getUsers`.
- Name collections with **plural nouns**, e.g., `/products`, `/orders`.
- Use nested resource paths to indicate hierarchy or relationship, e.g., `/users/{id}/orders`.
- Keep paths simple, consistent, and intuitive to understand without reading docs.

## Use Proper HTTP Methods
- Use the appropriate HTTP methods for CRUD operations:
  - `GET` to retrieve resources
  - `POST` to create new resources
  - `PUT` to update existing resources
  - `PATCH` for partial updates
  - `DELETE` to remove resources

## Use Standard HTTP Status Codes
- Respond with correct status codes to indicate the outcome:
  - `2xx` for success (200 OK, 201 Created)
  - `3xx` for redirection
  - `4xx` for client errors (404 Not Found, 400 Bad Request)
  - `5xx` for server errors (500 Internal Server Error)

## Use JSON for Request and Response Bodies
- Accept and respond with **JSON** as the default data format for wide compatibility and easy parsing.

## Implement Security Best Practices
- Use **HTTPS** to encrypt data in transit.
- Implement secure authentication and authorization methods, e.g., OAuth 2.0, JWT.
- Validate all inputs to avoid injection attacks.
- Use CORS policies suitably for cross-origin requests.

## Support Filtering, Sorting, and Pagination
- Allow clients to filter and sort data via query parameters.
- Implement pagination to handle large datasets efficiently and reduce payload sizes.

## Version Your API
- Use **semantic versioning** in the URL or headers, e.g., `/api/v1/products`.
- Clearly document version changes and maintain backward compatibility.

## Provide Comprehensive Documentation
- Use OpenAPI (Swagger) specifications for interactive, machine-readable API docs.
- Include examples of requests, responses, error codes, and descriptions.

## Optimize for Performance
- Implement caching to reduce redundant processing.
- Use response compression (gzip, Brotli).
- Design APIs to be stateless.
- Avoid returning unnecessary data; provide only what is requested.

## Test and Monitor APIs
- Conduct unit, integration, and security testing.
- Monitor API usage, errors, and performance metrics continuously.

## REST API Pagination Best Practices

Efficient pagination is essential for managing large datasets in REST APIs. It improves performance, enhances user experience, and reduces server load. This guide covers the most effective pagination patterns and practices.

### Common Pagination Techniques

#### 1. Offset and Limit Pagination
- Uses parameters like `offset` and `limit` to specify starting point and number of records.
- Example: `GET /items?offset=0&limit=10` fetches the first 10 records.
- Simple but can be inefficient and inconsistent with large or frequently changing datasets.

#### 2. Page-Based Pagination
- Uses `page` and `size` (or `limit`) parameters.
- Example: `GET /items?page=2&size=10` fetches the second page of 10 items.
- User-friendly and supports easy navigation but may suffer performance issues at scale.

#### 3. Cursor-Based Pagination (Keyset Pagination)
- Uses a cursor (unique identifier or token) to mark the last retrieved item.
- Example: `GET /items?cursor=abc123` fetches results after the cursor.
- Offers better performance, consistency, and stability for large or dynamic datasets.

#### 4. Time-Based Pagination
- Uses time range parameters like `start_time` and `end_time`.
- Suitable for datasets with temporal data, e.g., event logs.

### Best Practices for Pagination Output

- Use **standard parameter names** like `page` / `size` or `offset` / `limit` consistently.
- Always include **pagination metadata** in responses:
  - Total records
  - Current page
  - Total pages
  - Page size
- Provide **navigation links** (HATEOAS) for:
  - `next`
  - `prev`
  - `first`
  - `last`
- Return an **empty list with 200 OK** when no more data is available.
- Allow **customizable page sizes** within reasonable limits to improve flexibility.
- Set **maximum limits** on offset or page size to prevent performance issues.
- Implement **rate limiting** to protect the API from excessive requests.
- Support **sorting and filtering** parameters to improve data querying.
- Use cursor-based pagination for **large or frequently updated datasets** to improve consistency and reduce duplicate or missing records.

### Example Paginated Response (Page-Based)

- Paginate results following the pattern:

```json
{
  "data": [...],
  "pagination": {
    "current_page": 2,
    "page_size": 10,
    "total_pages": 100,
    "total_records": 1000,
    "links": {
      "first": "/items?page=1&size=10",
      "prev": "/items?page=1&size=10",
      "next": "/items?page=3&size=10",
      "last": "/items?page=100&size=10"
    }
  }
}
```

---

Following these practices will help build robust, scalable, and user-friendly REST APIs that stand the test of time.