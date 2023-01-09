package domain

// This is the root yaml config that contains the information needed to operate
type Config struct {
	Daemon      ConfigDaemon `yaml:"Daemon,omitempty"`
	Backup      BackupConfig `yaml:"Backup"`
	Alert       ConfigAlert  `yaml:"Alert,omitempty"`
	Destination ConfigDest   `yaml:"Destination,omitempty"`
}

type ConfigDaemon struct {
	Cron string `yaml:"Cron,omitempty"`
}

type BackupConfig struct {
	Docker []ContainerDocker `yaml:"Docker,omitempty"`
}

type ContainerDocker struct {
	Name      string              `yaml:"Name"`
	Directory string              `yaml:"Directory"`
	Tar       ConfigContainerTar  `yaml:"Tar"`
	Post      ConfigContainerPost `yaml:"Post,omitempty"`
}

type ConfigContainerTar struct {
	UseDate   bool   `yaml:"UseDate,omitempty"`
	Pattern   string `yaml:"Pattern,omitempty"`
	Directory string `yaml:"Directory,omitempty"`
}

type ConfigContainerPost struct {
	Reboot []string `yaml:"Reboot,omitempty"`
}

type ConfigDest struct {
	Retain ConfigRetain    `yaml:"Retain,omitempty"`
	Local  ConfigDestLocal `yaml:"Local,omitempty"`
	Sftp   ConfigDestSftp  `yaml:"Sftp,omitempty"`
}

// Defines how long backups should be retained
type ConfigRetain struct {
	Days int `yaml:"Days,omitempty"`
}

// Defines where and how to move data
type ConfigDestLocal struct {
	Path string `yaml:"Path,omitempty"`
}

type ConfigDestSftp struct {
	Path     string `yaml:"Path"`
	Server   string `yaml:"Server"`
	Username string `yaml:"Username"`
	Password string `yaml:"Password"`
}

type ConfigAlert struct {
	SendOnlyOnError bool               `yaml:"SendOnlyOnError,omitempty"`
	Discord         ConfigAlertDiscord `yaml:"Discord,omitempty"`
	Email           ConfigAlertEmail   `yaml:"Email,omitempty"`
}

type ConfigAlertDiscord struct {
	Username string   `yaml:"Username,omitempty"`
	Webhooks []string `yaml:"Webhooks,omitempty"`
}

type ConfigAlertEmail struct {
	Account ConfigAlertEmailAccount `yaml:"Account"`
	From    string                  `yaml:"From"`
	To      string                  `yaml:"To"`
}

type ConfigAlertEmailAccount struct {
	Username string `yaml:"Username"`
	Password string `yaml:"Password"`
	Host     string `yaml:"Host"`
	Port     int    `yaml:"Port"`
	UseTls   bool   `yaml:"UseTls"`
}
