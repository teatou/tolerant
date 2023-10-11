package withdraw

import (
	"net/http"

	"github.com/teatou/tolerant/pkg/mylogger"
)

type Withdrawer interface {
	Withdraw(from, sum int) error
}

func New(withdrawer Withdrawer, logger mylogger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
