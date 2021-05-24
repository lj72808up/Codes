package test

import (
	"fmt"
	"strings"
	"testing"
)

func TestScript(t *testing.T) {

	menuIds := []int{}
	for i := 0; i <= 24; i++ {
		menuIds = append(menuIds, i)
	}
	/*roles := []string{"波动总权限", "波动KA客户", "波动中小所有客户", "波动中小-华北客户", "波动中小-华东客户", "波动中小-华南客户", "波动中小-华中客户",
	"波动单个客户", "波动搜索渠道", "波动网盟渠道"}*/

	roles := []string{"波动客户&行业"}
	for _, menuId := range menuIds {
		for _, role := range roles {
			fmt.Printf("INSERT INTO bodong_menu_mapping (group_name, menu_id) VALUE ('%s',%d);\n", role, menuId)
		}
	}
}

func TestCasbinUser(t *testing.T) {
	content := `
1	lanyangyang@sogou-inc.com	兰阳阳
1	liujie02@sogou-inc.com	刘捷
1	wangkainuo@sogou-inc.com	王锴诺
1	wuzhonghai@sogou-inc.com	吴忠海
1	ligang@sogou-inc.com	李刚
1	fenghaiyu@sogou-inc.com	冯海钰
1	zhangshuguang@sogou-inc.com	张曙光
1	lijiguang@sogou-inc.com	李继光
1	lidanyang215017@sogou-inc.com	lidanyang215017
1	zhaitianzhi@sogou-inc.com	翟添智
1	kongyibo@sogou-inc.com	孔一博
1	zhangzhe@sogou-inc.com	张哲
1	yuhancheng@sogou-inc.com	于瀚程
1	wuzhaoren@sogou-inc.com	吴照人
1	liqian@sogou-inc.com	李倩
1	xuejingjing@sogou-inc.com	薛静静
1	zhouhui@sogou-inc.com	周慧
1	yangxiyu@sogou-inc.com	杨锡玉
1	hongtao@sogou-inc.com	洪涛
1	xuchenyang@sogou-inc.com	徐晨阳
1	songyu@sogou-inc.com	songyu
1	hexiaoyan@sogou-inc.com	何筱妍
1	guoxing@sogou-inc.com	张国兴
1	xulixian@sogou-inc.com	徐礼贤
1	wubo@sogou-inc.com	吴博
1	liyaodong@sogou-inc.com	李耀栋
1	pengchao@sogou-inc.com	彭超
1	hanmingyang@sogou-inc.com	韩明阳
1	chenyuhan@sogou-inc.com	陈禹含
1	yaowangqi@sogou-inc.com	yaowangqi
1	lizhaoyu@sogou-inc.com	李钊宇
1	sunao@sogou-inc.com	孙奥
1	zhengpengcheng@sogou-inc.com	郑鹏程
1	liuziwei@sogou-inc.com	liuziwei
1	huozhigang@sogou-inc.com	huozhigang
1	rensongchaosi2011@sogou-inc.com	rsc
1	jinghongpeng@sogou-inc.com	jinghongpeng
1	dengmengjiao@sogou-inc.com	邓梦娇
1	likang216651@sogou-inc.com	likang216651
1	wangzhulei@sogou-inc.com	王朱磊
1	guanweiqiu@sogou-inc.com	关维秋
1	zhaolu@sogou-inc.com	赵璐
1	yangguotao@sogou-inc.com	杨国涛
1	chenhuidu@sogou-inc.com	陈辉都
1	xuchunxiao@sogou-inc.com	许春晓
1	huangxiaoyan@sogou-inc.com	黄晓燕
1	caodi@sogou-inc.com	曹递
1	luyuan212521@sogou-inc.com	陆远
1	zhaoxiaole@sogou-inc.com	zhaoxiaole
1	yajuanwang@sogou-inc.com	王雅娟
1	lijing215119@sogou-inc.com	李静
2	huangguohua@sogou-inc.com	黄国华
2	liangying210795@sogou-inc.com	梁莹
2	lixiaoye@sogou-inc.com	李晓叶
2	lixiang206703@sogou-inc.com	李响
2	lixi@sogou-inc.com	李玺
2	zhuweiwei@sogou-inc.com	朱卫卫
2	lihui209397@sogou-inc.com	李慧
2	yangcaiyun@sogou-inc.com	杨彩云
2.1	wangdan215543@sogou-inc.com	wangdan215543
2.1	gongzhi215448@sogou-inc.com	龚志
2.1	sunyan@sogou-inc.com	孙岩
2.1	feisun@sogou-inc.com	孙霏
2.1	sunxue216911@sogou-inc.com	孙雪
2.1	wangbing@sogou-inc.com	王冰
2.1	liuxingya@sogou-inc.com	刘兴雅
2.1	huangyang213158@sogou-inc.com	黄杨
2.1	xueyingying@sogou-inc.com	薛莹莹
2.1	wanglinlin@sogou-inc.com	王琳琳
2.1	jessicazhang@sogou-inc.com	张琳
2.1	gaozhiguang@sogou-inc.com	高志广
2.1	yingkang@sogou-inc.com	康颖
2.1	zhangxinming@sogou-inc.com	张新明
2.1	taochen@sogou-inc.com	陈涛
2.1	huachaoyi@sogou-inc.com	华超逸
2.1	huangxiuhua@sogou-inc.com	黄秀花
2.1	liheng@sogou-inc.com	李恒
2.1	lining214305@sogou-inc.com	李宁
2.1	caodi@sogou-inc.com	曹递
2.1	jixin@sogou-inc.com	季鑫
2.1	zhoujiaxing@sogou-inc.com	周佳兴
2.2	yanghaipeng@sogou-inc.com	杨海朋
2.2	maxue@sogou-inc.com	马雪
2.2	weixin@sogou-inc.com	魏欣
2.2	wangyao214441@sogou-inc.com	王瑶
2.2	wangxi209754@sogou-inc.com	王希
2.2	renhongxiang@sogou-inc.com	任宏翔
2.2	nemoyang@sogou-inc.com	杨颖
2.2	baoshansun@sogou-inc.com	孙宝珊
2.2	limin@sogou-inc.com	李敏
2.2	wanghuanhuan@sogou-inc.com	王焕焕
2.2	liwei213341@sogou-inc.com	李伟
2.2	zhouhongjie@sogou-inc.com	周宏杰
2.2	xuhao212985@sogou-inc.com	徐浩
2.2	sunqi@sogou-inc.com	孙琪
2.2	liuye@sogou-inc.com	刘晔
3.1	liaorui205984@sogou-inc.com	廖锐
3.1	weneryan@sogou-inc.com	温二艳
3.1	qianziqiang@sogou-inc.com	钱自强
3.1	lihuabei@sogou-inc.com	李华北
3.1	wanghua@sogou-inc.com	王画
3.1	guowanyi@sogou-inc.com	郭婉怡
3.1	zhanglu214154@sogou-inc.com	张璐
3.1	zhaodansi4231@sogou-inc.com	赵丹
3.1	yangning@sogou-inc.com	杨宁
3.1	lixiumei@sogou-inc.com	李秀梅
3.1	fubing@sogou-inc.com	伏冰
3.1	mahanmei@sogou-inc.com	马寒梅
3.1	lianzhao@sogou-inc.com	连昭
3.1	zhouhui@sogou-inc.com	周慧
3.1	chennan@sogou-inc.com	陈楠
3.1	huangtao220022@sogou-inc.com	黄涛
3.1	zhangjinyi@sogou-inc.com	张晋毅
3.1	mazhisi2054@sogou-inc.com	马枝
3.1	houtingting@sogou-inc.com	侯婷婷
3.1	gengjingle@sogou-inc.com	耿京乐
3.1	ruanliuhui@sogou-inc.com	阮浏慧
3.1	zhufangfang@sogou-inc.com	朱芳芳
3.1	sunxiaolei@sogou-inc.com	孙晓磊
3.1	zhaohongxue@sogou-inc.com	赵宏雪
3.1	yiming@sogou-inc.com	易鸣
3.1	wudongyang@sogou-inc.com	吴东洋
3.1	lianghao@sogou-inc.com	梁浩
3.1	sunmaowei@sogou-inc.com	孙茂伟
3.1	dudong@sogou-inc.com	杜东
3.1	heyukai@sogou-inc.com	贺宇凯
3.1	maqian@sogou-inc.com	马茜
3.1	luyuan212521@sogou-inc.com	陆远
3.1	shupeng203672@sogou-inc.com	舒鹏
3.1	qipeng@sogou-inc.com	齐鹏
3.1	chenxudong@sogou-inc.com	陈绪东
3.1	zhaomeng@sogou-inc.com	赵萌
3.2	liaorui205984@sogou-inc.com	廖锐
3.2	qianziqiang@sogou-inc.com	钱自强
3.2	lihuabei@sogou-inc.com	李华北
3.2	wanghua@sogou-inc.com	王画
3.2	guowanyi@sogou-inc.com	郭婉怡
3.2	zhanglu214154@sogou-inc.com	张璐
3.2	zhaodansi4231@sogou-inc.com	赵丹
3.2	yangning@sogou-inc.com	杨宁
3.2	lixiumei@sogou-inc.com	李秀梅
3.2	fubing@sogou-inc.com	伏冰
3.2	mahanmei@sogou-inc.com	马寒梅
3.2	lianzhao@sogou-inc.com	连昭
3.2	zhouhui@sogou-inc.com	周慧
3.2	chennan@sogou-inc.com	陈楠
3.2	huangtao220022@sogou-inc.com	黄涛
3.2	mazhisi2054@sogou-inc.com	马枝
3.2	houtingting@sogou-inc.com	侯婷婷
3.2	gengjingle@sogou-inc.com	耿京乐
3.2	ruanliuhui@sogou-inc.com	阮浏慧
3.2	maqian@sogou-inc.com	maqian
3.2	sunmaowei@sogou-inc.com	sunmaowei
3.2	zhangjinyi@sogou-inc.com	zhangjinyi
3.2	wangzhulei@sogou-inc.com	王朱磊
3.2	shupeng203672@sogou-inc.com	shupeng203672
3.2	luyuan212521@sogou-inc.com	luyuan212521
3.2	yangjie01@sogou-inc.com	杨婕
3.2	chenxudong@sogou-inc.com	陈绪东
3.2	mayongbin@sogou-inc.com	马勇斌
3.3	liaorui205984@sogou-inc.com	廖锐
3.3	qianziqiang@sogou-inc.com	钱自强
3.3	lihuabei@sogou-inc.com	李华北
3.3	lijianhui@sogou-inc.com	李建辉
3.3	fubing@sogou-inc.com	伏冰
3.3	zhouhui@sogou-inc.com	周慧
3.3	huangtao220022@sogou-inc.com	黄涛
3.3	zhaoxiangfeng@sogou-inc.com	赵祥凤
3.3	houtingting@sogou-inc.com	侯婷婷
3.3	wangyi02@sogou-inc.com	王懿
3.3	wanglu214522@sogou-inc.com	王璐
3.3	mazhisi2054@sogou-inc.com	马枝
3.3	gengjingle@sogou-inc.com	耿京乐
3.3	limin208311@sogou-inc.com	李旻
3.3	gongqian@sogou-inc.com	巩茜
3.3	yangguang214314@sogou-inc.com	杨光
3.3	weneryan@sogou-inc.com	温二艳
3.4	liaorui205984@sogou-inc.com	廖锐
3.4	qianziqiang@sogou-inc.com	钱自强
3.4	lihuabei@sogou-inc.com	李华北
3.4	lijianhui@sogou-inc.com	李建辉
3.4	fubing@sogou-inc.com	伏冰
3.4	zhouhui@sogou-inc.com	周慧
3.4	huangtao220022@sogou-inc.com	黄涛
3.4	zhaoxiangfeng@sogou-inc.com	赵祥凤
3.4	houtingting@sogou-inc.com	侯婷婷
3.4	wangyi02@sogou-inc.com	王懿
3.4	wanglu214522@sogou-inc.com	王璐
3.4	mazhisi2054@sogou-inc.com	马枝
3.4	gengjingle@sogou-inc.com	耿京乐
3.4	limin208311@sogou-inc.com	李旻
3.4	gongqian@sogou-inc.com	巩茜
3.4	yangguang214314@sogou-inc.com	杨光
3.4	weneryan@sogou-inc.com	温二艳
3.5	liaorui205984@sogou-inc.com	廖锐
3.5	qianziqiang@sogou-inc.com	钱自强
3.5	lihuabei@sogou-inc.com	李华北
3.5	wuchuankai@sogou-inc.com	吴传凯
3.5	songyan@sogou-inc.com	宋妍
3.5	qibingyu@sogou-inc.com	祁秉钰
3.5	fubing@sogou-inc.com	伏冰
3.5	luyushan@sogou-inc.com	卢玉珊
3.5	chenshuchao@sogou-inc.com	陈树超
3.5	qixiaohe@sogou-inc.com	齐晓鹤
3.5	lianzhao@sogou-inc.com	连昭
3.5	liuyuansheng@sogou-inc.com	刘元生
3.5	zhouhui@sogou-inc.com	周慧
3.5	huangtao220022@sogou-inc.com	黄涛
3.5	yangdongxu207453@sogou-inc.com	yangdongxu
3.5	baijingjiao@sogou-inc.com	柏景娇
3.5	yangzhuoheng@sogou-inc.com	杨卓衡
3.5	gengjingle@sogou-inc.com	耿京乐
3.5	guowen@sogou-inc.com	郭雯
3.5	suyapeng@sogou-inc.com	苏亚鹏
3.5	liqiang214245@sogou-inc.com	李强
3.5	yangjie01@sogou-inc.com	杨婕
3.6	zhangguanglei@sogou-inc.com	张光磊
3.8	hemenglong212134@sogou-inc.com	何梦龙
3.8	zhaohongxue@sogou-inc.com	赵宏雪
3.8	chenfang@sogou-inc.com	陈芳
3.8	chenyaosi4587@sogou-inc.com	陈瑶
3.8	lihuabei@sogou-inc.com	李华北
3.8	qibingyu@sogou-inc.com	祁秉钰
3.8	chenshuchao@sogou-inc.com	陈树超
3.8	qixiaohe@sogou-inc.com	齐晓鹤
3.8	liuyuansheng@sogou-inc.com	刘元生
3.8	lianzhao@sogou-inc.com	连昭
3.8	luyushan@sogou-inc.com	卢玉珊
3.8	yiming@sogou-inc.com	易鸣
3.8	huangtao220022@sogou-inc.com	黄涛
3.8	liyangrui@sogou-inc.com	李杨瑞
3.8	xuxiaoting@sogou-inc.com	徐晓婷
3.8	sunxiaolei@sogou-inc.com	孙晓磊
3.8	liushukai212724@sogou-inc.com	刘树凯
3.8	wudongyang@sogou-inc.com	吴东洋
3.8	lianghao@sogou-inc.com	梁浩
3.8	luyuan212521@sogou-inc.com	陆远
3.8	zhaoxiaole@sogou-inc.com	zhaoxiaole
3.8	dushi@sogou-inc.com	dushi
3.8	huangboyang@sogou-inc.com	huangboyang
3.8	dudong@sogou-inc.com	杜东
3.8	heyukai@sogou-inc.com	贺宇凯
3.8	zhufangfang@sogou-inc.com	朱芳芳
3.8	wangya@sogou-inc.com	王丫
3.8	liuyutong211630@sogou-inc.com	刘宇同
3.8	shili@sogou-inc.com	史吏
3.8	dingxu@sogou-inc.com	丁旭
3.8	zhaomeng@sogou-inc.com	赵萌
3.8	yangjie01@sogou-inc.com	杨婕
4.1	niguochun@sogou-inc.com	倪国春
4.1	zhangyanou@sogou-inc.com	张艳欧
4.2	dongleijia@sogou-inc.com	贾东蕾
4.3	malin@sogou-inc.com	马蔺
4.3	zhanglijie@sogou-inc.com	张立婕
4.4	yuanyuandu@sogou-inc.com	杜媛媛
4.4	kangzhengguang@sogou-inc.com	亢正光
5	fuyanyan@sogou-inc.com	符艳燕
5	zhushunjie@sogou-inc.com	朱舜杰
5	zhangzhe214591@sogou-inc.com	张喆
5	maguoying@sogou-inc.com	马国营
5	liujiawen@sogou-inc.com	刘佳文
5	jiangzhen@sogou-inc.com	姜震
5	guyuxing@sogou-inc.com	顾宇星
5	sunfangyuan@sogou-inc.com	孙方园
5	menglingda@sogou-inc.com	孟令达
5	liuxianhe@sogou-inc.com	刘显赫
5	lipeng210100@sogou-inc.com	李鹏
5	lishujun@sogou-inc.com	李淑君
5	pandeng@sogou-inc.com	潘登
5	zhangchao216668@sogou-inc.com	张超216668
5	wangjie01@sogou-inc.com	王洁
5	liujiawei@sogou-inc.com	刘嘉炜
5	shiyan202878@sogou-inc.com	石岩
5	liqiaonan@sogou-inc.com	李乔楠
5	guojing@sogou-inc.com	郭静
5	zhanglinlin@sogou-inc.com	张琳琳
5	luoxiao@sogou-inc.com	罗霄
5	yangna@sogou-inc.com	杨娜
5	zhaodan207249@sogou-inc.com	赵丹
5	gaomin@sogou-inc.com	高敏
5	wangdan211213@sogou-inc.com	王丹
5	liweiwei@sogou-inc.com	黎维伟
5	yubinghe@sogou-inc.com	于秉禾
5	weishuqi@sogou-inc.com	魏淑琪
5	wuqi@sogou-inc.com	吴琪
5	lixiangnan@sogou-inc.com	李向南
5	zhaojingkai@sogou-inc.com	赵敬凯
5	liulin@sogou-inc.com	刘琳
5	wuxinyue@sogou-inc.com	吴欣悦
5	songsulin@sogou-inc.com	宋苏霖
5	zhangrui@sogou-inc.com	张锐
5	wangxiaoning@sogou-inc.com	王晓宁
5	huzhao@sogou-inc.com	胡昭
5	hening@sogou-inc.com	贺宁
5	zhina@sogou-inc.com	支娜
5	hansiwei@sogou-inc.com	韩思维
5	qianhai@sogou-inc.com	钱海
5	yaokai@sogou-inc.com	姚凯
5	zhoujielian@sogou-inc.com	周洁莲
5	hesonghai@sogou-inc.com	何宋海
5	xufeng213905@sogou-inc.com	徐峰
5	shenai@sogou-inc.com	沈爱
5	shengyu@sogou-inc.com	盛宇
5	qixiaoduo@sogou-inc.com	齐小朵
5	congchao@sogou-inc.com	丛超
5	liliang210741@sogou-inc.com	李亮
5	renna@sogou-inc.com	任娜
5	yangzheng@sogou-inc.com	杨征
5	liushuo@sogou-inc.com	刘硕
5	zhangyanou@sogou-inc.com	张艳欧
5	sunqi@sogou-inc.com	孙琪
5	sakurazhu@sogou-inc.com	朱佳
5	yanxiaokun@sogou-inc.com	闫小坤
5	yuanyuandu@sogou-inc.com	杜媛媛
5	xuhao212985@sogou-inc.com	徐浩
5	gaoxu@sogou-inc.com	高旭
5.1	kaqudao01@sogou-inc.com	KA渠道1
5.1	tinawang@sogou-inc.com	王琳
5.1	chenxichao@sogou-inc.com	陈西超
5.1	yutaoli@sogou-inc.com	李玉涛
5.1	wenwenyu@sogou-inc.com	于雯雯
5.1	zhuhong@sogou-inc.com	朱洪
5.1	yuntaocong@sogou-inc.com	丛云涛SH
5.1	chenruikai@sogou-inc.com	陈睿恺SH
5.1	dixuemei@sogou-inc.com	邸雪梅
5.1	shellxiao@sogou-inc.com	肖秉峰GZ
5.1	zhoudajiang@sogou-inc.com	周大江GZ
5.1	liumingyi@sogou-inc.com	刘明翼
5.1	qianzhen@sogou-inc.com	钱震
5.1	wangfuying@sogou-inc.com	王馥颖
5.1	zhoutong209047@sogou-inc.com	周桐
5.1	maxiaoyan@sogou-inc.com	马妍
5.1	zhusheng@sogou-inc.com	朱晟
5.1	xielujing@sogou-inc.com	谢鲁京
5.1	hewen@sogou-inc.com	何雯
5.1	surui@sogou-inc.com	苏瑞
5.1	liuzihui@sogou-inc.com	刘子辉
5.1	wanghongxu@sogou-inc.com	王宏旭
`
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		splits := strings.Split(line, "\t")
		role := ""
		if len(splits) > 2 {
			switch roleIdx := splits[0]; roleIdx {
			case "1": role = "波动总权限"
			case "2":role = "波动客户and行业"
			case "2.1":role = "波动KA客户"
			case "2.2":role = "波动中小所有客户"
			case "3.1":role = "波动搜索渠道"
			case "3.2":role = "波动搜索渠道"
			case "3.3":role = "波动网盟渠道"
			case "3.4":role = "波动网盟渠道"
			case "3.5":role = "波动搜索渠道"
			case "3.6":role = "波动搜索渠道"
			case "3.8":role = ""
			case "4.1":role = "波动中小-华北客户"
			case "4.2":role = "波动中小-华东客户"
			case "4.3":role = "波动中小-华南客户"
			case "4.4":role = "波动中小-华中客户"
			case "5":role = "波动单个客户"
			case "5.1":role = "波动单个客户"
			}
			if role != "" {
				userName := strings.Split(splits[1],"@")[0]
				//fmt.Print(role + ",")
				//fmt.Println(userName)
				fmt.Printf("INSERT INTO casbin_rule(p_type,v0,v1) VALUE('g','%s','%s'); \n", userName, role)
			}
		}

	}
}
