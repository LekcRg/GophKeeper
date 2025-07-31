package req

import (
	"context"

	"github.com/LekcRg/GophKeeper/internal/models"
	"github.com/LekcRg/GophKeeper/internal/routes"
)

func (r *Request) CreateVaultItem(
	ctx context.Context, item models.VaultCreateItemReq,
) (models.VaultItem, error) {
	var (
		resBody models.VaultItem
		resErrs map[string]string
	)

	res, err := r.client.R().
		SetHeader("Authorization", "Bearer "+r.config.Key).
		SetContext(ctx).
		SetBody(item).
		SetResult(&resBody).
		SetError(&resErrs).
		Post(r.config.Address + routes.VaultCreateItem)
	if err != nil {
		return models.VaultItem{}, err
	}

	if res.StatusCode() > 299 {
		return models.VaultItem{}, &ResError{Errors: resErrs}
	}

	return resBody, nil
}

func (r *Request) VaultGetAll(ctx context.Context) ([]models.VaultItem, error) {
	var (
		resBody []models.VaultItem
		resErrs map[string]string
	)

	res, err := r.client.R().
		SetHeader("Authorization", "Bearer "+r.config.Key).
		SetContext(ctx).
		SetResult(&resBody).
		SetError(&resErrs).
		Get(r.config.Address + routes.VaultGetAll)
	if err != nil {
		return resBody, err
	}

	if res.StatusCode() > minErrStatus {
		return resBody, &ResError{Errors: resErrs}
	}

	return resBody, nil
}
