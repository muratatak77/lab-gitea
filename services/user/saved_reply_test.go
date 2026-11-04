// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package user

import (
	"fmt"
	"strings"
	"testing"

	"gitea.dev/models/db"
	"gitea.dev/models/unittest"
	user_model "gitea.dev/models/user"
	"gitea.dev/modules/util"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateSavedReply(t *testing.T) {
	require.NoError(t, unittest.PrepareTestDatabase())

	reply, err := CreateSavedReply(t.Context(), 2, "  Thanks  ", "Thanks for the report!\n")
	require.NoError(t, err)
	assert.Equal(t, "Thanks", reply.Title)
	assert.Equal(t, "Thanks for the report!\n", reply.Content)
	unittest.AssertExistsAndLoadBean(t, &user_model.SavedReply{ID: reply.ID, UserID: 2})

	longASCII := strings.Repeat("a", user_model.SavedReplyTitleMaxLength)
	longNonASCII := strings.Repeat("ş", user_model.SavedReplyTitleMaxLength)
	for _, tc := range []struct {
		name, title, content string
		wantErr              error
	}{
		{name: "title at limit", title: longASCII, content: "x"},
		{name: "non-ascii title at limit", title: longNonASCII, content: "x"},
		{name: "title too long", title: longASCII + "a", content: "x", wantErr: util.ErrInvalidArgument},
		{name: "non-ascii title too long", title: longNonASCII + "ş", content: "x", wantErr: util.ErrInvalidArgument},
		{name: "blank title", title: "  ", content: "x", wantErr: util.ErrInvalidArgument},
		{name: "blank content", title: "Blank", content: " \n", wantErr: util.ErrInvalidArgument},
		{name: "non-ascii content at limit", title: "Long", content: strings.Repeat("ş", user_model.SavedReplyContentMaxLength)},
		{name: "content too long", title: "Too long", content: strings.Repeat("a", user_model.SavedReplyContentMaxLength+1), wantErr: util.ErrInvalidArgument},
		{name: "taken title", title: "Duplicate", content: "x", wantErr: util.ErrAlreadyExist},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := CreateSavedReply(t.Context(), 2, tc.title, tc.content)
			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}

	// another user can reuse a title
	_, err = CreateSavedReply(t.Context(), 1, "Duplicate", "x")
	require.NoError(t, err)
}

func TestCreateSavedReplyLimit(t *testing.T) {
	require.NoError(t, unittest.PrepareTestDatabase())

	for i := range user_model.SavedRepliesMaxPerUser {
		_, err := CreateSavedReply(t.Context(), 1, fmt.Sprintf("reply %d", i), "x")
		require.NoError(t, err)
	}
	_, err := CreateSavedReply(t.Context(), 1, "one too many", "x")
	require.ErrorIs(t, err, util.ErrInvalidArgument)
	assert.NotNil(t, util.ErrorAsTranslatable(err))

	count, err := db.Count[user_model.SavedReply](t.Context(), user_model.FindSavedRepliesOptions{UserID: 1})
	require.NoError(t, err)
	assert.EqualValues(t, user_model.SavedRepliesMaxPerUser, count)
}

func TestUpdateSavedReply(t *testing.T) {
	require.NoError(t, unittest.PrepareTestDatabase())

	reply := unittest.AssertExistsAndLoadBean(t, &user_model.SavedReply{ID: 1})
	require.NoError(t, UpdateSavedReply(t.Context(), reply, reply.Title, "Please add a stack trace."))
	reply = unittest.AssertExistsAndLoadBean(t, &user_model.SavedReply{ID: 1})
	assert.Equal(t, "Please add a stack trace.", reply.Content)

	require.ErrorIs(t, UpdateSavedReply(t.Context(), reply, "Duplicate", "x"), util.ErrAlreadyExist)
	reply = unittest.AssertExistsAndLoadBean(t, &user_model.SavedReply{ID: 1})
	assert.Equal(t, "Needs reproduction", reply.Title)
}
