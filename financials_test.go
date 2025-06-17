package ikentoo_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestFinancialsGet(t *testing.T) {
	req := client.NewFinancialsGetRequest()
	req.PathParams().BusinessID = businessID
	req.PathParams().From = time.Now().AddDate(0, 0, -1)
	req.PathParams().To = time.Now().AddDate(0, 0, -0)
	req.QueryParams().Include = append(req.QueryParams().Include, "payments")
	req.QueryParams().Include = append(req.QueryParams().Include, "payment")
	resp, err := req.Do(context.Background())
	if err != nil {
		t.Error(err)
	}

	b, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Println(string(b))
}
