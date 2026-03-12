package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"consul_demo/pkg/config"
	"consul_demo/pkg/service"
)

var (
	flagAll      bool
	flagServices string
	flagConfig   string
	flagStatus   bool
	flagPort     int
)

func init() {
	flag.BoolVar(&flagAll, "all", false, "启动所有服务")
	flag.StringVar(&flagServices, "services", "", "指定启动的服务 (逗号分隔)")
	flag.StringVar(&flagConfig, "config", "config.yaml", "配置文件路径")
	flag.BoolVar(&flagStatus, "status", false, "查看服务状态")
	flag.IntVar(&flagPort, "port", 0, "覆盖服务端口")
}

func main() {
	flag.Parse()

	// 加载配置
	cfg, err := config.Load(flagConfig)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 从环境变量加载
	cfg.LoadFromEnv()

	// 验证配置
	if err := cfg.Validate(); err != nil {
		log.Fatalf("配置验证失败: %v", err)
	}

	// 创建服务管理器
	manager := service.NewServiceManager()

	// 注册服务
	for name, svcCfg := range cfg.Services {
		// 覆盖端口
		if flagPort > 0 {
			svcCfg.Port = flagPort
		}

		svc := NewSimpleService(name, svcCfg)
		if err := manager.Register(svc); err != nil {
			log.Printf("注册服务 %s 失败: %v", name, err)
		}
	}

	// 查看状态模式
	if flagStatus {
		printStatus(manager)
		return
	}

	// 确定要启动的服务
	var servicesToStart []string
	if flagAll {
		servicesToStart = cfg.ListServices()
	} else if flagServices != "" {
		servicesToStart = strings.Split(flagServices, ",")
		for i, s := range servicesToStart {
			servicesToStart[i] = strings.TrimSpace(s)
		}
	} else {
		// 默认启动所有启用的服务
		servicesToStart = cfg.ListServices()
	}

	if len(servicesToStart) == 0 {
		log.Println("没有服务需要启动")
		printStatus(manager)
		return
	}

	// 创建上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 启动服务
	log.Println("==========================================")
	log.Println("启动服务...")
	log.Println("==========================================")

	for _, name := range servicesToStart {
		if !manager.Registry().Exists(name) {
			log.Printf("服务 %s 不存在，跳过", name)
			continue
		}

		log.Printf("启动服务: %s", name)
		if err := manager.Start(ctx, name); err != nil {
			log.Printf("启动服务 %s 失败: %v", name, err)
		} else {
			log.Printf("服务 %s 启动成功", name)
		}
	}

	// 打印状态
	log.Println("")
	printStatus(manager)
	log.Println("")
	log.Println("==========================================")
	log.Println("所有服务已启动，按 Ctrl+C 停止")
	log.Println("==========================================")

	// 等待信号
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Println("")
	log.Println("正在停止所有服务...")

	// 优雅关闭
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := manager.GracefulShutdown(10 * time.Second); err != nil {
		log.Printf("关闭服务时出错: %v", err)
	}

	<-shutdownCtx.Done()
	log.Println("所有服务已停止")
}

// printStatus 打印服务状态
func printStatus(manager *service.ServiceManager) {
	fmt.Println("\n服务状态:")
	fmt.Println("┌─────────────┬──────┬──────────┬─────────┐")
	fmt.Println("│ 服务名称    │ 端口 │ 类型     │ 状态    │")
	fmt.Println("├─────────────┼──────┼──────────┼─────────┤")

	for _, status := range manager.StatusAll() {
		running := "停止"
		if status.Running {
			running = "运行中"
		}

		fmt.Printf("│ %-11s │ %4d │ %-8s │ %-7s │\n",
			status.Name,
			status.Port,
			status.Type,
			running,
		)
	}

	fmt.Println("└─────────────┴──────┴──────────┴─────────┘")
}

// SimpleService 简单服务实现
type SimpleService struct {
	*service.BaseService
	cfg config.ServiceConfig
}

// NewSimpleService 创建简单服务
func NewSimpleService(name string, cfg config.ServiceConfig) *SimpleService {
	svcType := service.TypeHTTP
	if name == "order_srv" {
		svcType = service.TypeGRPC
	}

	s := &SimpleService{
		BaseService: service.NewBaseService(name, svcType, cfg.Port),
		cfg:         cfg,
	}

	// 设置启动函数
	s.SetStartFunc(func(ctx context.Context) error {
		log.Printf("[%s] 服务启动中... (端口: %d)", name, cfg.Port)
		return nil
	})

	// 设置停止函数
	s.SetStopFunc(func(ctx context.Context) error {
		log.Printf("[%s] 服务停止中...", name)
		return nil
	})

	// 设置健康检查函数
	s.SetHealthFunc(func() error {
		// TODO: 实现实际的健康检查
		return nil
	})

	return s
}
