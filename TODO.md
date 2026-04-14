## Goal
- 将当前 `main` 上已合并 upstream 的代码与本地业务改动发布为 `0.1.113-rc1`，重点回归计费与模型映射，随后提交、推分支、打 tag，并推送 `chengpengxiong/sub2api` 镜像。

## Todo
- 执行计费、模型映射、支付与构建相关回归。
- 回归通过后提交并推送 `main`。
- 创建并推送 `v0.1.113-rc1` tag。
- 构建并推送 Docker Hub 镜像 `chengpengxiong/sub2api:0.1.113-rc1` 与 `chengpengxiong/sub2api:latest`。

## Doing
- 更新版本号与发布执行状态。

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

## Validation
- 已执行：
- `cd backend && go test ./cmd/server ./internal/server/routes ./internal/handler/... ./internal/service/...`
- `cd backend && go test ./internal/service/...`
- `cd frontend && ./node_modules/.bin/eslint src/views src/components src/api src/composables src/utils src/i18n src/types --ext .ts,.vue`
- `cd frontend && ./node_modules/.bin/vue-tsc --noEmit`
- `cd frontend && ./node_modules/.bin/vitest run`
- 当前状态：验证通过；仓库已回到 `main`，原工作区改动已回放并保留。
- 待执行本轮发布专项回归、构建、镜像 smoke 与推送验证。

## Risks
- 当前 `main` 上仍保留本地修改与新增未跟踪文件，尚未形成新的本地提交。
- `stash@{0}` 仍保留，作为额外回退保险；如果后续确认不需要，可手动删除。
- upstream 主线引入的大量迁移、支付、通知、WebSearch 变更已落入本地 `main`，后续继续开发时应基于这些新结构继续改动。
- 这次会把 `0.1.113-rc1` 同时推成 `latest`，默认拉取 `latest` 的环境会直接拿到 RC 镜像。

## Next Steps
- 跑发布专项回归与构建验证。
- 回归通过后提交、推分支、打 tag、推镜像。
