// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package setting

import (
	"net/http"

	asymkey_model "gitea.dev/models/asymkey"
	"gitea.dev/models/db"
	"gitea.dev/services/context"
)

// KeysList returns the signed-in user's SSH keys for the deploy key picker, which lists them by name
func KeysList(ctx *context.Context) {
	keys, err := db.Find[asymkey_model.PublicKey](ctx, asymkey_model.FindPublicKeyOptions{
		ListOptions: db.ListOptionsAll,
		OwnerID:     ctx.Doer.ID,
		NotKeytype:  asymkey_model.KeyTypePrincipal,
	})
	if err != nil {
		ctx.ServerError("ListPublicKeys", err)
		return
	}
	ctx.JSON(http.StatusOK, keys)
}
