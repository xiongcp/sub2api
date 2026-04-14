## Goal
- 将官方 `Wei-Shaw/sub2api` 的 `upstream/main` 安全合并到本地仓库，并保留当前本地改动与回退路径。

## Todo
- 清理临时集成分支上的环境副作用和无关改动。
- 修复 upstream 合并后暴露的测试兼容问题。
- 将临时集成分支 fast-forward 回本地 `main`。
- 回放并整理原工作区 stash。

## Doing
- 处理 `merge/upstream-main-20260414` 分支上的测试失败与环境副作用。

## Done
- 创建备份分支 `backup/main-pre-upstream-20260414`。
- 导出补丁 `../sub2api-pre-upstream-20260414.patch`。
- stash 当前工作区改动 `pre-upstream-merge-20260414`。
- 新增 `upstream` 远端并拉取官方 `main`。
- 创建临时集成分支 `merge/upstream-main-20260414`。
- 将 `upstream/main` 合并到临时集成分支。

## Validation
- 已执行：
- `cd backend && go test ./cmd/server ./internal/server/routes ./internal/handler/... ./internal/service/...`
- `cd frontend && pnpm lint`
- `cd frontend && ./node_modules/.bin/vue-tsc --noEmit`
- `cd frontend && pnpm test:run`
- 当前状态：后端 `internal/service/...` 和部分前端测试仍需收口。

## Risks
- 当前仓库原本就是脏工作区，最终还需要处理 stash 回放冲突。
- `pnpm lint` 和 `pnpm test:run` 引入了部分环境副作用，需要甄别后保留最小必要改动。
- upstream 主线已引入较大范围功能变更，最终回放本地改动时可能出现第二轮冲突。

## Next Steps
- 回退无关改动，定位剩余失败测试。
- 完成临时集成分支验证后，将结果合回本地 `main`。
- `git stash pop` 回放本地改动并做第二轮验证。
