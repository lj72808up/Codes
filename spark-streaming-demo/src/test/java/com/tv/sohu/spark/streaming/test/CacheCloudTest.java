package com.tv.sohu.spark.streaming.test;

import com.tv.sohu.spark.streaming.dm3.online.RedisClusterOperator;
import org.apache.commons.lang.StringUtils;

import java.util.Scanner;

/**
 * Created by wentaoli213587 on 2021/4/13.
 */
public class CacheCloudTest {

    public static void main(String[] args) {

        RedisClusterOperator operator = new RedisClusterOperator(10684L);

        Scanner scanner = new Scanner(System.in);
        while (true) {
            System.out.print("command>");
            String str = scanner.nextLine();
            if (StringUtils.isBlank(str))
                continue;
            String[] arr = str.split("\\s+");
            if (arr[0].equals("exit"))
                break;
            try {
                switch (arr[0]) {
                    case "set":
                        System.out.println(operator.set(arr[1], arr[2]));
                        break;
                    case "setex":
                        System.out.println(operator.setex(arr[1], arr[2], Integer.parseInt(arr[3])));
                        break;
                    case "get":
                        System.out.println(operator.get(arr[1]));
                        break;
                    case "del":
                        System.out.println(operator.delete(arr[1]));
                        break;
                    case "exists":
                        System.out.println(operator.exists(arr[1]));
                        break;
                    case "ttl":
                        System.out.println(operator.ttl(arr[1]));
                        break;
                    case "help":
                        System.out.println("set <uid>\nget <uid>\ndel <uid>\nexists<uid>\nttl <uid>");
                        break;
                    default:
                        System.out.println("error command .....");
                }
            } catch (Exception e) {
                e.printStackTrace();
            }
        }
        operator.close();
    }


}
