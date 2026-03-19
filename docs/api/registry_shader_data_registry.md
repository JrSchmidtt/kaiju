# Package registry/shader_data_registry

**Import path:** `kaijuengine.com/registry/shader_data_registry`



## Constants


### const (

```go
ShaderDataStandardFlagOutline = StandardShaderDataFlags(1 << iota)
// Enable bit will be set anytime there are flags. This is needed because
// bits at the extremes of the float will be truncated to 0 otherwise. By
// setting this bit (largest exponent bit 2^1) this issue can be prevented.
ShaderDataStandardFlagEnable = 1 << 30
```

)

## Functions


### func Create(name string) rendering.DrawInstance

### func StandardShaderDataFlagsClear(target rendering.DrawInstance, flag StandardShaderDataFlags)

### func StandardShaderDataFlagsSet(target rendering.DrawInstance, flag StandardShaderDataFlags)

### func StandardShaderDataFlagsTest(target rendering.DrawInstance, flag StandardShaderDataFlags) bool


## Types


### type ShaderDataEdFrustumWire struct {

```go
rendering.ShaderDataBase `visible:"false"`
```


```go
Color             matrix.Color
FrustumProjection matrix.Mat4
```

}

### func (t ShaderDataEdFrustumWire) Size() int


### type ShaderDataEdThumbPreviewMesh struct {

```go
rendering.ShaderDataBase `visible:"false"`
```


```go
View       matrix.Mat4 `visible:"false"`
Projection matrix.Mat4 `visible:"false"`
```

}

### func (s *ShaderDataEdThumbPreviewMesh) SetCamera(view, projection matrix.Mat4)


### func (ShaderDataEdThumbPreviewMesh) Size() int


### type ShaderDataEdTransformWire struct {

```go
rendering.ShaderDataBase `visible:"false"`
```


```go
Color matrix.Color
```

}

### func (t ShaderDataEdTransformWire) Size() int


### type ShaderDataGrid struct {

```go
rendering.ShaderDataBase `visible:"false"`
```


```go
Color matrix.Color
```

}

### func (t ShaderDataGrid) Size() int


### type ShaderDataPBR struct {

```go
rendering.ShaderDataBase `visible:"false"`
```


```go
VertColors matrix.Color
MeRoEmAo   matrix.Vec4
Flags      StandardShaderDataFlags `visible:"false"`
LightIds   [4]int32                `visible:"false"`
```

}

### func (s *ShaderDataPBR) SelectLights(lights rendering.LightsForRender)


### func (t ShaderDataPBR) Size() int


### type ShaderDataParticle struct {

```go
rendering.ShaderDataBase `visible:"false"`
```


```go
Color matrix.Color
UVs   matrix.Vec4 `default:"0,0,1,1"`
```

}

### func (t ShaderDataParticle) Size() int


### type ShaderDataPbrSkinned struct {

```go
rendering.SkinnedShaderDataHeader `visible:"false"`
rendering.ShaderDataBase          `visible:"false"`
```


```go
VertColors matrix.Color
MeRoEmAo   matrix.Vec4
Flags      StandardShaderDataFlags `visible:"false"`
LightIds   [4]int32                `visible:"false"`
```

}

### func (t *ShaderDataPbrSkinned) BoundDataPointer() unsafe.Pointer


### func (t *ShaderDataPbrSkinned) InstanceBoundDataSize() int


### func (t ShaderDataPbrSkinned) Size() int


### func (t *ShaderDataPbrSkinned) SkinningHeader() *rendering.SkinnedShaderDataHeader


### func (t *ShaderDataPbrSkinned) UpdateBoundData() bool


### type ShaderDataStandard struct {

```go
rendering.ShaderDataBase `visible:"false"`
```


```go
Color matrix.Color
Flags StandardShaderDataFlags `visible:"false"`
```

}

### func (ShaderDataStandard) Size() int


### type ShaderDataStandardSkinned struct {

```go
rendering.SkinnedShaderDataHeader `visible:"false"`
rendering.ShaderDataBase          `visible:"false"`
```


```go
Color matrix.Color
Flags StandardShaderDataFlags `visible:"false"`
```

}

### func (t *ShaderDataStandardSkinned) BoundDataPointer() unsafe.Pointer


### func (t *ShaderDataStandardSkinned) InstanceBoundDataSize() int


### func (t ShaderDataStandardSkinned) Size() int


### func (t *ShaderDataStandardSkinned) SkinningHeader() *rendering.SkinnedShaderDataHeader


### func (t *ShaderDataStandardSkinned) UpdateBoundData() bool


### type ShaderDataUnlit struct {

```go
rendering.ShaderDataBase `visible:"false"`
```


```go
Color matrix.Color
UVs   matrix.Vec4             `default:"0,0,1,1"`
Flags StandardShaderDataFlags `visible:"false"`
```

}

### func (t ShaderDataUnlit) Size() int


### type StandardShaderDataFlags = uint32



