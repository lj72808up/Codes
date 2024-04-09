package testRegx.dfa;

/**
 * Deterministic finite automata for recognising language of strings with
 * following conditions:
 *
 * 1. It must contain only :, ), (, and _ characters.
 * 2. :has to be followed by ) or (.
 * 3. Has to have at least one smiley :) or sad :( face.
 * 4. If string has more than one face, then all of them have to be separated by at least one _ character.
 * 5. Each string can start and end with zero or more _ characters.
 *
 * 去这个地址看状态转换图表 http://www.dennis-grinch.co.uk/tutorial/dfa-in-java
 * @author Dennis Grinch
 */
// 那么怎么挖掘 testRegx.dfa.DFA 所有可能的状态呢？一个常用的技巧是，用 "当前处理到字符串的哪个部分" 当作状态的表述
// 构建 testRegx.dfa.DFA 最关键的是找出状态集合
public class DFA {

    /**
     * Enum representing set of states Q. Each state has boolean
     * variable "accept" which indicates if it is an accept
     * state or not. All states with "true" accept value
     * represents subset of accept states F.
     * <p>
     * Every state has an instance variable corresponding to the
     * symbol in the subset of alphabet Σ. Thus "eyes" stands for :,
     * smile - for ), sad - for (, and space for _.
     *
     */
    private enum States {
        // Q0: 初始态  ,  Q1: 笑脸,哭脸开始的第一个:
        // Q2: 当前是一个笑脸,哭脸  ,  Q3: 完整的笑脸或哭脸后跟着分隔符
        // Q4: 破坏规则的非法状态
        Q0(false), Q1(false), Q2(true), Q3(true), Q4(false);

        // transition table
        static {
            Q0.eyes = Q1;  Q0.space = Q0;  // 即状态转换表  f(Q0,:)=Q1 ; f(Q0,_)=Q0
            Q1.smile = Q2; Q1.sad = Q2;
            Q2.space = Q3;
            Q3.eyes = Q1; Q3.space = Q3;
        }

        States eyes;
        States smile;
        States sad;
        States space;

        final boolean accept;

        // Constructor for sate. Every state is either accepting or not.
        States(boolean accept) {this.accept = accept;}

        /**
         * Represents transition function δ : Q × Σ → Q.
         *
         * @param ch is a symbol from alphabet Σ (Unicode)
         * @return state from codomain Q
         *
         */
        States transition(char ch) {
            switch (ch) {
                case ':':
                    return this.eyes == null? Q4 : this.eyes;
                case ')':
                    return this.smile == null? Q4: this.smile;
                case '(':
                    return this.sad == null? Q4: this.sad;
                case '_':
                    return this.space == null? Q4: this.space;
                default:
                    return Q4;
            }
        }

    }

    /**
     * Checks if testRegx.dfa.DFA accepts supplied string.
     *
     * @param string to check
     * @return true if it does accept, false otherwise
     */
    public boolean accept(String string) {
        States state = States.Q0;
        for (int i = 0; i < string.length(); i++) {
            state = state.transition(string.charAt(i));
        }
        return state.accept;
    }

    /**
     * 测试方法
     */
    public static void main(String[] args) {
        String str1 = "__:)_:(__";
        String str2 = "__:):(__";
        String str3 = ":))__";
        DFA dfa = new DFA();
        System.out.println(str1);
        System.out.println(dfa.accept(str1));
        System.out.println(str2);
        System.out.println(dfa.accept(str2));
        System.out.println(str3);
        System.out.println(dfa.accept(str3));
    }
}