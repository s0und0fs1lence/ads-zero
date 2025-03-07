package configuration

import (
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
)

type Provider interface {
	ConfigFileUser() string
	Get(key string) interface{}
	GetBool(key string) bool
	GetDuration(key string) time.Duration
	GetFloat64(key string) float64
	GetInt(key string) int
	GetSizeInBytes(key string) uint
	GetString(key string) string
	GetStringSlice(key string) []string
	GetTime(key string) time.Time
	InConfig(key string) bool
	IsSet(key string) bool
	SetConfigFile(in string)
	ReadInConfig() error
	Set(key string, value interface{})
}

type providerImpl struct {
	v  *viper.Viper
	mx sync.RWMutex
}

// ConfigFileUser implements Provider.
func (p *providerImpl) ConfigFileUser() string {
	p.mx.RLock()
	defer p.mx.RUnlock()
	return p.v.ConfigFileUsed()
}

// Get implements Provider.
func (p *providerImpl) Get(key string) interface{} {
	p.mx.RLock()
	defer p.mx.RUnlock()
	return p.v.Get(key)
}

// GetBool implements Provider.
func (p *providerImpl) GetBool(key string) bool {
	p.mx.RLock()
	defer p.mx.RUnlock()
	return p.v.GetBool(key)
}

// GetDuration implements Provider.
func (p *providerImpl) GetDuration(key string) time.Duration {
	p.mx.RLock()
	defer p.mx.RUnlock()
	return p.v.GetDuration(key)
}

// GetFloat64 implements Provider.
func (p *providerImpl) GetFloat64(key string) float64 {
	p.mx.RLock()
	defer p.mx.RUnlock()
	return p.v.GetFloat64(key)
}

// GetInt implements Provider.
func (p *providerImpl) GetInt(key string) int {
	p.mx.RLock()
	defer p.mx.RUnlock()
	return p.v.GetInt(key)
}

// GetSizeInBytes implements Provider.
func (p *providerImpl) GetSizeInBytes(key string) uint {
	p.mx.RLock()
	defer p.mx.RUnlock()
	return p.v.GetSizeInBytes(key)
}

// GetString implements Provider.
func (p *providerImpl) GetString(key string) string {
	p.mx.RLock()
	defer p.mx.RUnlock()
	return p.v.GetString(key)
}

// GetStringSlice implements Provider.
func (p *providerImpl) GetStringSlice(key string) []string {
	p.mx.RLock()
	defer p.mx.RUnlock()
	return p.v.GetStringSlice(key)
}

// GetTime implements Provider.
func (p *providerImpl) GetTime(key string) time.Time {
	p.mx.RLock()
	defer p.mx.RUnlock()
	return p.v.GetTime(key)
}

// InConfig implements Provider.
func (p *providerImpl) InConfig(key string) bool {
	p.mx.RLock()
	defer p.mx.RUnlock()
	return p.v.InConfig(key)
}

// IsSet implements Provider.
func (p *providerImpl) IsSet(key string) bool {
	p.mx.RLock()
	defer p.mx.RUnlock()
	return p.v.IsSet(key)
}

// ReadInConfig implements Provider.
func (p *providerImpl) ReadInConfig() error {
	p.mx.RLock()
	defer p.mx.RUnlock()
	return p.v.ReadInConfig()
}

// Set implements Provider.
func (p *providerImpl) Set(key string, value interface{}) {
	p.mx.RLock()
	defer p.mx.RUnlock()
	p.v.Set(key, value)
}

// SetConfigFile implements Provider.
func (p *providerImpl) SetConfigFile(in string) {
	p.mx.RLock()
	defer p.mx.RUnlock()
	p.v.SetConfigFile(in)
}

var (
	defaultConfig     Provider
	defaultConfigOnce sync.Once
)

func defaultConfigOnceFunc() {
	defaultConfig = defaultViperConfig()
}

func Config() Provider {
	defaultConfigOnce.Do(defaultConfigOnceFunc)
	return defaultConfig
}

func defaultViperConfig() Provider {
	v := viper.New()
	v.SetDefault(SimpleWorkerTickInterval, 15) // minutes
	replacer := strings.NewReplacer("-", "_", ".", "_")
	v.SetEnvKeyReplacer(replacer)
	v.AutomaticEnv()
	return &providerImpl{
		v: v,
	}
}
