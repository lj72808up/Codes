import java.sql.Connection;
import java.sql.DriverManager;
import java.sql.SQLException;

public class TestJdbc {
    public static void main(String[] args) throws SQLException {
        Connection conn = DriverManager.getConnection(
                "jdbc:mysql://localhost:3306/adtl_test" ,
                "lj" ,
                "123456" );
        System.out.println(conn.isClosed());
    }
}
