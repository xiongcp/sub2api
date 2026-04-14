## Goal
- 将当前 `main` 上已合并 upstream 的代码与本地业务改动发布为 `0.1.113-rc1`，重点回归计费与模型映射，随后提交、推分支、打 tag，并推送 `chengpengxiong/sub2api` 镜像。

## Todo
- 无。

## Doing
- 无。

## Done
- 创建备份分支 `backup/main-pre-upstream-20260414`。
- 导出补丁 `../sub2api-pre-upstream-20260414.patch`。
- stash 当前工作区改动 `pre-upstream-merge-20260414`。
- 新增 `upstream` 远端并拉取官方 `main`。
- 创建临时集成分支 `merge/upstream-main-20260414`。
- 将 `upstream/main` 合并到临时集成分支。
- 清理 `packageManager` 等环境副作用。
- 适配 upstream 合并后的后端/前端测试基线。
- 将临时集成分支 fast-forward 回本地 `main`。
- 回放原工作区 stash，并解决 `AccountCapacityCell.vue`、`CreateAccountModal.vue`、`EditAccountModal.vue`、`QuotaLimitCard.vue` 四个冲突文件。
- 确认 Docker 本地环境可用，当前已登录 Docker Hub 用户 `chengpengxiong`。
- 将版本号提升到 `0.1.113-rc1`。
- 执行后端主干回归、前端静态检查、Vitest、专项计费/模型映射测试和整仓构建。
- 提交发布提交 `e1c94aee chore(release): prepare 0.1.113-rc1`。
- 推送 `origin/main` 并创建、推送标签 `v0.1.113-rc1`。
- 构建镜像 `chengpengxiong/sub2api:0.1.113-rc1`，本地 smoke 通过。
- 推送 Docker Hub 镜像 `chengpengxiong/sub2api:0.1.113-rc1` 与 `chengpengxiong/sub2api:latest`。

## Validation
- 已执行：
- `cd backend && go test ./cmd/server ./internal/server/routes ./internal/handler/... ./internal/service/...`
- `cd backend && go test ./internal/service/...`
- `cd frontend && ./node_modules/.bin/eslint src/views src/components src/api src/composables src/utils src/i18n src/types --ext .ts,.vue`
- `cd frontend && ./node_modules/.bin/vue-tsc --noEmit`
- `cd frontend && ./node_modules/.bin/vitest run`
- `cd backend && go test -count=1 ./cmd/server ./internal/server/routes ./internal/handler/... ./internal/service/... ./internal/repository/...`
- `cd frontend && pnpm lint`
- `cd frontend && pnpm typecheck`
- `cd frontend && pnpm test:run`
- `cd backend && go test -count=1 -run 'Test(OpenAIGatewayServiceRecordUsage|ResolveAccountStatsCost|CalculateCostUnified|UsageLogRepository)' ./internal/service ./internal/repository`
- `make build`
- `git diff --check`
- `docker build -f Dockerfile -t chengpengxiong/sub2api:0.1.113-rc1 --build-arg VERSION=0.1.113-rc1 --build-arg COMMIT=e1c94aee --build-arg DATE=2026-04-14T14:28:15Z --build-arg GOPROXY=https://goproxy.cn,direct --build-arg GOSUMDB=sum.golang.google.cn .`
- 分步 smoke：Postgres + Redis + `chengpengxiong/sub2api:0.1.113-rc1` 启动后，`curl http://127.0.0.1:38080/health` 返回 `{"status":"ok"}`。
- Docker Hub 推送结果：
- `chengpengxiong/sub2api:0.1.113-rc1` digest `sha256:4664ed97f9a6de2bed73fc868d75e0b27997b97e96d3339fb67797d4642b3b91`
- `chengpengxiong/sub2api:latest` digest `sha256:4664ed97f9a6de2bed73fc868d75e0b27997b97e96d3339fb67797d4642b3b91`

## Risks
- 当前仓库唯一未跟踪项是本地环境下的 `AGENTS.md`，未纳入提交。
- `stash@{0}` 仍保留，作为额外回退保险；如果后续确认不需要，可手动删除。
- upstream 主线引入的大量迁移、支付、通知、WebSearch 变更已落入本地 `main`，后续继续开发时应基于这些新结构继续改动。
- 这次会把 `0.1.113-rc1` 同时推成 `latest`，默认拉取 `latest` 的环境会直接拿到 RC 镜像。

## Next Steps
- 在线上环境验证一条真实计费请求和一条真实模型映射请求，确认 usage、额度、费用三处口径一致。
- 视需要删除 `stash@{0}` 与临时集成分支 `merge/upstream-main-20260414`。
