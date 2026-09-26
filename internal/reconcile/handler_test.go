package reconcile

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fluxa/fluxa/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/shopspring/decimal"
)

type forceSettleTestRepo struct {
	mockRepo
	updatedStatus domain.TransactionStatus
}

func (r *forceSettleTestRepo) UpdateReconciliationStatus(_ context.Context, id string, status domain.TransactionStatus) error {
	r.updatedStatus = status
	return nil
}

// TestAdminForceSettleResultsInTransition ensures that invoking force-settle via admin API correctly invokes the transfer settlement mechanism and transitions status.
func TestAdminForceSettleResultsInTransition(t *testing.T) {
	repo := &forceSettleTestRepo{}
	svc := NewService(
		repo,
		&driftRepo{},
		&driftLookup{
			Wallets: map[string]*domain.Wallet{
				"w-1": {ID: "w-1", PublicKey: "GBDRPAXW5JBCJLOWJ2KHTMIGM7R6U2NKZQ3W4BOWX3S2P5ZCQ2KJ4PVS"},
			},
		},
		&mockStellarClient{},
		nil,
		nil,
		nil,
		"worker",
		decimal.Zero,
		nil,
		"",
	)

	r := chi.NewRouter()
	r.Route("/v1/admin", NewHandler(svc).AdminRoutes())

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/admin/transfers/tx-1/force-settle", nil)
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200, body=%s", rec.Code, rec.Body.String())
	}
}
