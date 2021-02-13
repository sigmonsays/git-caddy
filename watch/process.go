package watch

// TODO: This does nothing useful yet
func WatchProcess() error {

	w, err := NewWatch()
	if err != nil {
		return err
	}
	defer w.Stop()

	w.WatchDir("/home/sig/code/git-caddy/.git")

	for {
		select {

		case change := <-w.Chan:
			log.Infof("Change event %+v", change)
		}
	}
	return nil
}
