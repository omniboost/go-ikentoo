package ikentoo_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
)

func TestPaymentMethodsGet(t *testing.T) {
	req := client.NewPaymentMethodsGetRequest()
	req.PathParams().BusinessID = businessID
	resp, err := req.Do(context.Background())
	if err != nil {
		t.Error(err)
	}

	b, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Println(string(b))
}
