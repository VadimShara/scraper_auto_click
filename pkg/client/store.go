package store

import (
	"context"
	"fmt"
	"sync"
	"time"
	"scraper-first/internal/config"
	"scraper-first/pkg/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type postgres struct {
	DB *pgxpool.Pool
}

var(
	pgInstance 	*postgres
	pgOnce 		sync.Once
)

// func NewPG(ctx context.Context, Username, Password, Host, Port, Database string) (*postgres, error) {

// 	connString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", Username, Password, Host, Port, Database)

// 	pgOnce.Do(func(){
// 		db, err := pgxpool.New(ctx, connString)
// 		if err != nil {
// 			panic(fmt.Errorf("unable to create connection pool: %w", err))
// 		}

// 		pgInstance = &postgres{db}
// 	})

// 	return pgInstance, nil
// }

func NewClient(ctx context.Context, maxAttampts int, sc config.StorageConfig) (pool *pgxpool.Pool, err error ){
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", sc.Username, sc.Password, sc.Host, sc.Port, sc.Database)

	err = repeatable.DoWithTries( func() error {
		ctx, cancel := context.WithTimeout(ctx, 5 * time.Second)
		defer cancel()

		pool, err = pgxpool.New(ctx, dsn)
		if err != nil {
			return err
		}

		return nil
	}, maxAttampts, 5 * time.Second)

	if err != nil{
		log.Fatal("error do with tries postgresql")
	}

	return pool, nil
}

// func (pg *postgres) Ping(ctx context.Context) error {
// 	return pg.DB.Ping(ctx)
// }

// func (pg *postgres) Close() {
// 	pg.DB.Close()
// }