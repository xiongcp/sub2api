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
