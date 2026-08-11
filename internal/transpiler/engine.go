package transpiler

import (
	"fmt"

	"github.com/AbHash-RixE/cloudinit-to-butane/internal/butane"
	"github.com/AbHash-RixE/cloudinit-to-butane/internal/cloudinit"
)

// stores all different translators for different config parts
// since each require its own translator
type Engine struct {
	translators []Translator
}

// returns a pointer for Engine struct
func NewEngine(t ...Translator) *Engine {
	return &Engine{
		translators: t,
	}
}

func (e *Engine) Run(in *cloudinit.Config) (*butane.Config, error) {
	//initialise butane config
	out := butane.NewDefaultConfig()

	for _, t := range e.translators {
		if err := t.Translate(in, out); err != nil {
			// Prints error + translator name
			return nil, fmt.Errorf("module %q failed: %w", t.Name(), err)
		}
	}

	return out, nil
}
