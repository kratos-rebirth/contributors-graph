# 贡献者图片生成器

## 这是什么

一个用于从 GitHub 拉取贡献者列表及对应的头像，并生成一张聚合图片的工具。

配套一个定期执行的 Actions 流水线，可以将生成的图片托管在 GitHub Pages 服务上。

## 如何使用

它有两个环境变量：

1. `REPO` 指定需要对什么仓库生成贡献者列表，例如默认的是 `Candinya/Kratos-Rebirth` 。
2. `TOKEN` 指定调用 GitHub API 时使用的访问令牌，如果是针对开源项目、在 GitHub Actions 执行的话，可以像这里的 CI 配置一样设置成 GitHub 自动临时分配的 `${{ github.token }}` 。

如果您觉得定期执行的策略不合适，您可以自行设定更为合适的策略。

## 开源授权

[MIT](./LICENSE)
