## Goal
- 修复 SMTP 配置只区分 `smtp_use_tls` 的问题，支持 `STARTTLS` / 隐式 TLS / 明文三种模式，解决 Outlook 587 握手失败；保持旧配置兼容，并完成代码提交、推送与 Docker 镜像发布。

## Todo
- 无。

## Doing
- 无。

## Done
- 已确认问题存在：当前 `smtp_use_tls=true` 会直接走隐式 TLS，连接 `smtp-mail.outlook.com:587` 时会报 `tls: first record does not look like a TLS handshake`。
- 已确认根因：Outlook 587 端口要求明文建连后执行 `STARTTLS`，不是隐式 TLS；当前“测试连接”和“真实发信”两套逻辑也不一致，容易出现测试结果与实际发信不一致。
- 已确定修复方案：增加 `smtp_security_mode` 显式安全模式，兼容旧字段 `smtp_use_tls`；统一 SMTP 拨号与认证逻辑；后台设置页改为显式模式选择，默认 `starttls`。

## Validation
- `git diff --check`
- `cd backend && go test ./internal/service/... ./internal/handler/... ./internal/server/...`
- `cd frontend && pnpm lint`
- `cd frontend && pnpm typecheck`
- `docker build -t chengpengxiong/sub2api:v0.1.113-rc1-hotfix1 -t chengpengxiong/sub2api:latest .`

## Risks
- 线上已保存的老配置只有 `smtp_use_tls`，需要兼容映射，避免保存前后行为漂移。
- 新镜像将按当前环境构建为单架构镜像；如果线上需要多架构发布，需要后续补齐 `buildx` / QEMU。

## Next Steps
- 完成代码改造、测试、提交与推送。
- 构建并推送 `chengpengxiong/sub2api:v0.1.113-rc1-hotfix1` 与 `chengpengxiong/sub2api:latest`。
