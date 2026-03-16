# IotManager - IoT设备管理应用

## 项目概述

IotManager是一个专为飞牛平台设计的IoT设备管理应用，提供设备监控、配置和管理功能。

## 项目信息

- **应用名称**: IotManager
- **版本**: 1.0.0
- **描述**: IoT设备管理应用
- **架构**: x86_64
- **维护者**: your-name
- **分发者**: your-name

## 项目结构

```
IotManager/
├── app/
│   ├── server/       # 后端服务
│   ├── ui/           # 前端界面
│   │   ├── images/   # 图标和图片资源
│   │   └── config    # UI配置文件
│   └── www/          # Web界面文件
├── cmd/              # 命令脚本
│   ├── config_callback   # 配置回调脚本
│   ├── config_init       # 配置初始化脚本
│   ├── install_callback  # 安装回调脚本
│   ├── install_init      # 安装初始化脚本
│   ├── main              # 主启动脚本
│   ├── uninstall_callback # 卸载回调脚本
│   ├── uninstall_init    # 卸载初始化脚本
│   ├── upgrade_callback  # 升级回调脚本
│   └── upgrade_init      # 升级初始化脚本
├── config/
│   ├── privilege     # 权限配置
│   └── resource      # 资源配置
├── ICON.PNG          # 应用图标
├── ICON_256.PNG      # 应用图标(256x256)
├── manifest          # 应用清单文件
└── README.md         # 项目说明文档
```

## 功能特性

- **设备管理**: 监控和管理IoT设备
- **Web界面**: 提供直观的用户界面
- **完整生命周期管理**: 支持安装、配置、升级和卸载流程
- **数据共享**: 支持应用间数据共享
- **进程管理**: 包含完整的进程启动、停止和状态检查功能
- **日志记录**: 详细的操作日志
- **网络监控**: 定时ping网关，监控网络连通性
- **自动检测**: 每30秒自动检测网关可达性
- **详细日志**: 记录ping结果和响应时间

## 系统要求

- 飞牛平台
- x86_64架构

## 安装说明

1. 将应用包上传到飞牛平台
2. 通过飞牛平台的应用管理界面安装
3. 安装过程会自动执行 `install_init` 和 `install_callback` 脚本

## 配置说明

1. 应用安装后，通过飞牛平台的配置界面进行配置
2. 配置过程会执行 `config_init` 和 `config_callback` 脚本
3. UI配置文件位于 `app/ui/config`，可根据需要修改

## 启动和停止

- 应用启动时会执行 `cmd/main start` 命令
- 应用停止时会执行 `cmd/main stop` 命令
- 查看应用状态可执行 `cmd/main status` 命令

## 升级和卸载

- **升级**: 通过飞牛平台上传新版本，会自动执行 `upgrade_init` 和 `upgrade_callback` 脚本
- **卸载**: 通过飞牛平台卸载应用，会自动执行 `uninstall_init` 和 `uninstall_callback` 脚本

## 开发指南

### 目录说明

- **app/server/**: 后端服务代码，负责设备管理逻辑
- **app/ui/**: 前端界面代码，提供用户交互界面
- **app/www/**: Web界面文件，包含HTML、CSS和JavaScript
- **cmd/**: 命令脚本，处理应用生命周期管理
- **config/**: 配置文件，包含权限和资源配置

### 自定义应用

1. 修改 `manifest` 文件中的应用信息
2. 更新 `app/ui/config` 中的UI配置
3. 修改 `config/resource` 中的资源配置
4. 在 `cmd/main` 中设置启动命令
5. 实现 `app/server` 中的后端逻辑
6. 开发 `app/ui` 和 `app/www` 中的前端界面

## 注意事项

- 确保 `cmd/main` 脚本中的 `CMD` 变量设置正确，指向实际的启动命令
- 应用运行时会在 `${TRIM_PKGVAR}` 目录下生成日志文件和PID文件
- 如需修改应用图标，替换 `ICON.PNG` 和 `ICON_256.PNG` 文件

## 故障排查

- 查看应用日志: `${TRIM_PKGVAR}/info.log`
- 查看网络监控日志: `${TRIM_PKGVAR}/app/server/ping.log`
- 检查应用状态: `cmd/main status`
- 手动启动/停止应用: `cmd/main start` / `cmd/main stop`

## 网络监控功能

### 功能说明
- 应用启动后会自动开始定时ping网关
- 默认每30秒ping一次默认网关
- 如果无法获取默认网关，会使用192.168.1.1作为默认值
- 记录ping结果和响应时间到ping.log文件

### 日志格式
```
2026-03-17 12:00:00 - 启动网关ping服务...
2026-03-17 12:00:00 - 默认网关: 192.168.1.1
2026-03-17 12:00:00 - ✅ 网关 192.168.1.1 可达 (耗时: 12ms)
2026-03-17 12:00:30 - ✅ 网关 192.168.1.1 可达 (耗时: 15ms)
```

## 许可证

本项目为飞牛应用模板项目，可根据实际需求进行修改和扩展。