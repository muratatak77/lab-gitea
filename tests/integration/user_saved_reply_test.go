// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package integration

import (
	"net/http"
	"testing"

	"gitea.dev/models/unittest"
	user_model "gitea.dev/models/user"
	api "gitea.dev/modules/structs"
	"gitea.dev/tests"

	"github.com/stretchr/testify/assert"
)

func TestUserSettingsSavedReplies(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	const link = "/user/settings/saved_replies"
	session := loginUser(t, "user2")

	req := NewRequestWithValues(t, "POST", link, map[string]string{"title": "Thanks", "content": "Thank you!"})
	session.MakeRequest(t, req, http.StatusOK)
	reply := unittest.AssertExistsAndLoadBean(t, &user_model.SavedReply{UserID: 2, Title: "Thanks"})

	req = NewRequestWithValues(t, "POST", link, map[string]string{"title": "Thanks", "content": "Again"})
	session.MakeRequest(t, req, http.StatusBadRequest)

	resp := session.MakeRequest(t, NewRequest(t, "GET", link), http.StatusOK)
	assert.Contains(t, resp.Body.String(), "Thank you!")

	replies := DecodeJSON(t, session.MakeRequest(t, NewRequest(t, "GET", link+"/list"), http.StatusOK), []*api.SavedReply{})
	assert.Len(t, replies, 3)
	assert.Empty(t, replies[0].Content)
	picked := DecodeJSON(t, session.MakeRequest(t, NewRequestf(t, "GET", "%s/%d", link, reply.ID), http.StatusOK), &api.SavedReply{})
	assert.Equal(t, "Thank you!", picked.Content)

	const editLink = link + "/1/edit"
	session.MakeRequest(t, NewRequest(t, "GET", editLink), http.StatusOK)
	req = NewRequestWithValues(t, "POST", editLink, map[string]string{"title": "Needs reproduction", "content": "Steps, please."})
	session.MakeRequest(t, req, http.StatusOK)
	assert.Equal(t, "Steps, please.", unittest.AssertExistsAndLoadBean(t, &user_model.SavedReply{ID: 1}).Content)

	// user4's reply
	session.MakeRequest(t, NewRequest(t, "GET", link+"/3"), http.StatusNotFound)
	session.MakeRequest(t, NewRequest(t, "GET", link+"/3/edit"), http.StatusNotFound)
	req = NewRequestWithValues(t, "POST", link+"/3/edit", map[string]string{"title": "x", "content": "y"})
	session.MakeRequest(t, req, http.StatusNotFound)
	session.MakeRequest(t, NewRequest(t, "POST", link+"/3/delete"), http.StatusNotFound)

	session.MakeRequest(t, NewRequestf(t, "POST", "%s/%d/delete", link, reply.ID), http.StatusOK)
	unittest.AssertNotExistsBean(t, &user_model.SavedReply{ID: reply.ID})
}
