package domain

import "errors"

var (
	ErrEmptyBatch           = errors.New("batch has no authorizations")
	ErrDuplicateDetected    = errors.New("duplicate authorization detected")
	ErrReconciliationFailed = errors.New("reconciliation failed")
)
