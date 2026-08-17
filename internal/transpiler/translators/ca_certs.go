package translators

import (
	"fmt"

	"github.com/AbHash-RixE/cloudinit-to-butane/internal/butane"
	"github.com/AbHash-RixE/cloudinit-to-butane/internal/cloudinit"
)

type CACertsTranslator struct{}

func NewCACertsTranslator() *CACertsTranslator {
	return &CACertsTranslator{}
}

func (c *CACertsTranslator) Name() string {
	return "ca_certs"
}

func (c *CACertsTranslator) Translate(in *cloudinit.Config, out *butane.Config) error {
	if in.CACerts == nil {
		return nil
	}

	for i, cert := range in.CACerts.Trusted {
		mode := 420 // 0644
		out.Storage.Files = append(out.Storage.Files, butane.File{
			Path: fmt.Sprintf("/etc/ssl/certs/c2b-ca-%d.pem", i),
			Mode: &mode,
			Contents: butane.Content{
				Inline: cert,
			},
		})
	}

	return nil
}
