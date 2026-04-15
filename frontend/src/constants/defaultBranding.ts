export const DEFAULT_BRANDING_HOME_CONTENT = `<div class="brand-google-home">
  <header class="brand-google-home__header">
    <a class="brand-google-home__brand" href="/">
      <span class="brand-google-home__brand-mark" aria-hidden="true">
        <span class="brand-google-home__brand-dot brand-google-home__brand-dot--blue"></span>
        <span class="brand-google-home__brand-dot brand-google-home__brand-dot--red"></span>
        <span class="brand-google-home__brand-dot brand-google-home__brand-dot--yellow"></span>
        <span class="brand-google-home__brand-dot brand-google-home__brand-dot--green"></span>
      </span>
      <span class="brand-google-home__brand-text">Sub2API</span>
    </a>
    <nav class="brand-google-home__nav">
      <a class="brand-google-home__nav-link" href="/login">登录</a>
      <a class="brand-google-home__nav-link brand-google-home__nav-link--primary" href="/register">注册</a>
    </nav>
  </header>

  <main class="brand-google-home__main">
    <div class="brand-google-home__badge">
      <span class="brand-google-home__badge-dot" aria-hidden="true"></span>
      统一接入 Claude、GPT、Gemini 与更多模型
    </div>
    <h1 class="brand-google-home__title">一个干净直接的 AI API 入口</h1>
    <p class="brand-google-home__subtitle">
      兼容常见调用方式，统一管理密钥、订阅、额度与路由策略，用更少步骤完成接入。
    </p>

    <div class="brand-google-home__search-shell" aria-hidden="true">
      <span class="brand-google-home__search-icon"></span>
      <span class="brand-google-home__search-text">OpenAI-compatible endpoint, subscriptions, billing, routing</span>
    </div>

    <div class="brand-google-home__actions">
      <a class="brand-google-home__button brand-google-home__button--primary" href="/login">进入控制台</a>
      <a class="brand-google-home__button" href="/register">创建账号</a>
    </div>

    <div class="brand-google-home__chips">
      <span class="brand-google-home__chip">统一入口</span>
      <span class="brand-google-home__chip">实时计费</span>
      <span class="brand-google-home__chip">订阅转 API</span>
    </div>
  </main>

  <footer class="brand-google-home__footer">
    <p>清晰的入口，稳定、可控的 AI API 管理体验。</p>
  </footer>
</div>
`

export const DEFAULT_BRANDING_CUSTOM_CSS = `:root {
  --brand-google-blue: #1a73e8;
  --brand-google-red: #ea4335;
  --brand-google-yellow: #fbbc04;
  --brand-google-green: #34a853;
  --brand-google-text: #202124;
  --brand-google-muted: #5f6368;
  --brand-google-muted-soft: #80868b;
  --brand-google-primary: #1a73e8;
  --brand-google-primary-hover: #1765cc;
  --brand-google-on-primary: #ffffff;
  --brand-google-page: #f8f9fa;
  --brand-google-surface-0: #ffffff;
  --brand-google-surface-1: #f8fafd;
  --brand-google-surface-2: #f1f5fb;
  --brand-google-outline: #dadce0;
  --brand-google-outline-strong: #c4c7c5;
  --brand-google-state-hover: rgba(26, 115, 232, 0.08);
  --brand-google-state-focus: rgba(26, 115, 232, 0.18);
  --brand-google-shadow-1: 0 1px 2px rgba(60, 64, 67, 0.16), 0 1px 3px 1px rgba(60, 64, 67, 0.08);
  --brand-google-shadow-2: 0 4px 12px rgba(60, 64, 67, 0.12), 0 2px 6px rgba(60, 64, 67, 0.08);
  --brand-google-font: ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
}

body {
  font-family: var(--brand-google-font);
  color: var(--brand-google-text);
  background: var(--brand-google-page);
}

.brand-google-home,
.brand-google-auth-note,
.brand-google-footer {
  font-family: var(--brand-google-font);
  color: var(--brand-google-text);
}

.brand-google-home {
  position: relative;
  min-height: 100vh;
  padding: 28px clamp(20px, 4vw, 40px);
  background: linear-gradient(180deg, #ffffff 0%, #f8fafd 100%);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.brand-google-home::before,
.brand-google-home::after {
  content: "";
  position: absolute;
  pointer-events: none;
}

.brand-google-home::before {
  inset: 0;
  background:
    radial-gradient(circle at 15% 20%, rgba(26, 115, 232, 0.08), transparent 28%),
    radial-gradient(circle at 85% 10%, rgba(251, 188, 4, 0.08), transparent 24%),
    radial-gradient(circle at 50% 100%, rgba(52, 168, 83, 0.08), transparent 24%);
}

.brand-google-home::after {
  top: 0;
  left: 50%;
  width: min(100%, 1080px);
  height: 1px;
  transform: translateX(-50%);
  background: linear-gradient(90deg, rgba(66, 133, 244, 0), rgba(66, 133, 244, 0.28), rgba(251, 188, 4, 0.22), rgba(52, 168, 83, 0));
}

.brand-google-home *,
.brand-google-auth-note *,
.brand-google-footer * {
  box-sizing: border-box;
}

.brand-google-home__header,
.brand-google-home__main,
.brand-google-home__footer {
  position: relative;
  z-index: 1;
}

.brand-google-home__header {
  width: 100%;
  max-width: 1120px;
  margin: 0 auto;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.brand-google-home__brand {
  display: inline-flex;
  align-items: center;
  gap: 12px;
  text-decoration: none;
  color: var(--brand-google-text);
  font-size: 14px;
  font-weight: 600;
  letter-spacing: -0.01em;
}

.brand-google-home__brand:focus-visible,
.brand-google-home__nav-link:focus-visible,
.brand-google-home__button:focus-visible {
  outline: none;
  box-shadow: 0 0 0 4px var(--brand-google-state-focus);
}

.brand-google-home__brand-mark {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.brand-google-home__brand-dot {
  width: 10px;
  height: 10px;
  border-radius: 999px;
  flex: 0 0 auto;
}

.brand-google-home__brand-dot--blue {
  background: var(--brand-google-blue);
}

.brand-google-home__brand-dot--red {
  background: var(--brand-google-red);
}

.brand-google-home__brand-dot--yellow {
  background: var(--brand-google-yellow);
}

.brand-google-home__brand-dot--green {
  background: var(--brand-google-green);
}

.brand-google-home__brand-text {
  line-height: 1;
}

.brand-google-home__nav {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.brand-google-home__nav-link {
  display: inline-flex;
  align-items: center;
  min-height: 40px;
  padding: 0 14px;
  border-radius: 999px;
  color: var(--brand-google-muted);
  text-decoration: none;
  font-size: 14px;
  line-height: 20px;
  transition: background-color 160ms ease, color 160ms ease;
}

.brand-google-home__nav-link:hover {
  color: var(--brand-google-text);
  background: rgba(95, 99, 104, 0.08);
}

.brand-google-home__nav-link--primary {
  color: var(--brand-google-primary);
  font-weight: 600;
}

.brand-google-home__nav-link--primary:hover {
  background: var(--brand-google-state-hover);
}

.brand-google-home__main {
  width: 100%;
  max-width: 760px;
  margin: auto;
  text-align: center;
}

.brand-google-home__badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  min-height: 36px;
  padding: 0 16px;
  border: 1px solid var(--brand-google-outline);
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.92);
  color: var(--brand-google-muted);
  font-size: 13px;
  line-height: 1;
  box-shadow: 0 1px 2px rgba(60, 64, 67, 0.08);
}

.brand-google-home__badge-dot {
  width: 8px;
  height: 8px;
  border-radius: 999px;
  background: linear-gradient(90deg, var(--brand-google-blue) 0%, var(--brand-google-green) 100%);
  flex: 0 0 auto;
}

.brand-google-home__title {
  max-width: 11ch;
  margin: 24px auto 14px;
  font-size: clamp(42px, 7vw, 72px);
  line-height: 1.04;
  font-weight: 500;
  letter-spacing: -0.05em;
}

.brand-google-home__subtitle {
  max-width: 640px;
  margin: 0 auto;
  color: var(--brand-google-muted);
  font-size: 17px;
  line-height: 1.76;
}

.brand-google-home__search-shell {
  display: flex;
  align-items: center;
  gap: 14px;
  width: 100%;
  max-width: 680px;
  min-height: 60px;
  margin: 34px auto 20px;
  padding: 0 22px;
  border: 1px solid var(--brand-google-outline);
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.98);
  box-shadow: var(--brand-google-shadow-1);
  transition: border-color 160ms ease, box-shadow 160ms ease, transform 160ms ease;
}

.brand-google-home__search-shell:hover {
  border-color: var(--brand-google-outline-strong);
  box-shadow: var(--brand-google-shadow-2);
  transform: translateY(-1px);
}

.brand-google-home__search-icon {
  width: 18px;
  height: 18px;
  border-radius: 999px;
  border: 2px solid var(--brand-google-blue);
  position: relative;
  flex: 0 0 auto;
}

.brand-google-home__search-icon::after {
  content: "";
  position: absolute;
  right: -6px;
  bottom: -5px;
  width: 8px;
  height: 2px;
  border-radius: 999px;
  background: var(--brand-google-red);
  transform: rotate(45deg);
}

.brand-google-home__search-text {
  width: 100%;
  color: var(--brand-google-muted-soft);
  font-size: 15px;
  text-align: left;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.brand-google-home__actions {
  display: flex;
  justify-content: center;
  gap: 14px;
  flex-wrap: wrap;
}

.brand-google-home__button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 144px;
  min-height: 46px;
  padding: 0 22px;
  border: 1px solid var(--brand-google-outline);
  border-radius: 999px;
  background: var(--brand-google-surface-0);
  color: var(--brand-google-text);
  text-decoration: none;
  font-size: 14px;
  font-weight: 500;
  box-shadow: 0 1px 2px rgba(60, 64, 67, 0.06);
  transition: background-color 160ms ease, border-color 160ms ease, box-shadow 160ms ease, transform 160ms ease;
}

.brand-google-home__button:hover {
  border-color: var(--brand-google-outline-strong);
  background: var(--brand-google-surface-1);
  box-shadow: 0 2px 6px rgba(60, 64, 67, 0.12);
  transform: translateY(-1px);
}

.brand-google-home__button:active {
  transform: translateY(0);
  box-shadow: 0 1px 2px rgba(60, 64, 67, 0.12);
}

.brand-google-home__button--primary {
  border-color: var(--brand-google-primary);
  background: var(--brand-google-primary);
  color: var(--brand-google-on-primary);
  box-shadow: 0 1px 2px rgba(26, 115, 232, 0.32);
}

.brand-google-home__button--primary:hover {
  border-color: var(--brand-google-primary-hover);
  background: var(--brand-google-primary-hover);
  box-shadow: 0 4px 10px rgba(26, 115, 232, 0.22);
}

.brand-google-home__chips {
  display: flex;
  justify-content: center;
  flex-wrap: wrap;
  gap: 10px;
  width: fit-content;
  max-width: 100%;
  margin: 24px auto 0;
  padding: 12px;
  border: 1px solid rgba(218, 220, 224, 0.72);
  border-radius: 22px;
  background: rgba(248, 250, 253, 0.88);
}

.brand-google-home__chip {
  display: inline-flex;
  align-items: center;
  min-height: 34px;
  padding: 0 14px;
  border: 1px solid rgba(218, 220, 224, 0.72);
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.88);
  color: var(--brand-google-muted);
  font-size: 13px;
}

.brand-google-home__footer {
  width: 100%;
  max-width: 760px;
  margin: 32px auto 0;
  padding-top: 24px;
  border-top: 1px solid rgba(218, 220, 224, 0.7);
  text-align: center;
  color: var(--brand-google-muted-soft);
  font-size: 13px;
  line-height: 1.7;
}

.auth-layout-shell {
  background: var(--brand-google-page);
}

.auth-layout-background {
  background:
    radial-gradient(circle at top center, rgba(66, 133, 244, 0.16), transparent 34%),
    linear-gradient(180deg, #ffffff 0%, #f8f9fa 100%) !important;
}

.auth-layout-decorations {
  opacity: 0.16;
}

.auth-layout-shell .card-glass {
  border: 1px solid var(--brand-google-outline);
  border-radius: 32px;
  background: rgba(255, 255, 255, 0.98);
  box-shadow: var(--brand-google-shadow-2);
  backdrop-filter: none;
}

.auth-layout-shell .text-gradient {
  background: none !important;
  -webkit-text-fill-color: var(--brand-google-text);
  color: var(--brand-google-text) !important;
}

.auth-layout-shell .auth-layout-subtitle,
.auth-layout-shell .auth-layout-footer,
.auth-layout-shell .auth-layout-copyright {
  color: var(--brand-google-muted) !important;
}

.auth-layout-shell .input {
  min-height: 50px;
  border: 1px solid var(--brand-google-outline);
  border-radius: 14px;
  background: var(--brand-google-surface-0);
  box-shadow: none;
  transition: border-color 160ms ease, box-shadow 160ms ease, background-color 160ms ease;
}

.auth-layout-shell .input:hover {
  border-color: var(--brand-google-outline-strong);
}

.auth-layout-shell .input:focus {
  border-color: var(--brand-google-primary);
  box-shadow: 0 0 0 4px var(--brand-google-state-focus);
}

.auth-layout-shell .btn.btn-primary {
  min-height: 48px;
  border-radius: 999px;
  background: var(--brand-google-primary);
  box-shadow: 0 1px 2px rgba(26, 115, 232, 0.32);
}

.auth-layout-shell .btn.btn-primary:hover {
  background: var(--brand-google-primary-hover);
  box-shadow: 0 4px 10px rgba(26, 115, 232, 0.22);
}

.brand-google-auth-note {
  border: 1px solid rgba(218, 220, 224, 0.78);
  border-radius: 24px;
  background: rgba(255, 255, 255, 0.92);
  box-shadow: 0 1px 2px rgba(60, 64, 67, 0.08);
  padding: 18px 20px;
  text-align: left;
}

.brand-google-auth-note__heading {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 10px;
}

.brand-google-auth-note__marker {
  width: 32px;
  height: 8px;
  border-radius: 999px;
  background: linear-gradient(90deg, var(--brand-google-blue) 0%, var(--brand-google-red) 33%, var(--brand-google-yellow) 66%, var(--brand-google-green) 100%);
  flex: 0 0 auto;
}

.brand-google-auth-note__eyebrow {
  margin: 0;
  color: var(--brand-google-text);
  font-size: 12px;
  font-weight: 600;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.brand-google-auth-note__text {
  margin: 0;
  color: var(--brand-google-muted);
  font-size: 14px;
  line-height: 1.75;
}

.brand-google-footer {
  padding: 8px 16px 0;
  text-align: center;
  background: transparent;
}

.brand-google-footer__divider {
  display: block;
  width: 56px;
  height: 1px;
  margin: 0 auto 14px;
  background: linear-gradient(90deg, rgba(66, 133, 244, 0), rgba(66, 133, 244, 0.44), rgba(52, 168, 83, 0.32), rgba(52, 168, 83, 0));
}

.brand-google-footer__text {
  margin: 0;
  color: var(--brand-google-muted);
  font-size: 13px;
  line-height: 1.7;
}

@media (max-width: 640px) {
  .brand-google-home {
    padding: 18px;
  }

  .brand-google-home__header {
    align-items: flex-start;
    flex-direction: column;
  }

  .brand-google-home__nav {
    width: 100%;
    justify-content: flex-start;
  }

  .brand-google-home__title {
    max-width: none;
    font-size: 38px;
  }

  .brand-google-home__subtitle {
    font-size: 15px;
    line-height: 1.7;
  }

  .brand-google-home__search-shell {
    min-height: 54px;
    padding: 0 16px;
  }

  .brand-google-home__search-text {
    font-size: 14px;
  }

  .brand-google-home__actions {
    flex-direction: column;
  }

  .brand-google-home__button {
    width: 100%;
  }

  .auth-layout-shell .card-glass {
    border-radius: 28px;
    padding: 24px;
  }
}
`

export const DEFAULT_BRANDING_LOGIN_EXTRA_HTML = `<div class="brand-google-auth-note">
  <div class="brand-google-auth-note__heading">
    <span class="brand-google-auth-note__marker" aria-hidden="true"></span>
    <p class="brand-google-auth-note__eyebrow">Sign in</p>
  </div>
  <p class="brand-google-auth-note__text">
    登录后即可进入控制台，查看密钥、订阅、额度与调用情况。
  </p>
</div>
`

export const DEFAULT_BRANDING_REGISTER_EXTRA_HTML = `<div class="brand-google-auth-note">
  <div class="brand-google-auth-note__heading">
    <span class="brand-google-auth-note__marker" aria-hidden="true"></span>
    <p class="brand-google-auth-note__eyebrow">Create account</p>
  </div>
  <p class="brand-google-auth-note__text">
    使用常用邮箱完成注册与验证，随后即可开始接入和管理你的 API 工作流。
  </p>
</div>
`

export const DEFAULT_BRANDING_GLOBAL_FOOTER_HTML = `<div class="brand-google-footer">
  <span class="brand-google-footer__divider" aria-hidden="true"></span>
  <p class="brand-google-footer__text">
    清晰的入口，稳定、可控的 AI API 管理体验。
  </p>
</div>
`
