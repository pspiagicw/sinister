package manage

type SyncOptions struct {
	ConfigPath string
	Update     UpdateOptions
	Download   DownloadOptions
}

func Sync(opts SyncOptions) {
	updateOpts := opts.Update
	updateOpts.ConfigPath = opts.ConfigPath

	downloadOpts := opts.Download
	downloadOpts.ConfigPath = opts.ConfigPath

	Update(updateOpts)
	Download(downloadOpts)
}
