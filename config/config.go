package config

import (
	"github.com/fsnotify/fsnotify"
	"gopkg.in/ini.v1"
)

const (
	LOCAL  = "local"
	SERVER = "server"
)

type Config struct {
	Name string
	// env local or server
	Env string `ini:"env"`
	// ssh host
	Host string `ini:"host"`
	// ssh user
	User string `ini:"user"`
	// ssh password
	Password string `ini:"password"`
	// ssh port
	Port int `ini:"port"`
	// 本地目录
	LocalPath string `ini:"local_path"`
	// 部署目录
	DeploymentPath string `ini:"deployment_path"`
	// 需要执行的命令
	Cmd []string `ini:"cmd,omitempty,allowshadow"`
	// 忽略前缀&后缀
	IgnorePrefix []string `ini:"ignore_prefix"`
	IgnoreSuffix []string `ini:"ignore_suffix"`

	// 监控事件(未使用)
	SyncEvent fsnotify.Op
}

func GetConfigsByIni(source string) map[string]*Config {
	cfg, err := ini.ShadowLoad(source)
	if err != nil {
		panic("配置文件路径不正确")
	}
	name := make(map[string]bool)
	configs := make(map[string]*Config)
	section := cfg.Sections()
	for _, s := range section {
		if s.Name() != ini.DefaultSection {
			n := &Config{}
			err := cfg.Section(s.Name()).MapTo(n)
			if err != nil {
				panic(err)
			}

			// 是否重复定义
			if _, ok := name[s.Name()]; ok {
				panic("配置文件[" + s.Name() + "]重复定义")
			}
			name[s.Name()] = true

			// 配置名称
			n.Name = s.Name()

			// 修改为全路径
			if len(n.IgnorePrefix) > 0 {
				for i, _ := range n.IgnorePrefix {
					n.IgnorePrefix[i] = n.LocalPath + "/" + n.IgnorePrefix[i]
				}
			}

			configs[n.LocalPath] = n
		}
	}

	return configs
}
