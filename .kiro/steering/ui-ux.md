---
inclusion: auto
name: ui-ux
description: Guide and best practices for UI/UX. Apply this skill when writing code in the frontend.
---

# UI/UX Best Practices for Frontend Development

This document defines **best practices for UI and UX design** when building a frontend application with **Next.js**. The goal is to create **responsive, accessible, and user-friendly interfaces** that deliver a consistent and intuitive experience across devices.

---

## 1. General Principles

- Prioritize **clarity and simplicity** in layouts and interactions.
- Follow **Dati's brand guidelines** for colors, typography, and visual identity.
- Maintain **consistency** across components and pages.
- Design with **accessibility** in mind (WCAG compliance).
- Use clear visual hierarchy and intuitive navigation patterns.
- Keep interfaces **minimalistic**, avoiding unnecessary visual clutter.
- **Do NOT use modals for creating or editing records.** These actions must occur on dedicated pages or well-structured sections to avoid usability and accessibility issues.

---

## 2. Responsive Design

- The UI must be fully **responsive** and functional across all screen sizes:

  - Mobile-first approach: start designing for small screens and scale up.
  - Use flexible layouts with CSS Grid and Flexbox.
  - Apply TailwindCSS responsive utilities (e.g., `sm:`, `md:`, `lg:`).
  - Ensure touch-friendly interactions on mobile devices.

- Test on multiple devices and viewport sizes to guarantee proper behavior.

Example Tailwind usage:

```html
<div class="p-4 sm:p-6 lg:p-8">
  <h1 class="text-lg sm:text-xl lg:text-2xl">Responsive Title</h1>
</div>
```

---

## 3. Navigation & Information Architecture

- Use a **clear and logical structure** for navigation.
- Group related items together and label them clearly.
- Keep the number of navigation levels minimal.
- Provide visual feedback for active states and hover effects.
- Implement a **responsive navigation menu**, such as a hamburger menu for mobile.
- In Next.js, leverage **`next/link`** for client-side navigation for performance.

---

## 4. Typography

- Use font sizes and weights that support **readability and hierarchy**.
- Maintain adequate line height and spacing.
- Ensure sufficient **color contrast** between text and background.
- Limit the number of font families and styles.
- Follow the brand typography rules defined in \[@rules/brand.md].

---

## 5. Color & Contrast

- Stick to the **brand color palette** for consistency.
- Use color to convey meaning (e.g., errors in red, success in green).
- Ensure contrast ratios meet accessibility standards (minimum AA).
- Avoid relying on color alone to convey information.

---

## 6. Accessibility

- Support keyboard navigation for all interactive elements.
- Provide descriptive **alt attributes** for images.
- Use semantic HTML tags (`<header>`, `<main>`, `<footer>`).
- Label form fields correctly and clearly.
- Test with screen readers and accessibility tools regularly.

---

## 7. Performance & Usability

- Optimize images for the web using **Next.js Image component (`next/image`)**.
- Minimize the use of heavy assets or unnecessary animations.
- Ensure **fast loading times** by lazy-loading content when possible.
- Provide **clear feedback** for loading, success, and error states.
- Prevent layout shifts (CLS) by reserving space for images and dynamic content.

---

## 8. Forms & Interactions

- Use clear labels and placeholders.
- Provide inline validation and error messages.
- Make form fields large enough for touch interaction.
- Group related form elements together.
- Use progress indicators for multi-step forms.

---

## 9. Testing UI/UX

- Test responsive behavior using browser dev tools and physical devices.
- Conduct usability tests to validate navigation and user flows.
- Use Playwright to verify UI behavior automatically.
- Check accessibility with Lighthouse and other auditing tools.

---

## 10. References

- [Next.js Image Component](https://nextjs.org/docs/api-reference/next/image)
- [TailwindCSS Documentation](https://tailwindcss.com/docs)
- [WCAG Accessibility Guidelines](https://www.w3.org/WAI/WCAG21/quickref/)
- [Flowbite Components](https://flowbite.com/)

---

**Following these best practices ensures a responsive, accessible, and user-friendly frontend in Next.js that aligns with brand identity and provides a seamless experience across devices.**