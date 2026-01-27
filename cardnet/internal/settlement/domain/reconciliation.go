package domain

type ReconciliationStatus string

const (
	Reconciled ReconciliationStatus = "Reconciled"
	Failed     ReconciliationStatus = "FAILED"
)

type ReconciliationIssue struct {
	BatchID string
	Reason  string
}
