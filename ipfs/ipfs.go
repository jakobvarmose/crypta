package ipfs

import (
	"context"
	"os"

	logging "github.com/ipfs/go-log"

	"github.com/ipfs/go-ipfs/core"
	repoconfig "github.com/ipfs/go-ipfs/repo/config"
	"github.com/ipfs/go-ipfs/repo/fsrepo"

	homedir "github.com/mitchellh/go-homedir"
)

func NewNode(repoPath string, offline bool) (*core.IpfsNode, error) {
	fileDescriptorCheck()

	if !fsrepo.IsInitialized(repoPath) {
		cfg, err := repoconfig.Init(os.Stdout, 2048)
		if err != nil {
			return nil, err
		}
		err = fsrepo.Init(repoPath, cfg)
		if err != nil {
			return nil, err
		}
	}
	repo, err := fsrepo.Open(repoPath)
	if err != nil {
		return nil, err
	}
	ncfg := &core.BuildCfg{
		Repo:   repo,
		Online: !offline,
		ExtraOpts: map[string]bool{
			"pubsub": true,
			"ipnsps": true,
			"mplex":  false,
		},
		Routing: core.DHTOption,
	}
	n, err := core.NewNode(context.Background(), ncfg)
	if err != nil {
		return nil, err
	}
	n.SetLocal(false)
	err = n.Bootstrap(core.DefaultBootstrapConfig)
	if err != nil {
		return nil, err
	}
	if offline {
		err = n.SetupOfflineRouting()
		if err != nil {
			return nil, err
		}
	}
	return n, nil
}

// Derived from go-ipfs

var log1 = logging.Logger("ipfs")

var fileDescriptorCheck = func() error { return nil }

// BestKnownPath returns the best known fsrepo path. If the ENV override is
// present, this function returns that value. Otherwise, it returns the default
// repo path.
func BestKnownPath() (string, error) {
	ipfsPath := "~/.crypta"
	if os.Getenv(repoconfig.EnvDir) != "" {
		ipfsPath = os.Getenv(repoconfig.EnvDir)
	}
	ipfsPath, err := homedir.Expand(ipfsPath)
	if err != nil {
		return "", err
	}
	return ipfsPath, nil
}
