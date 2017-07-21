package ae

import (
	"io"
	"net/http"
	"time"

	"github.com/clavoie/discord-pusher-deps/types"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

const webhookKind = "Webhook"

type hookContext struct {
	AppContext context.Context
}

func NewHookContext(r *http.Request) types.HookContext {
	return &hookContext{appengine.NewContext(r)}
}

func (hc *hookContext) Delete(encodedKey string) error {
	key, err := datastore.DecodeKey(encodedKey)

	if err != nil {
		return err
	}

	err = datastore.Delete(hc.AppContext, key)

	if err != nil {
		return err
	}

	time.Sleep(time.Second) // the server will refresh before the db deletes
	return nil
}

func (hc *hookContext) GetByHook(hookId string) (*types.HookDal, error) {
	q := datastore.NewQuery(webhookKind).Filter("Hook =", hookId).Limit(1)
	dals := make([]*types.HookDal, 0, 1)
	_, err := q.GetAll(hc.AppContext, &dals)

	if err != nil {
		return nil, err
	}

	if len(dals) < 1 {
		return nil, nil
	}

	return dals[0], nil
}

func (hc *hookContext) GetByTypeUrl(typeName, discordUrl string) (*types.HookDal, error) {
	q := datastore.NewQuery(webhookKind).Filter("DiscordHook = ", discordUrl).Filter("Type =", typeName).Limit(1)
	dals := make([]*types.HookDal, 0, 1)
	_, err := q.GetAll(hc.AppContext, &dals)

	if err != nil {
		return nil, err
	}

	if len(dals) > 0 {
		return dals[0], nil

	}

	return nil, nil
}

func (hc *hookContext) GetAll() ([]string, []*types.HookDal, error) {
	q := datastore.NewQuery(webhookKind)
	dals := make([]*types.HookDal, 0, 10)
	keys, err := q.GetAll(hc.AppContext, &dals)

	if err != nil {
		return nil, nil, err
	}

	keyStrs := make([]string, len(keys))

	for index, key := range keys {
		keyStrs[index] = key.Encode()
	}

	return keyStrs, dals, nil

}

func (hc *hookContext) Errorf(format string, args ...interface{}) {
	log.Errorf(hc.AppContext, format, args...)
}

func (hc *hookContext) Put(dal *types.HookDal) error {
	key := datastore.NewIncompleteKey(hc.AppContext, webhookKind, nil)
	_, err := datastore.Put(hc.AppContext, key, dal)
	time.Sleep(time.Second) // the page will refresh faster than the datastore will read / write
	return err
}

func (hc *hookContext) UrlPost(dal *types.HookDal, reader io.Reader) (*http.Response, error) {
	client := urlfetch.Client(hc.AppContext)
	return client.Post(dal.DiscordHook, "application/json", reader)
}
