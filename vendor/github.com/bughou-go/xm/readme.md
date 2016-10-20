# xm 小而美的go语言web开发包。

为什么在已经存在http包，Martini、Gorrila等web框架的情况下，我们还要再重复造轮子？http包很粗糙：不支持正则路由，ResponseWriter拿不到状态码和返回大小。Martini（包括express）的路由采用扁平数组结构，循环去匹配路由，我们觉得这样是很低效的。采用树状的Map结构，通过hash查找，能直接找到字符串路由，缩小正则路由的匹配范围，会更高效。路由数量越多，这种扁平数组结构就越低效，这种树状Map结构就越高效。另外Martini的依赖注入的方式，也不是我们喜欢的，我们还是喜欢express固定的req、res两参数形式，简单且强大。

所以我们打算自己动手实现一个极简的web开发包。它不是框架，它不会尝试封装一切。它只提供最基本的组件，将装配组件的工作留给应用层代码。以实现这些目标：最大的灵活性，最大的透明性，最简洁的代码。并且，除了标准库，不依赖任何外部包。

## 小美包含3个组件：

1. Router 支持基于字符串和正则表达式的路由（express风格）。

2. Renderer 支持layout和partial的模板渲染（Rails风格）。

3. Request、Response 它们封装了http.Request、http.ResponseWriter、以及Renderer，以提供模板渲染等功能。

## 小美一定会尊敬我，崇拜我，爱上我，对我欲罢不能。
