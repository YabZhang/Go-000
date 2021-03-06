### 学习笔记

本周课程是讲 Go 语言中的 `error`。

首先明确 `error` 在 Go 语言中只是一个接口声明，并和其他语言中的错误异常 `Excepton` 机制进行了对比。在 Go 语言中 `error` 是作为一个 *数值* 来传递，并且直接用来表示执行结果的正确性。尤其在处理多返回值时，首先需要检查 `error` 来确定结果是有效和正确的。此外也不应该基于 `error` 做任何隐藏的控制流 —— 出现 `error` 即意味着发生某种错误，结果值的有效性和正确性不再有任何保证。

然后罗列了集中 Go 语言中常见的错误类型的声明：

1. `Sentinel Error` (预定义错误), 类似于在包内预先声明的常量错误类型。通常要支持具体类型导出，便于等值判断，但这会制造严重的耦合。

2. `Error Types` (自定义错误类型), 用于提供更多的错误上下文，如代码行的文件和行号。相比于前者可以提供更多的上下文，但仍然需要支持错误类型导出和制造强耦合。

3. `Opaque errors` (模糊错误), 只返回错误而不假设内容，并提供方法来支持对错误类型的判断。这种方法的好处是降低耦合，不提供显示的类型或值，而是通过约定的接口或者方法来做行为判断。

接下来介绍了几种方法来更舒服地处理错误。

1. `Indented flow is for errors`, 减少错误判断的嵌套。
   
2. `Eliminate error handing by eliminating errors`, 减少冗余逻辑(有错误直接返回)，选择更好的代码实现和封装(如：Scanner, ErrWriter等)。

3. `Wrap errors`, 包裹错误。介绍了 `pkg/errors` 包通过巧妙地封装原始错误、调用栈和错误信息，然后来减少多处打log，提供完整错误信息等。

最后针对 `Wrap errors` 的使用场景总结了几条规则：

1. 高可重用的基础库不要包裹错误，而应该提供原始错误。

2. 如果不处理错误，那么就包裹错误并返回调用。

3. 一但错误被处理，那么久不能向上层调用栈返回错误；错误只应该被处理一次。

Go 1.13 借鉴了 `pkg/errors` 的实现，新增了错误包裹相关的特性，如新增函数`Is`和`As` 用于检查错误，新增`%w`谓词打印错误附加信息。但当前版本依然没有提供错误的调用栈信息。

此外介绍了 Go 2 对于 `error` 的一些改进的提议。

---


### 作业题目

> 我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？

DAO 层用于实现具体的数据访问细节，对于错误应该用 `Wrap` 包含调栈向上层返回，然后由上层来根据业务做错误处理。为了在上层能够在处理底层跑上来的错误时减轻耦合，可以通过 `Sentinel Error` 或者 `Opaque errors` 的方式来处理错误，具体就是定义新的结构来对错误做判定，这里使用重新定义 `RowNotFound` 的错误来取代原有的 `ErrNotRows`.