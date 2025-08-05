package req

import (
	"context"

	"github.com/LekcRg/GophKeeper/internal/models"
	"github.com/LekcRg/GophKeeper/internal/routes"
)

func (r *Request) getToken(
	ctx context.Context, body any, path string,
) (models.APIKeyRes, error) {
	var (
		resBody models.APIKeyRes
		resErrs map[string]string
	)

	res, err := r.client.R().
		SetContext(ctx).
		SetResult(&resBody).
		SetError(&resErrs).
		SetBody(body).
		Post(r.config.Address + path)
	if err != nil {
		return resBody, err
	}

	if res.StatusCode() > minErrStatus {
		return resBody, &ResError{Errors: resErrs}
	}

	return resBody, nil
}

func (r *Request) Register(
	ctx context.Context, body models.UserReq,
) (models.APIKeyRes, error) {
	return r.getToken(ctx, body, routes.UserRegister)
}

func (r *Request) UpdateAPIKey(
	ctx context.Context, body models.UserLogin,
) (models.APIKeyRes, error) {
	return r.getToken(ctx, body, routes.UserUpdateKey)
}

func (r *Request) GetCredentials(
	ctx context.Context, key string,
) (models.CryptoParamsRes, error) {
	var (
		resBody models.CryptoParamsRes
		resErrs map[string]string
	)

	res, err := r.client.R().
		SetHeader("Authorization", "Bearer "+key).
		SetContext(ctx).
		SetResult(&resBody).
		SetError(&resErrs).
		Get(r.config.Address + routes.UserGetCryptoParams)
	if err != nil {
		return resBody, err
	}

	if res.StatusCode() > minErrStatus {
		return resBody, &ResError{Errors: resErrs}
	}

	return resBody, nil
}
