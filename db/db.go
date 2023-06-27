package db

import (
	"final-project-backend/config"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db *gorm.DB
)

func getLogger() logger.Interface {
	recover()
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			LogLevel: logger.Info,
		},
	)
}

func Connect() (err error) {
	c := config.InitConfig().DBConfig
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=Asia/Jakarta",
		c.Host,
		c.User,
		c.Password,
		c.DBName,
		c.Port,
	)
	fmt.Println(dsn)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: getLogger(),
	})

	if err != nil {
		return err
	}

	err = db.AutoMigrate()
	if err != nil {
		return err
	}

	return
}

func Get() *gorm.DB {
	return db
}
