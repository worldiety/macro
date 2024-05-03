package domain

import (
	"encoding/json"
	"testing"
)

func TestComponent_MarshalJSON(t *testing.T) {
	compo := Component{}.WithButton(Button{Caption: "hello"})
	buf, err := json.Marshal(compo)
	if err != nil {
		t.Fatalf("%v", err)
	}

	var compo2 Component
	if err := json.Unmarshal(buf, &compo2); err != nil {
		t.Fatalf("%v", err)
	}

	if btn, _ := compo2.AsButton(); btn.Caption != "hello" {
		t.Fatalf("%v", btn.Caption)
	}
}
