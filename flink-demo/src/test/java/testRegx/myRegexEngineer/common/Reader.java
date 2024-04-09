package testRegx.myRegexEngineer.common;

public class Reader {
    private int cur = 0;
    private char[] chars;
    public Reader(String regex) {
        this.chars = regex.toCharArray();
    }
    public char peak() {
        if (cur == chars.length) {
            return '\0';
        }
        return chars[cur];
    }

    public char next() {
        if (cur == chars.length) {
            return '\0';
        }
        return chars[cur++];
    }

    public boolean hasNext() {
        return cur < chars.length;
    }

    public String getSubRegex(Reader reader) {
        int cntParem = 1;
        String regex = "";
        while (reader.hasNext()) {
            char ch = reader.next();
            if (ch == '(') {    // 对左右括号的计数：遇到左括号+1，遇到右括号-1
                cntParem++;
            } else if (ch == ')') {
                cntParem--;
                if (cntParem == 0) {  // 减到0说明左右括号全都匹配上了
                    break;
                } else {
                }
            }
            regex += ch;
        }
        return regex;
    }

    public String getRemainRegex(Reader reader) {
        String regex = "";
        while (reader.hasNext()) {
            regex += reader.next();
        }
        return regex;
    }
}
