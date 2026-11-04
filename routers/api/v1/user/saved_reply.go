// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package user

import (
	"net/http"

	"gitea.dev/models/db"
	user_model "gitea.dev/models/user"
	api "gitea.dev/modules/structs"
	"gitea.dev/modules/web"
	"gitea.dev/routers/api/v1/utils"
	"gitea.dev/services/context"
	"gitea.dev/services/convert"
	user_service "gitea.dev/services/user"
)

// ListSavedReplies lists the saved replies of the authenticated user
func ListSavedReplies(ctx *context.APIContext) {
	// swagger:operation GET /user/saved_replies user userListSavedReplies
	// ---
	// summary: List the authenticated user's saved replies
	// produces:
	// - application/json
	// parameters:
	// - name: page
	//   in: query
	//   description: page number of results to return (1-based)
	//   type: integer
	// - name: limit
	//   in: query
	//   description: page size of results
	//   type: integer
	// responses:
	//   "200":
	//     "$ref": "#/responses/SavedReplyList"
	//   default:
	//     "$ref": "#/responses/error"

	listOptions := utils.GetListOptions(ctx)
	replies, count, err := db.FindAndCount[user_model.SavedReply](ctx, user_model.FindSavedRepliesOptions{
		ListOptions: listOptions,
		UserID:      ctx.Doer.ID,
	})
	if err != nil {
		ctx.APIErrorInternal(err)
		return
	}

	ctx.SetLinkHeader(count, listOptions.PageSize)
	ctx.SetTotalCountHeader(count)
	ctx.JSON(http.StatusOK, convert.ToSavedReplyList(replies))
}

// GetSavedReply returns a saved reply of the authenticated user
func GetSavedReply(ctx *context.APIContext) {
	// swagger:operation GET /user/saved_replies/{id} user userGetSavedReply
	// ---
	// summary: Get a saved reply of the authenticated user
	// produces:
	// - application/json
	// parameters:
	// - name: id
	//   in: path
	//   description: id of the saved reply
	//   type: integer
	//   format: int64
	//   required: true
	// responses:
	//   "200":
	//     "$ref": "#/responses/SavedReply"
	//   default:
	//     "$ref": "#/responses/error"

	reply, err := user_model.GetSavedReply(ctx, ctx.Doer.ID, ctx.PathParamInt64("id"))
	if err != nil {
		ctx.APIErrorAuto(err)
		return
	}
	ctx.JSON(http.StatusOK, convert.ToSavedReply(reply))
}

// CreateSavedReply adds a saved reply for the authenticated user
func CreateSavedReply(ctx *context.APIContext) {
	// swagger:operation POST /user/saved_replies user userCreateSavedReply
	// ---
	// summary: Create a saved reply
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: body
	//   in: body
	//   schema:
	//     "$ref": "#/definitions/CreateSavedReplyOption"
	// responses:
	//   "201":
	//     "$ref": "#/responses/SavedReply"
	//   default:
	//     "$ref": "#/responses/error"

	form := web.GetForm[*api.CreateSavedReplyOption](ctx)
	reply, err := user_service.CreateSavedReply(ctx, ctx.Doer.ID, form.Title, form.Content)
	if err != nil {
		ctx.APIErrorAuto(err)
		return
	}
	ctx.JSON(http.StatusCreated, convert.ToSavedReply(reply))
}

// EditSavedReply updates a saved reply of the authenticated user
func EditSavedReply(ctx *context.APIContext) {
	// swagger:operation PATCH /user/saved_replies/{id} user userEditSavedReply
	// ---
	// summary: Edit a saved reply
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: id
	//   in: path
	//   description: id of the saved reply
	//   type: integer
	//   format: int64
	//   required: true
	// - name: body
	//   in: body
	//   schema:
	//     "$ref": "#/definitions/EditSavedReplyOption"
	// responses:
	//   "200":
	//     "$ref": "#/responses/SavedReply"
	//   default:
	//     "$ref": "#/responses/error"

	reply, err := user_model.GetSavedReply(ctx, ctx.Doer.ID, ctx.PathParamInt64("id"))
	if err != nil {
		ctx.APIErrorAuto(err)
		return
	}

	form := web.GetForm[*api.EditSavedReplyOption](ctx)
	title, content := reply.Title, reply.Content
	if form.Title != nil {
		title = *form.Title
	}
	if form.Content != nil {
		content = *form.Content
	}
	if err := user_service.UpdateSavedReply(ctx, reply, title, content); err != nil {
		ctx.APIErrorAuto(err)
		return
	}
	ctx.JSON(http.StatusOK, convert.ToSavedReply(reply))
}

// DeleteSavedReply removes a saved reply of the authenticated user
func DeleteSavedReply(ctx *context.APIContext) {
	// swagger:operation DELETE /user/saved_replies/{id} user userDeleteSavedReply
	// ---
	// summary: Delete a saved reply
	// parameters:
	// - name: id
	//   in: path
	//   description: id of the saved reply
	//   type: integer
	//   format: int64
	//   required: true
	// responses:
	//   "204":
	//     "$ref": "#/responses/empty"
	//   default:
	//     "$ref": "#/responses/error"

	if err := user_model.DeleteSavedReply(ctx, ctx.Doer.ID, ctx.PathParamInt64("id")); err != nil {
		ctx.APIErrorAuto(err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
