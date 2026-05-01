---
inclusion: auto
name: seo
description: Guide and best practices for applying SEO. Apply this skill when writing code in the frontend.
---

# SEO Best Practices for **Ranking** (Next.js)

This document defines SEO guidelines and technical best practices for the **Ranking** product, built with Next.js.  
The goal is to maximize organic traffic, improve discoverability of public rankings, and create a scalable foundation for search growth.

---

## 1. Core SEO Principles

The product must follow three fundamental SEO pillars:

1. **Indexability** – Search engines must be able to crawl and index all public content.
2. **Relevance** – Each ranking and item page must clearly communicate its topic.
3. **Performance** – Pages must be fast and mobile-friendly.

The site is content-driven. Each public ranking is a unique, valuable page that can rank in search engines.

---

## 2. URL Strategy

Use clean, semantic, and stable URLs.

Examples:

```

/rankings/019bc9c7-b857-72c5-a491-7f9686c7989a
/rankings/019bc9c7-b857-72c5-a491-7f9686c7989a/video-games

````

Best practices:

- Keep URLs lowercase.
- Avoid query parameters for core content.
- Use hyphen-separated slugs.
- Include the ranking name as a slug when possible.
- The UUID remains the canonical identifier.

---

## 3. Metadata with Next.js

Every public page must define:

- `<title>`
- `<meta name="description">`
- Open Graph tags
- Twitter Card tags
- Canonical URL

Example using Next.js App Router:

```ts
export const generateMetadata = ({ params }) => ({
  title: `Video Games Ranking – Compare Consoles by Graphics, Price & More`,
  description: `Compare video game consoles by graphics, sound, controls, game library and price. See community scores and add your own ratings.`,
  alternates: {
    canonical: `https://ranking.com/rankings/${params.id}/video-games`,
  },
  openGraph: {
    title: `Video Games Ranking`,
    description: `Community ranking of video game consoles.`,
    url: `https://ranking.com/rankings/${params.id}/video-games`,
    type: "website",
  },
});
```

Rules:

* Titles should be **50–60 characters**.
* Descriptions should be **140–160 characters**.
* Include the ranking name and main keywords.
* Never leave metadata empty.

---

## 4. Page Types and SEO Goals

### 4.1 Home Page

Purpose:

* Capture broad queries like:

  * “Create rankings online”
  * “Compare items with scores”
  * “Ranking tool”

Content:

* H1:
  `Create and Share Rankings Online`
* Short explanation of the product.
* List of latest public rankings.
* Internal links to popular rankings.

---

### 4.2 Public Ranking Page

Each public ranking is a landing page.

Structure:

```html
<h1>Video Games Ranking</h1>
<p>Compare video game consoles by graphics, sound, controls, game library and price.</p>

<h2>Items</h2>
<!-- Ranking table -->

<h2>About this Ranking</h2>
<p>Description written by the creator.</p>

<h2>Attributes</h2>
<ul>
  <li>Graphics – Visual quality</li>
  <li>Sound – Audio experience</li>
</ul>
```

SEO Rules:

* Exactly one `<h1>` per page.
* Use `<h2>` for sections.
* Always render textual descriptions.
* Do not rely only on tables for content.
* Make the ranking name the main keyword.

---

### 4.3 Search and Discovery Pages

* `/search?q=games`
* `/tags/video-games`

These pages must:

* Be indexable.
* Have dynamic titles:

  * `Search results for "games" – Ranking`
  * `Rankings tagged with "movies"`

---

## 5. Internal Linking

Internal links are critical.

Examples:

* Home → Latest rankings
* Ranking → Tags
* Tag → Related rankings
* “Meus Rankings” → Each user ranking

Rules:

* Use `<a href>` links, not only JS navigation.
* Always use descriptive anchor text:

  * ❌ “Click here”
  * ✅ “Video Games Ranking”

---

## 6. Structured Data (Schema.org)

Use JSON-LD for:

* `WebSite`
* `CollectionPage`
* `ItemList`

Example for a ranking:

```html
<script type="application/ld+json">
{
  "@context": "https://schema.org",
  "@type": "ItemList",
  "name": "Video Games Ranking",
  "itemListElement": [
    {
      "@type": "ListItem",
      "position": 1,
      "name": "Nintendo Entertainment System"
    },
    {
      "@type": "ListItem",
      "position": 2,
      "name": "Sega Mega Drive"
    }
  ]
}
</script>
```

This improves visibility in rich results.

---

## 7. Performance

SEO depends on speed.

Targets:

* LCP < 2.5s
* CLS < 0.1
* TTFB < 500ms

Use:

* Next.js Image Optimization
* Static Generation (SSG) for public rankings
* ISR (Incremental Static Regeneration)
* Edge caching (CDN)
* Minimal JS on ranking pages

Avoid:

* Blocking scripts
* Heavy client-only rendering
* Large unoptimized images

---

## 8. Crawl Control

Provide:

* `robots.txt`
* `sitemap.xml`

Rules:

* Only public rankings are indexable.
* Private rankings:

  * Must include `noindex`
  * Must not appear in sitemap

Example:

```
User-agent: *
Allow: /
Disallow: /private/
```

Sitemap:

* Include:

  * Home
  * Public rankings
  * Tag pages
* Update automatically when new rankings are created.

---

## 9. Content Strategy

SEO growth depends on volume and uniqueness.

Every public ranking is:

* A unique URL
* A unique topic
* A new indexable page

Encourage:

* Descriptions when creating rankings
* Meaningful attribute names
* Useful item names

Auto-generate supporting text:

* “This ranking compares 12 items using 5 attributes.”
* “Last updated on …”

This increases textual relevance.

---

## 10. Social Sharing = SEO

Sharing increases backlinks and traffic.

For every ranking:

* Use Open Graph images:

  * Ranking name
  * Top 3 items
* Ensure correct previews on:

  * WhatsApp
  * Facebook
  * X
  * Instagram

Preview example:

```
Video Games Ranking
Compare consoles by graphics, price and more.
```

---

## 11. Index Safety Rules

* Never expose private rankings in:

  * Sitemap
  * Search
  * Public listings
* Validate visibility on every request.
* Enforce canonical URLs.
* Avoid duplicate content:

  * One canonical URL per ranking.

---

## 12. Success Metrics

Track:

* Indexed pages
* Organic sessions
* Click-through rate (CTR)
* Rankings per keyword
* Time on page for ranking pages

Tools:

* Google Search Console
* Google Analytics
* Lighthouse

---

By following these guidelines, **Ranking** becomes a content engine where every public ranking is a potential landing page.
This creates a compounding SEO effect: more users → more rankings → more indexed pages → more traffic.