package testRegx.myRegexEngineer.nfa.strategy;

public class EpsilonMatchStrategy extends MatchStrategy {
    @Override
    public boolean isMatch(char c, String edge) {
        return true;
    }
}
