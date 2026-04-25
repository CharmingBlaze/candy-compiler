//go:build !enet || !cgo

package candy_enet

func tryNewCgoBackend() Backend {
	return nil
}

