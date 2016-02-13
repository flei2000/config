package config

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var yamlExample = []byte(`Hacker: true
name: steve
hobbies:
- skateboarding
- snowboarding
- go
clothing:
  jacket: leather
  trousers: denim
  pants:
    size: large
age: 35
eyes : brown
beard: true
`)

var yamlDst = []byte(`
Hacker: true
clothing:
  jacket: leather
  trousers: denim
`)

var yamlSrc = []byte(`
Hacker: false
clothing:
  jacket: textile
  pants:
    size: large
`)

func Test_readBuffer(t *testing.T) {
	cfg := Config{}
	err := cfg.readBuffer(yamlExample)

	assert.Nil(t, err, "Failed to convert buffer to data")
}

func Test_Get(t *testing.T) {
	cfg := Config{}
	cfg.readBuffer(yamlExample)
	fmt.Printf("%v\n", cfg.root)

	assert.True(t, cfg.GetBool("Hacker"))
	assert.False(t, cfg.GetBool("Nonexisting"))

	assert.Equal(t, 35, cfg.GetInt("age"))
	assert.Equal(t, 0, cfg.GetInt("Nonexisting"))

	assert.Equal(t, "large", cfg.GetString("clothing.pants.size"))
	assert.Equal(t, "", cfg.GetString("Nonexisting"))
}

func Test_Get_Global(t *testing.T) {
	cfg.readBuffer(yamlExample)

	assert.True(t, GetBool("Hacker"))
	assert.False(t, GetBool("Nonexisting"))

	assert.Equal(t, 35, GetInt("age"))
	assert.Equal(t, 0, GetInt("Nonexisting"))

	assert.Equal(t, "large", GetString("clothing.pants.size"))
	assert.Equal(t, "", GetString("Nonexisting"))
}

func Test_merge(t *testing.T) {
	dst := Config{}
	dst.readBuffer(yamlDst)

	src := Config{}
	src.readBuffer(yamlSrc)

	merge(&dst.root, &src.root)

	assert.False(t, dst.GetBool("Hacker"))
	assert.Equal(t, "textile", dst.GetString("clothing.jacket"))
	assert.Equal(t, "denim", dst.GetString("clothing.trousers"))
	assert.Equal(t, "large", dst.GetString("clothing.pants.size"))
}

func Test_Set(t *testing.T) {
	cfg.readBuffer(yamlExample)

	Set("key1", 1)
	assert.Equal(t, 1, GetInt("key1"))

	Set("key.2", true)
	assert.True(t, GetBool("key.2"))

	Set("key.3.3", "value3")
	assert.Equal(t, "value3", GetString("key.3.3"))

	Set("age", 36)
	assert.Equal(t, 36, GetInt("age"))

	Set("clothing.jacket", "textile")
	assert.Equal(t, "textile", GetString("clothing.jacket"))

	Set("clothing.pants.size", "small")
	assert.Equal(t, "small", GetString("clothing.pants.size"))
}

func Test_Sub(t *testing.T) {
	cfg := New()
	cfg.readBuffer(yamlExample)

	sub := cfg.Sub("clothing")
	assert.Equal(t, cfg.GetString("clothing.pants.size"), sub.GetString("pants.size"))

	sub = cfg.Sub("clothing.pants")
	assert.Equal(t, cfg.GetString("clothing.pants.size"), sub.GetString("size"))

	sub = cfg.Sub("clothing.pants.size.large")
	assert.Equal(t, sub, (*Config)(nil))
}

func Test_AllSettings(t *testing.T) {
	cfg := New()
	cfg.readBuffer(yamlExample)

	all_settings := cfg.AllSettings()

	assert.Equal(t, 9, len(all_settings))

	assert.Equal(t, true, all_settings["Hacker"])
	assert.Equal(t, "large", all_settings["clothing.pants.size"])
	assert.Equal(t, []interface{}{"skateboarding", "snowboarding", "go"}, all_settings["hobbies"])
}

func Test_BindEnvs(t *testing.T) {
	cfg.readBuffer(yamlExample)
	os.Setenv("TESTCFG_AGE", "36")
	os.Setenv("TESTCFG_CLOTHING_JACKET", "textile")
	BindEnvs("TESTCFG")

	assert.Equal(t, 36, GetInt("age"))
	assert.Equal(t, "textile", GetString("clothing.jacket"))
}
