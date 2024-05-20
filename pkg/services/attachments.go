package services

import (
	"context"
	"git.solsynth.dev/hydrogen/interactive/pkg/grpc"
	pcpb "git.solsynth.dev/hydrogen/paperclip/pkg/grpc/proto"
	"github.com/samber/lo"
)

func GetAttachmentByID(id uint) (*pcpb.Attachment, error) {
	return grpc.Attachments.GetAttachment(context.Background(), &pcpb.AttachmentLookupRequest{
		Id: lo.ToPtr(uint64(id)),
	})
}

func GetAttachmentByUUID(uuid string) (*pcpb.Attachment, error) {
	return grpc.Attachments.GetAttachment(context.Background(), &pcpb.AttachmentLookupRequest{
		Uuid: &uuid,
	})
}

func CheckAttachmentByIDExists(id uint, usage string) bool {
	_, err := grpc.Attachments.CheckAttachmentExists(context.Background(), &pcpb.AttachmentLookupRequest{
		Id:    lo.ToPtr(uint64(id)),
		Usage: &usage,
	})

	return err == nil
}
