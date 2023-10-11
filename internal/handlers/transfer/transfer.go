package transfer

import (
	"net/http"

	"github.com/teatou/tolerant/pkg/mylogger"
)

type Transferer interface {
	Transfer(from, to int) error
}

func New(transferer Transferer, logger mylogger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
