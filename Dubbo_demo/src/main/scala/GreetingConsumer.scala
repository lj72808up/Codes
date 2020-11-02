import org.apache.dubbo.config.ApplicationConfig
import org.apache.dubbo.config.ReferenceConfig
import org.apache.dubbo.config.RegistryConfig
import com.test.api.GreetingsService


object GreetingConsumer {
  def main(args: Array[String]): Unit = {

    val reference = new ReferenceConfig[GreetingsService]()
    reference.setApplication(new ApplicationConfig("first-dubbo-consumer"))
//    reference.setRegistry(new Nothing("zookeeper://local":2181"))
    reference.setUrl("dubbo://127.0.0.1:9090")
    reference.setInterface(classOf[GreetingsService])
    val service = reference.get
    val message = service.sayHi("dubbo")
    System.out.println(message)
  }
}
