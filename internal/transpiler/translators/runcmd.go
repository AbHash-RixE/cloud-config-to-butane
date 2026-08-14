package translators

import (
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

	var scriptLines []string
	scriptLines = append(scriptLines, "#!/bin/sh")
	// Fail immediately if any command fails.
	scriptLines = append(scriptLines, "set -e")
	for _, cmd := range in.RunCmd {
		scriptLines = append(scriptLines, string(cmd))
	}

	script := strings.Join(scriptLines, "\n")

	//write script to disk
	mode := 448 // 0700(decimal): owner can read, write, and execute
	out.Storage.Files = append(out.Storage.Files, butane.File{
		Path: "/opt/c2b/runcmd.sh",
		Mode: &mode,
		Contents: butane.Content{
			Inline: script,
		},
	})

	// Systemd unit string.
	unitContents := `[Unit]
Description=Execute cloud-init runcmd
Wants=network-online.target
After=network-online.target
ConditionPathExists=!/var/lib/c2b-runcmd.success

[Service]
Type=oneshot
ExecStart=/opt/c2b/runcmd.sh
ExecStartPost=/bin/touch /var/lib/c2b-runcmd.success

[Install]
WantedBy=multi-user.target`

	enabled := true

	out.Systemd.Units = append(out.Systemd.Units, butane.Unit{
		Name:     "c2b-runcmd.service",
		Enabled:  &enabled,
		Contents: unitContents,
	})

	return nil
}
