package cloudinit

type Config struct {
	Groups     []string     `yaml:"groups,omitempty"`
	Users      []User       `yaml:"users,omitempty"`
	WriteFiles []File       `yaml:"write_files,omitempty"`
	RunCmd     []RuncmdItem `yaml:"runcmd,omitempty"`
}

type User struct {
	Name              string   `yaml:"name"`
	SSHAuthorizedKeys []string `yaml:"ssh-authorized-keys,omitempty"`
	Sudo              string   `yaml:"sudo,omitempty"`
	Groups            string   `yaml:"groups,omitempty"`
}

type File struct {
	Path        string `yaml:"path"`
	Content     string `yaml:"content"`
	Owner       string `yaml:"owner,omitempty"`
	Permissions string `yaml:"permissions,omitempty"`
	Encoding    string `yaml:"encoding,omitempty"`
	Append      bool   `yaml:"append,omitempty"`
}

// custom type for runcmd: string AND list type
type RuncmdItem string
