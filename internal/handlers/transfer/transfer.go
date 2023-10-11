package transfer

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/teatou/tolerant/internal/storage"
	"github.com/teatou/tolerant/internal/storage/cache"
	"github.com/teatou/tolerant/pkg/mylogger"
)

type Transferer interface {
	Transfer(from, to, sum int) error
}

func New(transferer Transferer, cacher cache.Cacher, logger mylogger.Logger) func(c *gin.Context) {
	return func(c *gin.Context) {
		var t storage.Transaction
		c.Bind(&t)

		uid, err := cacher.Cache(t)
		cached := true
		if err != nil {
			cached = false
		}

		err = transferer.Transfer(t.From, t.To, t.Sum)
		if err != nil {
			if cached {
				c.JSON(http.StatusInsufficientStorage, gin.H{
					"status": "cached, not stored",
				})
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "not cached, not stored",
			})
		}

		// удалить кеш
		err = cacher.Delete(uid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "unable to delete cache",
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "stored, cache deleted",
		})
	}
}
