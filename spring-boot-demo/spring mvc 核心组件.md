spring boot 使用 `DispatcherServletAutoConfiguration` 自动注入 DispatcherServlet, 其 `doService()` 方法处理
 
#### 如何使用 `DispatcherServlet` 处理 request
1. 请求被 `DispatcherServlet` 捕获
2. `DispatcherServlet` 解析请求 URL, 得到资源标识符 URI, 根据 URI 获取处理器数组 `[]HandlerExecutionChain` 
    过程: 从 RequestMappingHandler 中获取 handler (从 AbstractHandlerMethodMapping 的 MappingRegistry 中获取. 
          启动时会根据 @RequestMapping 注解反射 url 和 HandlerMethod 注册到 MappingRegistry 中)
3. `DispatcherServlet` 根据获取的 Handler, 选择一个合适的 `HandlerAdapter`
4. `DispatcherServlet` 提取 request 中的参数, 填充 Handler(`Controller`) 的形参并执行, 然后返回结果