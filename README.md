# socks5-simulation
socks5 server with network simulation, like [clumsy](https://github.com/jagt/clumsy) but for Socks5 proxy.

支持网络模拟的 Socks5 服务器, 像 [clumsy](https://github.com/jagt/clumsy), 但是用于 Socks5 代理。

---

## 一、项目搭建记录
```bash
# 项目结构搭建
mkdir socks5-simulation
cd socks5-simulation/
go mod init socks5-simulation
# Go 包安装
go get github.com/urfave/cli/v3@latest
go get github.com/txthinking/socks5
go get github.com/txthinking/runnergroup
go get github.com/yanlinLiu0424/godivert
touch main.go
```

## 二、自行尝试
>  国内加速克隆:
> `git clone https://mirror.ghproxy.com/https://github.com/gsw945/socks5-simulation.git`

### (一)、克隆项目
```bash
git clone https://github.com/gsw945/socks5-simulation.git
cd socks5-simulation
```

### (二)、设置Go包源加速和安装Go包
```bash
# 启用 Go Mod
set GO111MODULE=on
# 使用国内源
set GOPROXY=https://proxy.golang.com.cn,direct
# 安装Go包
go mod tidy -v
```
### (三)、Windows 编译&运行
```bash
# 编译
go build -o socks5-simulation.exe main.go
# 运行
.\socks5-simulation.exe --listen 0.0.0.0:2060
```

### (四)、Linux 编译&运行
```bash
# 编译
go build -o socks5-simulation main.go
# 运行
./socks5-simulation --listen 0.0.0.0:2060
```

### 参考
- github.com/txthinking/socks5
- https://github.com/basil00/Divert
- https://github.com/yanlinLiu0424/godivert
- [Windivert ProcessId at NETWORK layer](https://stackoverflow.com/questions/58449491/windivert-processid-at-network-layer)
- [name of process in filter #169](https://github.com/basil00/Divert/issues/169#issuecomment-478176492)
- https://github.com/basil00/TorWall/blob/master/redirect.c

### 计划参考
- [proxychains4 with Go lang #199](https://github.com/rofl0r/proxychains-ng/issues/199)
- https://github.com/hmgle/graftcp
- https://github.com/mzz2017/gg
- [天朝局域网内go get的正确姿势](https://blog.scnace.me/post/%E4%B8%BAgo-get%E6%8A%A4%E8%88%AA-/)
- [WinDivert 2 Rust Wrapper](https://github.com/Rubensei/windivert-rust)
- https://github.com/dotpcap/sharppcap
- https://github.com/dotpcap/packetnet
- https://github.com/xljiulang/WindivertDotnet
- https://github.com/james-barrow/golang-ipc
