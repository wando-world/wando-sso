package test

import (
	"github.com/stretchr/testify/assert"
	apiModels "github.com/wando-world/wando-sso/internal/api/models"
	"github.com/wando-world/wando-sso/internal/mappers"
	"testing"
)

func TestLoginRequestToUser(t *testing.T) {
	authMapper := mappers.NewAuthMapper()

	expectedUserID, expectedVerifiedCode := "wando_test_id", "wando123"

	req := apiModels.LoginRequest{
		UserID:       expectedUserID,
		VerifiedCode: expectedVerifiedCode,
	}

	resultUser := authMapper.LoginRequestToUser(req)

	assert.Equal(t, expectedUserID, resultUser.UserID, "user id 검증")
	assert.Equal(t, expectedVerifiedCode, resultUser.VerifiedCode, "code 검증")
}
