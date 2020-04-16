// +build !plan9 !harvey linux darwin

package namespace

// DefaultNamespace is the default namespace
var DefaultNamespace = &unixnamespace{}

type unixnamespace struct{}

// Bind binds new on old.
func (ns *unixnamespace) Bind(new string, old string, flag int) error {
	panic("not implemented") // TODO: Implement
}

// Mount mounts servename on old.
func (ns *unixnamespace) Mount(servername string, old string, spec string, flag int) error {
	panic("not implemented") // TODO: Implement
}

// Unmount unmounts new from old, or everything mounted on old if new is missing.
func (ns *unixnamespace) Unmount(new string, old string) error {
	panic("not implemented") // TODO: Implement
}

// Clear clears the name space with rfork(RFCNAMEG).
func (ns *unixnamespace) Clear() error {
	panic("not implemented") // TODO: Implement
}

// Chdir changes the working directory to dir.
func (ns *unixnamespace) Chdir(dir string) error {
	panic("not implemented") // TODO: Implement
}

// Import imports a name space from a remote system
func (ns *unixnamespace) Import(host string, remotepath string, mountpoint string, flag int) error {
	panic("not implemented") // TODO: Implement
}
