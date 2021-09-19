package controller

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"

	"github.com/cockroachdb/errors"
	"github.com/gotway/gotway/pkg/log"
	"github.com/mmontes11/bankroach/internal/config"
	"github.com/mmontes11/bankroach/internal/model"
	"github.com/mmontes11/bankroach/internal/repository"
	"github.com/mmontes11/bankroach/pkg/cockroachdb"
)

type controller struct {
	repo   repository.AccountRepo
	db     *cockroachdb.CRDB
	config config.Config
	logger log.Logger
}

type Controller interface {
	CreateAccount(ctx context.Context, balance int64) (model.AccountID, error)
	PrintBalances(ctx context.Context) error
	Transfer(ctx context.Context) error
}

func (c *controller) CreateAccount(ctx context.Context, balance int64) (model.AccountID, error) {
	var account model.AccountID
	err := c.db.ExecuteTx(ctx, nil, func(tx *sql.Tx) error {
		var err error
		account, err = c.repo.Insert(ctx, tx, balance)
		return err
	})
	if err != nil {
		return 0, err
	}
	return account, nil
}

func (c *controller) PrintBalances(ctx context.Context) error {
	return c.db.ExecuteTx(ctx, nil, func(tx *sql.Tx) error {
		accounts, err := c.repo.List(ctx, tx)
		if err != nil {
			return err
		}
		c.logger.Infof("\nAccounts:\n")
		for _, account := range accounts {
			c.logger.Info(account)
		}
		return nil
	})
}

func (c *controller) Transfer(ctx context.Context) error {
	return c.db.ExecuteTx(ctx, nil, func(tx *sql.Tx) error {
		balance := rand.Int63n(int64(c.config.InitialBalance))

		source, err := c.repo.FindWithMinBalance(ctx, tx, balance)
		if err != nil {
			return errors.Wrap(
				err,
				fmt.Sprintf("error finding source acount with %d balance", balance),
			)
		}

		accounts, err := c.repo.List(ctx, tx)
		if err != nil {
			return errors.Wrap(
				err,
				"error listing accounts",
			)
		}
		destination := accounts[rand.Intn(len(accounts))]

		err = c.repo.Transfer(ctx, tx, repository.TransferParams{
			Source:      source,
			Destination: destination.ID,
			Amount:      balance,
		})
		if err != nil {
			return errors.Wrap(
				err,
				fmt.Sprintf(
					"error transferring %d from %d to %d",
					balance,
					source,
					destination.ID,
				),
			)
		}

		return nil
	})
}

func New(
	repo repository.AccountRepo,
	db *cockroachdb.CRDB,
	config config.Config,
	logger log.Logger,
) Controller {
	return &controller{
		repo:   repo,
		db:     db,
		config: config,
		logger: logger,
	}
}
