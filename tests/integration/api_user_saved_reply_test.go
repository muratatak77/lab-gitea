// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package integration

import (
	"fmt"
	"net/http"
	"testing"

	auth_model "gitea.dev/models/auth"
	"gitea.dev/models/unittest"
	user_model "gitea.dev/models/user"
	api "gitea.dev/modules/structs"
	"gitea.dev/tests"

	"github.com/stretchr/testify/assert"
)

func TestAPISavedReplies(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	session := loginUser(t, "user2")
	writeToken := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeWriteUser)
	readToken := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeReadUser)

	t.Run("List", func(t *testing.T) {
		req := NewRequest(t, "GET", "/api/v1/user/saved_replies?limit=1&page=2").AddTokenAuth(readToken)
		resp := MakeRequest(t, req, http.StatusOK)
		assert.Equal(t, "2", resp.Header().Get("X-Total-Count"))
		replies := DecodeJSON(t, resp, []*api.SavedReply{})
		if assert.Len(t, replies, 1) {
			assert.Equal(t, "Needs reproduction", replies[0].Title)
		}
	})

	t.Run("RequiresWriteScope", func(t *testing.T) {
		req := NewRequestWithJSON(t, "POST", "/api/v1/user/saved_replies", &api.CreateSavedReplyOption{Title: "t", Content: "c"}).AddTokenAuth(readToken)
		MakeRequest(t, req, http.StatusForbidden)
	})

	var created *api.SavedReply
	t.Run("Create", func(t *testing.T) {
		option := &api.CreateSavedReplyOption{Title: "Thanks", Content: "Thank you!"}
		req := NewRequestWithJSON(t, "POST", "/api/v1/user/saved_replies", option).AddTokenAuth(writeToken)
		created = DecodeJSON(t, MakeRequest(t, req, http.StatusCreated), &api.SavedReply{})
		assert.Equal(t, "Thanks", created.Title)
		unittest.AssertExistsAndLoadBean(t, &user_model.SavedReply{ID: created.ID, UserID: 2})

		req = NewRequestWithJSON(t, "POST", "/api/v1/user/saved_replies", option).AddTokenAuth(writeToken)
		MakeRequest(t, req, http.StatusConflict)
	})

	t.Run("Edit", func(t *testing.T) {
		content := "Thanks a lot!"
		req := NewRequestWithJSON(t, "PATCH", fmt.Sprintf("/api/v1/user/saved_replies/%d", created.ID), &api.EditSavedReplyOption{Content: &content}).AddTokenAuth(writeToken)
		edited := DecodeJSON(t, MakeRequest(t, req, http.StatusOK), &api.SavedReply{})
		assert.Equal(t, "Thanks", edited.Title)
		assert.Equal(t, content, edited.Content)
	})

	t.Run("OtherUsersReplyIsHidden", func(t *testing.T) {
		url := "/api/v1/user/saved_replies/3"
		MakeRequest(t, NewRequest(t, "GET", url).AddTokenAuth(readToken), http.StatusNotFound)
		MakeRequest(t, NewRequest(t, "DELETE", url).AddTokenAuth(writeToken), http.StatusNotFound)
		title := "stolen"
		MakeRequest(t, NewRequestWithJSON(t, "PATCH", url, &api.EditSavedReplyOption{Title: &title}).AddTokenAuth(writeToken), http.StatusNotFound)
		unittest.AssertExistsAndLoadBean(t, &user_model.SavedReply{ID: 3, Title: "Duplicate"})
	})

	t.Run("Delete", func(t *testing.T) {
		url := fmt.Sprintf("/api/v1/user/saved_replies/%d", created.ID)
		MakeRequest(t, NewRequest(t, "DELETE", url).AddTokenAuth(writeToken), http.StatusNoContent)
		MakeRequest(t, NewRequest(t, "GET", url).AddTokenAuth(readToken), http.StatusNotFound)
	})
}
