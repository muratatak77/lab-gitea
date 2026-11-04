// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package user

import (
	"context"

	"gitea.dev/models/db"
	"gitea.dev/modules/timeutil"
	"gitea.dev/modules/util"

	"xorm.io/builder"
)

const (
	SavedReplyTitleMaxLength   = 100
	SavedReplyContentMaxLength = 64 * 1024
	SavedRepliesMaxPerUser     = 100
)

// SavedReply is a reusable comment text owned by a user
type SavedReply struct {
	ID          int64              `xorm:"pk autoincr"`
	UserID      int64              `xorm:"UNIQUE(user_title) NOT NULL"`
	Title       string             `xorm:"VARCHAR(100) UNIQUE(user_title) NOT NULL"`
	Content     string             `xorm:"LONGTEXT NOT NULL"`
	CreatedUnix timeutil.TimeStamp `xorm:"created"`
	UpdatedUnix timeutil.TimeStamp `xorm:"updated"`
}

func init() {
	db.RegisterModel(new(SavedReply))
}

// FindSavedRepliesOptions lists the saved replies of a user ordered by title
type FindSavedRepliesOptions struct {
	db.ListOptions
	UserID int64
}

// ToConds limits the search to the user's saved replies
func (opts FindSavedRepliesOptions) ToConds() builder.Cond {
	return builder.Eq{"user_id": opts.UserID}
}

// ToOrders sorts saved replies by title
func (opts FindSavedRepliesOptions) ToOrders() string {
	return "title ASC, id ASC"
}

// GetSavedReply returns a saved reply only if it belongs to the given user
func GetSavedReply(ctx context.Context, userID, id int64) (*SavedReply, error) {
	reply := new(SavedReply)
	has, err := db.GetEngine(ctx).Where("id = ? AND user_id = ?", id, userID).Get(reply)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, util.NewNotExistErrorf("saved reply %d does not exist", id)
	}
	return reply, nil
}

// SavedReplyTitleInUse reports whether another saved reply of the user has the title
func SavedReplyTitleInUse(ctx context.Context, userID int64, title string, excludeID int64) (bool, error) {
	return db.Exist[SavedReply](ctx, builder.Eq{"user_id": userID, "title": title}.And(builder.Neq{"id": excludeID}))
}

// CreateSavedReply inserts a new saved reply
func CreateSavedReply(ctx context.Context, reply *SavedReply) error {
	return db.Insert(ctx, reply)
}

// UpdateSavedReply stores the title and content of an existing saved reply
func UpdateSavedReply(ctx context.Context, reply *SavedReply) error {
	_, err := db.GetEngine(ctx).ID(reply.ID).Cols("title", "content").Update(reply)
	return err
}

// DeleteSavedReply removes a saved reply only if it belongs to the given user
func DeleteSavedReply(ctx context.Context, userID, id int64) error {
	affected, err := db.GetEngine(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(new(SavedReply))
	if err != nil {
		return err
	} else if affected == 0 {
		return util.NewNotExistErrorf("saved reply %d does not exist", id)
	}
	return nil
}
