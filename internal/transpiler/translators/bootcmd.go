package translators

import (
	"strings"

	"github.com/AbHash-RixE/cloudinit-to-butane/internal/butane"
	"github.com/AbHash-RixE/cloudinit-to-butane/internal/cloudinit"
)

type BootCmdTranslator struct{}

func NewBootCmdTranslator() *BootCmdTranslator {
	return &BootCmdTranslator{}
}

func (b *BootCmdTranslator) Name() string {
	return "bootcmd"
}

func (b *BootCmdTranslator) Translate(in *cloudinit.Config, out *butane.Config) error {
	if len(in.BootCmd) == 0 {
		return nil
	}

	var lines []string
	lines = append(lines, "#!/bin/sh")
	lines = append(lines, "set -e")
	for _, cmd := range in.BootCmd {
		lines = append(lines, string(cmd))
	}

	script := strings.Join(lines, "\n")
	mode := 448 // 0700

	out.Storage.Files = append(out.Storage.Files, butane.File{
		Path: "/opt/c2b/bootcmd.sh",
		Mode: &mode,
		Contents: butane.Content{
			Inline: script,
		},
	})

	enabled := true
	out.Systemd.Units = append(out.Systemd.Units, butane.Unit{
		Name:    "c2b-bootcmd.service",
		Enabled: &enabled,
		Contents: `[Unit]
Description=Execute cloud-init bootcmd
DefaultDependencies=no
After=local-fs.target
Before=network.target

[Service]
Type=oneshot
ExecStart=/opt/c2b/bootcmd.sh

[Install]
WantedBy=multi-user.target`,
	})

	return nil
}
