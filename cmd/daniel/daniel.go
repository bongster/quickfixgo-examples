package daniel

import (
	"bufio"
	"fmt"
	"io"
	"time"

	"github.com/fatih/color"
	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/field"
	"github.com/quickfixgo/quickfix"
	"github.com/shopspring/decimal"

	// fix50sp1nos "github.com/quickfixgo/fix50sp1/newordersingle" // ordersingle has no legs configuration
	fix50sp1nom "github.com/quickfixgo/fix50sp1/newordermultileg"
	fix50sp1qr "github.com/quickfixgo/fix50sp1/quoterequest"

	"github.com/spf13/cobra"

	_ "github.com/lib/pq"
)

const (
	usage = "daniel"
	short = "Start an daniel"
	long  = "Start an daniel."
)

var (
	// Cmd is the quote command.
	Cmd = &cobra.Command{
		Use:     usage,
		Short:   short,
		Long:    long,
		Aliases: []string{"x"},
		Example: "qf daniel",
		RunE:    execute,
	}
)

func execute(cmd *cobra.Command, args []string) error {
	fmt.Println("calling daniel command")

	ordType := field.NewOrdType(enum.OrdType_MARKET)
	order := fix50sp1nom.New(
		field.NewClOrdID("TEST"),
		field.NewSide(enum.Side_BUY),
		field.NewTransactTime(time.Now()),
		ordType,
	)
	noLegs := fix50sp1nom.NewNoLegsRepeatingGroup()
	noLegs.Add().Set(field.NewLegSymbol("USD/JPY")).Set(field.NewLegSide("2")).Set(field.NewLegOrderQty(decimal.NewFromFloat(666.0), 0)).Set(field.NewLegSettlDate("20220814")).SetInt(quickfix.Tag(6120), 0)
	order.SetNoLegs(noLegs)
	// fmt.Println(order)
	// 8=FIXT.1.1 9=280 35=R 34=3 49=FIXPriceTaker 52=20220815-09:50:35.584 56=RETPriceMakerTransaction 369=2 131=422b9a9d-3570-4f06-9275-322e76592efa 146=1 55=USD/JPY 461=RCSXXX 167=FOR 15=USD 1=ThunesAccount 555=1 600=USD/JPY 624=2 588=20220817 685=100 6120=0 453=1 448=tier1@client 447=C 452=3 6111=1 10=235
	quoteReqId := field.NewQuoteReqID("422b9a9d-3570-4f06-9275-322e76592efa")
	quote := fix50sp1qr.New(quoteReqId)
	quote.SetSenderCompID("FIXPriceTaker")
	quote.SetSendingTime(time.Now())
	quote.SetMsgSeqNum(3)
	quote.SetLastMsgSeqNumProcessed(2)
	quote.SetTargetCompID("RETPriceMakerTransaction")

	noRelatedSym := fix50sp1qr.NewNoRelatedSymRepeatingGroup()
	noRelatedSym.Add().Set(field.NewSymbol("USD/JPY")).Set(field.NewCFICode("RCSXXX")).Set(field.NewSecurityType(enum.SecurityType(enum.SecurityType_FOREIGN_EXCHANGE_CONTRACT))).Set(field.NewCurrency("USD")).Set(field.NewAccount("ThunesAccount"))
	noRelatedSym.Get(0).SetGroup(noLegs)

	noPartyId := fix50sp1qr.NewNoPartyIDsRepeatingGroup()
	noPartyId.Add().Set(field.NewPartyID("tier1@client")).Set(field.NewPartyIDSource(enum.PartyIDSource_GENERALLY_ACCEPTED_MARKET_PARTICIPANT_IDENTIFIER)).Set(field.NewPartyRole(enum.PartyRole_CLIENT_ID))
	noRelatedSym.Add().SetNoPartyIDs(noPartyId)

	quote.SetNoRelatedSym(noRelatedSym)

	quote.SetCheckSum("235")
	fmt.Println(quote)
	return nil
}

func printConfig(reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	color.Set(color.Bold)
	fmt.Println("Starting FIX acceptor with config:")
	color.Unset()

	color.Set(color.FgHiMagenta)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}

	color.Unset()
}
