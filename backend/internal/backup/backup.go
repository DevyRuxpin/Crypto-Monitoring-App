package backup

import (
    "os/exec"
    "time"
    "fmt"
    "log"
)

type BackupConfig struct {
    DatabaseURL    string
    BackupPath    string
    RetentionDays int
}

func PerformBackup(config BackupConfig) error {
    timestamp := time.Now().Format("2006-01-02-15-04-05")
    filename := fmt.Sprintf("%s/backup-%s.sql", config.BackupPath, timestamp)

    // Execute pg_dump
    cmd := exec.Command("pg_dump",
        "-Fc",
        "-f", filename,
        config.DatabaseURL)

    if err := cmd.Run(); err != nil {
        log.Printf("Backup failed: %v", err)
        return err
    }

    // Cleanup old backups
    cleanup := exec.Command("find",
        config.BackupPath,
        "-name", "backup-*.sql",
        "-mtime", fmt.Sprintf("+%d", config.RetentionDays),
        "-delete")

    if err := cleanup.Run(); err != nil {
        log.Printf("Cleanup failed: %v", err)
        return err
    }

    log.Printf("Backup completed successfully: %s", filename)
    return nil
}

func NewBackupConfig() BackupConfig {
    return BackupConfig{
        DatabaseURL:    "postgresql://localhost:5432/crypto_monitor",
        BackupPath:    "/var/backups/crypto_monitor",
        RetentionDays: 7,
    }
}