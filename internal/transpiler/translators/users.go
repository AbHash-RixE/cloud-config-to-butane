package translators

import (
	"strings"

	"github.com/AbHash-RixE/cloudinit-to-butane/internal/butane"
	"github.com/AbHash-RixE/cloudinit-to-butane/internal/cloudinit"
)

type UserTranslator struct{}

func NewUserTranslator() *UserTranslator {
	return &UserTranslator{}
}

func (u *UserTranslator) Name() string {
	return "users"
}

func (u *UserTranslator) Translate(in *cloudinit.Config, out *butane.Config) error {
	//Translate top-level groups
	for _, g := range in.Groups {
		out.Passwd.Groups = append(out.Passwd.Groups, butane.Group{
			Name: g,
		})
	}

	//Translate users
	for _, cu := range in.Users {
		bu := butane.User{
			Name:              cu.Name,
			SSHAuthorizedKeys: cu.SSHAuthorizedKeys,
		}

		// Cloud-init groups -> comma separated string, but
		// Butane group -> []string array
		if cu.Groups != "" {
			rawGroups := strings.Split(cu.Groups, ",")
			for _, rg := range rawGroups {
				cleaned := strings.TrimSpace(rg)
				if cleaned != "" {
					bu.Groups = append(bu.Groups, cleaned)
				}
			}
		}

		out.Passwd.Users = append(out.Passwd.Users, bu)
	}

	return nil
}
