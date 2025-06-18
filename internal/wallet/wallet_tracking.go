package wallet

import (
	"context"
	"ton-utils-go/internal/app"
	"ton-utils-go/internal/structures"

	"github.com/sirupsen/logrus"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"
)

func StartTracking() error {
	if err := app.InitApp(); err != nil {
		logrus.Error(err)
		return err
	}

	liteclient := liteclient.NewConnectionPool()
	if err := liteclient.AddConnectionsFromConfig(context.Background(), app.CFG.MAINNET_CONFIG); err != nil {
		logrus.Error(err)
		return err
	}

	api := ton.NewAPIClient(liteclient)
 
	words := app.CFG.Wallet.SEED
	wall, err := wallet.FromSeed(api, words, wallet.V4R2)
	if err != nil {
		return err
	}

	block, err := api.CurrentMasterchainInfo(context.Background())
	if err != nil {
		return err
	}

	balance, err := wall.GetBalance(context.Background(), block)
	if err != nil {
		return err
	}

	logrus.Infof("💰 Wallet balance: %s TON", balance.TON())

	accountInfo, err := api.GetAccount(
		context.Background(),
		block,
		wall.Address(),
	)
	if err != nil {
		return err
	}

	transactions := make(chan *tlb.Transaction)
	go api.SubscribeOnTransactions(
		context.Background(),
		wall.Address(),
		accountInfo.LastTxLT, 
		transactions,
	)

	for {
		select {
		case newTx := <-transactions:
			if newTx.IO.In == nil || newTx.IO.In.MsgType != tlb.MsgTypeInternal {
				logrus.Info("🔸 not internal message")
				continue
			}

			internal := newTx.IO.In.AsInternal()
			if internal.Body == nil {
				logrus.Info("internal message body = nil")
				continue
			}

			slice := internal.Body.BeginParse()

			clonedSlice, err := slice.ToCell()
			if err != nil {
				logrus.Error("can't clone slice:", err)
				continue
			}
			s := clonedSlice.BeginParse()

			opcode, err := slice.LoadUInt(32)
			
			switch opcode {
				case 0: // TON + текст
					msg, err := slice.LoadStringSnake()
					if err != nil {
						logrus.Info("load string snake error")
						continue
					}
					logrus.Infof("📥 INCOMING TON TRANSFER: from %s, amount: %s TON, msg: %s", internal.SrcAddr, internal.Amount.TON(), msg)

				case 0x7362d09c: // Jetton transfer_notification$transfer_notification query_id:uint64 amount:(VarUInteger 16) sender:MsgAddress forward_payload:Cell = TransferNotification;
					var notif structures.JettonNotification
					if err := tlb.LoadFromCell(&notif, s, false); err != nil {
						logrus.Error("❌ Error decoding JettonNotification:", err)
						continue
					}
					logrus.Info("🪙 Found Jetton transfer!")
					logrus.Infof("Sender: %s", notif.Sender)
					logrus.Infof("Amount: %s", notif.Amount.String())
					logrus.Infof("QueryId: %d", notif.QueryId)

				default:
					logrus.Infof("⚠️ Unknown opcode: %d – skipping", opcode)
				}
		}
	}
}