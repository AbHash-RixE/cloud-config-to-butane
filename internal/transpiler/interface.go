package transpiler

import (
	"github.com/AbHash-RixE/cloudinit-to-butane/internal/butane"
	"github.com/AbHash-RixE/cloudinit-to-butane/internal/cloudinit"
)

type Translator interface {
	Name() string

	//this func. will translate struct from cloudconfig to butane
	Translate(in *cloudinit.Config, out *butane.Config) error
}
