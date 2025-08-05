package req

import (
	"context"
	"io"
	"strconv"

	"github.com/LekcRg/GophKeeper/internal/errs"
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

	if res.StatusCode() > minErrStatus {
		return models.VaultItem{}, &ResError{Errors: resErrs}
	}

	return resBody, nil
}

func (r *Request) CreateVaultBinaryItem(
	ctx context.Context, item models.VaultBinaryItemUploadReq,
) (models.VaultBinaryItemUploadRes, error) {
	var (
		resBody models.VaultBinaryItemUploadRes
		resErrs map[string]string
	)

	res, err := r.client.R().
		SetHeader("Authorization", "Bearer "+r.config.Key).
		SetContext(ctx).
		SetBody(item).
		SetResult(&resBody).
		SetError(&resErrs).
		Post(r.config.Address + routes.VaultCreateBinaryItem)
	if err != nil {
		return models.VaultBinaryItemUploadRes{}, err
	}

	if res.StatusCode() > minErrStatus {
		return models.VaultBinaryItemUploadRes{}, &ResError{Errors: resErrs}
	}

	return resBody, nil
}

func (r *Request) VaultUploadBinaryFile(ctx context.Context, url string, encFile []byte) error {
	resp, err := r.client.R().
		SetContext(ctx).
		SetBody(encFile).
		SetHeader("Content-Type", "application/octet-stream").
		Put(url)
	if err != nil {
		return err
	}

	if !resp.IsSuccess() {
		return errs.ErrBinaryFileUpload
	}

	return nil
}

func (r *Request) VaultConfirmCreateBinary(ctx context.Context, path string, id int) error {
	var resErrs map[string]string

	body := models.VaultConfirmBinaryUploadReq{
		VaultID: id,
		Path:    path,
	}

	res, err := r.client.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+r.config.Key).
		SetBody(body).
		SetError(&resErrs).
		Post(r.config.Address + routes.VaultBinaryConfirm)
	if err != nil {
		return err
	}

	if res.StatusCode() > minErrStatus {
		return &ResError{Errors: resErrs}
	}

	return nil
}

func (r *Request) VaultGetDowloadBidnaryURL(ctx context.Context, id int) (string, error) {
	var (
		resBody models.GetBinaryFileURLRes
		resErrs map[string]string
	)

	res, err := r.client.R().
		SetHeader("Authorization", "Bearer "+r.config.Key).
		SetContext(ctx).
		SetResult(&resBody).
		SetError(&resErrs).
		Get(r.config.Address + routes.VaultGetBinaryFile + strconv.Itoa(id))
	if err != nil {
		return "", err
	}

	if res.StatusCode() > minErrStatus {
		return "", &ResError{Errors: resErrs}
	}

	return resBody.URL, nil
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

func (r *Request) DownloadBinary(url string) ([]byte, error) {
	resp, err := r.client.R().Get(url)
	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return nil, errs.ErrDownloadBinary
	}

	file, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return file, nil
}
