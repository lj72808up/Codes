package testRegx.myRegexEngineer;

public class RegexTest {

//    private static Regex regex = Regex.compile("a(b|c)*");
//    private static String[] strs = {"ac", "acc", "a", "a   bcccdb", "ab", "abcd", "a3abcd", "a33333defd", "aabcabcabcabcabcabcdb",
//            "abbbbbbbbb", "acccccccbad", "acccccccccccccccccccccccccb", "abbbbbbbbbbbbbbbc"};
//
//    @Benchmark
//    @Measurement(iterations = 2)
//    @Threads(1)
//    @Fork(0)
//    @Warmup(iterations = 0)
//    public void nfa() {
//        for (String str : strs) {
//            regex.isMatch(str);
//        }
//    }
//
//    @Benchmark
//    @Measurement(iterations = 2)
//    @Threads(1)
//    @Fork(0)
//    @Warmup(iterations = 0)
//    public void dfaRecursion() {
//        for (String str : strs) {
//            regex.isMatch(str, 1);
//        }
//    }
//
//    @Benchmark
//    @Measurement(iterations = 2)
//    @Threads(1)
//    @Fork(0)
//    @Warmup(iterations = 0)
//    public void dfaNonRecursion() {
//        for (String str : strs) {
//            regex.isDfaMatch(str);
//        }
//    }

    public static void main(String[] args) {
//        test();     // 测试匹配结果
        printNFATest();   // 测试 NFA 打印
//        Options options = new OptionsBuilder().include(RegexTest.class.getSimpleName()).build();
//        try {
//            new Runner(options).run();
//        } catch (Exception e) {
//            System.out.println(e.fillInStackTrace());
//        } finally {
//            System.out.println("finshed");
//        }
    }

    private static void test() {
        String str = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaab";
        Regex regex1 = Regex.compile("a*aaaaaaaaaaaaaaaaaaaaaab");
        System.out.println(regex1.isDfaMatch(str)); // dfa 正则
        System.out.println("_________________");
        System.out.println(regex1.isMatch(str));    // nfa 正则
        System.out.println("_________________");
//        Regex regex = Regex.compile("a(b|c)*");
//        List<String> res = regex.match("aabacabbbcaccc");
//        regex.printNfa();
//        System.out.println("");
//        regex.printDfa();
//
//        System.out.println(regex.isMatch("ac"));
//        System.out.println(regex.isMatch("acc"));
//        System.out.println(regex.isMatch("a"));
//        System.out.println(regex.isMatch("a   bcccdb"));
//        System.out.println(regex.isMatch("ab"));
//        System.out.println(regex.isMatch("abcd"));
//        System.out.println(regex.isMatch("a3abcd"));
//        System.out.println(regex.isMatch("a33333defd"));
//        System.out.println(regex.isMatch("aabcabcabcabcabcabcdb"));
//
//        System.out.println("*********");
//        System.out.println(regex.isDfaMatch("ac"));
//        System.out.println(regex.isDfaMatch("acc"));
//        System.out.println(regex.isDfaMatch("a"));
//        System.out.println(regex.isDfaMatch("a   bcccdb"));
//        System.out.println(regex.isDfaMatch("ab"));
//        System.out.println(regex.isDfaMatch("abcd"));
//        System.out.println(regex.isDfaMatch("a3abcd"));
//        System.out.println(regex.isDfaMatch("a33333defd"));
//        System.out.println(regex.isDfaMatch("aabcabcabcabcabcabcdb"));
//
//        System.out.println("*********");
//        System.out.println(regex.isMatch("ac", 1));
//        System.out.println(regex.isMatch("acc", 1));
//        System.out.println(regex.isMatch("a", 1));
//        System.out.println(regex.isMatch("a   bcccdb", 1));
//        System.out.println(regex.isMatch("ab", 1));
//        System.out.println(regex.isMatch("abcd", 1));
//        System.out.println(regex.isMatch("a3abcd", 1));
//        System.out.println(regex.isMatch("a33333defd", 1));
//        System.out.println(regex.isMatch("aabcabcabcabcabcabcdb", 1));
    }

    /**
     * 测试打印正则对应的 NFA
     */
    private static void printNFATest(){
        Regex reg = Regex.compile("a(b|c)*");
        reg.printNfa();
    }
}
