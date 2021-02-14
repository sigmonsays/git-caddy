package watch

import "github.com/fsnotify/fsnotify"

func NewWatch() (*Watch, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	w := &Watch{
		watcher: watcher,
		Chan:    make(chan *ChangeEvent, 1),
	}
	go w.watchThread()
	return w, nil
}

type Watch struct {
	watcher *fsnotify.Watcher
	Chan    chan *ChangeEvent
}

func (me *Watch) WatchDir(path string) error {
	return me.watcher.Add(path)
}

func (me *Watch) Stop() error {
	return me.watcher.Close()
}
func (me *Watch) watchThread() {
	for {
		select {
		case event, ok := <-me.watcher.Events:
			if !ok {
				return
			}
			log.Infof("event: %+v", event)
			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Infof("modified file: %s", event.Name)
				ev := &ChangeEvent{
					Filename: event.Name,
				}
				me.Chan <- ev
			}
		case err, ok := <-me.watcher.Errors:
			if !ok {
				return
			}
			log.Warnf("error: %s", err)
		}
	}
}

type ChangeEvent struct {
	Filename string
}
