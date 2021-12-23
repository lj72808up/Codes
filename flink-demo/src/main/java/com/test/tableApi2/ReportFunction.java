package com.test.tableApi2;

import org.apache.flink.table.api.Table;
import org.apache.flink.table.expressions.TimeIntervalUnit;

import static org.apache.flink.table.api.Expressions.*;

public class ReportFunction {
    /**
     * 最基本的 report 方法，
     */
    public static Table reportBase(Table transactions) {
        return transactions.select(
                        $("account_id"),
                        $("transaction_time").floor(TimeIntervalUnit.HOUR).as("log_ts"),
                        $("amount"))
                .groupBy($("account_id"), $("log_ts"))
                .select(
                        $("account_id"),
                        $("log_ts"),
                        $("amount").sum().as("amount"));
    }
}
