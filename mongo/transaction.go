package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo/options"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	MaxTransactionRetryTimes = 20
	MaxCommitRetryTimes      = 5
)

// CommitWithRetry is an example function demonstrating transaction commit with retry logic.
func CommitWithRetry(sctx mongo.SessionContext) (err error) {
	for i := 0; i < MaxCommitRetryTimes; i++ {
		err = sctx.CommitTransaction(sctx)
		switch e := err.(type) {
		case nil:
			log.Infof("Transaction committed.")
			return nil
		case mongo.CommandError:
			if e.HasErrorLabel("UnknownTransactionCommitResult") {
				log.Infof("UnknownTransactionCommitResult, retrying commit operation...")
				continue
			}
			log.Errorf("Error during commit...")
			return e
		default:
			log.Errorf("Error during commit...")
			return e
		}
	}
	if err != nil {
		log.Errorf("Transaction commit failed after %v retries, transaction failed.", MaxCommitRetryTimes)
		return err
	}
	return nil
}

// RunTransactionWithRetry is an example function demonstrating transaction retry logic.
func RunTransactionWithRetry(ctx context.Context, client *mongo.Client, opts *options.SessionOptions, txnFn func(mongo.SessionContext) error) error {
	return client.UseSessionWithOptions(ctx, opts, func(sctx mongo.SessionContext) (err error) {
		for i := 0; i < MaxTransactionRetryTimes; i++ {
			err = txnFn(sctx) // Performs transaction.
			if err == nil {
				return nil
			}
			// If transient error, retry the whole transaction
			if cmdErr, ok := err.(mongo.CommandError); ok && cmdErr.HasErrorLabel("TransientTransactionError") {
				log.Println("TransientTransactionError, retrying transaction...")
				continue
			}
			return err
		}
		log.Errorf("Transaction failed after %v retries : %v", MaxTransactionRetryTimes, err)
		return err
	})
}
