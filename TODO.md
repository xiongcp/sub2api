## Goal
- 修复 scheduler cache 裁剪 `Account.Extra` 后丢失 OpenAI WS 能力字段，导致 `/openai/v1/responses` WebSocket 账号选择误判为不可用的问题。

## Todo
- 无。

## Doing
- 无。

## Done
- 已确认问题存在：scheduler snapshot 只保留少量 `Extra` 字段，而 OpenAI WS 协议解析和 scheduler transport 过滤依赖 `openai_*_responses_websockets_v2_*`、`responses_websockets_v2_enabled`、`openai_ws_enabled`、`openai_ws_force_http`。
- 已确认影响路径：`/openai/v1/responses` WebSocket 入口固定要求 `OpenAIUpstreamTransportResponsesWebsocketV2`，快照账号若丢失上述字段，会在账号选择阶段被错误过滤并报 `openai.websocket_account_select_failed`。
- 已补齐 scheduler cache 的 OpenAI WS 传输判定白名单，保持瘦快照策略不变，只保留 scheduler 选择阶段真正依赖的 `Extra` 字段。
- 已新增缓存层回归测试，确认快照会保留 OpenAI WS 所需字段，同时继续裁掉无关大字段。
- 已新增 scheduler 回归测试，确认走 snapshot 账号列表时，带有 WSv2 元数据的账号仍可被选中，并最终返回 hydration 后的完整账号。

## Validation
- `git diff --check`
- `cd backend && go test ./internal/service/...`
- `cd backend && go test -tags=integration ./internal/repository/... -run TestSchedulerCacheSnapshotUsesSlimMetadataButKeepsFullAccount`

## Risks
- 这是缓存快照问题，代码修复后如果线上继续命中旧 `sched:*` 快照，症状可能暂时还在；需要依赖重启后的 startup rebuild 或显式重建快照。
- 本次保持瘦快照策略，只补 transport 判定必需字段；如果后续 scheduler 再依赖新的 `Extra` 键，需要同步维护白名单。

## Next Steps
- 如果上线方式不是重启进程，需要补一条运维说明：发布后清理或重建 Redis 中的 scheduler snapshot。
