package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/kava-labs/kava/x/committee/types"
)

// RegisterInvariants registers all staking invariants
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {

	ir.RegisterRoute(types.ModuleName, "valid-committees",
		ValidCommitteesInvariant(k))
	ir.RegisterRoute(types.ModuleName, "valid-proposals",
		ValidProposalsInvariant(k))
	ir.RegisterRoute(types.ModuleName, "valid-votes",
		ValidVotesInvariant(k))
}

// ValidCommitteesInvariant verifies that all committees in the store are independently valid
func ValidCommitteesInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {

		var validationErr error
		var invalidCommittee types.Committee
		k.IterateCommittees(ctx, func(com types.Committee) bool {

			if err := com.Validate(); err != nil {
				validationErr = err
				invalidCommittee = com
				return true
			}
			return false
		})

		broken := validationErr != nil
		invariantMessage := sdk.FormatInvariant(
			types.ModuleName,
			"valid committees",
			fmt.Sprintf(
				"\tfound invalid committee, reason: %s\n"+
					"\tcommittee:\n\t%+v\n",
				validationErr, invalidCommittee),
		)
		return invariantMessage, broken
	}
}

// ValidProposalsInvariant verifies that all proposals in the store are valid
func ValidProposalsInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {

		var validationErr error
		var invalidProposal types.Proposal
		k.IterateProposals(ctx, func(proposal types.Proposal) bool {
			invalidProposal = proposal

			if err := proposal.PubProposal.ValidateBasic(); err != nil {
				validationErr = err
				return true
			}

			currentTime := ctx.BlockTime()
			if !currentTime.Equal(time.Time{}) { // this avoids a simulator bug where app.InitGenesis is called with blockTime=0 instead of the correct time
				if proposal.Deadline.Before(currentTime) {
					validationErr = fmt.Errorf("deadline after current block time (%s)", currentTime)
					return true
				}
			}

			com, found := k.GetCommittee(ctx, proposal.CommitteeID)
			if !found {
				validationErr = fmt.Errorf("proposal refers to non existant committee ID '%d'", proposal.CommitteeID)
				return true
			}

			if !com.HasPermissionsFor(proposal.PubProposal) {
				validationErr = fmt.Errorf("proposal not permitted for committee (%+v)", com)
				return true
			}

			return false
		})

		broken := validationErr != nil
		invariantMessage := sdk.FormatInvariant(
			types.ModuleName,
			"valid proposals",
			fmt.Sprintf(
				"\tfound invalid proposal, reason: %s\n"+
					"\tproposal:\n\t%s\n",
				validationErr, invalidProposal),
		)
		return invariantMessage, broken
	}
}

// ValidVotesInvariant verifies that all votes in the store are valid
func ValidVotesInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {

		voteIterator := sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), types.VoteKeyPrefix)
		defer voteIterator.Close()

		var validationErr error
		var invalidVote types.Vote
		for ; voteIterator.Valid(); voteIterator.Next() {
			var vote types.Vote
			k.cdc.MustUnmarshalBinaryLengthPrefixed(voteIterator.Value(), &vote)

			if _, found := k.GetProposal(ctx, vote.ProposalID); !found {
				validationErr = fmt.Errorf("vote refers to non existant proposal ID '%d'", vote.ProposalID)
				invalidVote = vote
				break
			}

			if vote.Voter.Empty() {
				validationErr = fmt.Errorf("empty voter address")
				invalidVote = vote
				break
			}

			// TODO check voter is a committee member?
		}

		broken := validationErr != nil
		invariantMessage := sdk.FormatInvariant(
			types.ModuleName,
			"valid votes",
			fmt.Sprintf(
				"\tfound invalid vote, reason: %s\n"+
					"\tvote:\n\t%+v\n",
				validationErr, invalidVote),
		)
		return invariantMessage, broken
	}
}
