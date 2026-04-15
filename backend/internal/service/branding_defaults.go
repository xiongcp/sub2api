package service

const defaultBrandingHomeContent = `<div class="brand-google-home">
  <header class="brand-google-home__header">
    <a class="brand-google-home__brand" href="/">
      <span class="brand-google-home__brand-dot brand-google-home__brand-dot--blue"></span>
      <span class="brand-google-home__brand-text">Sub2API</span>
    </a>
    <nav class="brand-google-home__nav">
      <a class="brand-google-home__nav-link" href="/login">登录</a>
      <a class="brand-google-home__nav-link brand-google-home__nav-link--primary" href="/register">注册</a>
    </nav>
  </header>

  <main class="brand-google-home__main">
    <div class="brand-google-home__badge">统一接入 Claude、GPT、Gemini 与更多模型</div>
    <h1 class="brand-google-home__title">一个干净直接的 AI API 入口</h1>
    <p class="brand-google-home__subtitle">
      兼容常见调用方式，管理密钥、订阅、额度与路由策略，用更少的步骤完成接入。
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
    <p>简洁的界面，直接的路径，面向稳定使用而设计。</p>
  </footer>
</div>
`

const defaultBrandingCustomCSS = `:root {
  --brand-google-blue: #4285f4;
  --brand-google-red: #ea4335;
  --brand-google-yellow: #fbbc05;
  --brand-google-green: #34a853;
  --brand-google-text: #202124;
  --brand-google-muted: #5f6368;
  --brand-google-border: #dadce0;
  --brand-google-surface: #ffffff;
  --brand-google-page: #f8f9fa;
  --brand-google-shadow: 0 1px 2px rgba(60, 64, 67, 0.18), 0 1px 3px 1px rgba(60, 64, 67, 0.1);
  --brand-google-font: ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
}

body {
  font-family: var(--brand-google-font);
  color: var(--brand-google-text);
}

.brand-google-home,
.brand-google-auth-note,
.brand-google-footer {
  font-family: var(--brand-google-font);
  color: var(--brand-google-text);
}

.brand-google-home {
  min-height: 100vh;
  padding: 24px;
  background: #ffffff;
  display: flex;
  flex-direction: column;
}

.brand-google-home *,
.brand-google-auth-note *,
.brand-google-footer * {
  box-sizing: border-box;
}

.brand-google-home__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.brand-google-home__brand {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  text-decoration: none;
  color: var(--brand-google-text);
  font-size: 14px;
  font-weight: 600;
}

.brand-google-home__brand-dot {
  width: 12px;
  height: 12px;
  border-radius: 999px;
  background: var(--brand-google-blue);
  box-shadow: 12px 0 0 var(--brand-google-red), 24px 0 0 var(--brand-google-yellow), 36px 0 0 var(--brand-google-green);
}

.brand-google-home__brand-text {
  margin-left: 40px;
}

.brand-google-home__nav {
  display: inline-flex;
  align-items: center;
  gap: 12px;
}

.brand-google-home__nav-link {
  color: var(--brand-google-muted);
  text-decoration: none;
  font-size: 14px;
  line-height: 20px;
}

.brand-google-home__nav-link:hover {
  color: var(--brand-google-text);
}

.brand-google-home__nav-link--primary {
  color: var(--brand-google-blue);
  font-weight: 600;
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
  min-height: 32px;
  padding: 0 14px;
  border: 1px solid #e8eaed;
  border-radius: 999px;
  background: #fff;
  color: var(--brand-google-muted);
  font-size: 13px;
  line-height: 1;
}

.brand-google-home__title {
  margin: 22px 0 14px;
  font-size: clamp(40px, 7vw, 72px);
  line-height: 1.08;
  font-weight: 500;
  letter-spacing: -0.04em;
}

.brand-google-home__subtitle {
  max-width: 620px;
  margin: 0 auto;
  color: var(--brand-google-muted);
  font-size: 16px;
  line-height: 1.7;
}

.brand-google-home__search-shell {
  display: flex;
  align-items: center;
  gap: 14px;
  width: 100%;
  max-width: 640px;
  min-height: 58px;
  margin: 32px auto 20px;
  padding: 0 22px;
  border: 1px solid transparent;
  border-radius: 999px;
  background: #fff;
  box-shadow: var(--brand-google-shadow);
}

.brand-google-home__search-shell:hover {
  box-shadow: 0 2px 8px rgba(60, 64, 67, 0.2), 0 1px 3px rgba(60, 64, 67, 0.15);
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
  color: #80868b;
  font-size: 15px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.brand-google-home__actions {
  display: flex;
  justify-content: center;
  gap: 12px;
  flex-wrap: wrap;
}

.brand-google-home__button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 132px;
  min-height: 42px;
  padding: 0 20px;
  border: 1px solid #dadce0;
  border-radius: 999px;
  background: #fff;
  color: var(--brand-google-text);
  text-decoration: none;
  font-size: 14px;
  font-weight: 500;
}

.brand-google-home__button:hover {
  border-color: #c6c6c6;
  box-shadow: 0 1px 1px rgba(0, 0, 0, 0.1);
}

.brand-google-home__button--primary {
  border-color: var(--brand-google-blue);
  background: var(--brand-google-blue);
  color: #fff;
}

.brand-google-home__button--primary:hover {
  border-color: #1a73e8;
  background: #1a73e8;
}

.brand-google-home__chips {
  display: flex;
  justify-content: center;
  flex-wrap: wrap;
  gap: 10px;
  margin-top: 24px;
}

.brand-google-home__chip {
  display: inline-flex;
  align-items: center;
  min-height: 32px;
  padding: 0 14px;
  border-radius: 999px;
  background: var(--brand-google-page);
  color: var(--brand-google-muted);
  font-size: 13px;
}

.brand-google-home__footer {
  padding-top: 24px;
  text-align: center;
  color: #80868b;
  font-size: 12px;
}

.auth-layout-shell {
  background: var(--brand-google-page);
}

.auth-layout-background {
  background: linear-gradient(180deg, #ffffff 0%, #f8f9fa 100%) !important;
}

.auth-layout-decorations {
  opacity: 0.28;
}

.auth-layout-shell .card-glass {
  border: 1px solid #e8eaed;
  border-radius: 28px;
  background: rgba(255, 255, 255, 0.94);
  box-shadow: 0 8px 24px rgba(60, 64, 67, 0.12);
  backdrop-filter: blur(10px);
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
  min-height: 48px;
  border-color: var(--brand-google-border);
  border-radius: 16px;
  background: #fff;
  box-shadow: none;
}

.auth-layout-shell .input:focus {
  border-color: var(--brand-google-blue);
  box-shadow: 0 0 0 3px rgba(66, 133, 244, 0.16);
}

.auth-layout-shell .btn.btn-primary {
  min-height: 48px;
  border-radius: 999px;
  background: #1a73e8;
  box-shadow: none;
}

.auth-layout-shell .btn.btn-primary:hover {
  background: #1765cc;
}

.brand-google-auth-note,
.brand-google-footer {
  border: 1px solid #e8eaed;
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.86);
  box-shadow: 0 1px 2px rgba(60, 64, 67, 0.08);
}

.brand-google-auth-note {
  padding: 16px 18px;
  text-align: left;
}

.brand-google-auth-note__eyebrow {
  margin: 0 0 8px;
  color: var(--brand-google-blue);
  font-size: 12px;
  font-weight: 600;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.brand-google-auth-note__text {
  margin: 0;
  color: var(--brand-google-muted);
  font-size: 14px;
  line-height: 1.7;
}

.brand-google-footer {
  padding: 12px 16px;
  text-align: center;
}

.brand-google-footer__text {
  margin: 0;
  color: var(--brand-google-muted);
  font-size: 13px;
  line-height: 1.6;
}

@media (max-width: 640px) {
  .brand-google-home {
    padding: 18px;
  }

  .brand-google-home__header {
    flex-direction: column;
    align-items: flex-start;
  }

  .brand-google-home__nav {
    width: 100%;
    justify-content: flex-start;
  }

  .brand-google-home__title {
    font-size: 36px;
  }

  .brand-google-home__search-shell {
    min-height: 52px;
    padding: 0 16px;
  }

  .brand-google-home__search-text {
    font-size: 14px;
  }

  .auth-layout-shell .card-glass {
    border-radius: 24px;
    padding: 24px;
  }
}
`

const defaultBrandingLoginExtraHTML = `<div class="brand-google-auth-note">
  <p class="brand-google-auth-note__eyebrow">Welcome back</p>
  <p class="brand-google-auth-note__text">
    登录后即可进入控制台，查看密钥、订阅、额度与调用情况。
  </p>
</div>
`

const defaultBrandingRegisterExtraHTML = `<div class="brand-google-auth-note">
  <p class="brand-google-auth-note__eyebrow">Create account</p>
  <p class="brand-google-auth-note__text">
    使用常用邮箱完成注册与验证，随后即可开始接入和管理你的 API 工作流。
  </p>
</div>
`

const defaultBrandingGlobalFooterHTML = `<div class="brand-google-footer">
  <p class="brand-google-footer__text">
    简洁、直接、可控的 AI API 管理体验。
  </p>
</div>
`

func defaultBrandingSettings() map[string]string {
	return map[string]string{
		SettingKeyHomeContent:       defaultBrandingHomeContent,
		SettingKeyCustomCSS:         defaultBrandingCustomCSS,
		SettingKeyLoginExtraHTML:    defaultBrandingLoginExtraHTML,
		SettingKeyRegisterExtraHTML: defaultBrandingRegisterExtraHTML,
		SettingKeyPaymentFooterHTML: "",
		SettingKeyGlobalFooterHTML:  defaultBrandingGlobalFooterHTML,
	}
}
