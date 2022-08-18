# study-go

## 项目目录

- Go 目录
  - /cmd 项目主干
  - /internal 私有应用程序和库代码
  - /pkg 外部应用程序可以使用的库代码
  - /vendor 应用程序依赖项
- 服务应用程序目录
  - /api 
- Web 应用程序目录
  - /web
- 通用应用目录
  - /configs 配置文件模板或默认配置
  - /init System init（systemd，upstart，sysv）和 process manager/supervisor（runit，supervisor）配置
  - /scripts 执行各种构建、安装、分析等操作的脚本
  - /build 打包和持续集成
  - /deployments IaaS、PaaS、系统和容器编排部署配置和模板
  - /test 额外的外部测试应用程序和测试数据
- 其他目录
  - /docs 设计和用户文档
  - /tools 支持工具
  - /examples 示例
  - /third_party 外部辅助工具
  - /githooks
  - /assets 与存储库一起使用的其他资产
  - /website 网站数据
