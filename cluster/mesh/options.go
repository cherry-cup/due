package mesh

import (
	"context"
	"github.com/dobyte/due/config"
	"github.com/dobyte/due/crypto"
	"github.com/dobyte/due/encoding"
	"github.com/dobyte/due/locate"
	"github.com/dobyte/due/registry"
	"github.com/dobyte/due/transport"
	"time"
)

const (
	defaultName    = "mesh"          // 默认节点名称
	defaultCodec   = "proto"         // 默认编解码器名称
	defaultTimeout = 3 * time.Second // 默认超时时间
)

const (
	defaultNameKey      = "config.cluster.mesh.name"
	defaultCodecKey     = "config.cluster.mesh.codec"
	defaultTimeoutKey   = "config.cluster.mesh.timeout"
	defaultEncryptorKey = "config.cluster.mesh.encryptor"
	defaultDecryptorKey = "config.cluster.mesh.decryptor"
)

type Option func(o *options)

type options struct {
	name        string                // 实例名称
	ctx         context.Context       // 上下文
	codec       encoding.Codec        // 编解码器
	timeout     time.Duration         // RPC调用超时时间
	locator     locate.Locator        // 用户定位器
	registry    registry.Registry     // 服务注册器
	transporter transport.Transporter // 消息传输器
	encryptor   crypto.Encryptor      // 消息加密器
	decryptor   crypto.Decryptor      // 消息解密器
}

func defaultOptions() *options {
	opts := &options{
		ctx:     context.Background(),
		name:    defaultName,
		codec:   encoding.Invoke(defaultCodec),
		timeout: defaultTimeout,
	}

	if name := config.Get(defaultNameKey).String(); name != "" {
		opts.name = name
	}

	if codec := config.Get(defaultCodecKey).String(); codec != "" {
		opts.codec = encoding.Invoke(codec)
	}

	if timeout := config.Get(defaultTimeoutKey).Int64(); timeout > 0 {
		opts.timeout = time.Duration(timeout) * time.Second
	}

	if encryptor := config.Get(defaultEncryptorKey).String(); encryptor != "" {
		opts.encryptor = crypto.InvokeEncryptor(encryptor)
	}

	if decryptor := config.Get(defaultDecryptorKey).String(); decryptor != "" {
		opts.decryptor = crypto.InvokeDecryptor(decryptor)
	}

	return opts
}

// WithName 设置实例名称
func WithName(name string) Option {
	return func(o *options) { o.name = name }
}

// WithCodec 设置编解码器
func WithCodec(codec encoding.Codec) Option {
	return func(o *options) { o.codec = codec }
}

// WithContext 设置上下文
func WithContext(ctx context.Context) Option {
	return func(o *options) { o.ctx = ctx }
}

// WithTimeout 设置RPC调用超时时间
func WithTimeout(timeout time.Duration) Option {
	return func(o *options) { o.timeout = timeout }
}

// WithLocator 设置定位器
func WithLocator(locator locate.Locator) Option {
	return func(o *options) { o.locator = locator }
}

// WithRegistry 设置服务注册器
func WithRegistry(r registry.Registry) Option {
	return func(o *options) { o.registry = r }
}

// WithTransporter 设置消息传输器
func WithTransporter(transporter transport.Transporter) Option {
	return func(o *options) { o.transporter = transporter }
}

// WithEncryptor 设置消息加密器
func WithEncryptor(encryptor crypto.Encryptor) Option {
	return func(o *options) { o.encryptor = encryptor }
}

// WithDecryptor 设置消息解密器
func WithDecryptor(decryptor crypto.Decryptor) Option {
	return func(o *options) { o.decryptor = decryptor }
}
