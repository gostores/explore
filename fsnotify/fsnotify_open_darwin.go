// +build darwin

/*=================================
* Copyright(c)2015-2016 gostores
* From github.com/howeyc/fsnotify
*=================================*/

package fsnotify

import "syscall"

const open_FLAGS = syscall.O_EVTONLY
