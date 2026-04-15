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
