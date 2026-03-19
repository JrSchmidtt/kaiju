# Package engine

**Import path:** `kaijuengine.com/engine`



## Constants


### const (

```go
DefaultWindowWidth  = 1280
DefaultWindowHeight = 720
```

)
### const InvalidFrameId = math.MaxUint64

```go
    InvalidFrameId can be used to indicate that a frame id is invalid
```



## Variables


### var DebugEntityDataRegistry = map[string]EntityData{}

### var LaunchParams = LaunchParameters{}


## Functions


### func LoadLaunchParams()

### func ReflectValueFromJson(v any, f reflect.Value)

### func RegisterEntityData(value EntityData) error


## Types


### type Entity struct {

```go
Transform matrix.Transform
Parent    *Entity
Children  []*Entity
```


```go
OnDestroy          events.Event
OnDestroyRequested events.Event
OnActivate         events.Event
OnDeactivate       events.Event
```


```go
// Has unexported fields.
```

}
```go
    Entity is a struct that represents an arbitrary object in the host system.
    It contains a 3D transformation and can be a parent of, or a child to,
    other entities. Entities can also contain arbitrary named data to make it
    easier to access data that is specific to the entity.
```


```go
    Child entities are unordered by default, you'll need to call
    #Entity.SetChildrenOrdered to make them ordered. It is recommended to leave
    children unordered unless you have a specific reason to order them.
```


### func NewEntity(workGroup *concurrent.WorkGroup) *Entity

```go
    NewEntity creates a new #Entity struct and returns it
```


### func (e *Entity) Activate()

```go
    Activate will set an active flag on the entity that can be queried with
    #Entity.IsActive. It will also set the active flag on all children of the
    entity. If the entity is already active, this function will do nothing.
```


### func (e *Entity) AddNamedData(key string, data any)

```go
    AddNamedData allows you to add arbitrary data to the entity that can be
    accessed by a string key. This is useful for storing data that is specific
    to the entity.
```


```go
    Named data is stored in a map of slices, so you can add multiple pieces of
    data to the same key. It is recommended to compile the data into a single
    structure so the slice length is 1, but sometimes that's not reasonable.
```


### func (e *Entity) ChildAt(idx int) *Entity

```go
    ChildAt returns the child entity at the specified index
```


### func (e *Entity) ChildCount() int

```go
    ChildCount returns the number of children the entity has
```


### func (e *Entity) Copy(other *Entity)


### func (e *Entity) Deactivate()

```go
    Deactivate will set an active flag on the entity that can be queried with
    #Entity.IsActive. It will also set the active flag on all children of the
    entity. If the entity is already inactive, this function will do nothing.
```


### func (e *Entity) DestroyShaderData()


### func (e *Entity) Duplicate(workGroup *concurrent.WorkGroup) *Entity


### func (e *Entity) FindByName(name string) *Entity

```go
    FindByName will search the entity and the tree of children for the first
    entity with the specified name. If no entity is found, nil will be returned.
```


### func (e *Entity) ForceCleanup()

```go
    ForceCleanup will force the full cleanup of the entity, typically this is to
    be called in very specific scenarios and not directly in game code. Unless
    there is a good reason (like this entity no longer bein gin the host).
```


### func (e *Entity) HasChildRecursive(child *Entity) bool


### func (e *Entity) HasChildren() bool

```go
    HasChildren returns true if the entity has any children
```


### func (e *Entity) HasParent(parent *Entity) bool

```go
    HasParent will loop through each parent and determine if any of them is the
    parent Entity supplied. If so, then it will return true, false otherwise.
```


### func (e *Entity) Id() EntityId

```go
    ID returns the unique identifier of the entity. The Id is only valid for
    entities that are not generated through template instantiation. The Id may
    also be stripped during game runtime if the entity is never externally
    referenced by any other part of the system.
```


### func (e *Entity) IndexOfChild(child *Entity) int


### func (e *Entity) Init(workGroup *concurrent.WorkGroup)


### func (e *Entity) IsActive() bool

```go
    IsActive will return true if the entity is active, false otherwise
```


### func (e *Entity) IsDestroyed() bool

```go
    IsDestroyed will return true if the entity is destroyed, false otherwise
```


### func (e *Entity) IsRoot() bool

```go
    IsRoot returns true if the entity is the root entity in the hierarchy
```


### func (e *Entity) Name() string

```go
    Name returns the name of the entity
```


### func (e *Entity) NamedData(key string) []any

```go
    NamedData will return the data associated with the specified key. If the key
    does not exist, nil will be returned.
```


### func (e *Entity) RemoveNamedData(key string, data any)

```go
    RemoveNamedData will remove the specified data from the entity's named data
    map. If the key does not exist, this function will do nothing.
```


```go
    *This will remove the entire slice and all of it's data*
```


### func (e *Entity) RemoveNamedDataByName(key string)

```go
    RemoveNamedDataByName will remove all of the stored named data that matches
    the given key on the entity
```


### func (e *Entity) Root() *Entity

```go
    Root will return the root entity of the entity's hierarchy
```


### func (e *Entity) SetActive(isActive bool)

```go
    SetActive will set the active flag on the entity that can be queried with
    #Entity.IsActive. It will also set the active flag on all children of the
    entity. If the entity is already active, this function will do nothing.
```


### func (e *Entity) SetChildrenOrdered()

```go
    SetChildrenOrdered sets the children of the entity to be ordered
```


### func (e *Entity) SetChildrenUnordered()

```go
    SetChildrenUnordered sets the children of the entity to be unordered
```


### func (e *Entity) SetName(name string)

```go
    SetName sets the name of the entity
```


### func (e *Entity) SetParent(newParent *Entity)

```go
    SetParent will set the parent of the entity. If the entity already has
    a parent, it will be removed from the parent's children list. If the new
    parent is nil, the entity will be removed from the hierarchy and will become
    the root entity. If the new parent is not nil, the entity will be added to
    the new parent's children list. If the new parent is not active, the entity
    will be deactivated as well.
```


```go
    This will also handle the transformation parenting internally
```


### func (e *Entity) ShaderData() rendering.DrawInstance


### func (e *Entity) StoreShaderData(sd rendering.DrawInstance)


### type EntityData interface {

```go
Init(entity *Entity, host *Host)
```

}

### type EntityId string

```go
    EntityId is a string that represents a unique identifier for an entity.
    The identifier is only valid for entities that are not generated through
    template instantiation. The identifier may also be stripped during game
    runtime if the entity is never externally referenced by any other part of
    the system.
```


### type FrameId = uint64

```go
    FrameId is a unique identifier for a frame
```


### type Host struct {

```go
Window    *windowing.Window
LogStream *logging.LogStream
```


```go
Cameras hostCameras
```


```go
Drawings rendering.Drawings
```


```go
Closing       bool
UIUpdater     Updater
UILateUpdater Updater
Updater       Updater
LateUpdater   Updater
```


```go
OnClose     events.Event
CloseSignal chan struct{}
```


```go
// Has unexported fields.
```

}
```go
    Host is the mediator to the entire runtime for the game/editor. It is the
    main entry point for the game loop and is responsible for managing all
    entities, the window, and the rendering context. The host can be used to
    create and manage entities, call update functions on the main thread,
    and access various caches and resources.
```


```go
    The host is expected to be passed around quite often throughout the program.
    It is designed to remove things like service locators, singletons, and other
    global state. You can have multiple hosts in a program to isolate things
    like windows and game state.
```


### func NewHost(name string, logStream *logging.LogStream, assetDb assets.Database) *Host

```go
    NewHost creates a new host with the given name and log stream. The log
    stream is the log handler that is used by the slog package functions. A Host
    that is created through NewHost has no function until #Host.Initialize is
    called.
```


```go
    This is primarily called from #host_container/New
```


### func (host *Host) AssetDatabase() assets.Database

```go
    AssetDatabase returns the asset database for the host
```


### func (host *Host) Audio() *audio.Audio

```go
    Audio returns the audio system for the host
```


### func (host *Host) Close()

```go
    Close will set the closing flag to true and signal the host to clean up
    resources and close the window.
```


### func (host *Host) CollisionManager() *collision_system.Manager

```go
    CollisionManager returns the collision manager for this host
```


### func (h *Host) Deadline() (time.Time, bool)

```go
    Deadline is here to fulfil context.Context and will return zero and false
```


### func (host *Host) DestroyEntity(entity *Entity)

```go
    DestroyEntity marks the given entity for destruction. The entity will be
    cleaned up at the beginning of the next frame.
```


### func (h *Host) Done() <-chan struct{}

```go
    Done is here to fulfil context.Context and will cose the CloseSignal channel
```


### func (h *Host) Err() error

```go
    Err is here to fulfil context.Context and will return nil or
    context.Canceled
```


### func (host *Host) FontCache() *rendering.FontCache

```go
    FontCache returns the font cache for the host
```


### func (host *Host) Frame() FrameId

```go
    Frame will return the current frame id
```


### func (host *Host) Game() any

```go
    Game will return the primary game mediator for the running application.
    In the editor, this would be *[editor.Editor], in the running game, this
    will be the *game_host.GameHost structure that is generated by the editor
    and filled out by the developer.
```


### func (host *Host) ImportPlugins(path string) error

```go
    ImportPlugins will read all of the plugins that are in the specified folder
    and prepare them for execution.
```


### func (host *Host) Initialize(width, height, x, y int, platformState any) error

```go
    Initializes the various systems and caches that are mediated through the
    host. This includes the window, the shader cache, the texture cache,
    the mesh cache, and the font cache, and the camera systems.
```


### func (host *Host) InitializeAudio() (err error)


### func (host *Host) InitializeRenderer() error


### func (host *Host) Lighting() *lighting.LightingInformation

```go
    Lighting returns a pointer to the internal lighting information
```


### func (host *Host) MaterialCache() *rendering.MaterialCache

```go
    MaterialCache returns the font cache for the host
```


### func (host *Host) MeshCache() *rendering.MeshCache

```go
    MeshCache returns the mesh cache for the host
```


### func (host *Host) Name() string

```go
    Name returns the name of the host
```


### func (host *Host) Physics() *StagePhysics

```go
    Physics returns the stage physics system
```


### func (host *Host) Plugins() []*plugins.LuaVM

```go
    Plugins returns all of the loaded plugins for the host
```


### func (host *Host) PrimaryCamera() cameras.Camera


### func (host *Host) Render()

```go
    Render will render the scene. This starts by preparing any drawings that
    are pending. It also creates any pending shaders, textures, and meshes
    before the start of the render. The frame is then readied, buffers swapped,
    and any transformations that are dirty on entities are then cleaned.
```


### func (host *Host) RunAfterFrames(wait int, call func())

```go
    RunAfterFrames will call the given function after the given number of frames
    have passed from the current frame
```


### func (host *Host) RunAfterNextUIClean(call func())

```go
    RunAfterNextUIClean will run the given function on the next frame.
```


### func (host *Host) RunAfterTime(wait time.Duration, call func())

```go
    RunAfterTime will call the given function after the given number of time has
    passed from the current frame
```


### func (host *Host) RunBeforeRender(call func())


### func (host *Host) RunNextFrame(call func())

```go
    RunNextFrame will run the given function on the next frame. This is the same
    as calling RunAfterFrames(0, func(){})
```


### func (host *Host) RunOnMainThread(call func())


### func (host *Host) Runtime() float64

```go
    Runtime will return how long the host has been running in seconds
```


### func (h *Host) SetFrameRateLimit(fps int64)

```go
    SetFrameRateLimit will set the frame rate limit for the host. If the frame
    rate is set to 0, then the frame rate limit will be removed.
```


```go
    If a frame rate is set, then the host will block until the desired frame
    rate is reached before continuing the update loop.
```


### func (host *Host) SetGame(game any)

```go
    SetGame is to be called by the engine in most cases. It is called by the
    editor when it first starts up to setup the editor game binding. For a game
    generated by the editor, it will be called when the game is bootstrapped and
    provide the *game_host.GameHost structure. You can call this function at any
    time you want, but you really only should need to for special cases.
```


### func (host *Host) ShaderCache() *rendering.ShaderCache

```go
    ShaderCache returns the shader cache for the host
```


### func (host *Host) StartPhysics()


### func (host *Host) Teardown()

```go
    Teardown will destroy the host and all of its resources. This will also
    execute the OnClose event. This will also signal the CloseSignal channel.
```


### func (host *Host) TextureCache() *rendering.TextureCache

```go
    TextureCache returns the texture cache for the host
```


### func (host *Host) Threads() *concurrent.Threads

```go
    Threads returns the long-running threads for this instance of host
```


### func (host *Host) UICamera() cameras.Camera


### func (host *Host) UIThreads() *concurrent.Threads

```go
    UIThreads returns the long-running threads for the UI
```


### func (host *Host) Update(deltaTime float64)

```go
    Update is the main update loop for the host. This will poll the window for
    events, update the entities, and render the scene. This will also check if
    the window has been closed or crashed and set the closing flag accordingly.
```


```go
    The update order is FrameRunner -> Update -> LateUpdate -> EndUpdate:
```


```go
    [-] FrameRunner: Functions added to RunAfterFrames [-] UIUpdate: Functions
    added to UIUpdater [-] UILateUpdate: Functions added to UILateUpdater
    [-] Update: Functions added to Updater [-] LateUpdate: Functions added to
    LateUpdater [-] EndUpdate: Internal functions for preparing for the next
    frame
```


```go
    Any destroyed entities will also be ticked for their cleanup. This will also
    tick the editor entities for cleanup.
```


### func (h *Host) Value(key any) any

```go
    Value is here to fulfil context.Context and will always return nil
```


### func (h *Host) WaitForFrameRate()

```go
    WaitForFrameRate will block until the desired frame rate limit is reached
```


### func (host *Host) WorkGroup() *concurrent.WorkGroup

```go
    WorkGroup returns the work group for this instance of host
```


### type LaunchParameters struct {

```go
Generate   string
StartStage string
Trace      bool
RecordPGO  bool
AutoTest   bool
```

}

### type StagePhysics struct {

```go
// Has unexported fields.
```

}

### func (p *StagePhysics) AddEntity(entity *Entity, body *physics.RigidBody)


### func (p *StagePhysics) Destroy()


### func (p *StagePhysics) FindCollision(hit physics.CollisionHit) (*StagePhysicsEntry, bool)


### func (p *StagePhysics) IsActive() bool


### func (p *StagePhysics) Start()


### func (p *StagePhysics) Update(threads *concurrent.Threads, deltaTime float64)


### func (p *StagePhysics) World() *physics.World


### type StagePhysicsEntry struct {

```go
Entity *Entity
Body   *physics.RigidBody
```

}

### type UpdateId int


### func (u *UpdateId) IsValid() bool


### type Updater struct {

```go
// Has unexported fields.
```

}
```go
    Updater is a struct that stores update functions to be called when the
    #Updater.Update function is called. This simply goes through the list from
    top to bottom and calls each function.
```


```go
    *Note that update functions are unordered, so don't rely on the order*
```


### func NewConcurrentUpdater(threads *concurrent.Threads) Updater

```go
    NewConcurrentUpdater creates a new concurrent #Updater struct and returns it
```


### func NewUpdater() Updater

```go
    NewUpdater creates a new #Updater struct and returns it
```


### func (u *Updater) AddUpdate(update func(float64)) UpdateId

```go
    AddUpdate adds an update function to the list of updates to be called when
    the #Updater.Update function is called. It returns the id of the update
    function that was added so that it can be removed later.
```


```go
    The update function is added to a back-buffer so it will not begin updating
    until the next call to #Updater.Update.
```


### func (u *Updater) Destroy()

```go
    Destroy cleans up the updater and should be called when the updater is no
    longer needed. It will close the pending and complete channels and clear the
    updates map.
```


### func (u *Updater) IsConcurrent() bool

```go
    IsConcurrent will return if this updater is a concurrent updater
```


### func (u *Updater) RemoveUpdate(id *UpdateId)

```go
    RemoveUpdate removes an update function from the list of updates to be
    called when the #Updater.Update function is called. It takes the id of the
    update function that was returned when the update function was added.
```


```go
    The update function is removed from a back-buffer so it will not be removed
    until the next call to #Updater.Update.
```


### func (u *Updater) Update(deltaTime float64)

```go
    Update calls all of the update functions that have been added to the
    updater. It takes a deltaTime parameter that is the approximate amount of
    time since the last call to #Updater.Update.
```



