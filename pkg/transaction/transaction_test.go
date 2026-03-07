package transaction

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
)

func TestFromContext(t *testing.T) {
	t.Parallel()

	t.Run("has tx in context", func(t *testing.T) {
		t.Parallel()

		tx := &testTx{}
		db := &testTx{}

		ctx := context.Background()
		ctx = context.WithValue(ctx, txKey, tx)

		result := FromContext(ctx, db)
		require.Same(t, tx, result)
		require.NotSame(t, db, result)
	})

	t.Run("has no tx in context", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		db := &testTx{}

		result := FromContext(ctx, db)
		require.Same(t, db, result)
	})
}

func TestExec(t *testing.T) {
	t.Parallel()

	expectedErr := fmt.Errorf("expected error")

	tests := []struct {
		name        string
		tx          postgres
		db          postgres
		callback    func(context.Context) error
		expectedErr error
	}{
		{
			name:        "tx already in context, success",
			tx:          &testTx{},
			callback:    func(_ context.Context) error { return nil },
			expectedErr: nil,
		},
		{
			name:        "tx already in context, callback error",
			tx:          &testTx{},
			callback:    func(_ context.Context) error { return expectedErr },
			expectedErr: expectedErr,
		},
		{
			name:     "no tx, success",
			db:       &testTx{},
			callback: func(_ context.Context) error { return nil },
		},
		{
			name:        "no tx, begin err",
			db:          &testTx{beginErr: expectedErr},
			callback:    func(_ context.Context) error { return nil },
			expectedErr: expectedErr,
		},
		{
			name:        "no tx, callback err",
			db:          &testTx{},
			callback:    func(_ context.Context) error { return expectedErr },
			expectedErr: expectedErr,
		},
		{
			name:        "no tx, commit err",
			db:          &testTx{commitErr: expectedErr},
			callback:    func(_ context.Context) error { return nil },
			expectedErr: expectedErr,
		},
		{
			name:        "no tx, rollback err",
			db:          &testTx{rollbackErr: expectedErr},
			callback:    func(_ context.Context) error { return nil },
			expectedErr: expectedErr,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			if tc.tx != nil {
				ctx = context.WithValue(ctx, txKey, tc.tx)
			}

			err := NewExecutor(tc.db).Exec(ctx, tc.callback)
			require.ErrorIs(t, err, tc.expectedErr)
		})
	}
}

type testTx struct {
	beginErr    error
	commitErr   error
	rollbackErr error
}

func (t *testTx) Begin(_ context.Context) (pgx.Tx, error) {
	if t.beginErr != nil {
		return nil, t.beginErr
	}

	return t, nil
}

func (t *testTx) Commit(_ context.Context) error {
	if t.commitErr != nil {
		return t.commitErr
	}

	return nil
}

func (t *testTx) Rollback(_ context.Context) error {
	if t.rollbackErr != nil {
		return t.rollbackErr
	}

	return nil
}

func (t *testTx) CopyFrom(_ context.Context, _ pgx.Identifier, _ []string, _ pgx.CopyFromSource) (int64, error) {
	return 0, nil
}

func (t *testTx) SendBatch(_ context.Context, b *pgx.Batch) pgx.BatchResults {
	return nil
}

func (t *testTx) LargeObjects() pgx.LargeObjects {
	return pgx.LargeObjects{}
}

func (t *testTx) Prepare(_ context.Context, name, sql string) (*pgconn.StatementDescription, error) {
	return nil, nil
}

func (t *testTx) Exec(_ context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error) {
	return pgconn.CommandTag{}, nil
}

func (t *testTx) Query(_ context.Context, sql string, args ...any) (pgx.Rows, error) {
	return nil, nil
}

func (t *testTx) QueryRow(_ context.Context, sql string, args ...any) pgx.Row {
	return nil
}

func (t *testTx) Conn() *pgx.Conn {
	return nil
}
