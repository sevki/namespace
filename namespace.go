package namespace

import (
	"fmt"
)

type syzcall int
type mountflag int

//go:generate stringer -type syzcall,mountflag
const (
	// REPL Replace the old file by the new one.
	// Henceforth, an evaluation of old will be translated to the new file.
	// If they are directories (for mount, this condition is true by definition),
	// old becomes a union directory consisting of one directory (the new file).
	REPL mountflag = 0x0000
	// BEFORE Both the old and new files must be directories.
	// Add the constituent files of the new directory to the
	// union directory at old so its contents appear first in the union.
	// After an BEFORE bind or mount, the new directory will be
	// searched first when evaluating file names in the union directory.
	BEFORE mountflag = 0x0001
	// AFTER Like MBEFORE but the new directory goes at the end of the union.
	AFTER mountflag = 0x0002
	// CREATE flag that can be OR'd with any of the above.
	// When a create system call (see open(2)) attempts to create in a union directory,
	// and the file does not exist, the elements of the union are searched in order until
	// one is found with CREATE set. The file is created in that directory;
	// if that attempt fails, the create fails.
	CREATE mountflag = 0x0004
	// CACHE flag, valid for mount only, turns on caching for files made available by the mount.
	// By default, file contents are always retrieved from the server.
	// With caching enabled, the kernel may instead use a local cache
	// to satisfy read(5) requests for files accessible through this mount point.
	CACHE mountflag = 0x0010
)

const (
	// These are copied over from the syscall pkg for plan9 https://go.plan9.io/pkg/syscall/

	// BIND is the plan9 bind syscall. https://9p.io/magic/man2html/2/bind
	BIND syzcall = 2
	// CHDIR is the plan9 bind syscall. https://9p.io/magic/man2html/2/chdir
	CHDIR syzcall = 3
	// UNMOUNT is the plan9 unmount syscall. https://9p.io/magic/man2html/2/bind
	UNMOUNT syzcall = 35
	// MOUNT is the plan9 MOUNT syscall. https://9p.io/magic/man2html/2/bind
	MOUNT syzcall = 46
	// RFORK is the plan9 rfork() syscall. https://9p.io/magic/man2html/2/fork
	// used to perform clear
	RFORK syzcall = 19
	// IMPORT is not a syscall. https://9p.io/magic/man2html/4/import
	IMPORT syzcall = 7
	// INCLUDE is not a syscall
	INCLUDE syzcall = 14
)

// File are a collection of calls that a namespace
// has to do as defined by the namespace files
type File []cmd

// Namespace is a plan9 namespace. It implmenets the bind(1)
// calls.
// Bind and mount modify the file name space of the current
// process and other processes in its name space group (see fork(2)).
// For both calls, old is the name of an existing file or directory
// in the current name space where the modification is to be made.
// The name old is evaluated as described in intro(2),
// except that no translation of the final path element is done.
type Namespace interface {
	// Bind binds new on old.
	Bind(new, old string, flag int) error
	// Mount mounts servename on old.
	Mount(servername, old, spec string, flag int) error
	// Unmount unmounts new from old, or everything mounted on old if new is missing.
	Unmount(new, old string) error
	// Clear clears the name space with rfork(RFCNAMEG).
	Clear() error
	// Chdir changes the working directory to dir.
	Chdir(dir string) error
	// Import imports a name space from a remote system
	Import(host, remotepath, mountpoint string, flag int) error
}

type cmd struct {
	syscall syzcall
	flag    mountflag

	args []string
}

func (c cmd) String() string { return fmt.Sprintf("%s(%v, %d)", c.syscall, c.args, c.flag) }

// NewNS builds a name space for user.
// It opens the file nsfile (/lib/namespace is used if nsfile is ""),
// copies the old environment, erases the current name space,
// sets the environment variables user and home, and interprets the commands in nsfile.
// The format of nsfile is described in namespace(6).
func NewNS(nsfile string) error {
	var ns Namespace
	ns = DefaultNamespace
	ns.Clear()
	r, err := NewBuilder()
	if err != nil {
		return err
	}
	if err := r.Parse(nsfile); err != nil {
		return err
	}
	return r.buildns(ns)
}

// AddNS also interprets and executes the commands in nsfile.
// Unlike newns it applies the command to the current name
// space rather than starting from scratch.
func AddNS(nsfile string) error {
	var ns Namespace
	ns = DefaultNamespace
	r, err := NewBuilder()
	if err != nil {
		return err
	}
	if err := r.Parse(nsfile); err != nil {
		return err
	}
	return r.buildns(ns)
}
