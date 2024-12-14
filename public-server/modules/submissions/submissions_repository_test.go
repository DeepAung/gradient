package submissions

import (
	"testing"

	"github.com/DeepAung/gradient/public-server/pkg/asserts"
)

func TestCreateSubmission(t *testing.T) {
	t.Run("invalid result percent", func(t *testing.T) {
		req := createReq
		req.Results = "PPPPPPPPPP"
		req.ResultPercent = -100
		_, err := submissionsRepo.CreateSubmission(req)
		asserts.EqualError(t, err, ErrInvalidResultPercent)
	})
}
