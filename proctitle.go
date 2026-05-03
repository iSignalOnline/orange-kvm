//go:build !orangepi

package kvm

import (
	"fmt"

	"github.com/erikdubbelboer/gspt"
)

func setProcTitle(status string) {
	if status != "" {
		status = " " + status
	}
	title := fmt.Sprintf("%s%s", procPrefix, status)
	gspt.SetProcTitle(title)
}
