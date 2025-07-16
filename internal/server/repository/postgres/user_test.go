package postgres

import (
	"context"
	"testing"

	"github.com/LekcRg/GophKeeper/internal/crypto"
	"github.com/LekcRg/GophKeeper/internal/models"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(t *testing.T) {

}

func TestCreateUser(t *testing.T) {
	pg, container := getPostgres(t)
	defer terminateContainer(t, container, pg)

	// type test struct{}
	err := goose.Up(pg.SQL, "../../../../migrations")
	require.NoError(t, err)

	login := "user"

	password := "TestP4$$word1234"

	hash, err := crypto.HashPassword(password)
	require.NoError(t, err)

	err = pg.UserRepo.CreateUser(context.Background(), models.UserReq{
		Login:        "user",
		PasswordHash: hash,
	})
	assert.NoError(t, err)

	user, err := pg.UserRepo.GetUserByLogin(context.Background(), login)
	assert.NoError(t, err)
	assert.Equal(t, login, user.Login)

	assert.True(t, crypto.CheckPasswordHash(password, user.PasswordHash),
		"hash from db is not valid")

	newPassword := "N3ewPa$sw0rd3d3d3d"
	newHash, err := crypto.HashPassword(newPassword)
	require.NoError(t, err)
	err = pg.UserRepo.UpdateUserPassword(context.Background(), models.User{
		Login:        login,
		PasswordHash: newHash,
	})
	require.NoError(t, err)

	changedPassUser, err := pg.UserRepo.GetUserByLogin(context.Background(), login)
	assert.NoError(t, err)
	assert.Equal(t, login, user.Login)

	assert.True(t, crypto.CheckPasswordHash(newPassword, changedPassUser.PasswordHash),
		"hash after UpdateUserPassword is not valid")
}
