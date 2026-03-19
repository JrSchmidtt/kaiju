# Package bootstrap

**Import path:** `kaijuengine.com/bootstrap`



## Functions


### func Main(game GameInterface, platformState any)


## Types


### type GameInterface interface {

```go
// Launch is used to bootstrap a game, the game should fill out this
// function's details to initialize itself. No updates are provided by the
// engine, so it is on the the implementing code to take care of registering
// any udpates with the supplied host.
Launch(*engine.Host)
```


```go
// PluginRegistry is used to expose types to be exported for use in Lua.
// Any type returned here will have it's members and functions mapped to
// be called by Lua. You can run the engine with the command line argument
// "generate=pluginapi" to dump a Lua API file and ensure your exposed
// types have been correctly inserted.
PluginRegistry() []reflect.Type
```


```go
// ContentDatabase must return the database interface for the engine to use
// when it is trying to access content. You can use exsiting types that
// implement [assets.Database], or you can create your own.
ContentDatabase() (assets.Database, error)
```

}
```go
    GameInterface is the primary interface to implement in order to bootstrap a
    game/application.
```



