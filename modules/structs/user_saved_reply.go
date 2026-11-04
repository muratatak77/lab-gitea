// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package structs

import "time"

// SavedReply represents a reusable comment text of a user
// swagger:model
type SavedReply struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	// swagger:strfmt date-time
	Created time.Time `json:"created_at"`
	// swagger:strfmt date-time
	Updated time.Time `json:"updated_at"`
}

// CreateSavedReplyOption options for creating a saved reply
type CreateSavedReplyOption struct {
	// required: true
	Title string `json:"title" binding:"Required;MaxSize(100)"`
	// required: true
	Content string `json:"content" binding:"Required"`
}

// EditSavedReplyOption options for editing a saved reply
type EditSavedReplyOption struct {
	Title   *string `json:"title" binding:"MaxSize(100)"`
	Content *string `json:"content"`
}
