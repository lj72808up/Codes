import java.sql.Connection;
import java.sql.DriverManager;
import java.sql.PreparedStatement;
import java.sql.ResultSet;

public class TestJdbc {
    public static void main(String[] args) throws Exception {
//        Class.forName("com.mysql.cj.jdbcz.Driver") ; // 这句话在 jdbc 4.0 以后已经可以省略
        String url = "jdbc:mysql://localhost:3306/test" ;
        String username = "root" ;
        String password = "root" ;
        Connection con = DriverManager.getConnection(url , username , password ) ;

        PreparedStatement pstmt = con.prepareStatement("select v0 fromm casbin_rule;") ;
        ResultSet rs = pstmt.executeQuery();
    }
}
