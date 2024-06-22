package services

import (
	"context"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/gap"
	"git.solsynth.dev/hydrogen/paperclip/pkg/proto"
	"github.com/samber/lo"
)

func CheckAttachmentByIDExists(id uint, usage string) bool {
	pc, err := gap.H.DiscoverServiceGRPC("Hydrogen.Paperclip")
	if err != nil {
		return false
	}
	_, err = proto.NewAttachmentsClient(pc).CheckAttachmentExists(context.Background(), &proto.AttachmentLookupRequest{
		Id:    lo.ToPtr(uint64(id)),
		Usage: &usage,
	})
	return err == nil
}
