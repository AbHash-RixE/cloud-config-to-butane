package translators

import (
	"fmt"
	"strings"

	"github.com/AbHash-RixE/cloudinit-to-butane/internal/butane"
	"github.com/AbHash-RixE/cloudinit-to-butane/internal/cloudinit"
)

type RunCmdTranslator struct{}

func NewRunCmdTranslator() *RunCmdTranslator {
	return &RunCmdTranslator{}
}

func (r *RunCmdTranslator) Name() string {
	return "runcmd"
}

func (r *RunCmdTranslator) Translate(in *cloudinit.Config, out *butane.Config) error {
	if len(in.RunCmd) == 0 {
		return nil
	}

	var commands []string
	for _, cmd := range in.RunCmd {
		commands = append(commands, string(cmd)) // Cast to string
	}

	script := strings.Join(commands, "\n")

	//Systemd unit string.
	unitContents := fmt.Sprintf(`[Unit]
Description=Execute cloud-init runcmd
Wants=network-online.target
After=network-online.target
ConditionPathExists=!/var/lib/c2b-runcmd.success

[Service]
Type=oneshot
ExecStart=/bin/sh -c %q
ExecStartPost=/bin/touch /var/lib/c2b-runcmd.success

[Install]
WantedBy=multi-user.target`, script)

	enabled := true

	out.Systemd.Units = append(out.Systemd.Units, butane.Unit{
		Name:     "kubeadm-runcmd.service",
		Enabled:  &enabled,
		Contents: unitContents,
	})

	return nil
}
