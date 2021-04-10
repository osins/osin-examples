package auth

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/osins/osin-storage/storage/pg"
	"github.com/osins/osin-storage/storage/pg/model"
)

func NewClient() Client {
	return &client{}
}

type Client interface {
	Init()
}

type client struct {
}

func (s *client) Init() {
	u := &model.User{
		Id:       uuid.UUID(uuid.New()),
		Username: "richard",
		EMail:    "296907@qq.com",
		Password: "123456",
	}

	if err := pg.NewUserStorage().Create(u); err != nil {
		fmt.Println(err)
	}

	client := &model.Client{
		Id:          uuid.UUID(uuid.New()),
		Secret:      "aabbccdd",
		RedirectUri: "http://localhost:14000/appauth",
		NeedLogin:   true,
		NeedRefresh: true,
	}

	// pg.NewClientManage().Delete(client.Id)
	if err := pg.NewClientStorage().Create(client); err != nil {
		fmt.Println(err)
	}

	f, err := pg.NewClientStorage().Get(client.Id.String())
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("\nclient manage first: %v\n", f)
	}

	b, err := json.Marshal(f)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("\nclient manage first: %s\n", b)
	}
}
