package zep_nakama

import (
	"context"
	"database/sql"
	"github.com/heroiclabs/nakama-common/runtime"
	"os"
	"time"
)

func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	initStart := time.Now()

	err := initializer.RegisterRpc("checksum", RpcCheckSum)
	if err != nil {
		return err
	}
	basePath := os.Getenv("FILE_BASE_PATH")
	logger.Info("basePath is %s", basePath)
	logger.Info("Module loaded in %dms", time.Since(initStart).Milliseconds())
	return nil
}
