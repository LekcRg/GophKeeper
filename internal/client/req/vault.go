package req

import (
	"context"
	"errors"
	"log"

	"github.com/LekcRg/GophKeeper/internal/models"
	"github.com/LekcRg/GophKeeper/internal/routes"
)

func (r *Request) CreateVaultItem(
	ctx context.Context, token string, item models.VaultCreateItemReq,
) (models.VaultItem, error) {
	var (
		res     models.VaultItem
		resErrs map[string]string
	)

	_, err := r.client.R().
		SetContext(ctx).
		SetBody(item).
		SetResult(&res).
		SetError(&resErrs).
		Post("http://localhost:8080/vault/create")
	if err != nil {
		return models.VaultItem{}, err
	}

	if len(resErrs) > 0 {
		log.Println(resErrs)
		return models.VaultItem{}, errors.New("req error")
	}

	return res, nil
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
