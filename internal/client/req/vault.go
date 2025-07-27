package req

import (
	"context"
	"errors"
	"log"

	"github.com/LekcRg/GophKeeper/internal/models"
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
