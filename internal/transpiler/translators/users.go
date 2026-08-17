package translators

import (
	"strconv"
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
	for _, g := range in.Groups {
		out.Passwd.Groups = append(out.Passwd.Groups, butane.Group{
			Name: g,
		})
	}

	for _, cu := range in.Users {
		bu := butane.User{
			Name:              cu.Name,
			SSHAuthorizedKeys: cu.SSHAuthorizedKeys,
		}

		if cu.Groups != "" {
			for _, rg := range strings.Split(cu.Groups, ",") {
				cleaned := strings.TrimSpace(rg)
				if cleaned != "" {
					bu.Groups = append(bu.Groups, cleaned)
				}
			}
		}

		if cu.UID != "" {
			uid, err := strconv.Atoi(cu.UID)
			if err == nil {
				bu.UID = &uid
			}
		}

		if cu.HashedPasswd != "" {
			bu.PasswordHash = cu.HashedPasswd
		}

		out.Passwd.Users = append(out.Passwd.Users, bu)
	}

	return nil
}
