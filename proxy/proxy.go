package proxy

import (
	"context"
	"io"
	"log"
	"net"
	"sync"

	"github.com/mikuta0407/tcpproxy/config"
)

type Proxy struct {
	// 設定ファイルの現在の状態を保持
	ActiveConfig *config.Config
	// 設定の更新時に使用するミューテックス
	ConfigMutex sync.RWMutex
	// 実行中のプロキシをキャンセルするための関数を保持
	CancelFuncs map[string]context.CancelFunc
	// cancelFuncsを保護するためのミューテックス
	CancelMutex sync.Mutex
}

func NewProxy(activeConfig *config.Config, cancelFuncs map[string]context.CancelFunc) *Proxy {
	return &Proxy{
		ActiveConfig: activeConfig,
		CancelFuncs:  cancelFuncs,
	}
}

func (p *Proxy) Start() {
	p.ConfigMutex.RLock()
	defer p.ConfigMutex.RUnlock()

	if p.ActiveConfig == nil {
		log.Fatal("no proxies setting")
	}

	log.Println("Starting proxies...")
	for _, proxy := range p.ActiveConfig.Proxies {
		// 各プロキシを個別のGoroutineで起動
		ctx, cancel := context.WithCancel(context.Background())

		p.CancelMutex.Lock()
		p.CancelFuncs[proxy.Name] = cancel
		p.CancelMutex.Unlock()

		go startProxy(ctx, proxy)
	}
}

// stopAllProxies は現在実行中のすべてのプロキシを停止します。
func (p *Proxy) StopAllProxies() {
	p.CancelMutex.Lock()
	defer p.CancelMutex.Unlock()

	log.Println("Stopping all active proxies...")
	for name, cancel := range p.CancelFuncs {
		cancel()
		delete(p.CancelFuncs, name)
	}
	// 少し待ってリスナーが閉じる時間を与える
	// time.Sleep(500 * time.Millisecond)
}

// startProxy は単一のポートフォワーディング処理を開始します。
func startProxy(ctx context.Context, p config.Proxy) {
	log.Printf("Starting proxy '%s': %s -> %s", p.Name, p.Source, p.Destination)

	// Sourceでリッスンを開始
	listener, err := net.Listen("tcp", p.Source)
	if err != nil {
		log.Printf("Error listening on %s for proxy '%s': %v", p.Source, p.Name, err)
		return
	}

	// コンテキストがキャンセルされたらリスナーを閉じる
	go func() {
		<-ctx.Done()
		log.Printf("Stopping proxy '%s' (%s)", p.Name, p.Source)
		listener.Close()
	}()

	// クライアントからの接続を待つループ
	for {
		conn, err := listener.Accept()
		if err != nil {
			// リスナーが閉じられた場合、ループを抜ける
			select {
			case <-ctx.Done():
				return
			default:
				log.Printf("Error accepting connection for proxy '%s': %v", p.Name, err)
			}
			continue
		}
		// 接続をハンドリングするGoroutineを起動
		go handleConnection(conn, p.Destination)
	}
}

// handleConnection はクライアントと接続先の間のデータ転送を処理します。
func handleConnection(sourceConn net.Conn, destinationAddr string) {
	defer sourceConn.Close()

	// Destinationに接続
	destConn, err := net.Dial("tcp", destinationAddr)
	if err != nil {
		log.Printf("Failed to connect to destination %s: %v", destinationAddr, err)
		return
	}
	defer destConn.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	// Source -> Destination のデータコピー
	go func() {
		defer wg.Done()
		if _, err := io.Copy(destConn, sourceConn); err != nil {
			// エラーロギングは冗長になる可能性があるので、必要に応じて有効化
			// log.Printf("Error copying from source to destination: %v", err)
		}
		// 反対側の接続も閉じる
		destConn.(*net.TCPConn).CloseWrite()
	}()

	// Destination -> Source のデータコピー
	go func() {
		defer wg.Done()
		if _, err := io.Copy(sourceConn, destConn); err != nil {
			// log.Printf("Error copying from destination to source: %v", err)
		}
		sourceConn.(*net.TCPConn).CloseWrite()
	}()

	wg.Wait()
}
