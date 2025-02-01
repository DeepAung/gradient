package submissions

import (
	"testing"

	"github.com/DeepAung/gradient/website-server/pkg/asserts"
)

func TestCreateSubmission(t *testing.T) {
	t.Run("invalid score", func(t *testing.T) {
		req := createReq
		req.Score = -100
		_, err := submissionsRepo.CreateSubmission(req)
		asserts.EqualError(t, err, ErrInvalidScore)
	})
}
