package watcher

import (
	"code_sync_tool/event"
	"code_sync_tool/log"
	"github.com/fsnotify/fsnotify"
	"os"
	"path/filepath"
)

type Watcher struct {
	watcher   *fsnotify.Watcher
	fileEvent chan event.FileEvent
}

func New(f chan event.FileEvent) *Watcher {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}

	watcher := &Watcher{watcher: w, fileEvent: f}
	go watcher.eventsHandler()

	return watcher
}

func (w *Watcher) AddDir(dir string) error {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			path, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			err = w.watcher.Add(path)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func (w *Watcher) eventsHandler() {
	for {
		select {
		case ev := <-w.watcher.Events:
			{
				if ev.Op&fsnotify.Create == fsnotify.Create {
					fi, err := os.Stat(ev.Name)
					// 如果是新文件夹加入到监测
					if err == nil && fi.IsDir() {
						err = w.watcher.Add(ev.Name)
						if err != nil {
							log.Warning(err, ev.Name)
						}
						w.fileEvent <- event.FileEvent{Event: fsnotify.Create, Name: ev.Name}
					}
				} else if ev.Op&fsnotify.Write == fsnotify.Write {
					w.fileEvent <- event.FileEvent{Event: fsnotify.Write, Name: ev.Name}
				} else if ev.Op&fsnotify.Remove == fsnotify.Remove {
					w.fileEvent <- event.FileEvent{Event: fsnotify.Remove, Name: ev.Name}

					fi, err := os.Stat(ev.Name)
					// 文件夹删除，移除监测
					if err == nil && fi.IsDir() {
						w.watcher.Remove(ev.Name)
					}
				} else if ev.Op&fsnotify.Rename == fsnotify.Rename {
					w.fileEvent <- event.FileEvent{Event: fsnotify.Rename, Name: ev.Name}
					w.watcher.Remove(ev.Name)
				}
			}
		case err := <-w.watcher.Errors:
			{
				log.Warning("监测异常：", err)
				return
			}
		}
	}
}
