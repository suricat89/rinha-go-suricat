package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suricat89/rinha-2024-q1/src/config"
	"github.com/suricat89/rinha-2024-q1/src/config/cache"
	"github.com/suricat89/rinha-2024-q1/src/config/database"
	"github.com/suricat89/rinha-2024-q1/src/controller"
	"github.com/suricat89/rinha-2024-q1/src/interfaces"
	"github.com/suricat89/rinha-2024-q1/src/repository"
	"github.com/suricat89/rinha-2024-q1/src/router"
	"go.mongodb.org/mongo-driver/mongo"
)

// Reference https://github.com/leorcvargas
func startProfiling() {
	conf := config.Env.Server.Profiling

	cpuProf, err := os.Create(conf.CpuFilePath)
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(cpuProf)

	memProf, err := os.Create(conf.MemoryFilePath)
	if err != nil {
		log.Fatal(err)
	}
	pprof.WriteHeapProfile(memProf)

	after := time.After(2 * time.Minute)

	go func() {
		<-after
		log.Info("Stopping CPU and Memory profiling")
		pprof.StopCPUProfile()
		cpuProf.Close()
		memProf.Close()
	}()
}

func main() {
	if config.Env.Server.Profiling.Enabled {
		startProfiling()
	}

	app := fiber.New(fiber.Config{
    Concurrency: 256,
    ReduceMemoryUsage: true,
  })

	app.Use(logger.New())
	log.SetLevel(config.Env.Server.LogLevel)

	var databaseConfig interfaces.DatabaseConfig

	if config.Env.Database.DatabaseSelected == "postgresql" {
		databaseConfig = database.NewPostgresDatabase()
	} else if config.Env.Database.DatabaseSelected == "mongodb" {
		databaseConfig = database.NewMongoDatabase()
	} else {
		log.Fatalf("Must select database 'postgresql' or 'mongodb' using env DATABASE_SELECTED")
	}

	db, err := databaseConfig.InitDb()
	if err != nil {
		log.Fatalf("Error opening database connection. Error: %s", err.Error())
		return
	}

	err = databaseConfig.PingDb()
	if err != nil {
		log.Fatalf("Error accessing database. Error: %s", err.Error())
		return
	}
	defer databaseConfig.Close()
	log.Info("Connected to database")

	cache.InitCache()
	err = cache.PingRedis()
	if err != nil {
		log.Panicf("Error accessing Redis. Error: %s", err.Error())
		return
	}
	defer cache.Rdb.Close()
	log.Info("Connected to cache")

	var databaseRepository interfaces.DatabaseRepository

	switch db := db.(type) {
	case *mongo.Database:
		databaseRepository = repository.NewMongoRepository(db)
	case *pgxpool.Pool:
		databaseRepository = repository.NewPostgresRepository(db)
	}

	cacheRepository := repository.NewCacheRepository(cache.Rdb)
	controller := controller.NewCustomerController(databaseRepository, cacheRepository)
	router := router.NewRouter(controller)
	router.Load(app)

	address := fmt.Sprintf(":%d", config.Env.Server.Port)
	app.Listen(address, fiber.ListenConfig{
		EnablePrefork: config.Env.Server.Prefork,
	})

	log.Infof("API started on port %s", address)
}
