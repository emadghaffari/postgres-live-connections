package app

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/gofiber/fiber"
	"github.com/stretchr/testify/assert"

	"github.com/emadghaffari/agileful/client/postgres"
	"github.com/emadghaffari/agileful/config"
)

func configPath() string {
	path, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		fmt.Println(err.Error())
	}

	return strings.TrimSpace(string(path)) + "/config.yaml"
}

type sqlMock struct {
	err error
}

func (s *sqlMock) Connect(config config.Config) error {
	return s.err
}

func (s sqlMock) DB() *pg.DB {
	return pg.Connect(&pg.Options{})
}

func (s sqlMock) Close() error {
	return nil
}

func TestInitPostgres(t *testing.T) {

	testCases := []struct {
		desc string
		err  error
	}{
		{
			desc: "first",
			err:  nil,
		},
		{
			desc: "second",
			err:  fmt.Errorf("some error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			test := sqlMock{}
			test.err = tc.err
			postgres.Storage = &test
			if err := Base.initPostgres(); err != nil {
				assert.Equal(t, err.Error(), tc.err.Error())
			}
		})
	}
}

type cfMock struct {
	err error
}

func (cf *cfMock) Get() config.Config {
	return config.Config{}
}
func (cf *cfMock) SetDebug(bool)       {}
func (cf *cfMock) Set(bts []byte) bool { return true }

func TestInitConfigs(t *testing.T) {
	testCases := []struct {
		desc string
		path string
		err  error
	}{
		{
			desc: "a",
			path: configPath(),
			err:  nil,
		},
		{
			desc: "a",
			path: "",
			err:  fmt.Errorf("init config read file err  open : no such file or directory"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			test := &cfMock{}
			test.err = tc.err
			err := Base.initConfigs(tc.path)
			if err != nil {
				assert.Equal(t, err, tc.err)
			}
		})
	}
}

func TestInitEndpoints(t *testing.T) {
	testCases := []struct {
		desc string
	}{
		{
			desc: "",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			Base.initEndpoints(fiber.New())
		})
	}
}

func TestStartApplication(t *testing.T) {
	testCases := []struct {
		desc string
	}{
		{
			desc: "",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			psql := sqlMock{}
			psql.err = nil
			postgres.Storage = &psql

			fbr := fiber.New()
			go func() {
				Base.StartApplication(fbr)
			}()
			time.Sleep(time.Second * 2)
			fbr.Shutdown()
		})
	}
}
