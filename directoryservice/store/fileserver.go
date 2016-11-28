package store

/*
 * Rootpath is '/' or '/test/' etc
 * Get file's path upto file name and search by that
 * Send back host
 * May be more than one record for each fileserver if they handle two or more folders
 */
type FileServer struct {
	FolderPath string
	Host       string
}
