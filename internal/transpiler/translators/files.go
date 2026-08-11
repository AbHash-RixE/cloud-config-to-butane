package translators

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/AbHash-RixE/cloudinit-to-butane/internal/butane"
	"github.com/AbHash-RixE/cloudinit-to-butane/internal/cloudinit"
)

type FileTranslator struct{}

func NewFileTranslator() *FileTranslator {
	return &FileTranslator{}
}

func (f *FileTranslator) Name() string {
	return "write_files"
}

func (f *FileTranslator) Translate(in *cloudinit.Config, out *butane.Config) error {
	for _, cf := range in.WriteFiles {
		bf := butane.File{
			Path: cf.Path,
		}

		// Handle Content and Base64 Encoding
		content := cf.Content
		if strings.ToLower(cf.Encoding) == "b64" || strings.ToLower(cf.Encoding) == "base64" {
			decoded, err := base64.StdEncoding.DecodeString(cf.Content)
			if err != nil {
				return fmt.Errorf("failed to decode base64 content for file %s: %w", cf.Path, err)
			}
			content = string(decoded)
		}

		// Handle appending OR overwriting
		if cf.Append {
			//since Append , expects []Append type
			bf.Append = []butane.Append{{Inline: content}}
		} else {
			//expects Content type
			bf.Contents = butane.Content{Inline: content}
		}

		// Handle Permissions (Cloud-init string -> Butane octal int)
		if cf.Permissions != "" {
			modeInt, err := strconv.ParseInt(cf.Permissions, 8, 32)
			if err != nil {
				return fmt.Errorf("invalid permissions %q for file %s: %w", cf.Permissions, cf.Path, err)
			}
			//int64 -> int
			m := int(modeInt)
			bf.Mode = &m
		}

		// Handle Ownership
		if cf.Owner != "" {
			parts := strings.Split(cf.Owner, ":")
			bf.User = &butane.NodeUser{Name: parts[0]}
		}

		out.Storage.Files = append(out.Storage.Files, bf)
	}

	return nil
}
