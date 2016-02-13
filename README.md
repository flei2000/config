# config

A GO module to parse config files.

## Why do I need a new config parsing module

Before creating this module, I tried [viper](https://github.com/spf13/viper) and [confer](https://github.com/jacobstr/confer).
Both are powerful with rich features. But lastly I decided to create a different module because neither are easy to implement what I need: a `Sub` function, which copies a branch of configuration and makes it a new config object.

For example: there is such a config representing:

```json
app:
  cache1:
    max-items: 100
    item-size: 64
  cache2:
    max-items: 200
    item-size: 80
```

After executing:
```go
cfg1 := config.Sub("app.cache1")
```

`cfg1` represents:

```json
max-items: 100
item-size: 64
```

Suppose we have:

```go
func NewCache(cfg *Viper) *Cache {...}
```

which creates a cache based on config information formatted as `cfg1`.
Now it's easy to create these 2 caches separately as:

```go
cfg1 := config.Sub("app.cache1")
cache1 := NewCache(cfg1)

cfg2 := config.Sub("app.cache2")
cache2 := NewCache(cfg2)
```

## How to make `Sub()` easier in this module

This module is carefully designed to make `Sub()` function easier in these aspects:

### Only-one internal data structure

`viper` has 6 internal data structures, one overlap another, to get the value of a config item lastly. `confer` has 3.
But in this module, we merge all config items into one single tree. Thus when you want to get a sub-tree, it's not sufferring to search all these internal data structure and merge them into one.

### Unified data structure

What makes the thing more complex is that, all these internal data structure may be inhomogenous.
It may be a tree, root level a `map[string]interface{}`, other levels `map[interface{}]interface{}`.
Or one level key-value store of `map[string]interface[]`.
Since GO is not a dynamic programming language, we have to switch for all these cases in the source code and deal with them one by one.

In this module, we use an unified data structure: a tree with each level a `map[interface{}][interface{}]`.

### Bind environments statically

Config items in environment are merged into the internal tree when `BindEnvs()` is called. If you change the env again after binding, the new value will not be in config.

### TODO
- Read remote config
- Refresh config periodically or triggered by a signal
- Bind pflag

## Usage
### Initialization

The package itself can be used as a global config:

```go
config.ReadFiles("/etc/app/config.yaml")
```
More than 1 file can be specified in this function:

```go
config.ReadFiles("/etc/app/config.yaml", "/home/user/.app/config.yaml")
```

The latter file will overwrite same config item in the former files.

Or you can create your own instance:

```go
cfg := config.NewConfig()
cfg.ReadFiles("/etc/app/config.yaml", "/home/user/.app/config.yaml")
```

### Set Key/Values

```go
config.Set("app.log.level", "debug")
```

### Get Values
There a group of ways to get a config vaule:
- `GetInt(key string) int`
- `GetString(key string) string`
- `GetBool(key string) bool`
- `Get(key string) interface{}`
- `AllSettings() map[string]interface{}`

`Get()` is common method to get any value: a scalar value, a sub-tree, or even an array.
It always returns `interface{}` so you need to convert it to right type.

`AllSettings()` returns a map of all `keys=>values`.

### Extract a sub-tree into a new Config object

```go
Sub(key string) *Config
```

Extract the data of a sub-tree, put it into a new Config object and return the object.

### Bind Environments

```go
BindEnvs(prefix string)
```

Read all envs and merge those having `prefix` to internal tree. For example:

```bash
bash$ export TESTCFG_APP_LOG_LEVEL=warning
```

```go
config.BindEnvs("TESTCFG")
config.Get("app.log.level")
```

It should returns `warning`.

There is no pre-defined priority of different config sources (file or env).
Usually env should be the first priority so call this function after all config files are read.
