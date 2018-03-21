// +build freebsd

package ipfs

// Derived from go-ipfs

import (
	"fmt"

	unix "golang.org/x/sys/unix"
)

func init() {
	fileDescriptorCheck = checkAndSetUlimit
}

func checkAndSetUlimit() error {
	var rLimit unix.Rlimit
	err := unix.Getrlimit(unix.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		return fmt.Errorf("error getting rlimit: %s", err)
	}

	ipfsFileDescNum := int64(ipfsFileDescNum)

	var setting bool
	if rLimit.Cur < ipfsFileDescNum {
		if rLimit.Max < ipfsFileDescNum {
			fmt.Println("adjusting max")
			rLimit.Max = ipfsFileDescNum
		}
		fmt.Printf("Adjusting current ulimit to %d...\n", ipfsFileDescNum)
		rLimit.Cur = ipfsFileDescNum
		setting = true
	}

	err = unix.Setrlimit(unix.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		return fmt.Errorf("error setting ulimit: %s", err)
	}

	if setting {
		fmt.Printf("Successfully raised file descriptor limit to %d.\n", ipfsFileDescNum)
	}

	return nil
}
