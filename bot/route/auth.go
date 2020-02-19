package route

import (
	"github.com/pingcap/community/bot/pkg/controller"
	"github.com/pingcap/community/bot/util"
)

func auth(ctl *controller.Controller, key string, secret string) bool {
	r := (*ctl).GetRepo(key)
	if r == nil {
		util.Event("repo not found")
		return false
	}
	util.Event("secret", secret, "webhook secret", r.WebhookSecret)
	if secret != r.WebhookSecret {
		return false
	}
	return true
}
