package namespace

import (
	"errors"
	"fmt"
	"syscall"
	"unsafe"
)

var DefaultNamespace = &namespace{}

type namespace struct{}

func (n *namespace) Bind(new string, old string, flag int) error { return syscall.Bind(new, old, flag) }
func (n *namespace) Chdir(dir string) error                      { return syscall.Chdir(dir) }

// Unmount unmounts
func (n *namespace) Unmount(new string, old string) error { return syscall.Unmount(new, old) }

// Clear clears the name space with rfork(RFCNAMEG).
func (n *namespace) Clear() error {
	r1, _, _ := syscall.RawSyscall(syscall.SYS_RFORK, uintptr(syscall.RFCNAMEG), 0, 0)
	if r1 != 0 {
		if int32(r1) == -1 {
			return errors.New(errstr())
		}
		// parent; return PID
		return nil
	}
	return nil
}

// Import imports a name space from a remote system
func (n *namespace) Import(host string, remotepath string, mountpoint string, f int) error {
	flag := mountflag(f)
	args := []string{host}
	if remotepath != "" {
		args = append(args, remotepath)
	}
	args = append(args, mountpoint)
	flg := ""
	if flag&AFTER != 0 {
		flg += "a"
	}
	if flag&BEFORE != 0 {
		flg += "b"
	}
	if flag&CREATE != 0 {
		flg += "c"
	}
	if len(flg) > 0 {
		args = append([]string{flg}, args...)
	}
	return syscall.Exec("import", args, nil)
}

// Mount opens a fd with the server name and mounts the open fd to
// old
func (n *namespace) Mount(servername string, old string, flag int) error {
	fd, err := syscall.Open(servername, syscall.O_RDWR)
	if err != nil {
		return fmt.Errorf("open failed: %v", err)
	}
	return syscall.Mount(fd, -1, old, flag, "")
}

func errstr() string {
	var buf [syscall.ERRMAX]byte

	syscall.RawSyscall(syscall.SYS_ERRSTR, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)), 0)

	buf[len(buf)-1] = 0
	return cstring(buf[:])
}

func cstring(s []byte) string {
	for i := range s {
		if s[i] == 0 {
			return string(s[0:i])
		}
	}
	return string(s)
}
