package ikentoo_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
)

func TestBusinessesGet(t *testing.T) {
	req := client.NewBusinessesGetRequest()
	resp, err := req.Do(context.Background())
	if err != nil {
		t.Error(err)
	}

	b, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Println(string(b))
}
