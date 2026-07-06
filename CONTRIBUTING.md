# 贡献者指南

感谢你对 **移动云盘管理系统** 的关注与贡献！这份指南将帮助你了解如何参与项目开发。

## 行为准则

- 尊重所有参与者，保持友好和包容
- 接受建设性批评，专注于项目改进
- 拒绝骚扰、歧视或不当行为
- 遵守开源协议和社区规范

## 如何贡献

### 报告 Bug

- 使用 [Bug Report](/.github/ISSUE_TEMPLATE/bug_report.md) 模板创建 Issue
- 提供清晰的标题和复现步骤
- 附上相关日志、截图和版本信息

### 提出新功能

- 使用 [Feature Request](/.github/ISSUE_TEMPLATE/feature_request.md) 模板
- 描述功能的价值和预期行为
- 如有原型/草图请一并附上

### 提交代码

1. Fork 本仓库
2. 创建功能分支：`git checkout -b feature/amazing-feature`
3. 提交更改：`git commit -m 'feat: add amazing feature'`
4. 推送分支：`git push origin feature/amazing-feature`
5. 开启 Pull Request

## 开发规范

### 后端（Go）

- 遵循 Clean Architecture：`domain/application/infrastructure/interfaces`
- 使用 `gofmt` 格式化代码
- 提交信息格式：`type(scope): description`
- 类型：`feat`、`fix`、`docs`、`style`、`refactor`、`test`、`chore`

### 前端（Vue3 + TypeScript）

- 使用 ESLint + Prettier 格式化
- 组件名使用 PascalCase，文件名使用 kebab-case
- Props/Events 需定义类型
- 提交信息格式同上

### 提交信息规范

| 类型 | 说明 |
|------|------|
| `feat` | 新功能 |
| `fix` | Bug 修复 |
| `docs` | 文档变更 |
| `style` | 代码格式（不影响功能） |
| `refactor` | 重构 |
| `test` | 测试相关 |
| `chore` | 构建/工具链 |

示例：`feat(auth): add JWT refresh token support`

## 开发环境

### 环境要求

- Go 1.25+
- Node.js 20+
- Docker 20.10+
- Docker Compose v2+

### 本地启动

```bash
# 后端
cd backend
go mod download
go build ./cmd/api
go build ./cmd/worker

# 前端
cd frontend
npm install
npm run dev
```

### 运行测试

```bash
# 后端测试
cd backend && go test ./...

# 前端测试
cd frontend && npm run lint && npm run typecheck && npm run test:unit

# E2E 测试
cd frontend && npm run e2e
```

## PR 检查清单

- [ ] 代码符合项目规范
- [ ] 添加/更新了相关测试
- [ ] 更新了文档（如需要）
- [ ] CI 检查全部通过
- [ ] 自测覆盖主要功能

## 发布流程

1. 更新 `CHANGELOG.md`
2. 打版本 tag：`git tag v1.0.0 && git push --tags`
3. GitHub Actions 会自动创建 Release

## 许可证

本项目采用 MIT 许可证。提交代码即表示你同意将代码以相同许可证发布。

## 联系方式

- Issues: https://github.com/xiaohai09/mobile-cloud-disk-management/issues
- Discussions: https://github.com/xiaohai09/mobile-cloud-disk-management/discussions
