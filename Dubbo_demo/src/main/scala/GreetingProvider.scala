import java.util.concurrent.CountDownLatch

import com.test.api.{GreetingsService, GreetingsServiceImpl}
import org.apache.dubbo.config.{ApplicationConfig, ProtocolConfig, RegistryConfig, ServiceConfig}

object GreetingProvider {
  def main(args: Array[String]): Unit = {
    val service = new ServiceConfig[GreetingsService]()
    service.setApplication(new ApplicationConfig("first-dubbo-provider"))

    val registryConfig = new RegistryConfig("zookeeper://localhost:2181")
    registryConfig.setRegister(false)
    service.setRegistry(registryConfig)

    val protocolConfig = new ProtocolConfig("dubbo",9090)
    service.setProtocol(protocolConfig)
    service.setInterface(classOf[GreetingsService])
    service.setRef(new GreetingsServiceImpl())
    service.export();

    System.out.println("dubbo service started")
    new CountDownLatch(1).await();
  }
}
