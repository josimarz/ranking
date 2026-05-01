# Guia: Screenshots Automatizadas com Playwright

## Problemas encontrados e soluções

### 1. Overlay intercepta clicks

**Problema**: Ao clicar num MatchCard, o MatchDetail abre com overlay `fixed inset-0 bg-black/50`. Qualquer click subsequente é interceptado pelo overlay — Playwright dá timeout.

**Solução**: Fechar via `page.keyboard.press('Escape')` ou recarregar a página.

```js
// ❌ NÃO funciona — overlay intercepta
await closeButton.click();

// ❌ NÃO funciona — overlay com pointer-events
await overlay.click({ force: true });

// ✅ Funciona — Escape ou reload
await page.keyboard.press('Escape');
await page.waitForTimeout(500);

// Se Escape não funcionar, reload
if (await page.locator('.fixed.inset-0').count()) {
  await page.goto(url, { waitUntil: 'networkidle' });
}
```

### 2. Globo demora pra carregar

**Problema**: O globo 3D (WebGL) demora 5-10s pra renderizar. Screenshot tirada cedo mostra "Carregando globo...".

**Solução**: Esperar mais tempo ou detectar o canvas.

```js
// ✅ Esperar tempo fixo (simples)
await page.waitForTimeout(8000);

// ✅ Esperar canvas existir (melhor)
await page.waitForSelector('canvas', { timeout: 15000 });
await page.waitForTimeout(2000); // extra pra texturas carregarem
```

### 3. Seletores com classes Tailwind escapadas

**Problema**: Classes como `bg-black/50` têm `/` que precisa ser escapado no CSS selector.

**Solução**: Usar `\\\\` pra escapar ou usar `locator` com texto.

```js
// ❌ Falha — / não escapado
await page.$('.bg-black/50');

// ✅ Funciona — escapar
await page.$('.bg-black\\\\/50');

// ✅ Melhor — usar locator com texto
await page.locator('text=Fase de Liga').first();

// ✅ Melhor ainda — usar data-testid (se disponível)
await page.locator('[data-testid="stage-filter"]');
```

### 4. Dynamic import + SSR false

**Problema**: Componentes com `dynamic(() => import(...), { ssr: false })` não existem no HTML inicial. Playwright pode não encontrá-los.

**Solução**: Usar `waitUntil: 'networkidle'` + timeout extra.

```js
await page.goto(url, { waitUntil: 'networkidle' });
await page.waitForTimeout(3000); // esperar hydration + dynamic imports
```

---

## Script completo de screenshots

```js
const { chromium } = require('@playwright/test');

async function takeScreenshots(baseUrl = 'http://localhost:3000') {
  const browser = await chromium.launch();
  const dir = '/tmp/visual-audit';
  
  // Desktop
  const desktop = await browser.newPage({ viewport: { width: 1920, height: 1080 } });
  
  // Helper: navegar e esperar
  async function go(page, path, waitMs = 3000) {
    await page.goto(`${baseUrl}${path}`, { waitUntil: 'networkidle', timeout: 15000 });
    await page.waitForTimeout(waitMs);
  }
  
  // Helper: screenshot com scroll
  async function screenshotWithScroll(page, name) {
    await page.screenshot({ path: `${dir}/${name}-top.png` });
    const scrollHeight = await page.evaluate(() => document.body.scrollHeight);
    if (scrollHeight > 1200) {
      await page.evaluate(() => window.scrollTo(0, document.body.scrollHeight));
      await page.waitForTimeout(1000);
      await page.screenshot({ path: `${dir}/${name}-bottom.png` });
      await page.evaluate(() => window.scrollTo(0, 0));
    }
  }
  
  // Helper: clicar e tirar print (com fallback)
  async function clickAndScreenshot(page, selector, name, closeAfter = true) {
    try {
      const el = page.locator(selector).first();
      if (await el.count() === 0) return;
      await el.click({ timeout: 3000 });
      await page.waitForTimeout(1500);
      await page.screenshot({ path: `${dir}/${name}.png` });
      if (closeAfter) {
        await page.keyboard.press('Escape');
        await page.waitForTimeout(500);
      }
    } catch {
      // Se falhar, continua sem screenshot dessa interação
    }
  }

  // ========== PÁGINAS PRINCIPAIS ==========
  
  // Home (globo precisa de mais tempo)
  await go(desktop, '/', 8000);
  await desktop.screenshot({ path: `${dir}/home.png` });
  
  // Home — interações
  await clickAndScreenshot(desktop, '[class*=cursor-pointer]', 'home-match-detail');
  
  // Classificação
  await go(desktop, '/classificacao');
  await screenshotWithScroll(desktop, 'classificacao');
  
  // Chaveamento
  await go(desktop, '/chaveamento');
  await desktop.screenshot({ path: `${dir}/chaveamento.png` });
  
  // Artilheiros
  await go(desktop, '/artilheiros');
  await screenshotWithScroll(desktop, 'artilheiros');
  
  // Times
  await go(desktop, '/times');
  await desktop.screenshot({ path: `${dir}/times.png` });
  
  // ========== SUBPÁGINAS DINÂMICAS ==========
  
  const teamIds = [86, 64, 81]; // Real Madrid, Liverpool, Barcelona
  for (const id of teamIds) {
    await go(desktop, `/times/${id}`);
    await screenshotWithScroll(desktop, `time-${id}`);
  }
  
  // ========== MOBILE ==========
  
  const mobile = await browser.newPage({ viewport: { width: 375, height: 812 } });
  
  await go(mobile, '/', 6000);
  await mobile.screenshot({ path: `${dir}/mobile-home.png` });
  
  await go(mobile, '/classificacao');
  await mobile.screenshot({ path: `${dir}/mobile-classificacao.png` });
  
  await go(mobile, '/times');
  await mobile.screenshot({ path: `${dir}/mobile-times.png` });
  
  // ========== CLEANUP ==========
  
  await browser.close();
  console.log('✅ Screenshots salvas em', dir);
}

takeScreenshots().catch(console.error);
```

## Como rodar

```bash
# Garantir que o servidor está rodando
pnpm dev --port 3000 &

# Rodar screenshots (de dentro do diretório e2e/)
cd e2e && node screenshot-audit.js
```

## Checklist pra evitar problemas

- [ ] Servidor rodando antes de tirar screenshots
- [ ] Usar `https://` nas URLs de texturas (não `//`)
- [ ] Esperar 8s+ pra páginas com WebGL (globo)
- [ ] Esperar 3s pra páginas normais
- [ ] Fechar modais com Escape, não com click
- [ ] Usar `{ force: true }` se precisar clicar em elemento coberto
- [ ] Usar `try/catch` em cada interação (não quebrar o script todo)
- [ ] Tirar screenshot do topo E do bottom (scroll) pra páginas longas
- [ ] Testar desktop (1920x1080) E mobile (375x812)
