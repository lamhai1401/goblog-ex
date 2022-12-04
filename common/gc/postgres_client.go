package gc

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

// GormPostgresClient linter
type GormPostgresClient interface {
	GetDB() *gorm.DB
	Check() bool
	SetupDB(sourceDNS string, replicas ...string) GormPostgresClient
}

// NewgormPostgresClient linter
func NewgormPostgresClient() GormPostgresClient {
	c := &gormPostgresClient{}
	return c
}

type gormPostgresClient struct {
	crDB *gorm.DB
}

// Check linter
func (gc *gormPostgresClient) Check() bool {
	return gc.crDB != nil
}

// SetupDB linter
func (gc *gormPostgresClient) SetupDB(sourceDNS string, replicas ...string) GormPostgresClient {
	fmt.Printf("Connecting with connection string: '%v' \n", sourceDNS)
	var err error

	gc.crDB, err = gorm.Open(postgres.Open(sourceDNS), &gorm.Config{
		SkipDefaultTransaction: true, // speed up 30% query turn off transaction by default
	})
	if err != nil {
		panic(err.Error())
	}

	// init logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,          // Disable color
		},
	)

	if replicas != nil {
		err = gc.addReplica(replicas...)
		if err != nil {
			panic(err.Error())
		}
	}

	gc.crDB.Logger = newLogger
	return gc
}

// addReplica to add replica for source db it was init when create new instance
func (gc *gormPostgresClient) addReplica(lst ...string) (err error) {
	temp := make([]gorm.Dialector, 0)
	cfg := dbresolver.Config{
		Policy: dbresolver.RandomPolicy{},
	}
	for _, dns := range lst {
		temp = append(temp, postgres.Open(dns))
	}
	cfg.Replicas = temp

	err = gc.crDB.Use(gc.setUpResolver(dbresolver.Register(cfg)))
	return
}

func (gc *gormPostgresClient) setUpResolver(cfg *dbresolver.DBResolver) *dbresolver.DBResolver {
	cfg = cfg.SetConnMaxIdleTime(time.Hour).
		SetConnMaxLifetime(24 * time.Hour).
		SetMaxIdleConns(100).
		SetMaxOpenConns(200)
	return cfg
}

// GetDB linter
func (gc *gormPostgresClient) GetDB() *gorm.DB {
	return gc.crDB
}

// Migration linter
func (gc *gormPostgresClient) Migration(lst ...interface{}) error {
	return gc.crDB.AutoMigrate(lst...)
}
