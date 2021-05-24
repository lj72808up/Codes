package SmallFlow

import (
	"AStar/models"
	"AStar/utils"
	"strings"
)

var ODS = "ods_ws_ad_pv_click_new"
var ODS_PC = "ods_pc_ad_pv_click_new"
var SOURCE = " AND source = 0\n"
var SOURCE_PC = " AND source = 4\n"

func GetABTestCase(flow models.FlowTaskDef) string {

	//var eRS string // Extend_reserve Result
	//if !flow.ExtendReserver.IsEmpty() {
	//	eRS = strings.Join([]string{" AND EXTEND_RESERVE_RESTRICTION(ext_reserve", "'" + flow.ExtendReserver.Rshift + "'", "'" + flow.ExtendReserver.And + "'", "'" + flow.ExtendReserver.Value + "') = true"}, ",")
	//} else {
	//	eRS = ""
	//}
	//
	//var sRS string // Style_reserve Result
	//if !flow.StyleReserver.IsEmpty() {
	//	sRS = strings.Join([]string{" AND STYLE_RESERVE_RESTRICTION(style_reserve", "'" + flow.StyleReserver.Pos + "'", "'" + flow.StyleReserver.Rshift + "'", "'" + flow.StyleReserver.And + "'", "'" + flow.StyleReserver.Value + "') = true"}, ",")
	//} else {
	//	sRS = ""
	//}

	return strings.Join([]string{"CASE WHEN FLOW_FLAG_RESTRICTION(flow_flag", "'" + flow.ExprFlow.Pos + "'", "'" + flow.ExprFlow.Value + "'", "'" + flow.ExprFlow.Type + "'" + ") = true"}, ",") +
	//eRS +
	//sRS +
		" THEN 'EXPR' " +
		strings.Join([]string{"WHEN FLOW_FLAG_RESTRICTION(flow_flag", "'" + flow.CtrlFlow.Pos + "'", "'" + flow.CtrlFlow.Value + "'", "'" + flow.CtrlFlow.Type + "'" + ") = true THEN 'CTRL' ELSE 'OTHER' END"}, ",")

	//CASE WHEN FLOW_FLAG_RESTRICTION(flow_flag,'0','0','1') = true
	// 		AND EXTEND_RESERVE_RESTRICTION(ext_reserve,'1','0','0',) = true
	// 		AND STYLE_RESERVE_RESTRICTION(style_reserve,'1','1','0','2',) = true
	// 		THEN 'EXPR'
	// WHEN FLOW_FLAG_RESTRICTION(flow_flag,'0','0','1') = true
	// 		THEN 'CTRL'
	// ELSE 'OTHER' END
	// AS ABTEST
}

func GetABTestWhere(flow models.FlowTaskDef) string {
	var eRS string
	if !flow.ExtendReserver.IsEmpty() {
		eRS = strings.Join([]string{" AND EXTEND_RESERVE_RESTRICTION(ext_reserve", "'" + flow.ExtendReserver.Rshift + "'", "'" + flow.ExtendReserver.And + "'", "'" + flow.ExtendReserver.Value + "'" + ") = true"}, ",")
	} else {
		eRS = ""
	}

	var sRS string

	if !flow.StyleReserver.IsEmpty() {
		sRS = strings.Join([]string{" AND STYLE_RESERVE_RESTRICTION(style_reserve", "'" + flow.StyleReserver.Pos + "'", "'" + flow.StyleReserver.Rshift + "'", "'" + flow.StyleReserver.And + "'", "'" + flow.StyleReserver.Value + "'" + ") = true"}, ",")
	} else {
		sRS = ""
	}

	return strings.Join([]string{"(FLOW_FLAG_RESTRICTION(flow_flag", "'" + flow.ExprFlow.Pos + "'", "'" + flow.ExprFlow.Value + "'", "'" + flow.ExprFlow.Type + "'" + ") = true"}, ",") +
		eRS +
		sRS +
		") OR " +
		strings.Join([]string{"(FLOW_FLAG_RESTRICTION(flow_flag", "'" + flow.CtrlFlow.Pos + "'", "'" + flow.CtrlFlow.Value + "'", "'" + flow.CtrlFlow.Type + "'" + ") = true)"}, ",")
}

func GetIPRestriction(condition string) string {
	if condition == "" {
		return ""
	}
	return "AND IP_RESTRICTION(flow_platform, '" + condition + "') = true "
}

func GetLocRestriction(condition string) string {
	if condition == "" {
		return "1 = 2"
	}
	return "LOCATION_RESTRICTION(reserved, '" + condition + "') = true"
}

// 全查询的SQL
func GetFullDataSQL(taskId string, taskType string, abTest string, abTestWhere string, ipRst string, locRst string) string {

	if taskType != "无线小流量任务" {
		ODS = ODS_PC
		SOURCE = SOURCE_PC
	}

	var sqlStr = "" +
		"SELECT " + ODS + ".DT                   AS RQ,\n" +
		"       '" + taskId + "'                 AS TaskId,\n" +
		"       ABTEST                           AS TAG,\n" +
		"       " + ODS + ".pvidpid as pv_id,\n" +
		"       IS_EXPR,\n" +
		"       PW,\n" +
		"       PV2,\n" +
		"       CHARGE,\n" +
		"       CLICK,\n" +
		"       ADPV2,\n" +
		"       ADCHARGE,\n" +
		"       ADCLICK,\n" +
		"       account_id                       AS KHID,\n" +
		"       CASE\n" +
		"         WHEN ISMOBILEQQ = 1 THEN 1\n" +
		"         ELSE 0\n" +
		"       END                              AS ISMOBILEQQ,\n" +
		"       CXC\n" +
		"FROM   (SELECT DT,\n" +
		"               pid,\n" +
		"               " + abTest + "           AS ABTEST,\n" +
		"               CASE\n" +
		"                 WHEN " + locRst + " THEN Pw(reserved)\n" +
		"                 ELSE 'other' \n" +
		"               END                      AS PW,\n" +
		"               pv_id,\n" +
		"               concat(pv_id,pid) as pvidpid,\n" +
		"               account_id,\n" +
		"               " + abTestWhere + "      AS IS_EXPR,\n" +
		"               query_word   AS CXC,\n" +
		"               pv2               AS ADPV2,\n" +
		"               charge_price AS ADCHARGE, \n" +
		"               charge_click AS ADCLICK\n" +
		"        FROM   adtl_biz." + ODS + "\n" +
		"        WHERE  " + ODS + ".DT = '{dt}'\n" +
		"               AND isruled = 1\n" +
		"               AND ad_id != '' \n" +
						SOURCE +
		"               )" + ODS + "\n" +
		"       INNER JOIN (SELECT concat(pv_id,pid) as pvidpid,\n" +
		"                          Sum(pv2)          AS pv2,\n" +
		"                          Sum(charge_price) AS CHARGE,\n" +
		"                          Sum(charge_click) AS CLICK\n" +
		"                   FROM   adtl_biz." + ODS + " tb1 \n" +
		"           WHERE  EXISTS (SELECT 1 \n" +
		"                          FROM   adtl_biz." + ODS + " \n" +
		"                          WHERE  dt = '{dt}' \n" +
		                                  SOURCE +
		"                                 AND ( " + abTestWhere + " )\n" +
		"                                 " + ipRst + " \n" +
		"                                 AND concat(tb1.pv_id,tb1.pid) = concat(pv_id,pid)) \n" +
		"                  AND dt = '{dt}' \n" +
		"                  AND isruled = 1\n" +
                           SOURCE +
		"           GROUP  BY pvidpid) pagepv \n" +
		"               ON " + ODS + ".pvidpid = pagepv.pvidpid\n" +
		"       LEFT JOIN (SELECT pid,\n" +
		"                         Str_like(second_channel, '手机QQ浏览器') AS ISMOBILEQQ\n" +
		"                  FROM   adtl_biz.dim_pid_da) dim_pid_da\n" +
		"              ON " + ODS + ".pid = dim_pid_da.pid\n"

	return sqlStr
}

// 同质查询的SQL
func GetSameQuerySQL(tb string) string {

	var sqlStr = "\n" +
		"SELECT *\n" +
		"FROM   " + tb + " a\n" +
		"WHERE  EXISTS (SELECT 1\n" +
		"               FROM   (SELECT CXC\n" +
		"                       FROM   " + tb + " \n" +
		"                       WHERE  TAG != 'OTHER'\n" +
		"                       GROUP  BY CXC\n" +
		"                       HAVING COUNT(DISTINCT TAG) = 2) b\n" +
		"               WHERE   a.CXC = b.CXC)" +
		"                       AND a.TAG != 'OTHER'"
	return sqlStr
}

func GetSmallFlowResSQL(taskId string, dateRange []string, isMobileQQ bool, isSame bool, isDtAvg bool, cusIds []string, query string) string {

	var mobileQQWhere = ""
	var tbName = "result_full"
	var adt = "A.dt,"
	var dt = "dt,"
	var dtJoin = " AND A.dt = B.dt"
	var acustomerId = "A.custom_id,"
	var customerId = "custom_id,"
	var customerJoin = " AND A.custom_id = B.custom_id"
	var customWhere = ""
	var styleWhere = ""
	var isPage = true
	var isStype = true

	if isMobileQQ {
		mobileQQWhere = " AND is_mobile_qq = true "
	}

	if isSame {
		tbName = "result_same"
	}

	if isDtAvg {
		dt = ""
		adt = ""
		dtJoin = ""
	}

	if len(cusIds) == 0 {
		tbName = tbName + "_only_c"
		customerId = ""
		customerJoin = ""
		acustomerId = ""
	}


	if len(cusIds) != 0 {
		customWhere = " AND custom_id IN ('" + strings.Join(cusIds, "','") + "')"
	}

	switch query {
	case "整体页面":
		isPage = true
		isStype = false
	case "实验样式-样式":
		isPage = false
		isStype = true
	case "实验样式-页面":
		isPage = true
		isStype = true
	}

	if isStype {
		styleWhere = " AND is_expr = true"
	} else {
		styleWhere = " AND is_expr = false"
	}

	if isPage {

		var sqlStr = "" +
			"SELECT  " + adt + "" +
			"        " + acustomerId + "" +
			"        A.tag," +
			"        pw," +
			"        sum_pv_page," +
			"        sum_pv_ad," +
			"        sum_charge," +
			"        sum_click," +
			"        avg_ctr_page," +
			"        avg_ctr_ad," +
			"        avg_acp," +
			"        avg_rpm_page" +
			" FROM (" +
			" SELECT " + dt + "" +
			"        task_id," +
			"        " + customerId + "" +
			"        CASE WHEN tag = 'CTRL' THEN '对照组' ELSE '实验组' END AS tag," +
			"        pw," +
			"        round(Sum(sum_pv_ad) / count(distinct dt),2)    AS sum_pv_ad," +
			"        round(Sum(sum_charge) / 100.0 / count(distinct dt),2)  AS sum_charge," +
			"        round(Sum(sum_click) / count(distinct dt),2)    AS sum_click," +
			"        100 * round(Sum(sum_click) / Sum(sum_pv_ad),4) AS avg_ctr_ad," +
			"        round(0.01 * Sum(sum_charge) / Sum(sum_click),4) AS avg_acp" +

			" FROM   " + tbName + "" +
			" WHERE  dt >= '" + dateRange[0] + "'" +
			"        AND tag != 'OTHER'" +
			"        AND dt <= '" + dateRange[1] + "'" +
			"        AND task_id = '" + taskId + "'" +
			"        " + mobileQQWhere + "" +
			"        " + customWhere + "" +
			"        " + styleWhere + "" +
			" GROUP  BY " + dt + "" +
			"           tag," + customerId +
			"           pw) A" +
			" JOIN " +
			" (SELECT " + dt + "" +
			"        task_id," +
			"        " + customerId + "" +
			"        CASE WHEN tag = 'CTRL' THEN '对照组' ELSE '实验组' END AS tag," +
			"        round(Sum(sum_pv_page) / count(distinct dt),2)  AS sum_pv_page," +
			"        100 * round(Sum(sum_click) / Sum(sum_pv_page),4) AS avg_ctr_page," +
			"        round(10.0 * Sum(sum_charge) / Sum(sum_pv_page),4) AS avg_rpm_page" +
			" FROM   " + tbName + "" +
			" WHERE  dt >= '" + dateRange[0] + "'" +
			"        AND tag != 'OTHER'" +
			"        AND pw = '全部'" +
			"        AND dt <= '" + dateRange[1] + "'" +
			"        AND task_id = '" + taskId + "'" +
			"        " + mobileQQWhere + "" +
			"        " + customWhere + "" +
			"        " + styleWhere + "" +
			" GROUP  BY " + dt + "" +
			"           " + customerId +
			"           tag) B" +
			" ON" +
			" A.task_id = B.task_id" +
			" AND A.tag = B.tag" + dtJoin + customerJoin
		return sqlStr
	} else {
		var sqlStr = "" +
		"SELECT " + dt + "" +
		"        task_id," +
			"        " + customerId + "" +
			"        CASE WHEN tag = 'CTRL' THEN '对照组' ELSE '实验组' END AS tag," +
			"        pw," +
			"        round(Sum(sum_pv_page) / count(distinct dt),2)  AS sum_pv_page," +
			"        round(Sum(sum_pv_ad) / count(distinct dt),2)    AS sum_pv_ad," +
			"        round(Sum(sum_charge) / 100.0 / count(distinct dt),2)  AS sum_charge," +
			"        round(Sum(sum_click) / count(distinct dt),2)    AS sum_click," +
			"        100 * round(Sum(sum_click) / Sum(sum_pv_page),4) AS avg_ctr_page," +
			"        100 * round(Sum(sum_click) / Sum(sum_pv_ad),4) AS avg_ctr_ad," +
			"        round(0.01 * Sum(sum_charge) / Sum(sum_click),4) AS avg_acp," +
			"        round(10.0 * Sum(sum_charge) / Sum(sum_pv_page),4) AS avg_rpm_page" +
			" FROM   " + tbName + "" +
			" WHERE  dt >= '" + dateRange[0] + "'" +
			"        AND tag != 'OTHER'" +
			"        AND dt <= '" + dateRange[1] + "'" +
			"        AND task_id = '" + taskId + "'" +
			"        " + mobileQQWhere + "" +
			"        " + customWhere + "" +
			"        " + styleWhere + "" +
			" GROUP  BY " + dt + "" +
			"           tag," + customerId +
			"           pw"
		return sqlStr
	}
}

func GetSmallFlowDownSQL(taskId string, dateRange []string, isMobileQQ bool, isSame bool, isDtAvg bool, cusIds []string, download string) (string, []string) {

	var mobileQQWhere = ""
	var sqlHeader = []string{""}
	var tbName = "adtl_john." + "small_flow_" + taskId

	//var adt = "A.dt,"
	//var dt = "dt,"
	var dtJoin = " AND A.dt = B.dt"

	//var acustomerId = "A.custom_id,"
	//var customerId = "custom_id,"
	//var customerJoin = " AND A.custom_id = B.custom_id"

	var cxc = ""
	var acxc = ""
	var cxcJoin = ""

	var isStyle = true

	var customWhere = ""
	var accountWhere = ""
	var styleWhere = ""
	var sameWhere = ""
	var dtWhere = " dt in ('" + strings.Join(utils.GetDateRange(dateRange[0], dateRange[1]), "','") + "')"

	var dataType = "'页面' AS data_type, "

	switch download {
	case "实验样式页面":
		isStyle = false
		isSame = false
	case "实验样式":
		isStyle = true
		isSame = false
	case "实验样式同质页面":
		isStyle = false
		isSame = true
	case "同质样式":
		isStyle = true
		isSame = true
	}

	if isMobileQQ {
		mobileQQWhere = " AND is_mobile_qq = true "
	}

	if isSame {
		sameWhere = " AND type='same' "
		cxc = " cxc, "
		acxc = " A.cxc, "
		cxcJoin = " AND A.cxc = B.cxc "
	} else {
		sameWhere = " AND type='full' "
	}
	/*

	if isDtAvg {
		dt = ""
		adt = ""
		dtJoin = ""
	}

	if len(cusIds) == 0 {
		customerId = ""
		customerJoin = ""
		acustomerId = ""
	}
*/
	if len(cusIds) != 0 {
		customWhere = " AND khid IN ('" + strings.Join(cusIds, "','") + "')"
		accountWhere = " AND account_id IN ('" + strings.Join(cusIds, "','") + "')"
	}


	if isStyle {
		styleWhere = " AND is_expr = true"
		dataType = "'样式' AS data_type, "
	}

	var sqlStr = "" +
		" SELECT " +
		"         A.dt, A.khid," +
		"         C.yjhy, C.ejhy, C.sjhy," + acxc +
		"         A.data_type as data_type_expr, A.ad_pv2 as ad_pv2_expr, A.ad_click as ad_click_expr, " +
		"         A.ad_charge as ad_charge_expr, A.ad_ctr as ad_ctr_expr, A.ad_acp as ad_acp_expr, A.ad_rpm as ad_rpm_expr," +
		"         B.data_type, B.ad_pv2, B.ad_click, B.ad_charge, B.ad_ctr, B.ad_acp, B.ad_rpm " +
		" FROM " +
		" (SELECT " +
		"         dt, khid, " + cxc + dataType +
		"         SUM(adpv2) AS ad_pv2," +
		"         SUM(adclick) AS ad_click," +
		"         SUM(adcharge) AS ad_charge," +
		"         FORMAT_NUMBER(100 * SUM(adclick) / SUM(adpv2), 4) AS ad_ctr," +
		"         FORMAT_NUMBER(0.01 * SUM(adcharge) / SUM(adclick), 4) AS ad_acp," +
		"         FORMAT_NUMBER(10.0 * SUM(adcharge) / SUM(adpv2), 4) AS ad_rpm " +
		" FROM " + tbName +
		" WHERE " + dtWhere + customWhere + styleWhere + mobileQQWhere + sameWhere +
		"         AND tag='EXPR' " +
		" GROUP BY " +
		"         dt," + cxc +
		"         khid) A" +
		" JOIN " +
		" (SELECT dt, khid, " + cxc + dataType +
		"         SUM(adpv2) AS ad_pv2," +
		"         SUM(adclick) AS ad_click," +
		"         SUM(adcharge) AS ad_charge," +
		"         FORMAT_NUMBER(100 * SUM(adclick) / SUM(adpv2), 4) AS ad_ctr," +
		"         FORMAT_NUMBER(0.01 * SUM(adcharge) / SUM(adclick), 4) AS ad_acp," +
		"         FORMAT_NUMBER(10.0 * SUM(adcharge) / SUM(adpv2), 4)  AS ad_rpm" +
		" FROM " + tbName + "" +
		" WHERE " + dtWhere + customWhere + styleWhere + mobileQQWhere + sameWhere +
		"         AND tag='CTRL' " +
		" GROUP  BY " +
		"         dt, " + cxc +
		"         khid) B" +
		" JOIN " +
		" (SELECT " +
		"         first_indus_name AS yjhy," +
		"         second_indus_name AS ejhy," +
		"         third_indus_name AS sjhy," +
		"         account_id AS khid " +
		" FROM adtl_biz.dim_account" +
		" WHERE " +
		"         1=1" + accountWhere +
		"         ) C" +
		" ON " +
		"         A.khid = B.khid AND A.khid = C.khid " + dtJoin + cxcJoin

	return sqlStr,sqlHeader
}

// about pc
func GetPcABTestCase(flow models.FlowTaskPcDef) string {

	return strings.Join([]string{" CASE WHEN EXPERIMENT_CONTROL_PC(style_reserve", "'" + flow.ExprExtendReserver.Pos + "'", "'" + flow.ExprExtendReserver.Rshift + "'", "'" + flow.ExprExtendReserver.And + "'", "'" + flow.ExprExtendReserver.Value + "'" + ") = true"}, "," ) +
		" THEN 'EXPR' " +
		strings.Join([]string{"WHEN EXPERIMENT_CONTROL_PC(style_reserve", "'" + flow.CtrlExtendReserver.Pos + "'", "'" + flow.CtrlExtendReserver.Rshift + "'", "'" + flow.ExprExtendReserver.And + "'", "'" + flow.ExprExtendReserver.Value + "'" + ") = true THEN 'CTRL' ELSE 'OTHER' END"}, ",")
}

func GetPcABTestWhere(flow models.FlowTaskPcDef) string {

	return strings.Join([]string{"(EXPERIMENT_CONTROL_PC(style_reserve", "'" + flow.ExprExtendReserver.Pos + "'", "'" + flow.ExprExtendReserver.Rshift + "'", "'" + flow.ExprExtendReserver.And + "'", "'" + flow.ExprExtendReserver.Value + "'" + ") = true"}, ",") +
		") OR " +
		strings.Join([]string{"(EXPERIMENT_CONTROL_PC(style_reserve", "'" + flow.CtrlExtendReserver.Pos + "'", "'" + flow.CtrlExtendReserver.Rshift + "'", "'" + flow.ExprExtendReserver.And + "'", "'" + flow.ExprExtendReserver.Value + "'" + ") = true)"}, ",")
}
