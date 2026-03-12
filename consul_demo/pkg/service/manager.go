package service

//111
import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Registry 服务注册表
type Registry struct {
	mu       sync.RWMutex
	services map[string]Service
}

// NewRegistry 创建服务注册表
func NewRegistry() *Registry {
	return &Registry{
		services: make(map[string]Service),
	}
}

// Register 注册服务
func (r *Registry) Register(svc Service) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := svc.Name()
	if _, exists := r.services[name]; exists {
		return fmt.Errorf("service %s already registered", name)
	}

	r.services[name] = svc
	return nil
}

// Unregister 注销服务
func (r *Registry) Unregister(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.services[name]; !exists {
		return fmt.Errorf("service %s not found", name)
	}

	delete(r.services, name)
	return nil
}

// Get 获取服务
func (r *Registry) Get(name string) (Service, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	svc, exists := r.services[name]
	if !exists {
		return nil, fmt.Errorf("service %s not found", name)
	}

	return svc, nil
}

// List 列出所有服务
func (r *Registry) List() []Service {
	r.mu.RLock()
	defer r.mu.RUnlock()

	services := make([]Service, 0, len(r.services))
	for _, svc := range r.services {
		services = append(services, svc)
	}

	return services
}

// ListByType 按类型列出服务
func (r *Registry) ListByType(svcType ServiceType) []Service {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var services []Service
	for _, svc := range r.services {
		if svc.Type() == svcType {
			services = append(services, svc)
		}
	}

	return services
}

// Exists 检查服务是否存在
func (r *Registry) Exists(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.services[name]
	return exists
}

// Count 返回注册服务数量
func (r *Registry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.services)
}

// ServiceManager 服务管理器
type ServiceManager struct {
	registry *Registry
	mu       sync.RWMutex
	running  map[string]bool
}

// NewServiceManager 创建服务管理器
func NewServiceManager() *ServiceManager {
	return &ServiceManager{
		registry: NewRegistry(),
		running:  make(map[string]bool),
	}
}

// Registry 返回服务注册表
func (m *ServiceManager) Registry() *Registry {
	return m.registry
}

// Register 注册服务
func (m *ServiceManager) Register(svc Service) error {
	return m.registry.Register(svc)
}

// Start 启动指定服务
func (m *ServiceManager) Start(ctx context.Context, name string) error {
	svc, err := m.registry.Get(name)
	if err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.running[name] {
		return fmt.Errorf("service %s is already running", name)
	}

	// 检查端口冲突
	if err := m.checkPortConflict(svc); err != nil {
		return err
	}

	if err := svc.Start(ctx); err != nil {
		return err
	}

	m.running[name] = true
	return nil
}

// StartAll 启动所有服务
func (m *ServiceManager) StartAll(ctx context.Context) error {
	services := m.registry.List()

	for _, svc := range services {
		if err := m.Start(ctx, svc.Name()); err != nil {
			return fmt.Errorf("failed to start service %s: %w", svc.Name(), err)
		}
	}

	return nil
}

// Stop 停止指定服务
func (m *ServiceManager) Stop(ctx context.Context, name string) error {
	svc, err := m.registry.Get(name)
	if err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.running[name] {
		return nil
	}

	if err := svc.Stop(ctx); err != nil {
		return err
	}

	delete(m.running, name)
	return nil
}

// StopAll 停止所有服务
func (m *ServiceManager) StopAll(ctx context.Context) error {
	services := m.registry.List()

	var errs []error
	for _, svc := range services {
		if err := m.Stop(ctx, svc.Name()); err != nil {
			errs = append(errs, fmt.Errorf("failed to stop service %s: %w", svc.Name(), err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors during shutdown: %v", errs)
	}

	return nil
}

// Restart 重启服务
func (m *ServiceManager) Restart(ctx context.Context, name string) error {
	if err := m.Stop(ctx, name); err != nil {
		return err
	}

	return m.Start(ctx, name)
}

// Status 获取服务状态
func (m *ServiceManager) Status(name string) (ServiceStatus, error) {
	svc, err := m.registry.Get(name)
	if err != nil {
		return ServiceStatus{}, err
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	status := ServiceStatus{
		Name:    svc.Name(),
		Type:    svc.Type(),
		Port:    svc.Port(),
		Running: m.running[name],
	}

	if status.Running {
		status.Healthy = svc.HealthCheck() == nil
	}

	return status, nil
}

// StatusAll 获取所有服务状态
func (m *ServiceManager) StatusAll() []ServiceStatus {
	services := m.registry.List()
	statuses := make([]ServiceStatus, 0, len(services))

	for _, svc := range services {
		status, _ := m.Status(svc.Name())
		statuses = append(statuses, status)
	}

	return statuses
}

// IsRunning 检查服务是否运行中
func (m *ServiceManager) IsRunning(name string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.running[name]
}

// GracefulShutdown 优雅关闭所有服务
func (m *ServiceManager) GracefulShutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return m.StopAll(ctx)
}

// checkPortConflict 检查端口冲突
func (m *ServiceManager) checkPortConflict(newSvc Service) error {
	for name, isRunning := range m.running {
		if !isRunning {
			continue
		}

		svc, err := m.registry.Get(name)
		if err != nil {
			continue
		}

		if svc.Port() == newSvc.Port() {
			return fmt.Errorf("port %d is already in use by service %s", newSvc.Port(), name)
		}
	}

	return nil
}

// ServiceStatus 服务状态
type ServiceStatus struct {
	Name    string
	Type    ServiceType
	Port    int
	Running bool
	Healthy bool
}
