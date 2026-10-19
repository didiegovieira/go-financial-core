package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/didiegovieira/go-financial-core/apps/core/internal/domain/entity"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/infrastructure/database/postgres/mapper"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/infrastructure/database/postgres/model"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/infrastructure/database/postgres/query"
	"github.com/didiegovieira/go-financial-core/pkg/telemetry"
)

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) Create(ctx context.Context, account entity.Account) (err error) {
	ctx, finish := telemetry.TraceWithError(ctx, "Repository.Account.Create", &err)
	defer finish()

	accountModel := mapper.AccountToModel(account)

	_, err = r.db.ExecContext(
		ctx,
		query.AccountCreate,
		accountModel.ID,
		accountModel.Name,
		accountModel.Email,
		accountModel.Password,
		accountModel.Phone,
		accountModel.DocumentType,
		accountModel.DocumentNumber,
		accountModel.DocumentCountry,
	)
	return err
}

func (r *AccountRepository) FindByDocument(ctx context.Context, account entity.Account) (foundAccount *entity.Account, err error) {
	ctx, finish := telemetry.TraceWithError(ctx, "Repository.Account.FindByDocument", &err)
	defer finish()

	accountModel := mapper.AccountToModel(account)
	var found model.AccountModel

	err = r.db.QueryRowContext(
		ctx,
		query.AccountFindByDocument,
		accountModel.DocumentNumber,
		accountModel.DocumentType,
		accountModel.DocumentCountry,
	).Scan(
		&found.ID,
		&found.Name,
		&found.Email,
		&found.Password,
		&found.Phone,
		&found.DocumentType,
		&found.DocumentNumber,
		&found.DocumentCountry,
		&found.CreatedAt,
		&found.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return mapper.AccountToEntity(found), nil
}

func (r *AccountRepository) FindByPhone(ctx context.Context, account entity.Account) (foundAccount *entity.Account, err error) {
	ctx, finish := telemetry.TraceWithError(ctx, "Repository.Account.FindByPhone", &err)
	defer finish()

	accountModel := mapper.AccountToModel(account)
	var found model.AccountModel

	err = r.db.QueryRowContext(
		ctx,
		query.AccountFindByPhone,
		accountModel.Phone,
	).Scan(
		&found.ID,
		&found.Name,
		&found.Email,
		&found.Password,
		&found.Phone,
		&found.DocumentType,
		&found.DocumentNumber,
		&found.DocumentCountry,
		&found.CreatedAt,
		&found.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return mapper.AccountToEntity(found), nil
}

func (r *AccountRepository) DeleteByPhone(ctx context.Context, account entity.Account) (err error) {
	ctx, finish := telemetry.TraceWithError(ctx, "Repository.Account.DeleteByPhone", &err)
	defer finish()

	accountModel := mapper.AccountToModel(account)

	_, err = r.db.ExecContext(
		ctx,
		query.AccountDeleteByPhone,
		accountModel.Phone,
	)

	return err
}
