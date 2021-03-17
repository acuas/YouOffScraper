package instance

import "testing"

func TestRun(t *testing.T) {
	err := Run("config.yaml", "https://www.youtube.com/channel/UCXMYnbGeoxCdhazKPDPk7DQ")
	if err != nil {
		t.Fatal(err)
	}
}
