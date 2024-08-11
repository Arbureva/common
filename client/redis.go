package client

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type Redis struct {
	Host     string `toml:"host"`
	Port     string `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`

	// 从网络连接中读取数据超时时间，可能的值：
	//  0 - 默认值，3秒
	// -1 - 无超时，无限期的阻塞
	// -2 - 不进行超时设置，不调用 SetReadDeadline 方法
	ReadTimeout int `toml:"read_timeout"`

	// 把数据写入网络连接的超时时间，可能的值：
	//  0 - 默认值，3秒
	// -1 - 无超时，无限期的阻塞
	// -2 - 不进行超时设置，不调用 SetWriteDeadline 方法
	WriteTimeout int `toml:"write_timeout"`

	// 连接池的类型，有 LIFO 和 FIFO 两种模式，
	// PoolFIFO 为 false 时使用 LIFO 模式，为 true 使用 FIFO 模式。
	// 当一个连接使用完毕时会把连接归还给连接池，连接池会把连接放入队尾，
	// LIFO 模式时，每次取空闲连接会从"队尾"取，就是刚放入队尾的空闲连接，
	// 也就是说 LIFO 每次使用的都是热连接，连接池有机会关闭"队头"的长期空闲连接，
	// 并且从概率上，刚放入的热连接健康状态会更好；
	// 而 FIFO 模式则相反，每次取空闲连接会从"队头"取，相比较于 LIFO 模式，
	// 会使整个连接池的连接使用更加平均，有点类似于负载均衡寻轮模式，会循环的使用
	// 连接池的所有连接，如果你使用 go-redis 当做代理让后端 redis 节点负载更平均的话，
	// FIFO 模式对你很有用。
	// 如果你不确定使用什么模式，请保持默认 PoolFIFO = false
	PoolFIFO bool `toml:"pool_fifo"`

	// 连接池最大连接数量，注意：这里不包括 pub/sub，pub/sub 将使用独立的网络连接
	// 默认为 10 * runtime.GOMAXPROCS
	PoolSize int `toml:"pool_size"`

	// PoolTimeout 代表如果连接池所有连接都在使用中，等待获取连接时间，超时将返回错误
	// 默认是 1秒 + ReadTimeout
	PoolTimeout int `toml:"pool_timeout"`

	// 连接池保持的最小空闲连接数，它受到PoolSize的限制
	// 默认为0，不保持
	MinIdleConns int `toml:"min_idle_conns"`

	// 连接池保持的最大空闲连接数，多余的空闲连接将被关闭
	// 默认为0，不限制
	MaxIdleConns int `toml:"max_idle_conns"`

	// ConnMaxIdleTime 是最大空闲时间，超过这个时间将被关闭。
	// 如果 ConnMaxIdleTime <= 0，则连接不会因为空闲而被关闭。
	// 默认值是30分钟，-1禁用
	ConnMaxIdleTime time.Duration `toml:"conn_max_idle_time"`

	// ConnMaxLifetime 是一个连接的生存时间，
	// 和 ConnMaxIdleTime 不同，ConnMaxLifetime 表示连接最大的存活时间
	// 如果 ConnMaxLifetime <= 0，则连接不会有使用时间限制
	// 默认值为0，代表连接没有时间限制
	ConnMaxLifetime time.Duration `toml:"conn_max_lifetime"`
}

func MustNewRedisClient(options Redis) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", options.Host, options.Port),
		Password: options.Password,
		Username: options.User,
	})
}
