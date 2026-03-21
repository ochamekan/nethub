package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ochamekan/nethub/backend/internal/dto"
	"github.com/ochamekan/nethub/backend/pkg/models"
)

type DeviceRepository interface {
	CreateDevice(ctx context.Context, ip, hostname string, is_active bool) (models.Device, error)
	GetDevices(ctx context.Context, data dto.GetDevicesRequest) ([]models.Device, error)
	GetDevice(ctx context.Context, id int64) (models.Device, error)
	UpdateDevice(ctx context.Context, id int64, data dto.UpdateDeviceRequest) (models.Device, error)
	DeleteDevice(ctx context.Context, id int64) (models.Device, error)
}

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) DeviceRepository {
	return &Repository{db}
}

var ErrNotFound = errors.New("not found")

func (r *Repository) CreateDevice(ctx context.Context, ip, hostname string, is_active bool) (models.Device, error) {
	rows, err := r.db.Query(ctx, "INSERT INTO devices (ip, hostname, is_active) VALUES ($1, $2, $3) RETURNING *", ip, hostname, is_active)
	if err != nil {
		return models.Device{}, err
	}
	defer rows.Close()
	d, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Device])
	if err != nil {
		return models.Device{}, err
	}
	return d, nil
}

func (r *Repository) GetDevices(ctx context.Context, data dto.GetDevicesRequest) ([]models.Device, error) {
	q := "SELECT * FROM devices"
	args, argsIdx := []any{}, 1
	var clauses []string

	if data.Hostname != nil {
		clauses = append(clauses, fmt.Sprintf("hostname ILIKE $%d", argsIdx))
		args = append(args, "%"+*data.Hostname+"%")
		argsIdx++
	}

	if data.IsActive != nil {
		clauses = append(clauses, fmt.Sprintf("is_active = $%d", argsIdx))
		args = append(args, *data.IsActive)
		argsIdx++
	}

	if len(clauses) > 0 {
		q += " WHERE" + strings.Join(clauses, " AND ")
	}

	rows, err := r.db.Query(ctx, q, args...)
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
	defer rows.Close()

	d, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Device])
	if err != nil {
		if err == pgx.ErrNoRows {
			return models.Device{}, ErrNotFound
		}
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
		if err == pgx.ErrNoRows {
			return models.Device{}, ErrNotFound
		}
		return models.Device{}, err
	}
	defer rows.Close()
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
	defer rows.Close()
	d, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Device])
	if err != nil {
		if err == pgx.ErrNoRows {
			return models.Device{}, ErrNotFound
		}
		return models.Device{}, err
	}

	return d, nil
}
