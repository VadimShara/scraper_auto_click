package postgresql

import (
	"context"
	"fmt"
	"reflect"
	"scraper-first/internal/service"
	"scraper-first/pkg/logging"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	repository struct {
		client *pgxpool.Pool
		logger *logging.Logger
	}
)

func NewRepository(client *pgxpool.Pool, logger *logging.Logger) service.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}

func (r *repository) Create(ctx context.Context, entity *service.CreateEntity) error {
	var newCar service.Entity

	err := r.client.QueryRow(
		ctx,
		`
			INSERT INTO cars (mark, model, volume, price, year, kilometers, power, transmission, fuel, owners, drive, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
			RETURNING id, mark, model, volume, price, year, kilometers, power, transmission, fuel, owners, drive, created_at, updated_at
		`,

		entity.Mark,
		entity.Model,
		entity.Volume,
		entity.Price,
		entity.Year,
		entity.Kilometers,
		entity.Power,
		entity.Transmission,
		entity.Fuel,
		entity.Owners,
		entity.Drive,
		time.Now(),
		time.Now(),
	).Scan(&newCar.ID, &newCar.Mark, &newCar.Model, &newCar.Volume, &newCar.Price, &newCar.Year, &newCar.Kilometers, &newCar.Power, &newCar.Transmission, &newCar.Fuel, &newCar.Owners, &newCar.Drive, &newCar.CreatedAt, &newCar.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) Check(ctx context.Context, entity *service.CreateEntity) (bool, error) {
	var exists bool

	if err := r.client.QueryRow(
		ctx,
		`
			SELECT EXISTS (
				SELECT 1 FROM cars WHERE mark = $1 AND model = $2 AND volume = $3 AND price = $4 AND year = $5 AND kilometers = $6 AND power = $7 AND transmission = $8 AND fuel = $9 AND owners = $10 AND drive = $11
			);
		`,
		entity.Mark,
		entity.Model,
		entity.Volume,
		entity.Price,
		entity.Year,
		entity.Kilometers,
		entity.Power,
		entity.Transmission,
		entity.Fuel,
		entity.Owners,
		entity.Drive,
	).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

func (r *repository) Find(ctx context.Context, entity *service.CreateEntity) (*[]service.Entity, error) {
	val := reflect.ValueOf(*entity)
	var conditions []string
	var args []interface{}
	var countArgs int
	if val.Kind() == reflect.Struct {
		for i := 0; i < val.NumField(); i++ {
			if val.Field(i).String() != "" {
				countArgs += 1
				conditions = append(conditions, fmt.Sprintf("%s = $%d", strings.ToLower(val.Type().Field(i).Name), countArgs))
				args = append(args, strings.ToLower(val.Field(i).String()))
			}
		}
	} else {
		return nil, fmt.Errorf("not struct")
	}

	q := `SELECT * FROM cars`

	if len(conditions) > 0 {
		q += ` WHERE ` + strings.Join(conditions, " AND ")
	}

	rows, err := r.client.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var res []service.Entity

	for rows.Next() {
		var e service.Entity
		err = rows.Scan(&e.ID, &e.Mark, &e.Model, &e.Volume, &e.Price, &e.CreatedAt, &e.UpdatedAt, &e.Year, &e.Kilometers, &e.Power, &e.Transmission, &e.Fuel, &e.Owners, &e.Drive)
		if err != nil {
			return nil, err
		}

		res = append(res, e)
	}

	return &res, err
}
