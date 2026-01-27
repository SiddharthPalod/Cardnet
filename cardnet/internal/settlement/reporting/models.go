package reporting

type BatchReport struct {
	BatchID     string
	BatchDate   string
	Status      string
	TotalCount  int
	TotalAmount int64
	Currency    string
}
