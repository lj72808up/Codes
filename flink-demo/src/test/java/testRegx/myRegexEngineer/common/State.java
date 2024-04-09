package testRegx.myRegexEngineer.common;

import java.util.HashMap;
import java.util.HashSet;
import java.util.Map;
import java.util.Set;

/**
 * State 表示 NFA 里的一个节点
 */
public class State {
    protected static int idCnt = 0;
    protected int id;
    protected StateType stateType;

    public State() {
        this.id = idCnt++;
        this.stateType = StateType.GENERAL;
    }

    /**
     * kv：<边，转以后的状态列表>
     * 边就是输入字符
     */
    public Map<String, Set<State>> next = new HashMap<>();

    /**
     * 为某个节点一个 <边,状态>
     */
    public void addNext(String edge, State nfaState) {
        Set<State> set = next.get(edge);
        if (set == null) {
            set = new HashSet<>();
            next.put(edge, set);
        }
        set.add(nfaState);
    }

    /**
     * 节点类型：判断是否是最终状态使用
     */
    public void setStateType(StateType stateType) {
        this.stateType = stateType;
    }

    public boolean isEndState() {
        return stateType == StateType.END;
    }

    /**
     * 状态id
     */
    public int getId() {
        return this.id;
    }
}
