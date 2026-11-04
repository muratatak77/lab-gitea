// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package user_test

import (
	"testing"

	"gitea.dev/models/db"
	"gitea.dev/models/unittest"
	user_model "gitea.dev/models/user"
	"gitea.dev/modules/util"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindSavedReplies(t *testing.T) {
	require.NoError(t, unittest.PrepareTestDatabase())

	replies, err := db.Find[user_model.SavedReply](t.Context(), user_model.FindSavedRepliesOptions{UserID: 2})
	require.NoError(t, err)
	require.Len(t, replies, 2)
	assert.Equal(t, "Duplicate", replies[0].Title)
	assert.Equal(t, "Needs reproduction", replies[1].Title)

	titles, err := user_model.FindSavedReplyTitles(t.Context(), 2)
	require.NoError(t, err)
	require.Len(t, titles, 2)
	assert.Equal(t, []int64{2, 1}, []int64{titles[0].ID, titles[1].ID})
	assert.Equal(t, "Duplicate", titles[0].Title)
	assert.Empty(t, titles[0].Content)
}

func TestSavedReplyOwnerScope(t *testing.T) {
	require.NoError(t, unittest.PrepareTestDatabase())

	reply, err := user_model.GetSavedReply(t.Context(), 2, 1)
	require.NoError(t, err)
	assert.Equal(t, "Needs reproduction", reply.Title)

	_, err = user_model.GetSavedReply(t.Context(), 4, 1)
	require.ErrorIs(t, err, util.ErrNotExist)

	require.ErrorIs(t, user_model.DeleteSavedReply(t.Context(), 4, 1), util.ErrNotExist)
	unittest.AssertExistsAndLoadBean(t, &user_model.SavedReply{ID: 1})

	require.NoError(t, user_model.DeleteSavedReply(t.Context(), 2, 1))
	unittest.AssertNotExistsBean(t, &user_model.SavedReply{ID: 1})
}

func TestSavedReplyTitleInUse(t *testing.T) {
	require.NoError(t, unittest.PrepareTestDatabase())

	inUse, err := user_model.SavedReplyTitleInUse(t.Context(), 2, "Duplicate", 0)
	require.NoError(t, err)
	assert.True(t, inUse)

	inUse, err = user_model.SavedReplyTitleInUse(t.Context(), 2, "Duplicate", 2)
	require.NoError(t, err)
	assert.False(t, inUse)

	inUse, err = user_model.SavedReplyTitleInUse(t.Context(), 1, "Duplicate", 0)
	require.NoError(t, err)
	assert.False(t, inUse)
}
