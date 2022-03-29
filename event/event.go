package event

import (
	"code_sync_tool/config"
	log "code_sync_tool/log"
	"code_sync_tool/sync"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"os/exec"
	"strings"
	"time"
)

var (
	// 发送CMD命令超时时间
	signalTimeOut = time.Millisecond * 200
)

type FileEvent struct {
	Event fsnotify.Op
	Name  string
}

type Event struct {
	fileEvent chan FileEvent
	configs   map[string]*config.Config
	sFTP      map[string]*sync.SFTP
}

func Run(cs map[string]*config.Config, f chan FileEvent) error {
	e := &Event{
		fileEvent: f,
		configs:   cs,
		sFTP:      make(map[string]*sync.SFTP, 2),
	}

	// 连接SFTP
	e.runSFTP()

	// 启动Event处理
	go e.eventsHandler()

	return nil
}

func (e *Event) runSFTP() {
	// 连接sftp
	// 连接服务器
	for _, c := range e.configs {
		if c.Env == config.LOCAL {
			s, err := sync.CreateSFTP(c.User, c.Password, c.Host, c.Port)
			if err != nil {
				panic(err)
			}
			e.sFTP[c.LocalPath] = s
		}
	}
}

func (e *Event) eventsHandler() {
	log.Println("文件事件处理……")

	var signal string
	sTime := time.Now()
	eTime := time.Now()
	for {
	EventsLoop:
		select {
		case <-time.After(signalTimeOut):
			if signal != "" {
				for _, cmd := range e.configs[signal].Cmd {
					if e.configs[signal].Env == config.LOCAL {
						e.sFTP[signal].Session.Run(cmd)
					} else if e.configs[signal].Env == config.SERVER {
						exec.Command("sh", "-c", cmd).Output()
					}
				}
				signal = ""
				log.Warning("Shell Done.")
				sub := eTime.Sub(sTime)
				fmt.Printf("总花费时间：%v\n", sub)
			}

			sTime = time.Now()
		case f := <-e.fileEvent:
			for localPath, cv := range e.configs {
				if strings.HasPrefix(f.Name, localPath) {
					// 文件夹or文件过滤
					for _, ig := range cv.IgnorePrefix {
						if strings.HasPrefix(f.Name, ig) {
							goto EventsLoop
						}
					}
					for _, ig := range cv.IgnoreSuffix {
						if strings.HasSuffix(f.Name, ig) {
							goto EventsLoop
						}
					}

					// 本地模式需要上传SFTP
					if cv.Env == config.LOCAL {
						df := strings.Replace(f.Name, cv.LocalPath, cv.DeploymentPath, 1)
						e.opEvent(e.sFTP[localPath], f, df)
					}
					log.Println(f.Name, f.Event)
					signal = localPath
					eTime = time.Now()
				}
			}
		}
	}
}

func (e *Event) opEvent(sftp *sync.SFTP, f FileEvent, df string) {
	if f.Event == fsnotify.Write || f.Event == fsnotify.Create {
		sftp.Upload(f.Name, df)
	} else if f.Event == fsnotify.Remove {
		sftp.Delete(f.Name, df)
	} else if f.Event == fsnotify.Rename {
		log.Println("rename:", f.Name)
	}
}
