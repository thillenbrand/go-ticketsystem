//2057008, 2624395, 9111696

package main

import "testing"

func TestGetMailQueue(t *testing.T) {
	err := getMailQueue()
	if err != nil {
		t.Error()
	}

}
