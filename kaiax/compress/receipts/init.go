package compress

import (
	"github.com/kaiachain/kaia/kaiax/compress"
	"github.com/kaiachain/kaia/storage/database"
)

var _ compress.CompressionModule = &ReceiptCompressModule{}

type ReceiptCompressModule struct {
	compress.InitOpts
}

func NewReceiptCompression() *ReceiptCompressModule {
	return &ReceiptCompressModule{}
}

func (rc *ReceiptCompressModule) GetChain() compress.BlockChain {
	return rc.Chain
}

func (rc *ReceiptCompressModule) GetDbm() database.DBManager {
	return rc.Dbm
}

func (rc *ReceiptCompressModule) Init(opts *compress.InitOpts) error {
	if opts == nil || opts.Chain == nil || opts.Dbm == nil {
		return errRCInitNil
	}
	rc.InitOpts = *opts
	return nil
}

func (rc *ReceiptCompressModule) Start() error {
	compress.Logger.Info("[Receipt Compression] Compression started")
	go rc.Compress()
	return nil
}

func (rc *ReceiptCompressModule) Stop() {
	compress.Logger.Info("[Receipt Compression] Compression Stopped")
}
