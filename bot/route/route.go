package route

import (
	"bytes"
	"io/ioutil"
	"strings"

	"github.com/pingcap/community/bot/pkg/controller"
	"github.com/pingcap/community/bot/util"

	"github.com/google/go-github/v29/github"
	"github.com/kataras/iris"
	"github.com/pkg/errors"
)

// HookBody for parsing webhook
type HookBody struct {
	Repository struct {
		FullName string `json:"full_name"`
	}
}

// Wrapper add webhook router
func Wrapper(app *iris.Application, ctl *controller.Controller) {
	// healthy test
	app.Get("/ping", func(ctx iris.Context) {
		ctx.JSON(iris.Map{
			"message": "pong",
		})
	})

	// Github webhook
	app.Post("/webhook", func(ctx iris.Context) {
		r := ctx.Request()
		body, err := ioutil.ReadAll(r.Body)

		// restore body for iris ReadJSON use
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		hookBody := HookBody{}
		if err := ctx.ReadJSON(&hookBody); err != nil {
			// body parse error
			util.Error(errors.Wrap(err, "webhook post request"))
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.WriteString(err.Error())
			return
		}

		key := strings.Replace(hookBody.Repository.FullName, "/", "-", 1)
		repo := (*ctl).GetRepo(key)
		if repo == nil {
			// repo not in config file
			// util.Error(errors.New("unsupported repo"))
			util.Event("unsupported repo", key)
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.WriteString("unsupported repo")
			return
		}

		// restore body for github ValidatePayload use
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		payload, err := github.ValidatePayload(r, []byte(repo.WebhookSecret))
		if err != nil {
			// invalid payload
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.WriteString(err.Error())
			util.Error(errors.Wrap(err, "webhook post request"))
			return
		}
		event, err := github.ParseWebHook(github.WebHookType(r), payload)
		if err != nil {
			// event parse err
			util.Error(errors.Wrap(err, "webhook post request"))
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.WriteString(err.Error())
			return
		}

		bot := (*ctl).GetBot(key)
		if bot == nil {
			util.Error(errors.New("bot not found, however config exist"))
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.WriteString("unsupported repo")
			return
		}
		ctx.WriteString("ok")
		go (*bot).Webhook(event)
	})

	// monthly check
	app.Get("/history/{owner:string}/{repo:string}", func(ctx iris.Context) {
		owner := ctx.Params().Get("owner")
		repo := ctx.Params().Get("repo")
		key := owner + "-" + repo
		secret := ctx.URLParam("secret")

		if !auth(ctl, key, secret) {
			// repo not in config file or auth fail
			util.Event("unsupported repo")
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.WriteString("unsupported repo")
			return
		}

		bot := (*ctl).GetBot(key)
		if bot == nil {
			util.Error(errors.New("bot not found, however config exist"))
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.WriteString("unsupported repo")
			return
		}

		res, err := (*bot).MonthlyCheck()
		if err != nil {
			util.Error(errors.Wrap(err, "monthly check"))
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.WriteString("check failed")
			return
		}

		ctx.JSON(res)
		return
	})

	// diaplay whitelist
	app.Get("/prlimit/whitelist/{owner:string}/{repo:string}", func(ctx iris.Context) {
		owner := ctx.Params().Get("owner")
		repo := ctx.Params().Get("repo")
		key := owner + "-" + repo
		secret := ctx.URLParam("secret")

		if !auth(ctl, key, secret) {
			// repo not in config file or auth fail
			util.Event("unsupported repo")
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.WriteString("unsupported repo")
			return
		}

		bot := (*ctl).GetBot(key)
		list, err := (*bot).GetMiddleware().Prlimit.GetWhiteList()
		if err != nil {
			util.Event(errors.Wrap(err, "get whitename list"))
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.WriteString("database query error")
			return
		}
		ctx.JSON(list)
	})

	// add whitename
	app.Post("/prlimit/whitelist/{owner:string}/{repo:string}/{username:string}", func(ctx iris.Context) {
		owner := ctx.Params().Get("owner")
		repo := ctx.Params().Get("repo")
		username := ctx.Params().Get("username")
		key := owner + "-" + repo
		secret := ctx.URLParam("secret")

		if !auth(ctl, key, secret) {
			// repo not in config file or auth fail
			util.Event("unsupported repo")
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.WriteString("unsupported repo")
			return
		}

		bot := (*ctl).GetBot(key)
		err := (*bot).GetMiddleware().Prlimit.AddWhiteList(username)
		if err != nil {
			util.Event(errors.Wrap(err, "get whitename list"))
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.WriteString("database query error")
			return
		}
		ctx.WriteString("ok")
	})

	// remove whitename
	app.Delete("/prlimit/whitelist/{owner:string}/{repo:string}/{username:string}", func(ctx iris.Context) {
		owner := ctx.Params().Get("owner")
		repo := ctx.Params().Get("repo")
		username := ctx.Params().Get("username")
		key := owner + "-" + repo
		secret := ctx.URLParam("secret")

		if !auth(ctl, key, secret) {
			// repo not in config file or auth fail
			util.Event("unsupported repo")
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.WriteString("unsupported repo")
			return
		}

		bot := (*ctl).GetBot(key)
		err := (*bot).GetMiddleware().Prlimit.RemoveWhiteList(username)
		if err != nil {
			util.Event(errors.Wrap(err, "get whitename list"))
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.WriteString("database query error")
			return
		}
		ctx.WriteString("ok")
	})

	// diaplay blacklist
	app.Get("/prlimit/blacklist/{owner:string}/{repo:string}", func(ctx iris.Context) {
		owner := ctx.Params().Get("owner")
		repo := ctx.Params().Get("repo")
		key := owner + "-" + repo
		secret := ctx.URLParam("secret")

		if !auth(ctl, key, secret) {
			// repo not in config file or auth fail
			util.Event("unsupported repo")
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.WriteString("unsupported repo")
			return
		}

		bot := (*ctl).GetBot(key)
		list, err := (*bot).GetMiddleware().Prlimit.GetBlackList()
		if err != nil {
			util.Event(errors.Wrap(err, "get blackname list"))
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.WriteString("database query error")
			return
		}
		ctx.JSON(list)
	})

	// add blacklist
	app.Post("/prlimit/blacklist/{owner:string}/{repo:string}/{username:string}", func(ctx iris.Context) {
		owner := ctx.Params().Get("owner")
		repo := ctx.Params().Get("repo")
		username := ctx.Params().Get("username")
		key := owner + "-" + repo
		secret := ctx.URLParam("secret")

		if !auth(ctl, key, secret) {
			// repo not in config file or auth fail
			util.Event("unsupported repo")
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.WriteString("unsupported repo")
			return
		}

		bot := (*ctl).GetBot(key)
		err := (*bot).GetMiddleware().Prlimit.AddBlackList(username)
		if err != nil {
			util.Event(errors.Wrap(err, "get blackname list"))
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.WriteString("database query error")
			return
		}
		ctx.WriteString("ok")
	})

	// remove blacklist
	app.Delete("/prlimit/blacklist/{owner:string}/{repo:string}/{username:string}", func(ctx iris.Context) {
		owner := ctx.Params().Get("owner")
		repo := ctx.Params().Get("repo")
		username := ctx.Params().Get("username")
		key := owner + "-" + repo
		secret := ctx.URLParam("secret")

		if !auth(ctl, key, secret) {
			// repo not in config file or auth fail
			util.Event("unsupported repo")
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.WriteString("unsupported repo")
			return
		}

		bot := (*ctl).GetBot(key)
		err := (*bot).GetMiddleware().Prlimit.RemoveBlackList(username)
		if err != nil {
			util.Event(errors.Wrap(err, "get blackname list"))
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.WriteString("database query error")
			return
		}
		ctx.WriteString("ok")
	})

	// diaplay whitelist
	app.Get("/merge/whitelist/{owner:string}/{repo:string}", func(ctx iris.Context) {
		owner := ctx.Params().Get("owner")
		repo := ctx.Params().Get("repo")
		key := owner + "-" + repo
		secret := ctx.URLParam("secret")

		if !auth(ctl, key, secret) {
			// repo not in config file or auth fail
			util.Event("unsupported repo")
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.WriteString("unsupported repo")
			return
		}

		bot := (*ctl).GetBot(key)
		list, err := (*bot).GetMiddleware().Merge.GetWhiteList()
		if err != nil {
			util.Event(errors.Wrap(err, "get whitename list"))
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.WriteString("database query error")
			return
		}
		ctx.JSON(list)
	})

	// add whitename
	app.Post("/merge/whitelist/{owner:string}/{repo:string}/{username:string}", func(ctx iris.Context) {
		owner := ctx.Params().Get("owner")
		repo := ctx.Params().Get("repo")
		username := ctx.Params().Get("username")
		key := owner + "-" + repo
		secret := ctx.URLParam("secret")

		if !auth(ctl, key, secret) {
			// repo not in config file or auth fail
			util.Event("unsupported repo")
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.WriteString("unsupported repo")
			return
		}

		bot := (*ctl).GetBot(key)
		err := (*bot).GetMiddleware().Merge.AddWhiteList(username)
		if err != nil {
			util.Event(errors.Wrap(err, "get whitename list"))
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.WriteString("database query error")
			return
		}
		ctx.WriteString("ok")
	})

	// remove whitename
	app.Delete("/merge/whitelist/{owner:string}/{repo:string}/{username:string}", func(ctx iris.Context) {
		owner := ctx.Params().Get("owner")
		repo := ctx.Params().Get("repo")
		username := ctx.Params().Get("username")
		key := owner + "-" + repo
		secret := ctx.URLParam("secret")

		if !auth(ctl, key, secret) {
			// repo not in config file or auth fail
			util.Event("unsupported repo")
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.WriteString("unsupported repo")
			return
		}

		bot := (*ctl).GetBot(key)
		err := (*bot).GetMiddleware().Merge.RemoveWhiteList(username)
		if err != nil {
			util.Event(errors.Wrap(err, "get whitename list"))
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.WriteString("database query error")
			return
		}
		ctx.WriteString("ok")
	})
}
