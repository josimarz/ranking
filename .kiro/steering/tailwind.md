---
inclusion: auto
name: tailwind
description: Guide and best practices for TailwindCSS. Apply this skill when writing code in the frontend.
---

# Tailwind CSS v4 + Next.js — Claude Code Memory

## Core Rules

- **Do NOT suggest or generate a `tailwind.config.js` file** by default.
- Tailwind v4 uses a **CSS-first configuration model**, meaning all customizations live in the global CSS file using directives like:
  - `@import "tailwindcss";`
  - `@theme { ... }`
  - `@source`
  - `@plugin`
- Use CSS variables to define design tokens for colors, spacing, typography, and breakpoints inside `@theme`.

---

## Next.js Integration

- Follow official Next.js setup guide for Tailwind v4: https://tailwindcss.com/docs/installation/framework-guides/nextjs
- The main Tailwind entry point should be in `app/globals.css` or `src/styles/globals.css`.
- Import Tailwind layers once at the top level:

```css
@import "tailwindcss";
```

- Do NOT import Tailwind into multiple files; global import only.

---

## Styling Practices

- Use semantic tokens via `@theme`:

```css
@theme {
  --color-primary: #1a73e8;
  --color-secondary: #9333ea;
  --font-sans: "Inter", system-ui, sans-serif;
  --spacing-7: 1.75rem;
  --breakpoint-lg: 64rem;
}
```

- Reference these tokens indirectly through Tailwind utilities (`bg-primary`, `text-primary`, `font-sans`).
- Avoid arbitrary values (`w-[137px]`, `text-[#ff0000]`) unless absolutely necessary.

---

## Responsive Design

- Mobile-first:
  - Base utilities apply to all viewports.
  - Use responsive prefixes like `sm:`, `md:`, `lg:`, `xl:`, etc., for progressive enhancement.
- Define custom breakpoints in `@theme` if needed.

---

## Component Design

- If the same utility combinations repeat across files, abstract them into React components (`Button`, `Card`, etc.).
- You may also use helper utilities like `tailwind-variants` or `tailwind-merge` for class composition.
- Avoid excessive use of `@apply` inside component CSS—it is okay for a few base patterns but not for everything.

---

## Performance

- Tailwind v4 automatically optimizes and tree-shakes unused styles.
- No need for a purge configuration or manual content paths.

---

## Readability

- Order class names logically:
  - Layout → Box Model → Typography → Color → State/Variants
  - Example: `flex flex-col items-center p-4 text-gray-900 hover:text-gray-700`
- Break long class strings across multiple lines in JSX for readability.

---

## Security

- Avoid inline styles for sensitive values.
- Keep global CSS minimal and rely on utility classes whenever possible.

---

## 2. Folder Structure

Here’s the suggested Next.js + Tailwind v4 structure:

```
my-app/
app/
layout.js
page.js
globals.css        # Tailwind entry point + theme customization
src/
components/
Button.jsx
Card.jsx
styles/
(optional) additional global styles
public/
favicon.ico
package.json
next.config.js
postcss.config.mjs   # Minimal PostCSS config for Tailwind

```

---

## 3. Example Code

### 3.1 Installing Dependencies

Run these commands to set up Tailwind v4 and Next.js:

```bash
npx create-next-app@latest my-app
cd my-app
npm install tailwindcss@latest postcss@latest autoprefixer@latest
```

---

### 3.2 Global Styles (`app/globals.css`)

> This file defines Tailwind layers and your design tokens with `@theme`.

```css
/* Tailwind CSS core layers */
@import "tailwindcss";

/* Theme configuration using CSS variables */
@theme {
  --color-primary: #1a73e8;
  --color-secondary: #9333ea;
  --color-neutral: #f3f4f6;

  --font-sans: "Inter", system-ui, sans-serif;
  --font-heading: "Poppins", system-ui, sans-serif;

  --spacing-7: 1.75rem;
  --spacing-9: 2.25rem;

  --breakpoint-sm: 40rem;
  --breakpoint-md: 48rem;
  --breakpoint-lg: 64rem;
}

/* Optional: base layer styles */
@layer base {
  html {
    font-family: var(--font-sans);
    background-color: var(--color-neutral);
  }
}

/* Optional: custom utility classes */
@layer utilities {
  .text-shadow {
    text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.2);
  }
}
```

---

### 3.3 Layout File (`app/layout.js`)

```jsx
export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body className="min-h-screen bg-neutral text-gray-900">
        <header className="p-4 bg-primary text-white">
          <h1 className="text-xl font-heading">My Next.js App</h1>
        </header>
        <main className="container mx-auto p-4">{children}</main>
        <footer className="p-4 text-center text-gray-500 text-sm">
          © {new Date().getFullYear()} My Company
        </footer>
      </body>
    </html>
  );
}
```

---

### 3.4 Example Page (`app/page.js`)

```jsx
import Button from "../src/components/Button";
import Card from "../src/components/Card";

export default function Home() {
  return (
    <div className="grid gap-4 sm:grid-cols-2 md:grid-cols-3">
      <Card
        title="Welcome to Tailwind v4"
        description="CSS-first configuration is powerful!"
      />
      <Button variant="primary">Get Started</Button>
    </div>
  );
}
```

---

### 3.5 Reusable Button Component (`src/components/Button.jsx`)

```jsx
import clsx from "clsx";

export default function Button({
  children,
  variant = "primary",
  className,
  ...props
}) {
  const base =
    "px-4 py-2 rounded-lg font-medium transition-colors duration-200";

  const variants = {
    primary: "bg-primary text-white hover:bg-blue-700",
    secondary: "bg-secondary text-white hover:bg-purple-700",
    neutral: "bg-neutral text-gray-800 hover:bg-gray-200",
  };

  return (
    <button className={clsx(base, variants[variant], className)} {...props}>
      {children}
    </button>
  );
}
```

---

### 3.6 Card Component (`src/components/Card.jsx`)

```jsx
export default function Card({ title, description }) {
  return (
    <div className="border rounded-xl p-6 bg-white shadow hover:shadow-lg transition-shadow">
      <h2 className="text-lg font-heading mb-2">{title}</h2>
      <p className="text-gray-600">{description}</p>
    </div>
  );
}
```

---

### 3.7 Minimal PostCSS Config (`postcss.config.mjs`)

> Required for Tailwind to process your CSS in Next.js.

```js
export default {
  plugins: {
    tailwindcss: {},
    autoprefixer: {},
  },
};
```

---

## 4. Benefits of This Setup

- **Latest Tailwind v4 syntax:**

  - No `tailwind.config.js`.
  - Theme and tokens declared directly in CSS.

- **Automatic optimization:**

  - Tree-shaking and purging are built-in.
  - Minimal setup, smaller bundle sizes.

- **Clear separation of concerns:**

  - Styling tokens in CSS.
  - Reusable UI components in React.

- **Scalable architecture:**

  - Ready for dark mode, theming, and responsive layouts.
  - Clean, semantic token definitions.