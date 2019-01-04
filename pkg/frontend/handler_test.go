package frontend

import (
	"os"
	"testing"
)

func TestOpenTickets(t *testing.T) {
	err := os.Chdir("../../")
	if err != nil {
		panic(err)
	}
	var tickets = OpenTickets()
	if len(tickets) == 0 {
		t.Error()
	}
}
