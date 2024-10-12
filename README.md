# tcpproxy

Port fowarder written by Golang

## ja

### これはなに?
- Go製のTCP用ポートフォワーダーです。
- iptablesやnginx、HAProxy、VMの設定でポートフォワードするのがめんどくさかったので、設定ファイルを読んでその通りにポートフォワードしてくれるシングルバイナリを作りました。

### 使い方
#### インストール
##### go install
```bash
go install github.com/mikuta0407/tcpproxy@latest
```
##### バイナリダウンロード
releasesページからダウンロードしてください。(linux/amd64, linux/arm64のみ準備しています)

#### コンフィグファイル作成
- デフォルトでは、`~/.config/tcpproxy/tcpproxy.yml`を読みに行きます。
  - Windowsでも`%userprofile%/tcpproxy/tcpproxy.yml`のはず
- `-c`オプションでymlファイルを渡してもOKです。
- config内に記載したフォワード設定はすべて並列に動作します。

例
```yaml
proxies:
  - name: SSH to VM1
    # KVMのNAT下にいるVMに外から2022でSSH接続できるようにする
    source: :2022
    destination: 192.168.122.2:22
  - name: Minecraft
    # VPS等外足を持っているサーバーでこれを動かし、VPNで繋いだ宅内Minecraftサーバーを外に出す
    source: :25565
    destination: 192.168.0.2:25565
```

### おすすめの使い方
- systemdのデーモンとして動作させておくと良いです。

### 既知の不具合
- 宛先に到達できない場合(no route to host)、クラッシュします。
  - systemdにまかせておけばある程度は自動復旧で助けられるかも……?

## en

### What is this?
- This is a TCP port forwarder written in Go.
- I created a single binary that reads a configuration file and forwards ports accordingly because setting up port forwarding with iptables, nginx, HAProxy, or VM settings was cumbersome.

### Usage
#### Installation
##### go install
```bash
go install github.com/mikuta0407/tcpproxy@latest
```
##### Binary download
Please download from the releases page. (Only linux/amd64 and linux/arm64 are prepared)

#### Creating a config file
- By default, it reads `~/.config/tcpproxy/tcpproxy.yml`.
  - On Windows, it should be `%userprofile%/tcpproxy/tcpproxy.yml`
- You can also pass the yml file with the `-c` option.
- All forwarding settings listed in the config will operate in parallel.

Example
```yaml
proxies:
  - name: SSH to VM1
    # Allow SSH connection to a VM under KVM NAT from outside on port 2022
    source: :2022
    destination: 192.168.122.2:22
  - name: Minecraft
    # Run this on a server with an external interface, and expose a home Minecraft server connected via VPN to the outside
    source: :25565
    destination: 192.168.0.2:25565
```

### Recommended usage
- It is recommended to run it as a systemd daemon.

### Known issues
- It crashes if the destination is unreachable (no route to host).
  - Leaving it to systemd might help with some automatic recovery...?
