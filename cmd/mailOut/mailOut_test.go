//2057008, 2624395, 9111696

package main

import (
	"os"
	"testing"
)

func init() {
	err := os.Chdir("../../")
	if err != nil {
		panic(err)
	}
}

func TestGetMailQueue(t *testing.T) {
	err := getMailQueue()
	if err != nil {
		t.Error()
	}

}
