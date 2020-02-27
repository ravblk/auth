package handlers_test

import (
	"bufio"
	"context"
	"fmt"

	"auth/model"
	"testing"
	"time"

	"github.com/valyala/fasthttp"

	"github.com/buaazp/fasthttprouter"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp/fasthttputil"
)

func TestRegistrationWrongData(t *testing.T) {
	var err error
	handlerTest, srvTest, err = InitHandlers()
	require.Nil(t, err)
	require.NotNil(t, handlerTest)
	require.NotNil(t, srvTest)
	regDataErr := []string{`
	{
		"email":"test@mail.com",
		"password": "dhfgffg43g"
	}`,
		`{
		"first_name":"Иван" ,
		"email":"test27@com",
		"password": "Vdhfgffg43g"
	}`,
		`{
		"first_name":"Иван",
		"last_name": "tutututut"
		"email":"test@com.ru",
		"password": "dhfg"
	}`,
		`{
		"first_name":"Иван",
		"last_name": "tutututut"
		"email":"test@com.ru",
		"password": "123"
	}`,
		`{
		"first_name":"Иван",
		"last_name": "tutututut"
		"email":"test@com.ru",
		"password": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAA1"
	}`,
		`{
		"last_name": "tutututut"
		"email":"test@com.ru",
		"password": "dhfgAAAAdfdf44",
		"registration_type": 1
	}`,
		`{
		"first_name":"Иван",
		"last_name": "tutututut"
		"email":"test@com.ru",
		"registration_type": 2
	}`}

	s := &fasthttp.Server{
		Handler: handlerTest.RegistrationLogin,
	}
	ln := fasthttputil.NewInmemoryListener()

	serverCh := make(chan struct{})
	go func() {
		if err := s.Serve(ln); err != nil {
			t.Logf("unexpected error: %s", err)
			return
		}
		close(serverCh)
	}()

	clientCh := make(chan struct{})
	go func() {
		for _, reg := range regDataErr {
			c, err := ln.Dial()
			if err != nil {
				t.Logf("unexpected error: %s", err)
				return
			}
			req := fmt.Sprintf("POST / HTTP/1.1\r\nHost: aa\r\nContent-Length: %d\r\nContent-Type: application/json\r\nX-Real-UserAgent: %s\r\nX-Real-IP: %s\r\n\r\n%s", len(reg), testUserAgent, testIP, reg)
			_, err = c.Write([]byte(req))
			if err != nil {
				t.Logf("unexpected error: %s", err)
				return
			}
			br := bufio.NewReader(c)
			var resp fasthttp.Response
			if err = resp.Read(br); err != nil {
				t.Logf("unexpected error: %s", err)
				return
			}
			if resp.StatusCode() == 200 {
				t.Logf("unexpected status code: %d. Expecting %s", resp.StatusCode(), "4XX")
				return
			}
			if err = c.Close(); err != nil {
				t.Logf("error when closing the connection: %s", err)
				return
			}
		}
		close(clientCh)
	}()
	select {
	case <-clientCh:
	case <-time.After(time.Second * 10):
		t.Fatal("timeout")
	}

	if err := ln.Close(); err != nil {
		t.Logf("unexpected error: %s", err)
	}

	select {
	case <-serverCh:
	case <-time.After(time.Second * 10):
		t.Fatal("timeout")
	}
}

func TestRegistrationValidData(t *testing.T) {
	regValidData := `{
		"first_name":"Иван" ,
		"last_name": "dert",
		"email":"testReg@mail.ru",
		"password": "dhfgfA4fg43g",
		"registration_type": 1
	}`

	s := &fasthttp.Server{
		Handler: handlerTest.RegistrationLogin,
	}

	ln := fasthttputil.NewInmemoryListener()

	serverCh := make(chan struct{})
	go func() {
		if err := s.Serve(ln); err != nil {
			t.Logf("unexpected error: %s", err)
			return
		}
		close(serverCh)
	}()

	clientCh := make(chan struct{})

	c, err := ln.Dial()
	if err != nil {
		t.Logf("unexpected error: %s", err)
	}
	req := fmt.Sprintf("POST / HTTP/1.1\r\nHost: aa\r\nContent-Length: %d\r\nContent-Type: application/json\r\nX-Real-UserAgent: %s\r\nX-Real-IP: %s\r\n\r\n%s", len(regValidData), testUserAgent, testIP, regValidData)
	_, err = c.Write([]byte(req))
	if err != nil {
		t.Logf("unexpected error: %s", err)
	}
	br := bufio.NewReader(c)
	var resp fasthttp.Response
	if err = resp.Read(br); err != nil {
		t.Logf("unexpected error: %s", err)
	}
	if resp.StatusCode() != 200 {
		t.Logf("unexpected status code: %d. Expecting %s", resp.StatusCode(), "200")
	}
	if err = c.Close(); err != nil {
		t.Logf("error when closing the connection: %s", err)
	}
	close(clientCh)
	select {
	case <-clientCh:
	case <-time.After(time.Second * 10):
		t.Logf("timeout")
	}

	if err := ln.Close(); err != nil {
		t.Logf("unexpected error: %s", err)
	}

	select {
	case <-serverCh:
	case <-time.After(time.Second * 10):
		t.Logf("timeout")
	}
}

func TestTwoFASecretCode(t *testing.T) {
	password := "gfgfgfh3111gF"
	registration := &request.Registration{
		FirstName:     "Name",
		LastName:      "Name",
		Email:         "testTwoFASecretCode@com.ru",
		Password:      password,
		AccountID:     "123657246",
		RegPlatformID: 1,
		RegType:       1,
	}
	_, _, user, errs := srvTest.UsrSvc.UserReg(context.Background(), registration)
	require.Nil(t, errs)
	require.NotNil(t, user)
	info := &model.SessionUserInfo{
		SessionDB: model.SessionDB{
			UserID:             user.ID,
			SessionHash:        "d403b19cc93919957225d02a362e828edd32798aac326965c5abea4f2f687003",
			IP:                 "127.0.0.1",
			UserAgent:          testUserAgent,
			LoginPlatformID:    1,
			LoginType:          1,
			Country:            "Russia",
			SessionRefreshHash: "d403b19cc933daf957225d02a362e828e0dd2798aac326965c5abea4f2f68703",
		},
	}
	info.SessionID = uuid.NewV4().String()
	info.ParentID = info.SessionID
	info.FirstName = user.FirstName
	info.LastName = user.LastName
	info.Nickname = user.Nickname
	info.Email = user.Email
	info.EmailConfirm = true
	info.Moderator = map[int]bool{1: true, 2: false, 3: false, 4: true, 5: false}
	info.Administrator = map[int]bool{1: false, 2: false, 3: false, 4: true, 5: false}
	errs = srvTest.SessSvc.SessionCreate(info)
	require.Nil(t, errs)

	router := fasthttprouter.Router{}
	h := router.Handler
	router.GET("/users/qr", handlerTest.TwoFASecretCode)

	s := &fasthttp.Server{
		Handler: h,
	}

	ln := fasthttputil.NewInmemoryListener()

	serverCh := make(chan struct{})
	go func() {
		if err := s.Serve(ln); err != nil {
			t.Logf("unexpected error: %s", err)
			return
		}
		close(serverCh)
	}()

	clientCh := make(chan struct{})

	c, err := ln.Dial()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	req := fmt.Sprintf("GET /users/qr HTTP/1.1\r\nHost: aa.ru\r\nauthorization: %s\r\nX-Real-UserAgent: %s\r\nX-Real-IP: %s\r\n\r\n", info.SessionHash, testUserAgent, testIP)
	_, err = c.Write([]byte(req))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	br := bufio.NewReader(c)
	var resp fasthttp.Response
	if err = resp.Read(br); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if resp.StatusCode() != 200 {
		t.Fatalf("unexpected status code: %d. Expecting %s", resp.StatusCode(), "200")
	}
	if err = c.Close(); err != nil {
		t.Fatalf("error when closing the connection: %s", err)
	}
	close(clientCh)
	select {
	case <-clientCh:
	case <-time.After(time.Second * 10):
		t.Fatalf("timeout")
	}

	if err := ln.Close(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	select {
	case <-serverCh:
	case <-time.After(time.Second * 10):
		t.Fatalf("timeout")
	}
}
