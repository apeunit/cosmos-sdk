package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	host "github.com/cosmos/cosmos-sdk/x/ibc/24-host"
)

type PacketDataI interface {
	GetCommitment() []byte    // GetCommitment returns (possibly non-recoverable) commitment bytes from its Data and Timeout
	GetTimeoutHeight() uint64 // GetTimeoutHeight returns the timeout height defined specifically for each packet data type instance

	ValidateBasic() sdk.Error // ValidateBasic validates basic properties of the packet data, implements sdk.Msg
	Type() string             // Type returns human readable identifier, implements sdk.Msg
}

// Packet defines a type that carries data across different chains through IBC
type Packet struct {
	PacketDataI `json:"data"` // opaque value which can be defined by the application logic of the associated modules.

	Sequence           uint64 `json:"sequence" yaml:"sequence"`                       // number corresponds to the order of sends and receives, where a Packet with an earlier sequence number must be sent and received before a Packet with a later sequence number.
	SourcePort         string `json:"source_port" yaml:"source_port"`                 // identifies the port on the sending chain.
	SourceChannel      string `json:"source_channel" yaml:"source_channel"`           // identifies the channel end on the sending chain.
	DestinationPort    string `json:"destination_port" yaml:"destination_port"`       // identifies the port on the receiving chain.
	DestinationChannel string `json:"destination_channel" yaml:"destination_channel"` // identifies the channel end on the receiving chain.
}

// NewPacket creates a new Packet instance
func NewPacket(
	data PacketDataI,
	sequence uint64, sourcePort, sourceChannel,
	destinationPort, destinationChannel string,
) Packet {
	return Packet{
		data,
		sequence,
		sourcePort,
		sourceChannel,
		destinationPort,
		destinationChannel,
	}
}

// GetSequence implements PacketI interface
func (p Packet) GetSequence() uint64 { return p.Sequence }

// GetSourcePort implements PacketI interface
func (p Packet) GetSourcePort() string { return p.SourcePort }

// GetSourceChannel implements PacketI interface
func (p Packet) GetSourceChannel() string { return p.SourceChannel }

// GetDestPort implements PacketI interface
func (p Packet) GetDestPort() string { return p.DestinationPort }

// GetDestChannel implements PacketI interface
func (p Packet) GetDestChannel() string { return p.DestinationChannel }

// ValidateBasic implements PacketI interface
func (p Packet) ValidateBasic() error {
	if err := host.DefaultPortIdentifierValidator(p.SourcePort); err != nil {
		return ErrInvalidPacket(DefaultCodespace, fmt.Sprintf("invalid source port ID: %s", err.Error()))
	}
	if err := host.DefaultPortIdentifierValidator(p.DestinationPort); err != nil {
		return ErrInvalidPacket(DefaultCodespace, fmt.Sprintf("invalid destination port ID: %s", err.Error()))
	}
	if err := host.DefaultChannelIdentifierValidator(p.SourceChannel); err != nil {
		return ErrInvalidPacket(DefaultCodespace, fmt.Sprintf("invalid source channel ID: %s", err.Error()))
	}
	if err := host.DefaultChannelIdentifierValidator(p.DestinationChannel); err != nil {
		return ErrInvalidPacket(DefaultCodespace, fmt.Sprintf("invalid destination channel ID: %s", err.Error()))
	}
	if p.Sequence == 0 {
		return ErrInvalidPacket(DefaultCodespace, "packet sequence cannot be 0")
	}
	if p.GetTimeoutHeight() == 0 {
		return ErrPacketTimeout(DefaultCodespace)
	}
	if len(p.GetCommitment()) == 0 {
		return ErrInvalidPacket(DefaultCodespace, "packet data cannot be empty")
	}
	return nil
}
