package simulation

// DONTCOVER

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

// Simulation parameter constants
const (
	SendEnabled = "send_enabled"
)

// GenBankGenesisState generates a random GenesisState for bank
func GenBankGenesisState(cdc *codec.Codec, r *rand.Rand, ap simulation.AppParams, genesisState map[string]json.RawMessage) {
	var sendEnabled bool
	ap.GetOrGenerate(cdc, SendEnabled, &sendEnabled, r,
		func(r *rand.Rand) {
			sendEnabled = r.Int63n(2) == 0
		})

	bankGenesis := bank.NewGenesisState(sendEnabled)

	fmt.Printf("Selected randomly generated bank parameters:\n%s\n", codec.MustMarshalJSONIndent(cdc, bankGenesis))
	genesisState[bank.ModuleName] = cdc.MustMarshalJSON(bankGenesis)
}
