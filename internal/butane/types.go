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
	GID  *int   `yaml:"gid,omitempty"`
}

type User struct {
	Name              string   `yaml:"name"`
	SSHAuthorizedKeys []string `yaml:"ssh_authorized_keys,omitempty"`
	Groups            []string `yaml:"groups,omitempty"`
	PasswordHash      string   `yaml:"password_hash,omitempty"`
	UID               *int     `yaml:"uid,omitempty"`
	HomeDir           string   `yaml:"home_dir,omitempty"`
	Shell             string   `yaml:"shell,omitempty"`
}

type Storage struct {
	Files []File `yaml:"files,omitempty"`
}

type File struct {
	Path      string     `yaml:"path"`
	Mode      *int       `yaml:"mode,omitempty"`
	User      *NodeUser  `yaml:"user,omitempty"`
	Group     *NodeGroup `yaml:"group,omitempty"`
	Contents  Content    `yaml:"contents"`
	Append    []Append   `yaml:"append,omitempty"`
	Overwrite *bool      `yaml:"overwrite,omitempty"`
}

type NodeUser struct {
	Name string `yaml:"name"`
	ID   *int   `yaml:"id,omitempty"`
}

type NodeGroup struct {
	Name string `yaml:"name"`
	ID   *int   `yaml:"id,omitempty"`
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
	Mask     *bool  `yaml:"mask,omitempty"`
	Contents string `yaml:"contents,omitempty"`
}

func NewDefaultConfig() *Config {
	return &Config{
		Variant: "flatcar",
		Version: "1.1.0",
	}
}
