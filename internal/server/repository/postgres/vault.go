package postgres

import (
	"context"

	"github.com/LekcRg/GophKeeper/internal/errs"
	"github.com/LekcRg/GophKeeper/internal/models"
	"github.com/jmoiron/sqlx"
)

type VaultRepo struct {
	db *sqlx.DB
}

func NewVaultRepo(db *sqlx.DB) *VaultRepo {
	return &VaultRepo{
		db: db,
	}
}

func (vr *VaultRepo) CreateItem(ctx context.Context, item models.VaultItem) (models.VaultItem, error) {
	res := models.VaultItem{}
	query := `INSERT INTO vault (user_id, name, type, encrypted_data)
	VALUES (:user_id, :name, :type, :encrypted_data)
	RETURNING id, user_id, name, type, encrypted_data, created_at, updated_at`

	rows, err := vr.db.NamedQueryContext(ctx, query, item)
	if err != nil {
		return models.VaultItem{}, err
	}
	defer rows.Close()

	if !rows.Next() {
		return models.VaultItem{}, errs.ErrRepoRowsNotFound
	}

	err = rows.StructScan(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (vr *VaultRepo) GetAllItems(
	ctx context.Context, userID int,
) ([]models.VaultItem, error) {
	var res []models.VaultItem

	query := `SElECT * FROM vault WHERE user_id = $1`

	err := vr.db.SelectContext(ctx, &res, query, userID)
	if err != nil {
		return res, err
	}

	return res, nil
}
