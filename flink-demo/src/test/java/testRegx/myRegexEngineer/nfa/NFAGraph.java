package testRegx.myRegexEngineer.nfa;

import testRegx.myRegexEngineer.common.Constant;

/**
 * 表示一个 NFAGraph 子图
 * 该类实现了，一个子图和另一个子图的“ |,*,?,() ”符号操作
 * 这部分参看 <a href="https://blog.csdn.net/xindoo/article/details/105875239">状态机下的正则表达式</a> 图例部分
 */
public class NFAGraph {
    public NFAState start;
    public NFAState end;

    public NFAGraph(NFAState start, NFAState end) {
        this.start = start;
        this.end = end;
    }

    // 2个 NFA 图组成'或'的关系
    public void addParallelGraph(NFAGraph NFAGraph) {
        NFAState newStart = new NFAState(); // 新图的 start 节点
        NFAState newEnd = new NFAState();   // 新图的 end 节点
        newStart.addNext(Constant.EPSILON, this.start);
        newStart.addNext(Constant.EPSILON, NFAGraph.start);
        this.end.addNext(Constant.EPSILON, newEnd);
        NFAGraph.end.addNext(Constant.EPSILON, newEnd);
        this.start = newStart;
        this.end = newEnd;
    }

    // 2个 NFA 子图串联
    public void addSeriesGraph(NFAGraph NFAGraph) {
        this.end.addNext(Constant.EPSILON, NFAGraph.start);
        this.end = NFAGraph.end;
    }

    // * 重复0-n次 (星号 = 问号 + 加号)
    public void repeatStar() {
        repeatPlus();
        addSToE(); // 重复0
    }

    // ? 重复0次 => 加一条 start 到 end 的 epsilon 边；不加新 start,end 节点
    public void addSToE() {
        start.addNext(Constant.EPSILON, end);
    }

    // + 重复1-n次： 加一条 end->start 的 epsilon 边; 最后连上新的 start,end 节点
    public void repeatPlus() {
        NFAState newStart = new NFAState();  // 新图的 start 节点
        NFAState newEnd = new NFAState();    // 新图的 end 节点
        newStart.addNext(Constant.EPSILON, this.start);
        end.addNext(Constant.EPSILON, newEnd);
        end.addNext(Constant.EPSILON, start);
        this.start = newStart;
        this.end = newEnd;
    }

}
