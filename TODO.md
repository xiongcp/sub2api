## Goal
- 提交当前工作区内已完成的密码重置与顶部横幅改动，推送 `main`，发布正式版 `v0.1.115`，并推送 `chengpengxiong/sub2api:0.1.115` 与 `chengpengxiong/sub2api:latest` 镜像。

## Todo
- 无。

## Doing
- 无。

## Done
- 已确认当前 `main`/`origin/main`/`v0.1.114` 指向同一已发布提交，新的未提交改动需要以新版本发布，不能覆盖旧 tag 与旧镜像。
- 已将 `backend/cmd/server/VERSION` 更新为 `0.1.115`。
- 已完成目标测试与脚本校验，并通过 `git diff --check`。
- 已创建发布提交 `8afbf362 feat: add top banner and password reset tooling` 并推送到 `origin/main`。
- 已创建并推送 git tag `v0.1.115`。
- 已构建并推送 `chengpengxiong/sub2api:0.1.115` 与 `chengpengxiong/sub2api:latest`，两者 digest 同为 `sha256:36084820fd5e7d37fb0488ec478370ed30aea6bff0bc4b95f7435ba51b58cf01`。

## Validation
- `git diff --check`
- `cd backend && go test -tags unit ./internal/service -run 'TestSettingService_(GetPublicSettings|UpdateSettings_)|TestAdminService_UpdateUser_(WithPasswordIncrementsTokenVersion|WithoutPasswordKeepsTokenVersion)'`
- `cd backend && go test -tags unit ./internal/server -run 'TestAPIContracts'`
- `cd frontend && pnpm test:run src/stores/__tests__/app.spec.ts src/components/layout/__tests__/AppHeader.spec.ts`
- `bash -n tools/reset_user_password.sh`
- `git push origin main`
- `git push origin v0.1.115`
- `docker build -t chengpengxiong/sub2api:0.1.115 -t chengpengxiong/sub2api:latest --build-arg VERSION=0.1.115 --build-arg COMMIT=8afbf362 --build-arg DATE=<UTC时间> --build-arg GOPROXY=https://goproxy.cn,direct --build-arg GOSUMDB=sum.golang.google.cn -f Dockerfile .`
- `docker push chengpengxiong/sub2api:0.1.115`
- `docker push chengpengxiong/sub2api:latest`

## Risks
- 当前工作区包含两组功能改动：管理员重置密码会失效旧会话，以及登录后顶部横幅；本轮将作为一次合并发布提交。
- 若直接沿用 `0.1.114` 发布，会造成版本号、git tag 和 Docker 镜像内容不一致；因此本轮统一提升到 `0.1.115`。

## Next Steps
- 无。

## Goal
- 为登录后页面增加可配置的顶部通知横幅，用于展示“充值联系 xxx”等固定运营消息；消息读取继续复用公开设置缓存，避免额外查询压力。

## Todo
- 无。

## Doing
- 无。

## Done
- 已为后端设置模型、公开设置 DTO、管理端设置更新入口新增 `top_banner_enabled` 与 `top_banner_text` 字段。
- 已让 `/settings/public` 和 `/api/v1/admin/settings` 读写链路贯通顶部横幅字段，并继续复用现有设置读取缓存。
- 已在管理端站点设置页增加顶部横幅开关与文案输入框。
- 已在登录后 `AppHeader` 增加顶部横幅展示、纯文本换行展示和本地关闭缓存。
- 已新增/更新后端缓存测试、公开设置测试、前端 store 测试和 `AppHeader` 横幅测试。
- 已同步更新后台接口契约测试，覆盖新增响应字段。

## Validation
- `git diff --check -- TODO.md backend/internal/service/domain_constants.go backend/internal/service/settings_view.go backend/internal/service/setting_service.go backend/internal/handler/dto/settings.go backend/internal/handler/setting_handler.go backend/internal/handler/admin/setting_handler.go backend/internal/service/setting_service_public_test.go backend/internal/service/setting_service_read_cache_test.go backend/internal/server/api_contract_test.go frontend/src/types/index.ts frontend/src/stores/app.ts frontend/src/api/admin/settings.ts frontend/src/views/admin/SettingsView.vue frontend/src/components/layout/AppHeader.vue frontend/src/stores/__tests__/app.spec.ts frontend/src/components/layout/__tests__/AppHeader.spec.ts frontend/src/i18n/locales/zh.ts frontend/src/i18n/locales/en.ts`
- `cd backend && go test -tags unit ./internal/service -run 'TestSettingService_(GetPublicSettings|UpdateSettings_)'`
- `cd backend && go test -tags unit ./internal/server -run 'TestAPIContracts'`
- `cd frontend && pnpm test:run src/stores/__tests__/app.spec.ts src/components/layout/__tests__/AppHeader.spec.ts`

## Risks
- 工作区当前为脏树，已有上一条密码重置任务的未提交改动；本轮只追加横幅相关改动，不回退无关文件。
- 顶部横幅复用 `/settings/public`，因此展示时效仍受现有公开设置缓存与强刷时机影响。
- 横幅关闭状态仅缓存在当前浏览器，并且以文案文本作为签名；若后台重新发布完全相同的文案，已关闭的浏览器仍会继续隐藏。

## Next Steps
- 无。

## Goal
- 为管理员提供安全的“按邮箱重置用户密码”方案：复用现有 Admin API，不查询数据库明文密码；同时修复管理员改密后旧会话未失效的问题。

## Todo
- 无。

## Doing
- 无。

## Done
- 已确认用户密码存储在 `users.password_hash`，仓库不存在可读取明文密码的实现。
- 已确认现有后台支持 `PUT /api/v1/admin/users/:id` 携带 `password` 修改用户密码。
- 已确认当前实现存在安全缺口：管理员改密时不会递增 `token_version`，旧 JWT / Refresh Token 不会立即失效。
- 已在 `backend/internal/service/admin_service.go` 中让管理员改密同步递增 `token_version`。
- 已新增 `backend/internal/service/admin_service_update_user_test.go`，覆盖管理员改密会使 `token_version` 递增，以及非改密更新不会误伤会话版本。
- 已新增 `tools/reset_user_password.sh`，支持用管理员 JWT 或 Admin API Key 按邮箱查找用户并重置密码。

## Validation
- `git diff --check -- backend/internal/service/admin_service.go backend/internal/service/admin_service_update_user_test.go tools/reset_user_password.sh TODO.md`
- `cd backend && go test -tags unit ./internal/service -run 'TestAdminService_UpdateUser_(WithPasswordIncrementsTokenVersion|WithoutPasswordKeepsTokenVersion)'`
- `bash -n tools/reset_user_password.sh`
- `tools/reset_user_password.sh --help`

## Risks
- 新脚本依赖 `curl` 和 `python3`，未额外兼容极简运行环境。
- 脚本通过 `search` 拉取用户列表后做邮箱精确匹配；如果后台检索策略未来变化，脚本需要同步调整。

## Next Steps
- 在管理机或本地终端用管理员 JWT / Admin API Key 执行 `tools/reset_user_password.sh`，不要再直接查数据库。

## Goal
- 提交 `UseKeyModal.vue` 的本地改动，忽略 `.playwright-mcp/` 调试产物，发布正式版 `v0.1.114`，并推送 `chengpengxiong/sub2api:0.1.114` 与 `latest` 镜像。

## Todo
- 验证工作区只包含预期文件。
- 运行 `UseKeyModal` 目标测试。
- 提交代码并推送 `origin/main`。
- 打并推送 git tag `v0.1.114`。
- 构建并推送 Docker 镜像标签 `0.1.114` 和 `latest`。

## Doing
- 正在落地 `.gitignore`、版本号和发布前验证。

## Done
- 已确认当前工作区只包含 `frontend/src/components/keys/UseKeyModal.vue` 和 `.playwright-mcp/` 调试产物。
- 已确认 `UseKeyModal.vue` 当前变更仅将两处 `model_context_window`/`model_auto_compact_token_limit` 配置从 `1000000/900000` 调整为 `200000/190000`。
- 已确认发布策略为正式版 `0.1.114`，Docker 同步推送 `chengpengxiong/sub2api:0.1.114` 与 `chengpengxiong/sub2api:latest`。

## Validation
- 待补充。

## Risks
- `.playwright-mcp/` 为本地调试产物，不应纳入本次提交。
- 当前仓库历史上已存在 `0.1.114-rc1` 版本文件和镜像，本轮需要确保正式版 tag 与镜像标签一致，避免 RC/正式版混淆。

## Next Steps
- 完成验证后提交代码、推送分支与 tag，并发布正式版 Docker 镜像。

## Goal
- 将版本号更新为 `0.1.114-rc1`，并构建、推送 `chengpengxiong/sub2api:0.1.114-rc1` Docker 镜像。

## Todo
- 无。

## Doing
- 无。

## Done
- 已确认当前版本文件为 `0.1.113-rc1`，最新 git tag 也是 `v0.1.113-rc1`。
- 已确认主 Dockerfile 支持通过构建参数或 `backend/cmd/server/VERSION` 注入版本号。
- 已决定 RC 版本只推送 `chengpengxiong/sub2api:0.1.114-rc1`，不覆盖 `latest`。
- 已将 `backend/cmd/server/VERSION` 更新为 `0.1.114-rc1`。
- 已重新执行 `go generate`，确认 `OAuthRefreshAPI` 的 Wire provider 缺口已修复，`wire_gen.go` 可重新生成。
- 已完成镜像构建并推送 `chengpengxiong/sub2api:0.1.114-rc1`，推送 digest 为 `sha256:7e47d8cd47b3230c36f03590ad559da7592b0e04e96c747511caba2c4b1caf12`。

## Validation
- `cd backend/cmd/server && go generate`
- `cd backend && go test ./cmd/server`
- `cd backend && go test -tags unit ./internal/service -run 'TestSettingService_(GetAPIKeyUsageGuide|GetPublicSettings|GetFrontendURL|UpdateSettings_)'`
- `docker build -t chengpengxiong/sub2api:0.1.114-rc1 --build-arg VERSION=0.1.114-rc1 --build-arg COMMIT=66630cb6 --build-arg DATE=2026-04-16T05:03:35Z --build-arg GOPROXY=https://goproxy.cn,direct --build-arg GOSUMDB=sum.golang.google.cn -f Dockerfile .`
- `docker push chengpengxiong/sub2api:0.1.114-rc1`

## Risks
- 工作区当前为脏树，本轮只更新版本文件并执行镜像构建/推送，不处理其他未提交改动。
- 未额外创建 git tag 或推送代码；当前只完成了本地版本更新和 Docker 镜像发布。

## Next Steps
- 如需让仓库版本和镜像版本对外一致，下一步应提交版本变更并按仓库约定创建/推送对应 git tag。

## Goal
- 修复 `Wire` 无法为 `OAuthRefreshAPI` 自动生成依赖注入代码的问题，消除 `no provider found for []time.Duration` 报错，同时保持现有默认刷新锁 TTL 逻辑不变。

## Todo
- 为 `OAuthRefreshAPI` 增加固定参数的 Wire provider 包装函数。
- 替换 `service.ProviderSet` 中对 variadic 构造函数的直接注册。
- 重新生成 `wire_gen.go` 并验证编译。

## Doing
- 正在修复 `OAuthRefreshAPI` 的 Wire provider 缺口并验证 `go generate`。

## Done
- 已确认根因是 `NewOAuthRefreshAPI(accountRepo, tokenCache, lockTTL ...time.Duration)` 的 variadic 参数会被 Wire 解析成缺失的 `[]time.Duration` provider。
- 已确认当前运行逻辑仍依赖两参默认调用，问题集中在代码生成阶段而非运行时行为。

## Validation
- 待补充。

## Risks
- 工作区当前为脏树，本轮只处理 `OAuthRefreshAPI` 的 Wire 注入缺口，不扩修无关依赖链。

## Next Steps
- 修复完成后执行 `go generate`、`go test ./cmd/server` 并更新本节状态。

## Goal
- 为设置类只读接口补充 `L1 + Redis L2` 两级缓存，覆盖公开设置、API Key 使用说明和前端 URL，减少数据库读取压力，并在设置更新后主动失效。

## Todo
- 无。

## Doing
- 无。

## Done
- 已确认仓库已有 Redis、singleflight 和多实例缓存失效先例，当前设置类读取尚未接入统一缓存。
- 已确认本次最小闭环适合只覆盖 `GetPublicSettings`、`GetAPIKeyUsageGuide` 和 `GetFrontendURL`。
- 已为 `SettingService` 增加可选 `SettingReadCache` 注入、进程内 L1 缓存、Redis L2 读取以及基于 Pub/Sub 的本地缓存失效订阅。
- 已让 `GetPublicSettings`、`GetAPIKeyUsageGuide` 和 `GetFrontendURL` 优先命中 L1，再查 Redis，最后才回源数据库。
- 已在 `UpdateSettings` 成功写库后主动清理本机 L1、删除 Redis bundle key，并发布失效消息，避免只靠 TTL 被动过期。
- 已新增 `backend/internal/repository/setting_read_cache.go` 并接入 Wire，同时手动同步 `backend/cmd/server/wire_gen.go` 中的 `SettingService` 注入参数。
- 已新增设置读取缓存单测，覆盖首次回源、Redis 命中和更新后失效重载场景。

## Validation
- `cd backend && go test -tags unit ./internal/service -run 'TestSettingService_(GetAPIKeyUsageGuide|GetPublicSettings|GetFrontendURL|UpdateSettings_)'`
- `cd backend && go test ./cmd/server`
- `cd backend && go test ./internal/repository -run '^$'`
- `git diff --check -- backend/internal/service/setting_service.go backend/internal/service/wire.go backend/internal/repository/setting_read_cache.go backend/internal/repository/wire.go backend/internal/service/setting_service_read_cache_test.go backend/cmd/server/wire_gen.go TODO.md`

## Risks
- 工作区当前为脏树，需严格限制改动范围，只处理设置读取缓存相关代码。
- `wire` 生成当前仍受工作区内另一条既有 DI 缺口影响：`initializeApplication` 缺少 `[]time.Duration` provider；本轮仅手动同步了 `SettingService` 相关的 `wire_gen.go` 调用，不扩修无关依赖链。

## Next Steps
- 如果后续要继续降低数据库压力，下一步优先把其他 `SettingService` 的只读入口按同样 bundle 方式接入，而不是直接做全站通用缓存框架。

## Goal
- 在关闭 Turnstile 后重启线上 `cc.taylor-link.xyz` 的应用容器，并确认公网配置与健康检查恢复正常。

## Todo
- 无。

## Doing
- 无。

## Done
- 已通过 `ssh taylor@34.92.180.210` 确认线上为 `/mnt/data/sub2api/docker-compose.yml` 的 Docker 部署，服务名为 `sub2api`。
- 已执行 `docker compose restart sub2api` 重启应用容器。
- 已确认重启后容器重新进入 `healthy` 状态。
- 已确认公网 `https://cc.taylor-link.xyz/api/v1/settings/public` 返回 `turnstile_enabled=false`。

## Validation
- `ssh taylor@34.92.180.210 'cd /mnt/data/sub2api && docker compose restart sub2api && docker compose ps sub2api'`
- `curl -fsS 'https://cc.taylor-link.xyz/api/v1/settings/public?timezone=Asia%2FShanghai'`
- `curl -fsS https://cc.taylor-link.xyz/health`

## Risks
- 当前虽然 `turnstile_enabled=false`，但公网设置里仍保留了 `turnstile_site_key` 字段；只要前端严格以 `turnstile_enabled` 为开关就不会生效，但若后续有错误依赖 site key 的逻辑，仍需继续检查。

## Next Steps
- 如需进一步确认大陆可用性，下一步应实测登录/注册页是否已完全不再加载 Turnstile 相关脚本或组件。

## Goal
- 优化默认 branding 的自定义 CSS 与辅助文案结构，在保持简洁风格的前提下强化 Google 风格的层级、节奏与品牌识别。

## Todo
- 无。

## Doing
- 无。

## Done
- 已将默认 branding 的首页品牌标识从单点阴影改为四色标记组合，并重写首页 badge、search shell、按钮、chips 与 footer 的视觉 token，强化 Google 风格但保留简洁结构。
- 已调整认证页默认 branding 覆盖：弱化 glass 效果、收敛背景装饰、统一输入框与主按钮状态，并为登录/注册说明块新增 `marker + eyebrow` 结构。
- 已同步更新 `frontend/src/constants/defaultBranding.ts` 与 `backend/internal/service/branding_defaults.go`，确保前端“恢复默认品牌”和后端初始化默认值保持一致。
- 已补强 `frontend/src/constants/__tests__/defaultBranding.spec.ts`，新增对新 token 与新 HTML 结构的断言。

## Validation
- `git diff --check -- TODO.md frontend/src/constants/defaultBranding.ts frontend/src/constants/__tests__/defaultBranding.spec.ts backend/internal/service/branding_defaults.go`
- `cd frontend && pnpm test:run src/constants/__tests__/defaultBranding.spec.ts`
- `cd backend && go test -tags unit ./internal/service -run 'TestSettingService_InitializeDefaultSettings_(FillsMissingBrandingDefaults|PreservesExplicitEmptyBrandingValues)'`

## Risks
- 默认 branding 常量当前在前后端各维护一份，若不同步修改会导致“恢复默认品牌”前后不一致。
- 工作区存在与本任务无关的现有变更：`frontend/package.json` 已被 Corepack 写入 `packageManager` 字段，本轮不应混入该改动。

## Next Steps
- 如需继续微调视觉效果，下一步应优先做登录页与首页的实际页面截图对比，再决定是否继续收紧色彩和间距。

## Goal
- 将 Cloudflare 控制台已开启的全局 Authenticated Origin Pulls（AOP）真正落地到源站 Nginx，要求源站只接受带 Cloudflare 共享客户端证书的 HTTPS 请求。

## Todo
- 无。

## Doing
- 无。

## Done
- 已确认 Cloudflare 控制台侧的全局 AOP 已开启，但源站 Nginx 之前仅配置了服务端证书 `ssl_certificate` / `ssl_certificate_key`，未配置 `ssl_client_certificate` / `ssl_verify_client`，因此 AOP 当时未真正生效。
- 已按 Cloudflare 官方共享证书模式下载并部署源站校验用 CA：`/etc/nginx/ssl/cloudflare-origin-pull-ca.pem`。
- 已新增 Nginx `log_format main_cf`，将 `ssl_client_verify`、客户端证书主体和签发者写入访问日志，便于验证 AOP 行为。
- 已在 `/etc/nginx/sites-enabled/cc.taylor-link.xyz` 启用 AOP，先以 `ssl_verify_client optional` 观测，再切换为 `ssl_verify_client on` 强制校验。
- 已清理误放入 `/etc/nginx/sites-enabled/` 的站点备份文件，避免 Nginx 因重复加载备份文件产生 `conflicting server name` 告警。

## Validation
- `curl -fsSL https://developers.cloudflare.com/ssl/static/authenticated_origin_pull_ca.pem | openssl x509 -noout -subject -issuer -fingerprint -sha256`
- `ssh ... 'sudo nginx -t && sudo systemctl reload nginx'`
- `curl -fsS https://cc.taylor-link.xyz/api/v1/settings/public?timezone=Asia%2FShanghai&aop_cf_enforced=1`
- `ssh ... \"curl -sk --resolve cc.taylor-link.xyz:443:127.0.0.1 https://cc.taylor-link.xyz/api/v1/settings/public?...\"`
- `ssh ... \"sudo grep -nE 'aop_cf_enforced=1|aop_local_enforced=1' /var/log/nginx/cc.taylor-link.xyz.access.log\"`

## Risks
- 当前使用的是 Cloudflare 全局共享 AOP 证书，只能证明请求来自 Cloudflare 网络，不能证明只来自你自己的 Cloudflare 账号；若要更严格隔离，需要改为 zone-level 或 per-hostname AOP 自有证书模式。
- 如果后续替换 Nginx 站点文件或重建服务器，必须同时恢复 `ssl_client_certificate`、`ssl_verify_client on` 和 `cloudflare-origin-pull-ca.pem`，否则会退回到“Cloudflare 控制台开了 AOP、但源站未校验”的假生效状态。

## Next Steps
- 如果要继续收紧，可考虑将目前仍保留的 3 个源站 `80/443` 白名单 IP 清理掉，只保留 Cloudflare 入站。
- 如果要把安全强度从“来自 Cloudflare 网络”提升到“来自你自己的 Cloudflare zone/hostname”，下一步改用 zone-level 或 per-hostname AOP 自有证书。

## Goal
- 收紧服务器 SSH 入口，只允许固定管理 IP 登录，并将当前执行主机公网出口 `104.225.146.113` 纳入白名单。

## Todo
- 无。

## Doing
- 无。

## Done
- 已确认当前执行主机公网出口 IP 为 `104.225.146.113`，且与现有源站白名单 IP 之一一致。
- 已在线上为 `22/tcp` 新增 3 条显式 UFW 放行规则，仅允许以下 IPv4 来源访问 SSH：
  - `212.50.248.185`
  - `154.44.21.124`
  - `104.225.146.113`
- 已删除原有 `22/tcp ALLOW IN Anywhere` 与对应 IPv6 全网放行规则，SSH 不再对全网开放。

## Validation
- `curl -4 https://ifconfig.me`
- `ssh ... 'sudo ufw status numbered'`
- `ssh ... 'echo ssh-ok'`
- `ssh ... 'sudo fail2ban-client status'`
- `curl -fsS https://cc.taylor-link.xyz/health`

## Risks
- 当前 SSH 白名单只包含 3 个 IPv4 出口；如果你的管理出口 IP 变更，新的来源将无法直接登录，需要从现有白名单 IP 之一进入后再调整 UFW。
- 当前未额外放开 IPv6 SSH；如果未来你的管理链路切到 IPv6，需要单独补充白名单。

## Next Steps
- 如果确认不再需要那 3 个 IP 直连源站 `80/443`，下一步可以把它们从 Web 白名单里移除，只保留 Cloudflare 入站。

## Goal
- 开启应用自身的 `security.url_allowlist` 以补齐 SSRF 防护，并收紧服务器对 `20201/20202` 的防火墙与 fail2ban 配置。

## Todo
- 无。

## Doing
- 无。

## Done
- 已确认应用侧 URL allowlist 功能已经存在，但线上 `.env` 仍显式配置为 `SECURITY_URL_ALLOWLIST_ENABLED=false`，且允许 `http://` 与私网主机。
- 已确认 `20201/20202` 分别由 Google Cloud Ops Agent 的 `otelopscol` 与 `fluent-bit` 监听在公网地址，但当前 UFW 默认策略为 `deny incoming`，属于“默认被拦截但缺少显式规则”的状态。
- 已确认 fail2ban 已安装并启用，当前只对 `sshd` 和 `nginx-http-auth` 生效。
- 已收紧部署示例配置，将 URL allowlist 的生产默认值改为启用，并将 `allow_insecure_http` / `allow_private_hosts` 默认值改为 `false`。
- 已在线上启用 URL allowlist，并显式限制上游/定价下载主机，禁止私网主机和明文 HTTP。
- 已在线上为 `20201/tcp` 与 `20202/tcp` 增加显式 UFW deny 规则。
- 已在线上强化 fail2ban，保留 `sshd`，新增 `nginx-botsearch` 与 `recidive` jail，并统一通过 UFW 执行封禁。

## Validation
- `git diff --check`
- `ssh ... 'sudo ufw status numbered'`
- `ssh ... 'sudo fail2ban-client status'`
- `ssh ... 'docker compose restart sub2api && docker logs sub2api | grep url_allowlist'`
- `curl -fsS https://cc.taylor-link.xyz/health`

## Risks
- URL allowlist 开启后，若后续新增自定义上游域名、CRS 源或定价源，但未同步加入白名单，会被应用主动拒绝。
- `20201/20202` 目前仍监听在公网地址，只是被防火墙显式阻断；如果未来希望进一步缩小暴露面，还应继续把 Google Ops Agent 监听地址改为 `127.0.0.1`。

## Next Steps
- 如果后续启用自定义上游或私有网络上游，先评估风险，再按需增加白名单，不要重新放开 `allow_private_hosts`。
- 后续可继续收紧 SSH，只允许固定管理 IP，并把额外 `80/443` 白名单 IP 清理掉。

## Goal
- 修复 SMTP 测试连接与真实发信仅支持 `AUTH PLAIN` 的兼容性问题，支持 Outlook / Microsoft 365 这类会拒绝 `AUTH PLAIN` 但接受 `AUTH LOGIN` 的服务端。

## Todo
- 无。

## Doing
- 无。

## Done
- 已确认 `backend/internal/service/email_service.go` 中测试连接与真实发信都硬编码使用 `smtp.PlainAuth(...)`，这是当前 `504 5.7.4 Unrecognized authentication type` 的直接触发点。
- 已新增 SMTP AUTH 机制探测逻辑，按服务端 `EHLO` 返回的 `AUTH` 能力构造认证尝试。
- 已新增 `AUTH LOGIN` 支持，并在服务端返回“不识别认证类型”时，重新建连后从 `PLAIN` 回退到 `LOGIN`，避免同一连接重试导致的 broken pipe。
- 已将测试连接与真实发信统一复用同一套 SMTP 建连与认证逻辑，避免行为分叉。
- 已补充 SMTP 单元测试，覆盖 `PLAIN -> LOGIN` 回退和不支持认证机制的报错场景。

## Validation
- `cd backend && go test -tags unit ./internal/service -run 'TestEmailService_(TestSMTPConnectionWithConfig|GetSMTPConfig|NormalizeSMTPSecurityModeWithLegacy)'`
- `git diff --check`

## Risks
- 本轮只补 `AUTH PLAIN` 与 `AUTH LOGIN`，未扩展 `XOAUTH2`；如果未来接入强制 OAuth2 的 SMTP 服务，还需要单独实现。
- 若邮件服务端宣告的 `AUTH` 能力与实际行为不一致，当前只会在“服务端明确拒绝认证类型”时做重连回退，不会无条件穷举所有机制。

## Next Steps
- 重新在管理后台执行一次 SMTP 连接测试，确认 Outlook `587 + STARTTLS` 场景从 `504` 转为成功。
- 如果测试连接成功，再实际发送一封验证码邮件，确认真实发信链路也正常。

## Goal
- 修复生产环境 Cloudflare -> Nginx -> Docker bridge -> Sub2API 的真实 IP 信任链，补齐 `trusted_proxies` / `cors.allowed_origins`，修正 Cloudflare IP 自动更新脚本的中断风险，并把 Nginx `client_max_body_size` 统一为 256MB；同步更新 `/root/gcp-sub2api` 部署文档。

## Todo
- 无。

## Doing
- 无。

## Done
- 已确认线上实际部署不是 host network，而是 `Cloudflare -> Nginx -> 127.0.0.1:8080 -> Docker bridge(172.18.0.1 -> 172.18.0.4)`。
- 已确认当前 Nginx 未配置 `real_ip_header CF-Connecting-IP` / `set_real_ip_from`，公网 access log 记录的是 Cloudflare 边缘 IP。
- 已确认当前应用未配置 `server.trusted_proxies` 与 `cors.allowed_origins`，启动日志存在对应 warning。
- 已确认当前 Cloudflare IP 更新脚本存在高风险：先删旧规则再拉新 IP、删除逻辑不可靠、且不会同步更新 Nginx `set_real_ip_from`。
- 已确认公网响应头存在 `X-Frame-Options: DENY` 与 `X-Frame-Options: SAMEORIGIN` 双写冲突。
- 已完成远端 Nginx 修正：新增 Cloudflare `real_ip` 白名单、统一代理头为 `$remote_addr`、移除重复安全头、将 `client_max_body_size` 调整为 `256m`。
- 已完成远端应用修正：`/mnt/data/sub2api/config.yaml` 新增 `server.trusted_proxies` 与 `cors.allowed_origins`，容器重启后 warning 消失。
- 已完成远端脚本修正：`/opt/sub2api-backup/scripts/update_cloudflare_ips.sh` 改为先拉取校验、再增量同步 UFW 与 Nginx `real_ip`，并成功生成 `/opt/sub2api-backup/state/cloudflare-ips-v4.txt` 与 `cloudflare-ips-v6.txt`。
- 已完成部署文档修正：`/root/gcp-sub2api` 中相关文档已从 host network/危险 cron 示例更新为 bridge 网络、可信代理与安全更新脚本方案。

## Validation
- `git diff --check`
- `nginx -t`
- `docker compose restart sub2api`
- `docker logs sub2api | grep trusted_proxies`
- `docker logs sub2api | grep allowed_origins`
- `curl -I https://cc.taylor-link.xyz/`
- `tail /var/log/nginx/cc.taylor-link.xyz.access.log`
- `ufw status numbered`
- `tail /var/log/cloudflare-ip-update.log`

## Risks
- Docker bridge 网关如果后续重建变成非 `172.18.0.1`，需要同步更新 `server.trusted_proxies`。
- 目前仍保留少量源站白名单 IP；它们只能保留在 UFW 中，不能进入 trusted proxy / real_ip 信任链。
- 文档包含历史敏感信息，本轮只修正部署事实，不扩散改动范围。

## Next Steps
- 如果后续重建 `docker compose` 网络，先用 `docker network inspect sub2api_sub2api-network` 确认 gateway，再同步更新 `server.trusted_proxies`。
- 如果未来取消源站直连白名单，可再收紧 UFW，仅保留 Cloudflare 段与 SSH 管理入口。

## Goal
- 审查上游 PR `#1668`（`https://github.com/Wei-Shaw/sub2api/pull/1668/changes`）是否应合并到当前分支。

## Todo
- 无。

## Doing
- 无。

## Done
- 已读取本地 `TODO.md`、确认当前工作区为脏树，不切换当前分支、不覆盖现有改动。
- 已只读拉取 `upstream/pr-1668`，并对比 `upstream/main...upstream/pr-1668` 的服务层与测试层差异。
- 已在隔离 worktree `/tmp/sub2api-pr1668` 跑相关单测，确认 PR 分支可通过目标测试集。
- 已在当前工作区跑同一组目标单测，确认主线当前也能通过。
- 已确认该 PR 的核心行为不是局部修正，而是移除了多处 `codex extra / snapshot -> RateLimitResetAt` 的回写路径；同时保留了真实 429 路径中的 `calculateOpenAI429ResetTime(...)` 逻辑，因此若目标是“7d 可用时不因 5h 为 0 回写 429”，该 PR 不是完整修复。

## Validation
- `git fetch upstream pull/1668/head:refs/remotes/upstream/pr-1668`
- `git diff --unified=20 upstream/main...upstream/pr-1668 -- backend/internal/service/account_usage_service_test.go backend/internal/service/admin_service.go`
- `cd /tmp/sub2api-pr1668/backend && go test -tags unit ./internal/service -run 'Test(AccountTestService_OpenAI|AccountUsageService_|ExtractOpenAICodex|OpenAIGatewayService_UpdateCodexUsageSnapshot|OpenAIGatewayService_GetSchedulableAccount|AdminService_ListAccounts_ExhaustedCodexExtra|OpenAIGatewayService_ProxyResponsesWebSocketFromClient_ErrorEventUsageLimitReachedMarksAccountRateLimited)'`
- `cd /root/sub2api/backend && go test -tags unit ./internal/service -run 'Test(AccountTestService_OpenAI|AccountUsageService_|ExtractOpenAICodex|OpenAIGatewayService_UpdateCodexUsageSnapshot|OpenAIGatewayService_GetSchedulableAccount|AdminService_ListAccounts_ExhaustedCodexExtra|OpenAIGatewayService_ProxyResponsesWebSocketFromClient_ErrorEventUsageLimitReachedMarksAccountRateLimited)'`

## Risks
- 如果直接合并 PR，会改变当前“通过 snapshot/extra 提前补齐运行时限流状态”的策略，可能让真正已耗尽的账号继续进入调度，直到下一次真实 429 才被动落限流。
- 如果不修正 `backend/internal/service/ratelimit_service.go` 里的真实 429 判定逻辑，而只合并这个 PR，标题描述的问题仍可能在真实 429 路径存在。

## Next Steps
- 若要修标题里的问题，建议做更小的定向修复：明确哪些场景允许 `5h` 窗口耗尽提升为 `RateLimitResetAt`，哪些场景只更新展示用 `extra`，并同步修改真实 429 路径与对应测试。

## Goal
- 分析 `cache creation` 统计长期为 0 的原因，并评估提交 `1b33ccd27705294a3d14aac2d7002f089ea4df75` 是否已解决。

## Todo
- 无。

## Doing
- 无。

## Done
- 已确认 `usage_log`、聚合 SQL 和前端展示都已支持 `cache_creation_tokens`，因此“总是 0”不是聚合层或前端展示层缺字段。
- 已确认当前 OpenAI 网关链路在 [openai_gateway_service.go](/root/sub2api/backend/internal/service/openai_gateway_service.go#L4016) 与 [openai_gateway_service.go](/root/sub2api/backend/internal/service/openai_gateway_service.go#L3998) 只解析 `input_tokens`、`output_tokens` 和 `input_tokens_details.cached_tokens`，没有解析 `cache_creation_input_tokens`，因此写库时 [openai_gateway_service.go](/root/sub2api/backend/internal/service/openai_gateway_service.go#L4637) 会持续写入 0。
- 已确认项目里另一条通用网关链路 [gateway_service.go](/root/sub2api/backend/internal/service/gateway_service.go#L5115) 已经解析 `cache_creation_input_tokens` 及 `cache_creation.ephemeral_5m_input_tokens` / `ephemeral_1h_input_tokens`，说明问题集中在 OpenAI 专用链路而不是全局设计不支持。
- 已确认提交 `1b33ccd27705294a3d14aac2d7002f089ea4df75` 只修改 `backend/internal/service/openai_gateway_messages.go`，内容是为 `messages -> responses` 兼容路径补写 `prompt_cache_key`，未修改任何 usage 解析或 usage log 落库逻辑。

## Validation
- `rg -n "cache_creation_input_tokens|CacheCreationInputTokens|cache_creation_tokens" backend frontend`
- `git fetch https://github.com/FjlI5/sub2api.git 1b33ccd27705294a3d14aac2d7002f089ea4df75`
- `git show --stat --summary 1b33ccd27705294a3d14aac2d7002f089ea4df75`
- `nl -ba backend/internal/service/openai_gateway_service.go | sed -n '3998,4032p'`
- `nl -ba backend/internal/service/openai_gateway_service.go | sed -n '4527,4642p'`
- `nl -ba backend/internal/service/gateway_service.go | sed -n '5038,5076p'`
- `nl -ba backend/internal/service/gateway_service.go | sed -n '5110,5128p'`

## Risks
- 如果上游本身没有返回 `cache_creation_input_tokens`，即使修了解析，本地统计仍会是 0；但当前代码已能确定“至少存在本地漏解析”这一层问题。
- 你给的那个提交可能改善某些 `messages -> responses` 路径的上游缓存命中稳定性，但它不会直接修复“统计字段一直为 0”的落库问题。

## Next Steps
- 若要真正修复，应该在 OpenAI 网关的 JSON/SSE usage 解析里补齐 `cache_creation_input_tokens`，并在上游返回 TTL 明细时一并解析 `cache_creation.ephemeral_5m_input_tokens` / `ephemeral_1h_input_tokens`。

## Goal
- 将 `upstream/main@be7551b9` 合入当前本地 `main`，保留本地既有修复与 branding 改动，并吸收上游支付、账号成本、SSE 关闭与 OpenAI 限流修复。

## Todo
- 无。

## Doing
- 无。

## Done
- 已在隔离 worktree 上执行真实 `git merge upstream/main`，并将结果快进到当前本地 `main`，当前 `HEAD` 为合并提交 `66630cb6`。
- 已吸收上游支付充值倍率 / 手续费率、账号成本展示、账号测试 SSE 可关闭，以及 OpenAI 账号限流回流修复相关改动。
- 已保留本地 `backend/cmd/server/VERSION` 为 `0.1.113-rc1`，未跟随上游改成正式版 `0.1.113`。
- 已保留本地既有 branding 与其他历史修复；自动合并未产生除 `VERSION` 以外的文本冲突。
- 已补一个最小稳定性修复：`frontend/src/views/admin/DashboardView.vue` 的 `formatCost()` 现在可处理 `undefined/null`，避免新引入的 `account_cost` 字段在旧 mock/旧数据下触发渲染异常。
- 已清理临时合并 worktree 与临时分支，未把验证用 `node_modules` 软链混入仓库。

## Validation
- `git fetch upstream main`
- `git merge --ff-only codex/merge-upstream-main`
- `cd /tmp/sub2api-merge/backend && go test -tags unit ./internal/service/... ./internal/handler/... ./internal/server/...`
- `cd /tmp/sub2api-merge/frontend && COREPACK_ENABLE_AUTO_PIN=0 pnpm test:run src/components/account src/views/admin src/views/user src/constants`
- `cd /tmp/sub2api-merge/frontend && COREPACK_ENABLE_AUTO_PIN=0 pnpm lint:check`
- `cd /tmp/sub2api-merge/frontend && COREPACK_ENABLE_AUTO_PIN=0 pnpm typecheck`
- `cd /tmp/sub2api-merge && make test` 失败于环境缺少 `golangci-lint`，不是代码测试失败。

## Risks
- 工作区仍保留与本任务无关的现有改动：`TODO.md`、`frontend/package.json`、`.playwright-mcp/`；提交本次合并时需要注意不要误混入。
- `Antigravity-Manager` gitlink 目前随 `upstream/main` 一并进入本地分支；上游树当前就是这个状态，本轮未额外清理。
- `make test` 依赖本机安装 `golangci-lint`；当前仅能确认 Go 测试、前端测试、前端 lint/typecheck 通过，不能声称完整 CI 基线已在本机跑通。

## Next Steps
- 如要提交本次结果，先明确是否保留 `Antigravity-Manager` gitlink，再决定是否额外做一次清理提交。
- 如要补全本地 CI 验证，先安装 `golangci-lint` 后重新执行 `make test`。

## Goal
- 修复用户端无法查看可用分组的问题，让用户可以直接看到自己可使用的分组，并且只暴露基础信息。

## Todo
- 无。

## Doing
- 无。

## Done
- 已确认现有后端已具备“按用户权限过滤可用分组”的能力，但仅被 `KeysView` 内部下拉框使用，用户侧没有独立入口。
- 已新增用户侧摘要接口 `GET /api/v1/groups/available/summary`，复用现有权限判断，仅返回 `id/name/description/platform/rate_multiplier/subscription_type/access_scope`。
- 已保持现有 `GET /api/v1/groups/available` 不变，避免影响 `frontend/src/views/user/KeysView.vue` 现有分组切换与 `allow_messages_dispatch` 相关逻辑。
- 已新增用户页 [GroupsView.vue](/root/sub2api/frontend/src/views/user/GroupsView.vue)，用户可以在侧边栏直接查看当前账号可使用的分组、分组类型、平台、倍率和可用资格说明。
- 已在 [AppSidebar.vue](/root/sub2api/frontend/src/components/layout/AppSidebar.vue) 和 [router/index.ts](/root/sub2api/frontend/src/router/index.ts) 增加用户侧“可用分组”入口与路由。
- 已补充后端合约测试和前端视图/导航测试，覆盖摘要接口的权限过滤和用户入口可见性。

## Validation
- `cd backend && go test -tags unit ./internal/server -run 'TestAPIContracts'`
- `cd frontend && COREPACK_ENABLE_AUTO_PIN=0 pnpm test:run src/views/user/__tests__/GroupsView.spec.ts src/components/layout/__tests__/AppSidebar.spec.ts`
- `cd frontend && COREPACK_ENABLE_AUTO_PIN=0 pnpm typecheck`
- `cd frontend && COREPACK_ENABLE_AUTO_PIN=0 pnpm lint:check`

## Risks
- 当前新增的是“摘要接口 + 独立页面”，未把老的 `/groups/available` 一并收窄；这是刻意保留兼容性，意味着短期内用户侧仍存在一条更宽的历史接口供 `KeysView` 使用。
- `TODO.md`、`frontend/package.json` 与 `.playwright-mcp/` 在本轮开始前就已经是脏工作区状态，提交时仍需注意不要误混入无关改动。

## Next Steps
- 如果后续希望进一步收紧用户侧数据暴露，可以再评估把 `KeysView` 逐步迁移到更细粒度的类型，而不是继续复用宽 `Group` DTO。
- 如果希望用户更容易理解“为什么我能用这个分组”，下一步可以在页面上补一段资格说明文案，但本轮先保持最小可用实现。

## Goal
- 为网关请求体读取失败补最小可观测性，并把这类失败统一收敛成可重试的 JSON 错误提示，避免 API 用户看到笼统的 body read 失败文案。

## Todo
- 无。

## Doing
- 无。

## Done
- 已新增统一的请求体读取失败处理 helper，覆盖 OpenAI / Anthropic / Claude-compatible / Gemini 入口，响应中统一附带 `Retry-After: 1`，并返回明确的“请重试” JSON 文案。
- 已为 OpenAI / Anthropic / Claude-compatible 错误体补充稳定的 `error.code=request_body_read_failed`；Gemini 兼容错误体补充 `details.reason=REQUEST_BODY_READ_FAILED` 与 `retry_after_seconds=1`。
- 已为这类失败补充结构化日志，记录底层读流错误、路径、方法、`content_length`、`transfer_encoding`、`content_type`、`user_agent` 与客户端 IP，不记录原始请求体。
- 已让 `ops_error_logs` 识别这类错误码并标记为 `is_retryable=true`，同时把响应头里的 `Retry-After` 落入 `retry_after_seconds`。
- 已将仓库内 `deploy/Caddyfile` 的 API 错误页示例改为 JSON，避免受管 Caddy 部署在代理层返回纯文本错误页。
- 已补充 handler 单测，覆盖返回格式、`Retry-After`、错误码提取和 retryable 分类。

## Validation
- `cd backend && go test -tags unit ./internal/handler -run 'Test(HandleRetryableRequestBodyReadError|ParseOpsErrorResponse_ExtractsStructuredCode|ClassifyOpsIsRetryable_RequestBodyReadFailed|ParseRetryAfterSeconds|ReadRequestBodyWithPrealloc|ReadRequestBodyWithPrealloc_MaxBytesError|OpenAIHandleStreamingAwareError_)'`
- `cd backend && go test -tags unit ./internal/repository -run 'Test(DoesNotExist)'`
- `cd backend && go test -tags unit ./internal/service -run 'Test(DoesNotExist)'`
- `git diff --check -- backend/internal/handler/request_body_read_error.go backend/internal/handler/request_body_read_error_test.go backend/internal/handler/openai_gateway_handler.go backend/internal/handler/openai_chat_completions.go backend/internal/handler/gateway_handler.go backend/internal/handler/gateway_handler_responses.go backend/internal/handler/gateway_handler_chat_completions.go backend/internal/handler/gemini_v1beta_handler.go backend/internal/handler/ops_error_logger.go backend/internal/handler/ops_error_logger_test.go backend/internal/service/ops_port.go backend/internal/repository/ops_repo.go deploy/Caddyfile`
- `caddy` 本地二进制缺失；尝试用 Docker `caddy:2.10.0` 做 `caddy adapt` 校验时因拉镜像网络超时失败，Caddyfile 仅完成静态检查，未完成运行时语法验证。

## Risks
- 当前仓库不包含线上 Nginx 配置，因此这轮只能保证“应用自己处理到的错误”返回 JSON；如果请求在外层 Nginx 就被拒绝，线上仍需要单独做 API 错误页 JSON 化。
- `ops_error_logs.retry_after_seconds` 已开始写入，但现有管理端列表/详情还没有专门展示这个字段；当前主要用于分类与后续排查。
- 工作区里仍有本任务之外的脏改动，提交时要避免混入 `frontend/package.json`、`.playwright-mcp/` 以及上一轮未提交的用户分组改动。

## Next Steps
- 如需彻底消灭线上 HTML 错误页，需要把同样的 API 错误页 JSON 化策略落到实际使用的 Nginx 配置。
- 如需让客户端更容易自动重试，可以后续再评估是否补充额外的显式重试头或 SDK 侧重试规则，但本轮先保持 `Retry-After + stable error.code`。

## Goal
- 将“使用 API 密钥”按钮对应的说明从前端静态文案改为管理员可编辑，并让用户点击时动态请求最新配置，避免每次改文案都重新打包上线。

## Todo
- 无。

## Doing
- 无。

## Done
- 已在后端新增专用接口 `GET /api/v1/keys/usage-guide`，返回最新 `api_base_url` 与结构化的 API Key 使用说明文案。
- 已在设置域模型中新增 `api_key_usage_guide_content`，支持初始化默认值、管理后台保存、读取解析与审计差异记录。
- 已在管理后台“站点设置”中新增 API Key 使用说明编辑区，管理员可直接编辑通用说明、未分组提示、OpenAI / Gemini / Antigravity / OpenCode 的基础文案。
- 已把用户侧 [UseKeyModal.vue](/root/sub2api/frontend/src/components/keys/UseKeyModal.vue) 改为优先使用后端动态文案，并在缺失字段时回退到现有 i18n 默认值。
- 已把用户侧 [KeysView.vue](/root/sub2api/frontend/src/views/user/KeysView.vue) 改为点击“使用 API 密钥”后动态请求最新说明，并在关闭弹窗或组件卸载时中止未完成请求。
- 已保持生成配置文件/代码片段的逻辑仍在前端，只把需要频繁调整的说明文案与 `api_base_url` 动态化，控制改动范围。
- 已补充后端单测与接口合约测试，并补充前端弹窗测试覆盖服务端文案覆盖和动态 base URL 生效。

## Validation
- `cd backend && go test -tags unit ./internal/service -run 'TestSettingService_(GetAPIKeyUsageGuide|ParseAPIKeyUsageGuideContent_InvalidJSONReturnsEmpty|GetPublicSettings)'`
- `cd backend && go test -tags unit ./internal/server -run 'TestAPIContracts'`
- `cd backend && go test -tags unit ./...`
- `cd frontend && pnpm test:run src/components/keys/__tests__/UseKeyModal.spec.ts`
- `cd frontend && pnpm typecheck`
- `cd frontend && pnpm exec eslint src/api/admin/settings.ts src/api/keys.ts src/components/keys/UseKeyModal.vue src/components/keys/__tests__/UseKeyModal.spec.ts src/types/index.ts src/views/admin/SettingsView.vue src/views/user/KeysView.vue src/i18n/locales/en.ts src/i18n/locales/zh.ts`
- `cd frontend && pnpm build`

## Risks
- 当前后端存储的是一份统一文案，不区分多语言；如果后续需要中英文分别维护，需要再扩展数据结构和管理界面。
- 本轮刻意没有把这份说明并入 `GET /api/v1/settings/public`，而是只通过已登录用户专用接口返回，以减少无关暴露面。
- 工作区仍存在本任务之外的既有脏改动，如用户分组页、请求体读流错误处理、`frontend/package.json`、`.playwright-mcp/` 等，提交时需要谨慎选择文件。

## Next Steps
- 如果后续希望支持多语言管理员文案，下一步应先确定是否采用“按 locale 存整份文案”还是“字段级多语言对象”的数据模型。
- 如果希望用户打开弹窗时有更明确的加载反馈，可以再补一个轻量 loading 态，但本轮先保持最小改动。
