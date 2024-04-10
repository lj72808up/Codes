#### DFA
dfa 包下的 `DFA.java` 是一个外国小哥写的 blog，是一个简单的匹配 smail，sad 表情的有限状态机 demo
* 为什么状态机能用作字串匹配？ 
    因为一个 dfa 就定义了一种状态转移模式。如果把输入的字符串看做字符流，当前字符和下一个字符就构成了状态转移。

#### NFA 
myRegexEngineer 包下是一个正则表达式引擎  
* [https://blog.csdn.net/xindoo/article/details/105875239](https://blog.csdn.net/xindoo/article/details/105875239)
1. nfa 图例符号：
   * -> 加 ‘单圆环’ => 起始态
   * 双圆环 => 接收态
   * ϵ闭包：不需要任何条件就能转移状态的边
2. 正则表达式的每个循环，或，嵌套组匹配， 都有固定的 nfa 图例模式表示。所以只要写出这个正则表达式对应的 NFA 类，就能实现正则匹配。   
  这个 regexExp -> NFA 的转换，在这里是人工执行的    

#### DFA
编程的角度,NFA和DFA都是有向图. 所有的 NFA 都有一个其等价的 DFA. 因为 NFA 在一个状态输入后, 可以转换成多个不同的状态, 所以 NFA 在匹配失败时要回溯,
