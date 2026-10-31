// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package setting

import (
	"errors"
	"net/http"

	"gitea.dev/models/db"
	user_model "gitea.dev/models/user"
	"gitea.dev/modules/setting"
	"gitea.dev/modules/templates"
	"gitea.dev/modules/util"
	"gitea.dev/services/context"
	"gitea.dev/services/convert"
	"gitea.dev/services/forms"
	user_service "gitea.dev/services/user"
)

const tplSettingsSavedReplies templates.TplName = "user/settings/saved_replies"

func savedRepliesLink() string {
	return setting.AppSubURL + "/user/settings/saved_replies"
}

func prepareSavedRepliesPage(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("settings.saved_replies")
	ctx.Data["PageIsSettingsSavedReplies"] = true
	ctx.Data["SavedReplyTitleMaxLength"] = user_model.SavedReplyTitleMaxLength
	ctx.Data["SavedReplyContentMaxLength"] = user_model.SavedReplyContentMaxLength

	replies, err := db.Find[user_model.SavedReply](ctx, user_model.FindSavedRepliesOptions{
		ListOptions: db.ListOptionsAll,
		UserID:      ctx.Doer.ID,
	})
	if err != nil {
		ctx.ServerError("Find", err)
		return
	}
	ctx.Data["SavedReplies"] = replies
}

// SavedReplies renders the saved replies settings page
func SavedReplies(ctx *context.Context) {
	prepareSavedRepliesPage(ctx)
	if ctx.Written() {
		return
	}
	ctx.HTML(http.StatusOK, tplSettingsSavedReplies)
}

// SavedRepliesList returns the saved replies for the markdown editor's menu,
// the API can't be used because its page size may be smaller than the per-user limit
func SavedRepliesList(ctx *context.Context) {
	replies, err := db.Find[user_model.SavedReply](ctx, user_model.FindSavedRepliesOptions{
		ListOptions: db.ListOptionsAll,
		UserID:      ctx.Doer.ID,
	})
	if err != nil {
		ctx.ServerError("Find", err)
		return
	}
	ctx.JSON(http.StatusOK, convert.ToSavedReplyList(replies))
}

func respondSavedReplyError(ctx *context.Context, title string, err error) {
	switch {
	case errors.Is(err, util.ErrAlreadyExist):
		ctx.JSONErrorWithField(ctx.Tr("settings.saved_replies.title_taken", title), "title")
	default:
		ctx.JSONErrorAuto(err)
	}
}

// SavedRepliesPost adds a saved reply
func SavedRepliesPost(ctx *context.Context) {
	form := context.GetFetchActionForm[*forms.SavedReplyForm](ctx)
	if ctx.Written() {
		return
	}

	reply, err := user_service.CreateSavedReply(ctx, ctx.Doer.ID, form.Title, form.Content)
	if err != nil {
		respondSavedReplyError(ctx, form.Title, err)
		return
	}
	ctx.Flash.Success(ctx.Tr("settings.saved_replies.add_success", reply.Title))
	ctx.JSONRedirect(savedRepliesLink())
}

func loadSavedReply(ctx *context.Context) *user_model.SavedReply {
	reply, err := user_model.GetSavedReply(ctx, ctx.Doer.ID, ctx.PathParamInt64("id"))
	if errors.Is(err, util.ErrNotExist) {
		ctx.NotFound(err)
		return nil
	} else if err != nil {
		ctx.ServerError("GetSavedReply", err)
		return nil
	}
	return reply
}

// SavedReplyEdit renders the settings page with the form filled for one saved reply
func SavedReplyEdit(ctx *context.Context) {
	reply := loadSavedReply(ctx)
	if ctx.Written() {
		return
	}
	prepareSavedRepliesPage(ctx)
	if ctx.Written() {
		return
	}
	ctx.Data["SavedReplyToEdit"] = reply
	ctx.HTML(http.StatusOK, tplSettingsSavedReplies)
}

// SavedReplyEditPost saves changes to a saved reply
func SavedReplyEditPost(ctx *context.Context) {
	form := context.GetFetchActionForm[*forms.SavedReplyForm](ctx)
	if ctx.Written() {
		return
	}
	reply := loadSavedReply(ctx)
	if ctx.Written() {
		return
	}

	if err := user_service.UpdateSavedReply(ctx, reply, form.Title, form.Content); err != nil {
		respondSavedReplyError(ctx, form.Title, err)
		return
	}
	ctx.Flash.Success(ctx.Tr("settings.saved_replies.update_success", reply.Title))
	ctx.JSONRedirect(savedRepliesLink())
}

// SavedReplyDelete removes a saved reply
func SavedReplyDelete(ctx *context.Context) {
	if err := user_model.DeleteSavedReply(ctx, ctx.Doer.ID, ctx.PathParamInt64("id")); err != nil {
		ctx.JSONErrorAuto(err)
		return
	}
	ctx.Flash.Success(ctx.Tr("settings.saved_replies.delete_success"))
	ctx.JSONRedirect(savedRepliesLink())
}
