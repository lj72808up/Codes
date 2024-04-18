package testRegx.myRegexEngineer;

public class RegexTest {
    public static void main(String[] args) {
        test();     // 测试匹配结果
//        printNFATest();   // 测试 NFA 打印
    }

    private static void test() {
        // 这里的匹配不同于 java 的正则表达式,只要匹配出子串就返回 true. 这里要匹配到最后一个字符满足才返回 true
        String str = "aacba";
        Regex regex1 = Regex.compile("a*(b|c)");
        System.out.println(regex1.isDfaMatch(str)); // dfa 正则
        System.out.println("_________________");
        System.out.println(regex1.isMatch(str));    // nfa 正则
        System.out.println("_________________");
    }

    /**
     * 测试打印正则对应的 NFA
     */
    private static void printNFATest(){
        Regex reg = Regex.compile("a(b|c)*");
        reg.printNfa();
    }
}
