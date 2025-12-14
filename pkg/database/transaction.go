package database

import (
	"context"
	"database/sql"
)

// contextKey 用於在 context 中存儲 transaction
type contextKey string

const txKey contextKey = "db_transaction"

// WithTransaction 執行一個事務操作
// 如果 fn 返回 error,會自動 rollback
// 如果 fn 成功執行完畢,會自動 commit
func WithTransaction(db *sql.DB, fn func(context.Context) error) error {
	// 1. 開始事務
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	// 2. 將 transaction 放入 context
	ctx := context.WithValue(context.Background(), txKey, tx)

	// 3. 使用 defer 確保事務一定會被處理(commit 或 rollback)
	defer func() {
		if p := recover(); p != nil {
			// 如果發生 panic,rollback 並繼續 panic
			tx.Rollback()
			panic(p)
		} else if err != nil {
			// 如果有錯誤,rollback
			tx.Rollback()
		} else {
			// 成功則 commit
			err = tx.Commit()
		}
	}()

	// 4. 執行業務邏輯
	err = fn(ctx)
	return err
}

// GetTx 從 context 中取得 transaction
// 如果 context 中沒有 transaction,返回 nil, false
func GetTx(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(txKey).(*sql.Tx)
	return tx, ok
}
