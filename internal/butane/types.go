package butane

type Config struct {
	Variant string  `yaml:"variant"`
	Version string  `yaml:"version"`
	Passwd  Passwd  `yaml:"passwd,omitempty"`
	Storage Storage `yaml:"storage,omitempty"`
	Systemd Systemd `yaml:"systemd,omitempty"`
}

type Passwd struct {
	Groups []Group `yaml:"groups,omitempty"`
	Users  []User  `yaml:"users,omitempty"`
}

type Group struct {
	Name string `yaml:"name"`
}

type User struct {
	Name              string   `yaml:"name"`
	SSHAuthorizedKeys []string `yaml:"ssh_authorized_keys,omitempty"`
	Groups            []string `yaml:"groups,omitempty"`
}

type Storage struct {
	Files []File `yaml:"files,omitempty"`
}

type File struct {
	Path     string    `yaml:"path"`
	Mode     *int      `yaml:"mode,omitempty"`
	User     *NodeUser `yaml:"user,omitempty"`
	Contents Content   `yaml:"contents"`
	Append   []Append  `yaml:"append,omitempty"`
}

type NodeUser struct {
	Name string `yaml:"name"`
}

type Content struct {
	Inline string `yaml:"inline,omitempty"`
}

type Append struct {
	Inline string `yaml:"inline"`
}

type Systemd struct {
	Units []Unit `yaml:"units,omitempty"`
}

type Unit struct {
	Name     string `yaml:"name"`
	Enabled  *bool  `yaml:"enabled,omitempty"`
	Contents string `yaml:"contents,omitempty"`
}

// initializes a Butane config
func NewDefaultConfig() *Config {
	//default data
	return &Config{
		Variant: "flatcar",
		Version: "1.1.0",
	}
}
