package translators

import (
	"fmt"
	"strings"

	"github.com/AbHash-RixE/cloudinit-to-butane/internal/butane"
	"github.com/AbHash-RixE/cloudinit-to-butane/internal/cloudinit"
)

type NTPTranslator struct{}

func NewNTPTranslator() *NTPTranslator {
	return &NTPTranslator{}
}

func (n *NTPTranslator) Name() string {
	return "ntp"
}

func (n *NTPTranslator) Translate(in *cloudinit.Config, out *butane.Config) error {
	if in.NTP == nil {
		return nil
	}

	var timesyncLines []string
	timesyncLines = append(timesyncLines, "[Time]")
	for _, s := range in.NTP.Servers {
		timesyncLines = append(timesyncLines, fmt.Sprintf("NTP=%s", s))
	}
	for _, p := range in.NTP.Pools {
		timesyncLines = append(timesyncLines, fmt.Sprintf("Pool=%s", p))
	}

	dropinContent := strings.Join(timesyncLines, "\n")

	out.Systemd.Units = append(out.Systemd.Units, butane.Unit{
		Name:     "systemd-timesyncd.service",
		Contents: dropinContent,
	})

	return nil
}
