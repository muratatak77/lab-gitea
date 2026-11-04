// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package v29

import (
	"context"

	"gitea.dev/modelmigration/base"
	"gitea.dev/modules/timeutil"
)

func AddSavedReplyTable(_ context.Context, x base.EngineMigration) error {
	type SavedReply struct {
		ID          int64              `xorm:"pk autoincr"`
		UserID      int64              `xorm:"UNIQUE(user_title) NOT NULL"`
		Title       string             `xorm:"VARCHAR(100) UNIQUE(user_title) NOT NULL"`
		Content     string             `xorm:"LONGTEXT NOT NULL"`
		CreatedUnix timeutil.TimeStamp `xorm:"created"`
		UpdatedUnix timeutil.TimeStamp `xorm:"updated"`
	}

	return x.Sync(new(SavedReply))
}
