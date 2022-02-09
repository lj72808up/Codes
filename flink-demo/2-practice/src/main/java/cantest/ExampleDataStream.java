package cantest;

import org.apache.flink.api.common.functions.FilterFunction;
import org.apache.flink.streaming.api.datastream.DataStreamSource;
import org.apache.flink.streaming.api.datastream.SingleOutputStreamOperator;
import org.apache.flink.streaming.api.environment.StreamExecutionEnvironment;

public class ExampleDataStream {
    public static void main(String[] args) throws Exception {
        StreamExecutionEnvironment env = StreamExecutionEnvironment.getExecutionEnvironment();
        DataStreamSource<Person> source = env.fromElements(
                new Person("zhangsan", 23),
                new Person("lisi", 24),
                new Person("wangwu",2)
        );
        SingleOutputStreamOperator<Person> filterStream = source.filter(new FilterFunction<Person>() {
            @Override
            public boolean filter(Person person) throws Exception {
                return person.getAge() >= 18;
            }
        });

        filterStream.print();   // 执行
        env.execute();
    }
}

