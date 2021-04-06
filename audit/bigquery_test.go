package audit

import "testing"

func TestDatasetFormat(t *testing.T) {
	bq := NewBigQueryClient("dataset-with_hyphen")

	if bq.dataset != "dataset_with_hyphen" {
		t.Errorf("Wrong format for BQ Client: expected 'dataset-with_hyphen' got '%s'", bq.dataset)
	}
}
