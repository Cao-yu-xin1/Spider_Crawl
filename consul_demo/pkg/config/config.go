package config

//111

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

// Config 总配置
type Config struct {
	Consul   ConsulConfig             `yaml:"consul"`
	Services map[string]ServiceConfig `yaml:"services"`
}

// ConsulConfig Consul 配置
type ConsulConfig struct {
	Address string `yaml:"address"`
}

// ServiceConfig 服务配置
type ServiceConfig struct {
	Enabled  bool                   `yaml:"enabled"`
	Port     int                    `yaml:"port"`
	Host     string                 `yaml:"host"`
	Database *DatabaseConfig        `yaml:"database,omitempty"`
	Consul   *ServiceConsulConfig   `yaml:"consul,omitempty"`
	Order    *OrderConfig           `yaml:"order,omitempty"`
	Extra    map[string]interface{} `yaml:"extra,omitempty"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

// ServiceConsulConfig 服务 Consul 配置
type ServiceConsulConfig struct {
	Address     string `yaml:"address"`
	ServiceName string `yaml:"service_name"`
	ServiceID   string `yaml:"service_id"`
}

// OrderConfig Order 服务配置
type OrderConfig struct {
	ServiceName string `yaml:"service_name"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Consul: ConsulConfig{
			Address: "115.190.154.22:8500",
		},
		Services: map[string]ServiceConfig{
			"h5_bff": {
				Enabled: true,
				Port:    8080,
				Host:    "127.0.0.1",
				Consul: &ServiceConsulConfig{
					Address:     "115.190.154.22:8500",
					ServiceName: "order-service",
				},
				Order: &OrderConfig{
					ServiceName: "order-service",
				},
			},
			"order_srv": {
				Enabled: true,
				Port:    50051,
				Host:    "127.0.0.1",
				Database: &DatabaseConfig{
					Host:     "115.190.154.22",
					Port:     3306,
					User:     "root",
					Password: "4ay1nkal3u8ed77y",
					DBName:   "zy0226",
				},
				Consul: &ServiceConsulConfig{
					Address:     "115.190.154.22:8500",
					ServiceName: "order-service",
					ServiceID:   "order-service-1",
				},
			},
		},
	}
}

// Load 从文件加载配置
func Load(path string) (*Config, error) {
	cfg := DefaultConfig()

	if path == "" {
		return cfg, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return cfg, nil
}

// LoadFromEnv 从环境变量加载配置
func (c *Config) LoadFromEnv() {
	// Consul 地址
	if addr := os.Getenv("CONSUL_ADDRESS"); addr != "" {
		c.Consul.Address = addr
	}

	// 遍历服务配置
	for name, svc := range c.Services {
		// 服务端口
		if port := os.Getenv(fmt.Sprintf("%s_PORT", upper(name))); port != "" {
			if p, err := strconv.Atoi(port); err == nil {
				svc.Port = p
			}
		}

		// 服务启用状态
		if enabled := os.Getenv(fmt.Sprintf("%s_ENABLED", upper(name))); enabled != "" {
			svc.Enabled = enabled == "true" || enabled == "1"
		}

		// 服务ID
		if serviceID := os.Getenv("SERVICE_ID"); serviceID != "" {
			if svc.Consul != nil {
				svc.Consul.ServiceID = serviceID
			}
		}

		c.Services[name] = svc
	}
}

// Validate 验证配置
func (c *Config) Validate() error {
	// 检查端口冲突
	ports := make(map[int]string)
	for name, svc := range c.Services {
		if !svc.Enabled {
			continue
		}

		if existing, exists := ports[svc.Port]; exists {
			return fmt.Errorf("port conflict: service %s and %s both use port %d", name, existing, svc.Port)
		}
		ports[svc.Port] = name
	}

	return nil
}

// GetService 获取服务配置
func (c *Config) GetService(name string) (ServiceConfig, error) {
	svc, exists := c.Services[name]
	if !exists {
		return ServiceConfig{}, fmt.Errorf("service %s not found in config", name)
	}
	return svc, nil
}

// ListServices 列出所有启用的服务
func (c *Config) ListServices() []string {
	services := make([]string, 0)
	for name, svc := range c.Services {
		if svc.Enabled {
			services = append(services, name)
		}
	}
	return services
}

// upper 转大写（用于环境变量）
func upper(s string) string {
	result := make([]byte, len(s))
	for i, c := range s {
		if c >= 'a' && c <= 'z' {
			result[i] = byte(c - 32)
		} else {
			result[i] = byte(c)
		}
	}
	return string(result)
}
