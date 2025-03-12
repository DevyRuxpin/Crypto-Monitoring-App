package health

import (
    "context"
    "database/sql"
    "github.com/go-redis/redis/v8"
    "syscall"
)

type HealthChecker struct {
    db    *sql.DB
    redis *redis.Client
}

type HealthStatus struct {
    Database    bool    `json:"database"`
    Redis       bool    `json:"redis"`
    DiskSpace   bool    `json:"disk_space"`
    Memory      bool    `json:"memory"`
    DiskUsage   float64 `json:"disk_usage"`
    MemoryUsage float64 `json:"memory_usage"`
}

func NewHealthChecker(db *sql.DB, redis *redis.Client) *HealthChecker {
    return &HealthChecker{
        db:    db,
        redis: redis,
    }
}

func (h *HealthChecker) Check() HealthStatus {
    status := HealthStatus{}

    // Check database
    status.Database = h.db.Ping() == nil

    // Check Redis
    status.Redis = h.redis.Ping(context.Background()).Err() == nil

    // Check disk space
    diskUsage, diskOk := checkDiskSpace()
    status.DiskSpace = diskOk
    status.DiskUsage = diskUsage

    // Check memory
    memUsage, memOk := checkMemoryUsage()
    status.Memory = memOk
    status.MemoryUsage = memUsage

    return status
}

func checkDiskSpace() (float64, bool) {
    var stat syscall.Statfs_t
    err := syscall.Statfs("/", &stat)
    if err != nil {
        return 0, false
    }

    // Calculate disk usage percentage
    total := float64(stat.Blocks) * float64(stat.Bsize)
    free := float64(stat.Bfree) * float64(stat.Bsize)
    used := total - free
    usagePercent := (used / total) * 100

    return usagePercent, usagePercent < 90
}

func checkMemoryUsage() (float64, bool) {
    var info syscall.Sysinfo_t
    err := syscall.Sysinfo(&info)
    if err != nil {
        return 0, false
    }

    // Calculate memory usage percentage
    total := float64(info.Totalram)
    free := float64(info.Freeram)
    used := total - free
    usagePercent := (used / total) * 100

    return usagePercent, usagePercent < 90
}