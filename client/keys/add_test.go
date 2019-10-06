package keys

import (
	"fmt"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/tests"
)

func Test_runAddCmdBasic(t *testing.T) {
	runningOnServer := isRunningOnServer()
	cmd := addKeyCommand()
	assert.NotNil(t, cmd)
	mockIn, mockOut, _ := tests.ApplyMockIO(cmd)

	kbHome, kbCleanUp := tests.NewTestCaseDir(t)
	assert.NotNil(t, kbHome)
	defer kbCleanUp()
	viper.Set(flags.FlagHome, kbHome)
	viper.Set(cli.OutputFlag, OutputFormatText)

	if runningOnServer {
		mockIn.Reset("testpass1\ntestpass1\n")
	} else {
		mockIn.Reset("y\n")
		kb := NewKeyring(mockIn)
		defer func() {
			kb.Delete("keyname1", "", false)
			kb.Delete("keyname2", "", false)
		}()
	}
	assert.NoError(t, runAddCmd(cmd, []string{"keyname1"}))

	if runningOnServer {
		mockIn.Reset("testpass1\nN\n")
	} else {
		mockIn.Reset("N\n")
	}
	assert.Error(t, runAddCmd(cmd, []string{"keyname1"}))

	if runningOnServer {
		mockIn.Reset("testpass1\ny\ntestpass1\n")
	} else {
		mockIn.Reset("y\n")
	}
	err := runAddCmd(cmd, []string{"keyname1"})
	fmt.Println(mockOut.String())
	assert.NoError(t, err)

	viper.Set(cli.OutputFlag, OutputFormatJSON)
	if runningOnServer {
		mockIn.Reset("testpass1\n")
	} else {
		mockIn.Reset("y\n")
	}
	assert.NoError(t, runAddCmd(cmd, []string{"keyname2"}))
}
