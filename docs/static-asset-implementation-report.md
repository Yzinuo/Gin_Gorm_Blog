# 静态资源治理实施报告

日期：2026-09-19
基线：`e49867fe67d9f0a3c965a1d70ce06b104e21aa51`

## 已交付

- Go 端新增 provider-neutral ObjectStore、R2 S3 适配器、严格上传校验、受管资源版本表、公开 manifest、ETag、发布/回滚事务、最近三版保护和延迟 GC。
- Admin 新增首页成对图片、Resume GLB、8 张贴纸的上传、校验、图片/Three.js GLB 预览、预览确认、发布与回滚入口。
- 前台首页先显示静态 Before；桌面在可见且 idle 后增强，移动端、Save-Data、慢网和 reduced-motion 默认不请求 After。About 的 GLB 同样按设备/网络条件加载并支持 Meshopt。
- 通用上传与飞书导入统一走 ObjectStore；上传文件使用内容寻址 key，保留 source，静态图片生成按用途限宽的 WebP delivery。
- 提供固定依赖版本的图片/GLB 优化脚本、迁移/GC CLI、数据库 migration、Nginx/Docker 配置和部署/回滚 runbook。

## 性能结果

具体逐文件字节数见 `docs/static-asset-performance.md`。首页两张桌面投放图合计 134,650 B；GLB 从 28,128,092 B 降至 3,858,836 B；favicon 从 1,485,199 B 降至 2,725 B。生产镜像不再复制两张原始 PNG、原始 GLB 和原始 favicon，按静态文件净载荷计算减少约 27.5 MiB。

## 验证结果

| 检查 | 结果 |
|---|---|
| `go test ./...` | 通过；无 Redis 时两项既有集成测试显式 skip |
| `front` ESLint | 通过；仅 UnoCSS 图标解析提示 |
| `front` production build | 通过 |
| `admin` ESLint | 通过；保留一项既有 `console` warning |
| `admin` production build | 通过 |
| Compose config | 通过 |
| 干净环境 Web Docker build | 通过；97,008,799 B，原始 PNG/GLB/favicon 已排除 |
| Server Docker build | 通过；92,723,078 B |
| R2/CDN/Lighthouse | 待所有者轮换凭据、激活域名并部署后执行 |

## 迁移失败清单

开发环境 dry-run 报告为 `docs/asset-migration-dry-run.json`：共盘点 24 个仓库资源、44,154,464 B、0 个读取失败。该报告不包含生产数据库或旧远程 URL，因此不能代替生产迁移报告。生产 apply 会把每个不可读、超限或上传失败的引用以 `status=failed` 和错误原因写入 JSON，且保持原数据库引用不变。执行命令与报告保存位置见 `docs/asset-migration.md`。

## Definition of Done

- [x] D01–D18 的工程实现已落入代码、migration、脚本与 runbook。
- [x] Admin 可动态管理首页图片、GLB 和 8 个贴纸，不需要重新构建前端。
- [x] 对象 URL 版本化且不可变；发布/回滚只切换数据库指针。
- [x] 当前版本与最近两个历史发布版本受保护。
- [x] 首页移动端不自动请求 After；About 移动端不自动请求 GLB。
- [x] 仓库投放 GLB 小于 6 MiB，结构契约有自动测试；真实视觉交互列入部署验收。
- [x] 通用上传及飞书导入默认走 ObjectStore，响应 envelope 保持兼容。
- [ ] 已使用生产数据库生成迁移 dry-run 报告并人工核对失败项。
- [x] 仓库 diff 未包含 R2 token、Access Key 或 Secret。
- [x] Go 测试、两个前端 lint/build、Compose 配置验证通过。
- [x] 部署和回滚 runbook 已更新。
- [ ] `assets.heliar.top` 已激活且完成缓存头、CORS、三次 Lighthouse 中位数及 24 小时观察（2026-09-19 的只读探测仍在 TLS 握手阶段失败）。

未勾选项都是 PRD 第 15 节允许的生产账号、域名或部署前置操作，不代表可在代码中伪造完成。

## Ambiguity Report

```text
Ambiguity Report:
  Goals:        0.0   ✓ clear
  Acceptance:   0.25  ✓ production checks remain
  Boundaries:   0.0   ✓ clear
  Alternatives: 0.25  ✓ main tradeoffs resolved
  Assumptions:  0.25  ✓ external state identified
  ──────────────────────────────
  Aggregate:    0.15  ✓ below threshold (0.2 spec)

Push lightly on: production migration evidence, CDN/RUM acceptance.
```
