package grid

import (
	"testing"
)

func TestValidateAddress(t *testing.T) {
	resp, err := ValidateAddress("TJLLCqUWQN3LPgm1AwJm4EMFnwe28uMNAS")
	if err != nil {
		t.Error(err)
	}
	t.Logf("response: %+v", resp)
}
