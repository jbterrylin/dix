# Go DI Container

A lightweight, type-safe, concurrency-safe Dependency Injection container for Go.

## 🚀 Features

- 🔑 Type + Key based registration & resolution
- 🧠 Auto injection with `InjectFunc` and `InjectStruct`
- 🌀 Provider functions with or without `context.Context`
- 🧼 Optional safe deletion with ref counter
- 📌 Lifecycle hooks for registration, first access, and more
- 🧵 Full thread-safe map-based storage with minimal lock granularity

---

## 🧰 Quick Start

```go
const TestKey = dix.ValueKey("test")
const TestInterfaceKey = dix.ValueKey("test interface")

var _ ITestInterface = NewTestInterface("test interface")

type ITestInterface interface {
	Print()
}

type Test struct {
	count int
	name  string
}

func NewTest(name string) *Test {
	return &Test{
		name: name,
	}
}

func NewTestInterface(name string) ITestInterface {
	return &Test{
		name: name,
	}
}

func (v *Test) Print() {
	v.count += 1
	fmt.Printf("hello, name: %s, count: %d\n", v.name, v.count)
}

func main() {
	test := NewTest("test")
	dix.Add(TestKey, test, dix.WithValueSetDefault())

	testInterface := NewTestInterface("test interface")
	dix.Add(TestInterfaceKey, testInterface, dix.WithValueSetDefault())

	testFromC, err := dix.Get[*Test]()
	if err != nil {
		fmt.Println(err)
		return
	}
	testFromC.Print() // hello, name: test, count: 1

	testInterfaceFromC, err := dix.Get[ITestInterface]()
	if err != nil {
		fmt.Println(err)
		return
	}
	testInterfaceFromC.Print() // hello, name: test interface, count: 1
}
```

---

## 🧩 Core Concepts

### Value
- A Value is a fixed instance registered in the container.
- You can attach an `OnCloseHook`, which will be triggered on [Delete](#func-delete) or [Reset](#func-reset).
### Provider
- A Provider is a factory function that can optionally accept a `context.Context` as a parameter.
- By default, the return value is cached after the first successful call. To disable caching, use WithProviderNoCache when calling [AddProvider](#func-addprovider) or [AddCtxProvider](#func-addctxprovider).
### ValueKey / ProviderKey
- All values and providers are uniquely identified by their type and a string key.
- This design encourages using constants over magic strings for better safety and maintainability.
### DefaultKey
- Each type can have a "default" instance under the `DefaultValueKey`.
- By default, the key is an empty string. You can customize it using [SetDefaultValueKey](#func-setdefaultvaluekey) and [SetDefaultProviderKey](#func-setdefaultproviderkey), although this is generally not recommended.
### Lifecycle Hooks
- [AfterAdd](#func-afteradd), [AfterProviderRun](#func-afterproviderrun), [AfterFirstAccess](#func-afterfirstaccess) and [BeforeDuplicateRegister](#func-beforeduplicateregister).
### Safe Delete
- When enabled via [SetSafeDelete](#func-setsafedelete), the container ensures no active users remain before invoking the `OnCloseHook`.
- Since the container cannot track external usage, users must explicitly signal end-of-use via [DeductRefCount](#func-deductrefcount) or [DeductRefCountByKey](#func-deductrefcountbykey).
- Enabling this feature is generally not recommended, as deletion is not a common practice in DI container design. It also introduces performance overhead due to atomic operations on the reference counter.
### InjectFunc && InjectStruct
- Uses reflection to automatically resolve and inject dependencies into functions or struct fields.

---

## 🛠️ API Index

### Value
#### Add
- [Add](#func-add)
###### ValueAddOption
- [WithValueOnClose](#func-withvalueonclose)
- [WithValueSetDefault](#func-withvaluesetdefault)
- [WithValueTag](#func-withvaluetag)
#### Get
- [Get](#func-get)
- [GetByKey](#func-getbykey)
- [MustGet](#func-mustget)
- [MustGetByKey](#func-mustgetbykey)
#### Check Exist
- [Exist](#func-exist)
- [ExistByKey](#func-existbykey)
#### Delete
- [Delete](#func-delete)
- [DeleteByKey](#func-deletebykey)
##### ValueDeleteOption
- [WithValueSkipOnClose](#func-withvalueskiponclose)
#### Other
- [ListKeys](#func-listkeys)
- [GetAll](#func-getall)
#### Safe Delete
- [DeductRefCount](#func-deductrefcount)
- [DeductRefCountByKey](#func-deductrefcountbykey)

### Provider
#### Add
- [AddProvider](#func-addprovider)
- [AddCtxProvider](#func-addctxprovider)
##### ProviderAddOption
- [WithProviderSetDefault](#func-withprovidersetdefault)
- [WithProviderNoCache](#func-withprovidernocache)
- [WithProviderTag](#func-withprovidertag)
#### Get
- [GetProvider](#func-getprovider)
- [GetProviderWithCtx](#func-getproviderwithctx)
- [GetProviderByKey](#func-getproviderbykey)
- [GetProviderByKeyWithCtx](#func-getproviderbykeywithctx)
- [MustGetProvider](#func-mustgetprovider)
- [MustGetProviderWithCtx](#func-mustgetproviderwithctx)
- [MustGetProviderByKey](#func-mustgetproviderbykey)
- [MustGetProviderByKeyWithCtx](#func-mustgetproviderbykeywithctx)
- [ProviderExist](#func-providerexist)
- [ProviderExistByKey](#func-providerexistbykey)
- [DeleteProvider](#func-deleteprovider)
- [DeleteProviderByKey](#func-deleteproviderbykey)
- [ListProviderKeys](#func-listproviderkeys)
- [GetAllProvider](#func-getallprovider)
##### ProviderGetOption
- [WithProviderReload](#func-withproviderreload)
### Inject
- [InjectStruct](#func-injectstruct)
- [InjectStructWithCtx](#func-injectstructwithctx)
- [InjectFunc](#func-injectfunc)
- [InjectFuncWithCtx](#func-injectfuncwithctx)
##### InjectFuncOption
- [WithInjectFuncProvider](#func-withinjectfuncprovider)
- [WithInjectFuncKey](#func-withinjectfunckey)
- [WithInjectFuncReload](#func-withinjectfuncreload)
- [WithInjectFuncOptional](#func-withinjectfuncoptional)
### Hook
- [AfterAdd](#func-afteradd)
- [AfterProviderRun](#func-afterproviderrun)
- [AfterFirstAccess](#func-afterfirstaccess)
- [BeforeDuplicateRegister](#func-beforeduplicateregister)
### Global
- [SetDefaultValueKey](#func-setdefaultvaluekey)
- [SetDefaultProviderKey](#func-setdefaultproviderkey)
- [SetSafeDelete](#func-setsafedelete)
- [SetResetMaxConcurrent](#func-setresetmaxconcurrent)
- [Reset](#func-reset)

## 🧪 API Reference

### Value
#### Add 
<a id="func-add"></a>

```go
	func Add[T any](key ValueKey, val T, opts ...ValueAddOption) error
```
##### ValueAddOption
<a id="func-withvalueonclose"></a>

```go
	func WithValueOnClose(f func()) ValueAddOption
```
<a id="func-withvaluesetdefault"></a>

```go
	func WithValueSetDefault() ValueAddOption
```
<a id="func-withvaluetag"></a>

```go
	func WithValueTag(tagMap map[string]any) ValueAddOption
```
#### Get
If a get function does not include `ByKey`, it retrieves the value added with the [WithValueSetDefault](#func-withvaluesetdefault) option.

<a id="func-get"></a>

```go
	func Get[T any]() (T, error)
```
<a id="func-getbykey"></a>

```go
	func GetByKey[T any](key ValueKey) (T, error)
```
<a id="func-mustget"></a>

```go
	func MustGet[T any]() T
```
<a id="func-mustgetbykey"></a>

```go
	func MustGetByKey[T any](key ValueKey) T
```
#### Check Exist
<a id="func-exist"></a>

```go
	func Exist[T any]() bool
```
<a id="func-existbykey"></a>

```go
	func ExistByKey[T any](key ValueKey) bool
```
#### Delete
<a id="func-delete"></a>

```go
	func Delete[T any](opts ...ValueDeleteOption)
```
<a id="func-deletebykey"></a>

```go
	func DeleteByKey[T any](key ValueKey, opts ...ValueDeleteOption) error
```
##### ValueDeleteOption
<a id="func-withvalueskiponclose"></a>

```go
	func WithValueSkipOnClose() ValueDeleteOption 
```
#### Other
<a id="func-listkeys"></a>

```go
	func ListKeys[T any]() []ValueKey
```
<a id="func-getall"></a>

```go
	func GetAll[T any]() []T
```
#### Safe Delete
⚠️ The functions below are only relevant when SetSafeDelete is set to true.
<a id="func-deductrefcount"></a>

```go
	func DeductRefCount[T any]() error
```
<a id="func-deductrefcountbykey"></a>

```go
	func DeductRefCountByKey[T any](key ValueKey) error
```
### Provider
#### Add
<a id="func-addprovider"></a>

```go
	func AddProvider[T any](key ProviderKey, value func() (T, error), opts ...ProviderAddOption) error
```
<a id="func-addctxprovider"></a>

```go
	func AddCtxProvider[T any](key ProviderKey, valueWithCtx func(context.Context) (T, error), opts ...ProviderAddOption) error
```
##### ProviderAddOption
<a id="func-withprovidersetdefault"></a>

```go
func WithProviderSetDefault() ProviderAddOption
```
<a id="func-withprovidernocache"></a>

```go
func WithProviderNoCache() ProviderAddOption
```
<a id="func-withprovidertag"></a>

```go
func WithProviderTag(tagMap map[string]any) ProviderAddOption
```
#### Get
If a get function does not include `ByKey`, it retrieves the value added with the [WithValueSetDefault](#func-withvaluesetdefault) option.

<a id="func-getprovider"></a>

```go
	func GetProvider[T any](opts ...ProviderGetOption) (T, error)
```
<a id="func-getproviderwithctx"></a>

```go
	func GetProviderWithCtx[T any](ctx context.Context, opts ...ProviderGetOption) (T, error)
```
<a id="func-getproviderbykey"></a>

```go
	func GetProviderByKey[T any](key ProviderKey, opts ...ProviderGetOption) (T, error)
```
<a id="func-getproviderbykeywithctx"></a>

```go
	func GetProviderByKeyWithCtx[T any](ctx context.Context, key ProviderKey, opts ...ProviderGetOption) (T, error)
```
<a id="func-mustgetprovider"></a>

```go
	func MustGetProvider[T any](opts ...ProviderGetOption) T
```
<a id="func-mustgetproviderwithctx"></a>

```go
	func MustGetProviderWithCtx[T any](ctx context.Context, opts ...ProviderGetOption) T
```
<a id="func-mustgetproviderbykey"></a>

```go
	func MustGetProviderByKey[T any](key ProviderKey, opts ...ProviderGetOption) T
```
<a id="func-mustgetproviderbykeywithctx"></a>

```go
	func MustGetProviderByKeyWithCtx[T any](ctx context.Context, key ProviderKey, opts ...ProviderGetOption) T
```
##### ProviderGetOption
<a id="func-withproviderreload"></a>

```go
	func WithProviderReload() ProviderGetOption
```
This option forces the factory function to run again, ignoring any existing cached value.

#### Check Exist
<a id="func-providerexist"></a>

```go
	func ProviderExist[T any]() bool
```
<a id="func-providerexistbykey"></a>

```go
	func ProviderExistByKey[T any](key ProviderKey) bool
```
#### Delete
<a id="func-deleteprovider"></a>

```go
	func DeleteProvider[T any]()
```
<a id="func-deleteproviderbykey"></a>

```go
	func DeleteProviderByKey[T any](key ProviderKey) error
```
#### Other
<a id="func-listproviderkeys"></a>

```go
	func ListProviderKeys[T any]() []ProviderKey
```
<a id="func-getallprovider"></a>

```go
	func GetAllProvider[T any](opts ...ProviderGetOption) ([]T, error)
```
### Inject
<a id="func-injectstruct"></a>

```go
	func InjectStruct(target any) error
```
<a id="func-injectstructwithctx"></a>

```go
	func InjectStructWithCtx(ctx context.Context, target any) error
```
You can inject dependencies into struct fields using the `di:"..."` tag.<br>
<br>
Supported tag options:<br>
| Field     | Type               | Default | Description								|
|-----------|--------------------|---------|--------------------------------------------|
| `type`    | `"provider"` / `"value"` | `"value"` | The source type to inject from.	|
| `key`     | `string`           | `""`    | A string key used for lookup.				|
| `reload`  | `true` / `false`   | `false` | Only for providers.						|
| `optional`| `true` / `false`   | `false` | If true, injection is optional.			|
<br>
To skip injection for a field, use `di:"-"`.<br>
`Only exported (public) fields can be injected.`

<a id="func-injectfunc"></a>

```go
	func InjectFunc[T any](fn any, opts ...InjectFuncOption)
```
<a id="func-injectfuncwithctx"></a>

```go
	func InjectFuncWithCtx[T any](ctx context.Context, fn any, opts ...InjectFuncOption)
```

##### InjectFuncOption

<a id="func-withinjectfuncprovider"></a>

```go
	func WithInjectFuncProvider(variable string) InjectFuncOption
```

<a id="func-withinjectfunckey"></a>

```go
	func WithInjectFuncKey(variable string, key string) InjectFuncOption
```

<a id="func-withinjectfuncreload"></a>

```go
	func WithInjectFuncReload(variable string) InjectFuncOption
```

<a id="func-withinjectfuncoptional"></a>

```go
	func WithInjectFuncOptional(variable string) InjectFuncOption
```

Variable can be params index or variable name.

### Hook
<a id="func-afteradd"></a>

```go
	func AfterAdd(f AfterAddFunc)

	type AfterAddCtx struct {
		Type              reflect.Type
		ValueKey          *ValueKey
		ContainerValue    *containerValue
		ProviderKey       *ProviderKey
		ContainerProvider *containerProvider
	}

	type AfterAddFunc func(ctx AfterAddCtx)
```
This hook is triggered after a successful addition.<br>
You can determine whether the added item is a value or a provider by checking whether `ValueKey` or `ProviderKey` is non-nil.

<a id="func-afterproviderrun"></a>

```go
	func AfterProviderRun(f AfterProviderRunFunc)

	type AfterProviderRunCtx struct {
		Type              reflect.Type
		Key               ProviderKey
		ContainerProvider *containerProvider
		Value             any
	}

	type AfterProviderRunFunc func(ctx AfterProviderRunCtx)
```
This hook is triggered after the provider's factory function successfully returns a value—not just on the first invocation.<br>
If `NoCache` is set, the factory function runs every time, so this hook will be triggered on every `GetProvider` call.

<a id="func-afterfirstaccess"></a>

```go
	func AfterFirstAccess(f AfterFirstAccessFunc)

	type AfterFirstAccessCtx struct {
		Type              reflect.Type
		ValueKey          *ValueKey
		ContainerValue    *containerValue
		ProviderKey       *ProviderKey
		ContainerProvider *containerProvider
	}

	type AfterFirstAccessFunc func(ctx AfterFirstAccessCtx)
```
This hook is triggered on the first successful access of a value or provider.<br>
You can determine whether the added item is a value or a provider by checking whether `ValueKey` or `ProviderKey` is non-nil.

<a id="func-beforeduplicateregister"></a>

```go
	func BeforeDuplicateRegister(f BeforeDuplicateRegisterFunc)

	type BeforeDuplicateRegisterCtx struct {
		Type                 reflect.Type
		ValueKey             *ValueKey
		OldContainerValue    *containerValue
		NewContainerValue    *containerValue
		ProviderKey          *ProviderKey
		OldContainerProvider *containerProvider
		NewContainerProvider *containerProvider
		IsDefault            bool
	}

	type BeforeDuplicateRegisterFunc func(ctx BeforeDuplicateRegisterCtx) error
```
Invoked before a value or provider with the same type and key is overwritten.<br>
If an error is returned, the registration is aborted.<br>
You can determine whether the added item is a value or a provider by checking whether `ValueKey` or `ProviderKey` is non-nil.

### Global
<a id="func-setdefaultvaluekey"></a>

```go
	func SetDefaultValueKey(defaultValueKey string)
```
Default is empty string.
<a id="func-setdefaultproviderkey"></a>

```go
	func SetDefaultProviderKey(defaultProviderKey string)
```
Default is empty string.
<a id="func-setsafedelete"></a>

```go
	func SetSafeDelete(safeDelete bool)
```
Default is false. When set to true, every Value retrieval will increment the internal reference counter.<br>
Users must explicitly call [DeductRefCount](#func-deductrefcount) or [DeductRefCountByKey](#func-deductrefcountbykey) to decrement it when done.<br>
When safeDelete is enabled, Reset and Delete operations will block until the reference count reaches zero.
<a id="func-setresetmaxconcurrent"></a>

```go
	func SetResetMaxConcurrent(resetMaxConcurrent int)
```
Default is 100.
<a id="func-reset"></a>

```go
	func Reset(opts ...ResetOption) []error
```
Clears all registered values and providers.<br>
If `WithResetSkipOnClose` is not provided, `OnCloseHook` will be executed for each value/provider.<br>
Cleanup runs in parallel to speed up the process.

##### ResetOption

```go
	func WithResetSkipOnClose() ResetOption
```
---

## 🤔 Q&A

**Q: Can I stop duplicate register same value / component?**<br>
A: Yes. You may stop it using [BeforeDuplicateRegister](#func-beforeduplicateregister) by return error.<br>

**Q: Will Value's OnCloseHook trigger when replace?**<br>
A: No. But you can trigger through [BeforeDuplicateRegister](#func-beforeduplicateregister).

**Q: Are all actions transactional?**<br>
A: Not entirely.<br>
[Add](#func-add), [AddProvider](#func-addprovider), and [AddCtxProvider](#func-addctxprovider) are transactional — both the typed key and the default key will be added together or not at all.<br>
However, [GetAllProvider](#func-getallprovider) is **not** transactional. For example, if four uncached providers are triggered by [GetAllProvider](#func-getallprovider), and the first two succeed while the third fails, the first two will still be cached. There is no rollback mechanism.

**Q: Will Reset take a long time?**<br>
A: It depends. The duration is determined by the execution time of each OnCloseHook, but they run in parallel — and it won’t overload your CPU, so relax.

**Q: Why doesn't OnCloseHook receive the value as a parameter?**<br>
A: Because in practice, you already know what you need to handle when the `OnCloseHook` is triggered.<br>
Your code usually looks like this:

```go
dix.Add(TestKey, test, dix.WithValueOnClose(func() {
	// do something with test
}))
```
The variable is right next to the hook — you can just capture it directly.<br>
Also, it's tricky to pass the value into the hook cleanly, since the container doesn’t know the type and would have to pass it as `any`. 🥹

**Q: Why is the documentation a bit messy?**<br>
A: Because… I’ve never written documentation before. 😂 Be kind.

**Q: Why choose go version 1.20 not lower?**<br>
A: reflect.Type to satisfy comparable requires go1.20 or latercompilerInvalidTypeArg

## 💭 Future:
1. Should containerValue and containerProvider store their own Type and Key?
2. More test case.