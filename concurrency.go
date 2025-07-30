package wtoken

import (
	"context"
	"sync"
	"time"
)

// RWMutexWithTimeout 带超时的读写锁
type RWMutexWithTimeout struct {
	mu      sync.RWMutex
	timeout time.Duration
}

// NewRWMutexWithTimeout 创建带超时的读写锁
func NewRWMutexWithTimeout(timeout time.Duration) *RWMutexWithTimeout {
	return &RWMutexWithTimeout{
		timeout: timeout,
	}
}

/**
 * TryLockWithTimeout 尝试获取写锁（带超时）
 * @param {context.Context} ctx 上下文
 * @returns {bool} 是否成功获取锁
 */
func (rw *RWMutexWithTimeout) TryLockWithTimeout(ctx context.Context) bool {
	done := make(chan bool, 1)
	
	go func() {
		rw.mu.Lock()
		select {
		case done <- true:
			// 成功发送信号
		default:
			// 如果无法发送（超时），释放锁
			rw.mu.Unlock()
		}
	}()

	select {
	case <-done:
		return true
	case <-ctx.Done():
		return false
	case <-time.After(rw.timeout):
		return false
	}
}

/**
 * TryRLockWithTimeout 尝试获取读锁（带超时）
 * @param {context.Context} ctx 上下文
 * @returns {bool} 是否成功获取锁
 */
func (rw *RWMutexWithTimeout) TryRLockWithTimeout(ctx context.Context) bool {
	done := make(chan bool, 1)
	
	go func() {
		rw.mu.RLock()
		select {
		case done <- true:
			// 成功发送信号
		default:
			// 如果无法发送（超时），释放锁
			rw.mu.RUnlock()
		}
	}()

	select {
	case <-done:
		return true
	case <-ctx.Done():
		return false
	case <-time.After(rw.timeout):
		return false
	}
}

/**
 * Unlock 释放写锁
 */
func (rw *RWMutexWithTimeout) Unlock() {
	rw.mu.Unlock()
}

/**
 * RUnlock 释放读锁
 */
func (rw *RWMutexWithTimeout) RUnlock() {
	rw.mu.RUnlock()
}

// WorkerPool 工作池
type WorkerPool struct {
	workerCount int
	jobQueue    chan func()
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
	stopped     bool
	mu          sync.RWMutex
}

// NewWorkerPool 创建工作池
func NewWorkerPool(workerCount int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPool{
		workerCount: workerCount,
		jobQueue:    make(chan func(), workerCount*2),
		ctx:         ctx,
		cancel:      cancel,
	}
}

/**
 * Start 启动工作池
 */
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workerCount; i++ {
		wp.wg.Add(1)
		go wp.worker()
	}
}

/**
 * Submit 提交任务
 * @param {func()} job 任务函数
 * @returns {bool} 是否成功提交
 */
func (wp *WorkerPool) Submit(job func()) bool {
	defer func() {
		if r := recover(); r != nil {
			// 捕获向已关闭channel发送数据的panic
		}
	}()
	
	select {
	case <-wp.ctx.Done():
		return false
	default:
	}
	
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	if wp.stopped {
		return false
	}
	
	select {
	case wp.jobQueue <- job:
		return true
	case <-wp.ctx.Done():
		return false
	default:
		return false
	}
}

/**
 * Stop 停止工作池
 */
func (wp *WorkerPool) Stop() {
	wp.mu.Lock()
	wp.stopped = true
	wp.mu.Unlock()
	
	// 先关闭jobQueue，让worker自然退出
	close(wp.jobQueue)
	// 等待所有worker完成
	wp.wg.Wait()
	// 最后取消context
	wp.cancel()
}

/**
 * worker 工作协程
 */
func (wp *WorkerPool) worker() {
	defer wp.wg.Done()
	for {
		select {
		case job, ok := <-wp.jobQueue:
			if !ok {
				return
			}
			job()
		case <-wp.ctx.Done():
			return
		}
	}
}

// ConcurrentMap 并发安全的映射
type ConcurrentMap[K comparable, V any] struct {
	shards   []*mapShard[K, V]
	shardNum int
}

// mapShard 映射分片
type mapShard[K comparable, V any] struct {
	mu   sync.RWMutex
	data map[K]V
}

// NewConcurrentMap 创建并发安全的映射
func NewConcurrentMap[K comparable, V any](shardNum int) *ConcurrentMap[K, V] {
	if shardNum <= 0 {
		shardNum = 32
	}
	cm := &ConcurrentMap[K, V]{
		shards:   make([]*mapShard[K, V], shardNum),
		shardNum: shardNum,
	}
	for i := 0; i < shardNum; i++ {
		cm.shards[i] = &mapShard[K, V]{
			data: make(map[K]V),
		}
	}
	return cm
}

/**
 * getShard 获取分片
 * @param {K} key 键
 * @returns {*mapShard[K, V]} 分片
 */
func (cm *ConcurrentMap[K, V]) getShard(key K) *mapShard[K, V] {
	hash := cm.hash(key)
	return cm.shards[hash%uint32(cm.shardNum)]
}

/**
 * hash 计算哈希值
 * @param {K} key 键
 * @returns {uint32} 哈希值
 */
func (cm *ConcurrentMap[K, V]) hash(key K) uint32 {
	// 简单的哈希函数，实际应用中可以使用更复杂的哈希算法
	switch k := any(key).(type) {
	case string:
		h := uint32(0)
		for _, c := range k {
			h = h*31 + uint32(c)
		}
		return h
	case int:
		return uint32(k)
	case uint:
		return uint32(k)
	default:
		return 0
	}
}

/**
 * Set 设置键值对
 * @param {K} key 键
 * @param {V} value 值
 */
func (cm *ConcurrentMap[K, V]) Set(key K, value V) {
	shard := cm.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()
	shard.data[key] = value
}

/**
 * Get 获取值
 * @param {K} key 键
 * @returns {V, bool} 值和是否存在
 */
func (cm *ConcurrentMap[K, V]) Get(key K) (V, bool) {
	shard := cm.getShard(key)
	shard.mu.RLock()
	defer shard.mu.RUnlock()
	value, exists := shard.data[key]
	return value, exists
}

/**
 * Delete 删除键值对
 * @param {K} key 键
 * @returns {bool} 是否删除成功
 */
func (cm *ConcurrentMap[K, V]) Delete(key K) bool {
	shard := cm.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()
	_, exists := shard.data[key]
	if exists {
		delete(shard.data, key)
	}
	return exists
}

/**
 * Size 获取映射大小
 * @returns {int} 大小
 */
func (cm *ConcurrentMap[K, V]) Size() int {
	size := 0
	for _, shard := range cm.shards {
		shard.mu.RLock()
		size += len(shard.data)
		shard.mu.RUnlock()
	}
	return size
}

/**
 * Keys 获取所有键
 * @returns {[]K} 键列表
 */
func (cm *ConcurrentMap[K, V]) Keys() []K {
	keys := make([]K, 0)
	for _, shard := range cm.shards {
		shard.mu.RLock()
		for key := range shard.data {
			keys = append(keys, key)
		}
		shard.mu.RUnlock()
	}
	return keys
}

/**
 * Range 遍历所有键值对
 * @param {func(K, V) bool} fn 遍历函数，返回false时停止遍历
 */
func (cm *ConcurrentMap[K, V]) Range(fn func(K, V) bool) {
	for _, shard := range cm.shards {
		shard.mu.RLock()
		for key, value := range shard.data {
			if !fn(key, value) {
				shard.mu.RUnlock()
				return
			}
		}
		shard.mu.RUnlock()
	}
}

/**
 * Clear 清空所有数据
 */
func (cm *ConcurrentMap[K, V]) Clear() {
	for _, shard := range cm.shards {
		shard.mu.Lock()
		shard.data = make(map[K]V)
		shard.mu.Unlock()
	}
}