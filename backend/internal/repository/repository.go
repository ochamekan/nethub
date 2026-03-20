package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ochamekan/nethub/backend/internal/dto"
	"github.com/ochamekan/nethub/backend/pkg/models"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) Repository {
	return Repository{db}
}

func (r *Repository) CreateDevice(ctx context.Context, ip, hostname string) (models.Device, error) {
	rows, err := r.db.Query(ctx, "INSERT INTO devices (ip, hostname) VALUES ($1, $2) RETURNING *", ip, hostname)
	if err != nil {
		return models.Device{}, err
	}
	rows.Close()
	d, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Device])
	if err != nil {
		return models.Device{}, err
	}
	return d, nil
}

func (r *Repository) GetDevices(ctx context.Context) ([]models.Device, error) {
	rows, err := r.db.Query(ctx, "SELECT * FROM devices")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	devices, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Device])
	if err != nil {
		return nil, err
	}
	return devices, nil
}

func (r *Repository) GetDevice(ctx context.Context, id int64) (models.Device, error) {
	rows, err := r.db.Query(ctx, "SELECT * FROM devices WHERE id = $1", id)
	if err != nil {
		return models.Device{}, err
	}
	rows.Close()

	d, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Device])
	if err != nil {
		return models.Device{}, err
	}

	return d, nil
}

func (r *Repository) UpdateDevice(ctx context.Context, id int64, data dto.UpdateDeviceRequest) (models.Device, error) {
	var clauses []string
	args, argIdx := []any{}, 1

	if data.Hostname != nil {
		clauses = append(clauses, fmt.Sprintf("hostname = $%d", argIdx))
		args = append(args, *data.Hostname)
		argIdx++
	}

	if data.IP != nil {
		clauses = append(clauses, fmt.Sprintf("ip = $%d", argIdx))
		args = append(args, *data.IP)
		argIdx++
	}

	if data.IsActive != nil {
		clauses = append(clauses, fmt.Sprintf("is_active = $%d", argIdx))
		args = append(args, *data.IsActive)
		argIdx++
	}

	if len(clauses) == 0 {
		return models.Device{}, fmt.Errorf("at least one field required")
	}

	args = append(args, id)

	q := fmt.Sprintf("UPDATE devices SET %s WHERE id = $%d RETURNING *", strings.Join(clauses, ", "), argIdx)
	rows, err := r.db.Query(ctx, q, args...)
	if err != nil {
		return models.Device{}, err
	}
	rows.Close()
	d, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Device])
	if err != nil {
		return models.Device{}, err
	}

	return d, nil
}

func (r *Repository) DeleteDevice(ctx context.Context, id int64) (models.Device, error) {
	rows, err := r.db.Query(ctx, "UPDATE devices SET is_deleted = true WHERE id = $1 RETURNING *", id)
	if err != nil {
		return models.Device{}, err
	}
	rows.Close()
	d, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Device])
	if err != nil {
		return models.Device{}, err
	}

	return d, nil
}
