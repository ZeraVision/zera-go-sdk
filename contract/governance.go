package contract

import (
	"fmt"

	pb "github.com/ZeraVision/go-zera-network/grpc/protobuf"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GovernanceType int16

const (
	Staged    GovernanceType = 0
	Cycle     GovernanceType = 1
	Adaptive  GovernanceType = 2
	Staggared GovernanceType = 3
)

type GovernanceTypeHelper struct {
	Type           GovernanceType
	ProposalPeriod *ProposalPeriod
	Stages         []*Stage
	StartTimestamp *timestamppb.Timestamp // only present for staged and cycle, represents the first start time of the cycle with truncated hour
}

type ProposalPeriodType int16

const (
	Days   ProposalPeriodType = 0
	Months ProposalPeriodType = 1
)

// Present with staged, cycle, and staggered types -- not adaptive
// For staged this represents the SUM of all periods specified in StageLength
type ProposalPeriod struct {
	PeriodType   ProposalPeriodType // days or months
	VotingPeriod uint32             // how long is it voted on
}

// Only present in staged governance type
type Stage struct {
	PeriodType  ProposalPeriodType // days or months
	Length      uint32             // how long does the period last
	Break       bool               // if true, no voting occurs, if false, voting occurs
	MaxApproved uint32             // max number of approved proposals that pass through this stage
}

func CreateGovernance(govType GovernanceTypeHelper, regularQuorum float64, fastQuorum *float64, allowedProposalContracts []string, allowedVotingContracts []string, votingThreshold float64, alwaysWinner *bool, allowMultiChoice bool) (*pb.Governance, error) {

	// Convert regularQuorum to a whole number between 0-9999
	if regularQuorum < 0 || regularQuorum > 100 {
		return nil, fmt.Errorf("regularQuorum must be between 0 and 100")
	}

	regularQuorumScaled := uint32(regularQuorum * 100)

	// Convert fastQuorum to a whole number between 0-9999 (if provided)
	var fastQuorumScaled *uint32
	if fastQuorum != nil {
		if *fastQuorum < 0 || *fastQuorum > 100 {
			return nil, fmt.Errorf("fastQuorum must be between 0 and 100")
		}
		scaled := uint32(*fastQuorum * 100)
		fastQuorumScaled = &scaled
	}

	// Convert voting threshold to a whole number between 0-1000
	if votingThreshold < 0 || votingThreshold > 100 {
		return nil, fmt.Errorf("votingThreshold must be between 0 and 100")
	}

	votingThresholdScaled := uint32(regularQuorum * 10)

	gov := &pb.Governance{
		Type:                      pb.GOVERNANCE_TYPE(govType.Type),
		RegularQuorum:             regularQuorumScaled,
		FastQuorum:                fastQuorumScaled,
		AllowedProposalInstrument: allowedProposalContracts,
		VotingInstrument:          allowedVotingContracts,
		Threshold:                 votingThresholdScaled,
		ChickenDinner:             alwaysWinner,
		AllowMulti:                allowMultiChoice,
	}

	if govType.Type == Staged || govType.Type == Cycle || govType.Type == Staggared {
		if govType.ProposalPeriod == nil {
			return nil, fmt.Errorf("proposalPeriod is required for staged, cycle, and staggered governance types")
		}

		period := pb.PROPOSAL_PERIOD(govType.ProposalPeriod.PeriodType)
		gov.ProposalPeriod = &period
		gov.VotingPeriod = &govType.ProposalPeriod.VotingPeriod
	} else {
		if govType.ProposalPeriod != nil {
			return nil, fmt.Errorf("proposalPeriod is not allowed for adaptive governance type -- have you made a mistake?")
		}
	}

	// TODO verify this
	if govType.Type == Staged || govType.Type == Cycle {
		if gov.StartTimestamp == nil {
			return nil, fmt.Errorf("startTimestamp is required for staged and cycle governance types")
		}

		gov.StartTimestamp = govType.StartTimestamp
	} else {
		if govType.StartTimestamp != nil {
			return nil, fmt.Errorf("startTimestamp is not allowed for adaptive or staggared governance types -- have you made a mistake?")
		}
	}

	// Add staged parameters
	if govType.Type == Staged {
		for _, stage := range govType.Stages {
			if stage.Length < 1 {
				return nil, fmt.Errorf("stage length must be greater than 0")
			}
			if stage.MaxApproved < 1 {
				return nil, fmt.Errorf("max approved must be greater than 0")
			}

			gov.StageLength = append(gov.StageLength, &pb.Stage{
				Period:      pb.PROPOSAL_PERIOD(stage.PeriodType),
				Length:      stage.Length,
				Break:       stage.Break,
				MaxApproved: stage.MaxApproved,
			})
		}
	}

	return gov, nil

}
