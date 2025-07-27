package req

import (
	"context"

	"github.com/LekcRg/GophKeeper/internal/models"
	"github.com/LekcRg/GophKeeper/internal/routes"
)

func (r *Request) UserRegister(
	ctx context.Context, body models.UserReq,
) (models.APIKeyRes, error) {
	var (
		minErrStatus = 299
		resBody      models.APIKeyRes
		resErrs      map[string]string
	)

	res, err := r.client.R().
		SetContext(ctx).
		SetResult(&resBody).
		SetError(&resErrs).
		SetBody(body).
		Post("http://localhost:8080" + routes.UserRegister)
	if err != nil {
		return resBody, err
	}

	if res.StatusCode() > minErrStatus {
		return resBody, &ResError{Errors: resErrs}
	}

	return resBody, nil
}
