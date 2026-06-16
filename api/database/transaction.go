package database

import "gorm.io/gorm"

// Tx 执行无返回值事务
func Tx(db *gorm.DB, txFunc func(tx *gorm.DB) error) error {
	return db.Transaction(txFunc)
}

// TxResult 执行带返回值的事务（泛型）
func TxResult[T any](db *gorm.DB, txFunc func(tx *gorm.DB) (T, error)) (T, error) {
	var zero T
	var result T
	err := db.Transaction(func(tx *gorm.DB) error {
		var txErr error
		result, txErr = txFunc(tx)
		return txErr
	})
	if err != nil {
		return zero, err
	}
	return result, nil
}