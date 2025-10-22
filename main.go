package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"concierge-be/config"
	"concierge-be/database"
	"concierge-be/internal/amenities"
	"concierge-be/internal/amenities_categories"
	"concierge-be/internal/facilities"
	"concierge-be/internal/roles"
	"concierge-be/internal/tenants"
	"concierge-be/internal/users"
	"concierge-be/router"
	"concierge-be/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	// 打印启动 banner
	utils.PrintBanner()

	// 解析命令行参数
	env := flag.String("e", "development", "运行环境 (development, production, test)")
	flag.Parse()

	// 加载配置
	config.LoadConfig(*env)

	// 设置 Gin 模式
	gin.SetMode(config.AppConfig.Server.Mode)

	// 初始化数据库
	database.InitDB()

	// 自动迁移数据库表
	if err := database.GetDB().AutoMigrate(
		&roles.Role{},
		&users.User{}, 
		&users.UserTenant{}, 
		&tenants.Tenant{},
		&amenities_categories.AmenityCategory{},
		&amenities.Amenity{},
		&facilities.Facility{},
		&facilities.Booking{},
	); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// 运行角色迁移脚本
	if err := runRoleMigration(); err != nil {
		log.Printf("Warning: Failed to run role migration: %v", err)
	}

	// 种子数据：创建默认角色
	roleService := roles.NewService()
	if err := roleService.SeedRoles(); err != nil {
		log.Printf("Warning: Failed to seed roles: %v", err)
	} else {
		log.Println("Roles seeded successfully")
	}

	// 为现有用户分配默认角色
	userService := users.NewService()
	if err := userService.AssignDefaultRoleToUsersWithoutRole(); err != nil {
		log.Printf("Warning: Failed to assign default roles to existing users: %v", err)
	} else {
		log.Println("Default roles assigned to existing users successfully")
	}

	// 设置路由
	r := router.SetupRouter()

	// 启动服务
	addr := fmt.Sprintf(":%s", config.AppConfig.Server.Port)
	log.Printf("Server is running on http://localhost%s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// runRoleMigration runs the role migration SQL script
func runRoleMigration() error {
	// Read the SQL file
	sqlFile := filepath.Join("scripts", "role-migration.sql")
	sqlContent, err := ioutil.ReadFile(sqlFile)
	if err != nil {
		return fmt.Errorf("failed to read role migration file: %v", err)
	}

	// Execute the SQL
	db := database.GetDB()
	if err := db.Exec(string(sqlContent)).Error; err != nil {
		return fmt.Errorf("failed to execute role migration: %v", err)
	}

	log.Println("Role migration completed successfully")
	return nil
}
