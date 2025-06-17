package ikentoo_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	ikentoo "github.com/omniboost/go-ikentoo"
)

func TestDailyFinancialsGet(t *testing.T) {
	req := client.NewDailyFinancialsGetRequest()
	req.PathParams().BusinessID = businessID
	req.QueryParams().Date = ikentoo.Date{time.Now().AddDate(0, 0, -1)}
	req.QueryParams().Include = append(req.QueryParams().Include, "payments")
	req.QueryParams().Include = append(req.QueryParams().Include, "payment")
	resp, err := req.Do(context.Background())
	if err != nil {
		t.Error(err)
	}

	b, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Println(string(b))
}
