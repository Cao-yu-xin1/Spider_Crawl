package service

//111
import (
	"context"
	"fmt"
)

// ServiceType 服务类型
type ServiceType string

const (
	TypeHTTP ServiceType = "http"
	TypeGRPC ServiceType = "grpc"
)

// Service 服务接口定义
type Service interface {
	// Name 返回服务名称
	Name() string

	// Type 返回服务类型
	Type() ServiceType

	// Port 返回服务端口
	Port() int

	// Start 启动服务
	Start(ctx context.Context) error

	// Stop 停止服务
	Stop(ctx context.Context) error

	// IsRunning 检查服务是否运行中
	IsRunning() bool

	// HealthCheck 健康检查
	HealthCheck() error
}

// BaseService 基础服务实现，可嵌入其他服务实现
type BaseService struct {
	name       string
	svcType    ServiceType
	port       int
	running    bool
	stopFunc   func(context.Context) error
	startFunc  func(context.Context) error
	healthFunc func() error
}

// NewBaseService 创建基础服务
func NewBaseService(name string, svcType ServiceType, port int) *BaseService {
	return &BaseService{
		name:    name,
		svcType: svcType,
		port:    port,
		running: false,
	}
}

// Name 返回服务名称
func (s *BaseService) Name() string {
	return s.name
}

// Type 返回服务类型
func (s *BaseService) Type() ServiceType {
	return s.svcType
}

// Port 返回服务端口
func (s *BaseService) Port() int {
	return s.port
}

// IsRunning 检查服务是否运行中
func (s *BaseService) IsRunning() bool {
	return s.running
}

// SetRunning 设置运行状态
func (s *BaseService) SetRunning(running bool) {
	s.running = running
}

// SetStartFunc 设置启动函数
func (s *BaseService) SetStartFunc(fn func(context.Context) error) {
	s.startFunc = fn
}

// SetStopFunc 设置停止函数
func (s *BaseService) SetStopFunc(fn func(context.Context) error) {
	s.stopFunc = fn
}

// SetHealthFunc 设置健康检查函数
func (s *BaseService) SetHealthFunc(fn func() error) {
	s.healthFunc = fn
}

// Start 启动服务
func (s *BaseService) Start(ctx context.Context) error {
	if s.running {
		return fmt.Errorf("service %s is already running", s.name)
	}

	if s.startFunc != nil {
		if err := s.startFunc(ctx); err != nil {
			return fmt.Errorf("failed to start service %s: %w", s.name, err)
		}
	}

	s.running = true
	return nil
}

// Stop 停止服务
func (s *BaseService) Stop(ctx context.Context) error {
	if !s.running {
		return nil
	}

	if s.stopFunc != nil {
		if err := s.stopFunc(ctx); err != nil {
			return fmt.Errorf("failed to stop service %s: %w", s.name, err)
		}
	}

	s.running = false
	return nil
}

// HealthCheck 健康检查
func (s *BaseService) HealthCheck() error {
	if !s.running {
		return fmt.Errorf("service %s is not running", s.name)
	}

	if s.healthFunc != nil {
		return s.healthFunc()
	}

	return nil
}
