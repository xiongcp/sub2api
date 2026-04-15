## Goal
- 为现有白牌 HTML/CSS 能力补充一套默认的 Google 风极简配置，并让新部署、缺失 key 补齐、后台“恢复默认”三条链路保持一致。

## Todo
- 无。

## Doing
- 无。

## Done
- 后端新增 Google 风默认白牌模板常量，覆盖 `home_content`、`custom_css`、`login_extra_html`、`register_extra_html`、`global_footer_html`，并保持 `payment_footer_html` 默认为空。
- `InitializeDefaultSettings` 改为“只补齐缺失 key，不覆盖已有值”，显式空字符串视为管理员有意清空，不自动回填。
- `ProvideSettingService` 启动时会执行默认 setting 补齐，确保默认白牌模板真正落库，而不是只存在于代码里。
- 后台“恢复默认”按钮改为恢复这套默认模板，并把首页 `home_content` 一起纳入恢复范围。
- 认证布局补了稳定 class，默认 CSS 只对首页和认证布局生效，避免粗暴影响后台页面。
- 新增后端初始化测试和前端默认模板测试。

## Validation
- `git diff --check`
- `gofmt -w backend/internal/service/branding_defaults.go backend/internal/service/setting_service.go backend/internal/service/wire.go backend/cmd/server/wire_gen.go backend/internal/service/setting_service_initialize_test.go`
- `cd backend && go test ./internal/service/... ./internal/server/... ./cmd/server/...`
- `pnpm --dir frontend test:run`
- `pnpm --dir frontend typecheck`
- `pnpm --dir frontend lint:check`

## Risks
- 默认首页采用完整 `home_content` 整页替换，后续如果首页产品结构再改，这套默认模板需要同步维护。
- 这套默认风格接近 Google 的极简视觉，但仍需避免后续运营内容直接使用 Google 商标、Logo 或官方品牌素材。
- `custom_css` 仍是全局注入，虽然本次已尽量收敛到首页/认证布局选择器，但后续自定义样式仍可能影响未验证页面。

## Next Steps
- 如果要继续做品牌模板能力，下一步更适合抽成后台“一键套用模板”，而不是继续在前端手写多个默认字符串副本。
