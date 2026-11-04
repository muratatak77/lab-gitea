// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package convert

import (
	user_model "gitea.dev/models/user"
	api "gitea.dev/modules/structs"
)

// ToSavedReply converts a saved reply to its API format
func ToSavedReply(reply *user_model.SavedReply) *api.SavedReply {
	return &api.SavedReply{
		ID:      reply.ID,
		Title:   reply.Title,
		Content: reply.Content,
		Created: reply.CreatedUnix.AsTime(),
		Updated: reply.UpdatedUnix.AsTime(),
	}
}

// ToSavedReplyList converts a list of saved replies to their API format
func ToSavedReplyList(replies []*user_model.SavedReply) []*api.SavedReply {
	result := make([]*api.SavedReply, len(replies))
	for i, reply := range replies {
		result[i] = ToSavedReply(reply)
	}
	return result
}
