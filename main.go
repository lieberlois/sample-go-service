package main

import (
	"log"

	"go-rest-api/api"
	"go-rest-api/cfg"
	"go-rest-api/db"
	"go-rest-api/storage"

	"github.com/go-sql-driver/mysql"
)

func main() {
	env := cfg.InitConfig()

	cfg := mysql.Config{
		User:                 env.DBUser,
		Passwd:               env.DBPassword,
		Addr:                 env.DBAddress,
		DBName:               env.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	sql := db.NewMySQLStorage(cfg)
	db, err := sql.Init()

	if err != nil {
		log.Fatal(err)
	}

	store := storage.NewDbStorage(db)

	api := api.NewAPIServer(":3000", store)
	api.Serve()
}
