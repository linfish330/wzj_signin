# 微助教签到助手 (wzj_signin)

一个运行在本机或服务器上的微助教 (TeacherMate) 辅助工具。通过精致的 Web 页面管理微信 OpenID 监控，自动轮询并自动完成可签到列表。

本项目基于上游仓库演进（上游地址：[Azuka753/wzj_sign_public](https://github.com/Azuka753/wzj_sign_public)），并进行了深度的体验升级与功能重构。

---

## 🎨 视觉与交互重构 (Zhijian Design 风格)

本项目引入了全新的 **Zhijian Design** 视觉体系，为原本单调的签到监控工具带来了极具现代感和高级感的界面：
- **卡片式玻璃拟态**：界面元素采用毛玻璃与深色渐变微光，色彩搭配高雅和谐。
- **可视化轮询与延迟机制**：在设置页内嵌直观的三步骤流程图（扫描探测 ➔ 拟真延迟 ➔ 防刷提交），带有多态动效（雷达扫描、沙漏、呼吸灯）以及一键运行测试动画。
- **表格化 GPS 标签管理**：废弃了凌乱的徽章标签，改用响应式的精致数据表格，支持按行清晰查看经纬度并快速删除或管理。
- **全局微动效**：按钮悬停抖动、输入框外发光、运行状态呼吸灯，让操作更加富有生命力。

---

## ✨ 核心功能

- **自动化监控轮询**：周期性扫描 OpenID 待签到列表，轮询间隔可精准调节（1-3600秒）。
- **多种签到类型支持**：
  - **普通签到 & GPS 签到**：自动完成，并在 GPS 提交时对经纬度施加微小的物理随机抖动（Jittering），防止因定位千篇一律而被风控判定为防刷。
  - **二维码签到**：检测到二维码签到后，自动生成本地轮询刷新的二维码页面，并通过邮件通知您，扫码后自动提交。
- **人性化拟真延迟**：支持在探测到签到任务后，模拟等待设定秒数再提交，贴合真人手动操作逻辑。
- **数据本地化与安全隔离**：敏感配置（如邮箱 SMTP 密码、Redis 密码）自动隔离开并存储于本地的 `data/secrets.json`，不会进入版本控制。

---

## 🛠 运行前准备

- **Go 1.20+**（若直接在本机编译运行需要）
- **Redis**（用于存储轮询队列与用户状态数据）

---

## 🚀 快速开始

### 方式一：Windows 本机运行

1. **启动 Redis**
   - 默认读取 `localhost:6379`（可用环境变量 `REDIS_ADDRESS` 覆盖）。
2. **启动服务**
   ```powershell
   go run .
   ```
3. **访问网页**
   - 打开浏览器访问：[http://localhost:8080/home](http://localhost:8080/home)

> **端口冲突？**
> 如果默认的 `8080` 被占用，可使用如下命令更换端口启动：
> ```powershell
> $env:PORT=8081
> $env:SERVER_ADDRESS="http://localhost:8081"
> go run .
> ```
> 也可以直接使用提供的运行脚本：
> ```powershell
> powershell -ExecutionPolicy Bypass -File .\scripts\run-local.ps1 -Port 8081
> ```

---

### 方式二：macOS 本机运行 (推荐)

项目内提供了一个便捷的一键管理脚本。首先安装快捷指令：
```bash
chmod +x scripts/wzj_signin
ln -sf "$(pwd)/scripts/wzj_signin" ~/.local/bin/wzj_signin
```

之后，您可以在系统的任意目录下运行以下命令来管理服务：
- **启动服务**（自动检查 Redis、编译代码、后台运行并自动打开主页）：
  ```bash
  wzj_signin
  ```
- **查看状态**：`wzj_signin status`
- **查看日志**：`wzj_signin logs`
- **重启服务**：`wzj_signin restart`
- **停止服务**：`wzj_signin stop`

---

### 方式三：Docker / Docker Compose 部署

#### 1. 使用 docker-compose（推荐，自带独立 Redis 容器）
```bash
docker compose up -d --build
```
- **Web 访问地址**：[http://localhost:18080/home](http://localhost:18080/home)
- **数据持久化**：默认将 `./data` 目录挂载到容器中，方便持久化配置。

#### 2. 使用纯 Docker 容器（连接外部 Redis）
```bash
docker build -t wzj-sign:local .

docker run -d \
  --name wzj_sign \
  -p 18080:8080 \
  -e PORT=8080 \
  -e SERVER_ADDRESS=http://localhost:18080 \
  -e REDIS_ADDRESS=host.docker.internal:6379 \
  -v %CD%/data:/app/data \
  wzj-sign:local
```

---

## ⚙️ 配置与数据管理

### 1. 非敏感配置：`local/config.yml`
- 提供服务的基本参数设置（如轮询间隔、服务器地址等）。
- 若本地不存在此文件，启动时会自动以 `examples/config.example.yml` 为模板生成 `local/config.yml`（已加入 `.gitignore`，防止泄露个人配置）。

### 2. 敏感密码配置：`data/secrets.json`
- 存放 Redis 密码、邮件 SMTP 授权码等。
- 格式参考 `examples/secrets.example.json`。在设置页面更新密码时，会安全地直接写入该文件，无需暴露到配置文件中。

### 3. 数据持久化目录：`data/`
- `appconfig.json`：前端/服务端的非敏感运行机制配置（轮询间隔、延迟时间等）。
- `frontend_settings.json`：默认邮箱、已保存的 GPS 标签列表等。
- `secrets.json`：敏感密码。

---

## 📖 Web 页面指南

- **`/home` (运行概览)**：展示当前服务的运行状态、正在监控的 OpenID 数量、最近签到成功统计，并提供快速提交 OpenID 的快捷入口。
- **`/submit` (提交监控)**：用于粘贴 OpenID 或包含 OpenID 的链接，可选择保存的 GPS 标签进行绑定提交。
- **`/history` (监控中心)**：展示当前正在监控的用户列表，支持一键“开始/暂停轮询”、“查看日志”、“重新打开二维码”。
- **`/settings` (控制中心)**：
  - **标签管理**：以精致表格形式添加、删除、检索自定义 GPS 定位。
  - **邮件提醒**：配置签到通知收件箱与 SMTP 发信设置。
  - **轮询与延迟控制**：直观配置轮询扫描频率与拟真提交等待时间，并可通过运行动画观察流程。

---

## ⚠️ 免责声明

1. 本项目为开源辅助工具，与“微助教 / TeacherMate”官方及其母公司无任何关联。
2. 请在遵守学校、班级及相关平台规则的前提下，合理、合规地使用本工具。因非正常使用导致的任何后果由使用者自行承担。

---

## 🙏 致谢

- 感谢上游项目的开源实现：[Azuka753/wzj_sign_public](https://github.com/Azuka753/wzj_sign_public)。
