package add

import (
	"net/http"

	"github.com/teatou/tolerant/pkg/mylogger"
)

type Adder interface {
	Add(to, sum int) error
}

func New(adder Adder, logger mylogger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
