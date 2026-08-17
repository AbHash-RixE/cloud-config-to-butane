package cloudinit

type Config struct {
	Groups     []string     `yaml:"groups,omitempty"`
	Users      []User       `yaml:"users,omitempty"`
	WriteFiles []File       `yaml:"write_files,omitempty"`
	RunCmd     []RuncmdItem `yaml:"runcmd,omitempty"`
	BootCmd    []RuncmdItem `yaml:"bootcmd,omitempty"`
	NTP        *NTP         `yaml:"ntp,omitempty"`
	CACerts    *CACerts     `yaml:"ca_certs,omitempty"`
}

type User struct {
	Name              string   `yaml:"name"`
	SSHAuthorizedKeys []string `yaml:"ssh-authorized-keys,omitempty"`
	Sudo              string   `yaml:"sudo,omitempty"`
	Groups            string   `yaml:"groups,omitempty"`
	LockPasswd        *bool    `yaml:"lock_passwd,omitempty"`
	HashedPasswd      string   `yaml:"hashed_passwd,omitempty"`
	UID               string   `yaml:"uid,omitempty"`
}

type File struct {
	Path        string `yaml:"path"`
	Content     string `yaml:"content"`
	Owner       string `yaml:"owner,omitempty"`
	Permissions string `yaml:"permissions,omitempty"`
	Encoding    string `yaml:"encoding,omitempty"`
	Append      bool   `yaml:"append,omitempty"`
}

type RuncmdItem string

type NTP struct {
	Servers []string `yaml:"servers,omitempty"`
	Pools   []string `yaml:"pools,omitempty"`
}

type CACerts struct {
	RemoveDefaults bool     `yaml:"remove_defaults,omitempty"`
	Trusted        []string `yaml:"trusted,omitempty"`
}
