package instance

import "testing"

func TestRun(t *testing.T) {
	err := Run("config.yaml", "https://www.youtube.com/channel/UCdgUTNVvHrcJq7K9nyEQ6qg")
	if err != nil {
		t.Fatal(err)
	}
}
