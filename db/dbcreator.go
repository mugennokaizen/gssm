package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/samber/do"
	"github.com/testcontainers/testcontainers-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"path/filepath"
	"runtime"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // used by migrator
	_ "github.com/golang-migrate/migrate/v4/source/file"       // used by migrator
	_ "github.com/jackc/pgx/v4/stdlib"                         // used by migrator
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	DbName = "gssn"
	DbUser = "user"
	DbPass = "password"
)

type TestDatabase struct {
	DbInstance *gorm.DB
	DbAddress  string
	container  testcontainers.Container
}

func CreateInjectorWithDb() *do.Injector {
	testDB := SetupTestDatabase()
	testDbInstance := testDB.DbInstance

	injector := do.New()
	do.ProvideValue(injector, testDbInstance)

	return injector
}

func SetupTestDatabase() *TestDatabase {

	// setup db container
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	container, dbInstance, dbAddr, err := createContainer(ctx)
	if err != nil {
		log.Fatal("failed to setup test", err)
	}

	// migrate db schema
	err = migrateDb(dbAddr)
	if err != nil {
		log.Fatal("failed to perform db migration", err)
	}
	cancel()

	return &TestDatabase{
		container:  container,
		DbInstance: dbInstance,
		DbAddress:  dbAddr,
	}
}

func (tdb *TestDatabase) TearDown() {
	db, err := tdb.DbInstance.DB()
	if err != nil {
		return
	}
	_ = db.Close()
	// remove test container
	_ = tdb.container.Terminate(context.Background())
}

func createContainer(ctx context.Context) (testcontainers.Container, *gorm.DB, string, error) {
	var env = map[string]string{
		"POSTGRES_PASSWORD": DbPass,
		"POSTGRES_USER":     DbUser,
		"POSTGRES_DB":       DbName,
	}
	var port = "5432/tcp"

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:15",
			ExposedPorts: []string{port},
			Env:          env,
			WaitingFor: wait.ForAll(
				wait.ForLog("database system is ready to accept connections"),
				wait.ForListeningPort("5432/tcp"),
				wait.ForExec([]string{"apt", "install", "/pgx-ulid.deb"}),
			),
			Files: []testcontainers.ContainerFile{
				{
					HostFilePath:      "../lib/pgx-ulid.deb",
					ContainerFilePath: "/pgx-ulid.deb",
					FileMode:          0o775,
				},
			},
		},
		Started: true,
	}
	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return container, nil, "", fmt.Errorf("failed to start container: %v", err)
	}

	p, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return container, nil, "", fmt.Errorf("failed to get container external port: %v", err)
	}

	log.Println("postgres container ready and running at port: ", p.Port())

	time.Sleep(time.Second)
	dbAddr := fmt.Sprintf("localhost:%s", p.Port())
	db, err := gorm.Open(postgres.Open(fmt.Sprintf("host=localhost user=%v password=%v dbname=%v port=%s sslmode=disable", DbUser, DbPass, DbName, p.Port())), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		return container, db, dbAddr, fmt.Errorf("failed to establish database connection: %v", err)
	}

	return container, db, dbAddr, nil
}

func migrateDb(dbAddr string) error {

	// get location of test
	_, path, _, _ := runtime.Caller(0)
	pathToMigrationFiles := filepath.Dir(path) + "/migrations"

	databaseURL := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", DbUser, DbPass, dbAddr, DbName)
	m, err := migrate.New(fmt.Sprintf("file:%s", pathToMigrationFiles), databaseURL)
	if err != nil {
		return err
	}
	defer m.Close()

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	log.Println("migration done")

	return nil
}
