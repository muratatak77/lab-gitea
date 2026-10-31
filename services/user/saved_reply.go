// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package user

import (
	"context"
	"strings"
	"unicode/utf8"

	"gitea.dev/models/db"
	user_model "gitea.dev/models/user"
	"gitea.dev/modules/util"
)

func validateSavedReply(ctx context.Context, userID, replyID int64, title, content string) error {
	if title == "" {
		return util.NewInvalidArgumentErrorf("saved reply title must not be empty")
	}
	if utf8.RuneCountInString(title) > user_model.SavedReplyTitleMaxLength {
		return util.NewInvalidArgumentErrorf("saved reply title must not exceed %d characters", user_model.SavedReplyTitleMaxLength)
	}
	if strings.TrimSpace(content) == "" {
		return util.NewInvalidArgumentErrorf("saved reply content must not be empty")
	}
	if len(content) > user_model.SavedReplyContentMaxLength {
		return util.NewInvalidArgumentErrorf("saved reply content must not exceed %d bytes", user_model.SavedReplyContentMaxLength)
	}

	inUse, err := user_model.SavedReplyTitleInUse(ctx, userID, title, replyID)
	if err != nil {
		return err
	}
	if inUse {
		return util.NewAlreadyExistErrorf("saved reply title %q already exists", title)
	}
	return nil
}

// CreateSavedReply adds a saved reply for the user, the title is unique per user
func CreateSavedReply(ctx context.Context, userID int64, title, content string) (*user_model.SavedReply, error) {
	title = strings.TrimSpace(title)
	return db.WithTx2(ctx, func(ctx context.Context) (*user_model.SavedReply, error) {
		count, err := db.Count[user_model.SavedReply](ctx, user_model.FindSavedRepliesOptions{UserID: userID})
		if err != nil {
			return nil, err
		}
		if count >= user_model.SavedRepliesMaxPerUser {
			return nil, util.NewInvalidArgumentErrorf("a user can have at most %d saved replies", user_model.SavedRepliesMaxPerUser)
		}
		if err := validateSavedReply(ctx, userID, 0, title, content); err != nil {
			return nil, err
		}
		reply := &user_model.SavedReply{UserID: userID, Title: title, Content: content}
		if err := user_model.CreateSavedReply(ctx, reply); err != nil {
			return nil, err
		}
		return reply, nil
	})
}

// UpdateSavedReply changes the title and content of an existing saved reply
func UpdateSavedReply(ctx context.Context, reply *user_model.SavedReply, title, content string) error {
	title = strings.TrimSpace(title)
	return db.WithTx(ctx, func(ctx context.Context) error {
		if err := validateSavedReply(ctx, reply.UserID, reply.ID, title, content); err != nil {
			return err
		}
		reply.Title, reply.Content = title, content
		return user_model.UpdateSavedReply(ctx, reply)
	})
}
