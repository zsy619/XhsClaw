// Package main 提供数据库迁移命令
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"xiaohongshu/internal/config"
	"xiaohongshu/internal/repository"
)

func main() {
	// 解析命令行参数
	migrateCmd := flag.NewFlagSet("migrate", flag.ExitOnError)
	configPath := migrateCmd.String("config", "config.yaml", "配置文件路径")

	// 解析子命令
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	subCommand := os.Args[1]

	switch subCommand {
	case "user-configs":
		// 执行 user_configs 迁移
		migrateCmd.Parse(os.Args[2:])

		// 加载配置
		cfg, err := config.LoadConfig(*configPath)
		if err != nil {
			log.Fatalf("加载配置文件失败: %v", err)
		}

		// 初始化数据库
		if err := repository.InitDatabase(&cfg.Database); err != nil {
			log.Fatalf("初始化数据库失败: %v", err)
		}

		// 执行迁移
		fmt.Println("开始迁移 user_configs 表...")
		fmt.Println("")

		if err := repository.MigrateUserConfigsToSeparateTables(); err != nil {
			log.Fatalf("迁移失败: %v", err)
		}

		fmt.Println("")
		fmt.Println("迁移完成！")

	case "backup":
		// 创建备份
		migrateCmd.Parse(os.Args[2:])

		// 加载配置
		cfg, err := config.LoadConfig(*configPath)
		if err != nil {
			log.Fatalf("加载配置文件失败: %v", err)
		}

		// 初始化数据库
		if err := repository.InitDatabase(&cfg.Database); err != nil {
			log.Fatalf("初始化数据库失败: %v", err)
		}

		// 创建备份
		backupTable, err := repository.BackupUserConfigs()
		if err != nil {
			log.Fatalf("备份失败: %v", err)
		}

		fmt.Printf("备份成功！备份表名: %s\n", backupTable)

	case "restore":
		// 恢复备份
		if len(os.Args) < 3 {
			fmt.Println("用法: migrate restore <备份表名>")
			os.Exit(1)
		}

		backupTable := os.Args[2]
		migrateCmd.Parse(os.Args[3:])

		// 加载配置
		cfg, err := config.LoadConfig(*configPath)
		if err != nil {
			log.Fatalf("加载配置文件失败: %v", err)
		}

		// 初始化数据库
		if err := repository.InitDatabase(&cfg.Database); err != nil {
			log.Fatalf("初始化数据库失败: %v", err)
		}

		// 恢复备份
		if err := repository.RestoreUserConfigsFromBackup(backupTable); err != nil {
			log.Fatalf("恢复失败: %v", err)
		}

		fmt.Printf("恢复成功！从备份表: %s\n", backupTable)

	case "clean-backup":
		// 清理过期备份
		migrateCmd.Parse(os.Args[2:])

		// 加载配置
		cfg, err := config.LoadConfig(*configPath)
		if err != nil {
			log.Fatalf("加载配置文件失败: %v", err)
		}

		// 初始化数据库
		if err := repository.InitDatabase(&cfg.Database); err != nil {
			log.Fatalf("初始化数据库失败: %v", err)
		}

		// 清理备份
		if err := repository.CleanBackupTables(); err != nil {
			log.Fatalf("清理失败: %v", err)
		}

		fmt.Println("清理完成！")

	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("数据库迁移工具")
	fmt.Println("")
	fmt.Println("用法:")
	fmt.Println("  migrate <command> [options]")
	fmt.Println("")
	fmt.Println("可用命令:")
	fmt.Println("  user-configs    迁移 user_configs 表到独立表")
	fmt.Println("  backup          备份 user_configs 表")
	fmt.Println("  restore <table>  从备份表恢复")
	fmt.Println("  clean-backup    清理过期的备份表")
	fmt.Println("")
	fmt.Println("全局选项:")
	fmt.Println("  -config string  配置文件路径 (默认: config.yaml)")
	fmt.Println("")
	fmt.Println("示例:")
	fmt.Println("  migrate user-configs")
	fmt.Println("  migrate user-configs -config /path/to/config.yaml")
	fmt.Println("  migrate backup")
	fmt.Println("  migrate restore user_configs_backup_1234567890")
	fmt.Println("  migrate clean-backup")
}
