package model

import (
	"fmt"
)

type SituationProgress string

const (
	SituationProgressDraft           SituationProgress = "draft"
	SituationProgressPendingApproval SituationProgress = "pendingApproval"
	SituationProgressApprovedDraft   SituationProgress = "approvedDraft"
	SituationProgressOpend           SituationProgress = "open"
	SituationProgressPublished       SituationProgress = "published"
	SituationProgressClosing         SituationProgress = "closing"
	SituationProgressClosed          SituationProgress = "closed"
)

func (progress *SituationProgress) FromString(s string) error {
	switch SituationProgress(s) {
	case SituationProgressDraft:
		fallthrough
	case SituationProgressPendingApproval:
		fallthrough
	case SituationProgressApprovedDraft:
		fallthrough
	case SituationProgressOpend:
		fallthrough
	case SituationProgressPublished:
		fallthrough
	case SituationProgressClosing:
		fallthrough
	case SituationProgressClosed:
		*progress = SituationProgress(s)
		return nil
	}
	return fmt.Errorf("invalid progress %s", s)
}
