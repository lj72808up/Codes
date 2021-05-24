package controllers

import (
	"AStar/models"
	"AStar/utils"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type XiaopController struct {
	beego.Controller
}

// @Title 登录时添加用户
// @Summary  登录时添加用户
// @Description  登录时添加用户
// @Accept json
// @router /auth/user/postUser [post]
func (this *XiaopController) PostOrUpdateUser() {
	var userNameMapping models.UserNameMapping
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &userNameMapping)
	userNameMapping.AddOrUpdate()
}

// @Title Robot
// @Summary 小P权限申请机器人接口
// @Description 小P权限申请机器人接口
// @Success 200 {array[string]}
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /robot/ [post]
func (this *XiaopController) Robot() {
	uid := this.Ctx.Input.Query("uid")
	q := this.Ctx.Input.Query("q")
	fmt.Println(q)
	if q == "申请权限" {
		ts := strconv.FormatInt(time.Now().Unix()*1000, 10)
		token := utils.GetMd5("adtechdatacenter:b8851867c9f7f0e3d719d634957d19f9:" + ts)
		urlStr := "https://puboa.sogou-inc.com/moa/sylla/open/v1/pns/send"
		title := "[待办] 申请单填写"
		content := "请填写申请单"
		to := uid
		println(uid)
		http.PostForm(urlStr,
			url.Values{"content": {content}, "to": {to}, "token": {token}, "title": {title},
				"pubid": {"adtechdatacenter"}, "ts": {ts}, "tp": {"3"}, "url": {"http://datasearch.sogou:8081/#/apply?uid=" + uid}})
		ans := models.RobotAns{2, "2", "OK"}
		res, _ := json.Marshal(ans)
		this.Ctx.WriteString(string(res))
	} else {
		ans := models.RobotAns{2, "2", "OK"}
		res, _ := json.Marshal(ans)
		this.Ctx.WriteString(string(res))
	}
}

// @Title Apply
// @Summary 申请单填写
// @Description 申请单填写
// @Param para body models.XiaopApply true "提交信息"
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /apply/ [post]
func (this *XiaopController) Apply() {
	var applyJSON = models.XiaopApply{}
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &applyJSON)
	if err != nil {
		models.Logger.Error(err.Error())
		this.Ctx.ResponseWriter.WriteHeader(400)
		return
	}
	o := orm.NewOrm()
	_, err = o.Insert(&applyJSON)
	fmt.Println(err)

	ts := strconv.FormatInt(time.Now().Unix()*1000, 10)
	fmt.Println(applyJSON)
	token := utils.GetMd5("adtechdatacenter:b8851867c9f7f0e3d719d634957d19f9:" + ts)
	urlStr := "https://puboa.sogou-inc.com/moa/sylla/open/v1/pns/send"
	title := "[待办] 权限审批 " + applyJSON.Name
	content := "请完成" + applyJSON.Name + "的权限审批单"
	//to := "kongyibo,liujie02"
	to := beego.AppConfig.String("xiaopApproval")
	println(token)
	urlx := fmt.Sprintf("http://datasearch.sogou:8081/#/approval?uid=%s&name=%s&boss=%s&apartment=%s&position=%s&reason=%s", applyJSON.Uid, applyJSON.Name, applyJSON.Boss, applyJSON.Apartment, applyJSON.Position, applyJSON.Reason)
	http.PostForm(urlStr,
		url.Values{"content": {content}, "to": {to}, "token": {token}, "title": {title},
			"pubid": {"adtechdatacenter"}, "ts": {ts}, "tp": {"3"}, "url": {urlx}})

	this.Ctx.WriteString("")
}

// @Title 查看该用户的申请状态  (0:申请中,1:已审核)
// @Summary 查看该用户的申请状态
// @Description 查看该用户的申请状态
// @Accept json
// @router /applyStatus/:uid [get]
func (this *XiaopController) GetApplyStatus() {
	user := this.GetString(":uid")
	o := orm.NewOrm()
	var xiaopApply models.XiaopApply
	_ = o.Raw("SELECT * FROM xiaop_apply WHERE uid=?", user).QueryRow(&xiaopApply)
	res, _ := json.Marshal(xiaopApply)
	this.Ctx.WriteString(string(res))
}

// @Title Auth
// @Summary 获取所有角色
// @Description 获取所有角色
// @Param para body models.XiaopApply true "提交信息"
// @Success 200 {array[string]}
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /roles/ [get]
func (this *XiaopController) Roles() {
	jsonStr, _ := json.Marshal(models.Enforcer.GetAllRoles())
	this.Ctx.WriteString(string(jsonStr))
}

// @Title Auth
// @Summary 设置用户角色
// @Description 设置用户角色
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /role/add [post]
func (this *XiaopController) AddRole() {
	var applyJSON = models.XiaopApply{}
	o := orm.NewOrm()
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &applyJSON)
	if err != nil {
		models.Logger.Error(err.Error())
		this.Ctx.ResponseWriter.WriteHeader(400)
		return
	}
	uid := applyJSON.Uid
	role := applyJSON.Role
	roles := models.Enforcer.GetRolesForUser(uid) // 该用户已存在的角色
	sql := "UPDATE xiaop_apply SET status = ?, detail = ? WHERE uid = ?"
	if len(roles) == 0 {
		models.Enforcer.AddRoleForUser(uid, role)
		_, _ = o.Raw(sql, 1, "审核通过", uid).Exec()
		successMsg := "审核通过"
		pushBack(successMsg, uid)
		this.Ctx.WriteString(utils.VueResponse{
			Res:  true,
			Info: successMsg,
		}.String())
	} else {
		failMsg := "该用户已存在角色: " + roles[0]
		_, _ = o.Raw(sql, -1, failMsg, uid).Exec()
		pushBack(failMsg, uid)
		this.Ctx.WriteString(utils.VueResponse{
			Res:  false,
			Info: failMsg,
		}.String())
	}
}

// 审核后通知该用户
func pushBack(msg string, uid string) {
	ts := strconv.FormatInt(time.Now().Unix()*1000, 10)
	token := utils.GetMd5("adtechdatacenter:b8851867c9f7f0e3d719d634957d19f9:" + ts)
	urlStr := "https://puboa.sogou-inc.com/moa/sylla/open/v1/pns/send"
	title := "请查看审核结果"
	content := msg
	to := uid
	http.PostForm(urlStr,
		url.Values{"pubid": {"adtechdatacenter"}, "token": {token},
			"ts": {ts}, "to": {to}, "content": {content}, "title": {title}, "tp": {"3"},})

}

// @Title Auth
// @Summary 查询用户角色
// @Description 查询用户角色
// @Success 200 string
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /role/:uid [get]
func (this *XiaopController) GetRole() {
	uid := this.Ctx.Input.Params()[":uid"]
	roles := models.Enforcer.GetRolesForUser(uid)
	if len(roles) == 0 {
		this.Ctx.WriteString("")
	} else {
		this.Ctx.WriteString(roles[0])
	}

}

// @Title 解析pToken
// @Summary  解析pToken
// @Description
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /lvpi/parseToken [post]
func (this *XiaopController) ParseLvpiToken() {
	var params map[string]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &params)
	pToken := params["pToken"]

	verifyURL := "https://login.sogou-inc.com/backend/verify/"
	defer func() {
		if r := recover(); r != nil {
			models.Logger.Error("%v", r)
			bytes, _ := json.Marshal(utils.VueResponse{
				Res:  false,
				Info: fmt.Sprintf("%v", r),
			})
			this.Ctx.WriteString(string(bytes))
		}
	}()

	//client := http.Client{Timeout: 30 * time.Second}
	formData := url.Values{"appid": {"1899"}, "ptoken": {pToken}}
	resp, err := http.PostForm(verifyURL, formData)

	if err != nil {
		panic(err.Error())
	} else {
		entity := LogInEntity{}
		body, _ := ioutil.ReadAll(resp.Body)
		_ = json.Unmarshal(body, &entity)

		models.Logger.Info("SSO登录验证结果:%v", entity)

		logInData := "接口token验证失败"
		if entity.Status == 1 {
			bytes,_ := json.Marshal(entity.Data)
			logInData = string(bytes)
		} else {
			logInData = logInData + ":" + entity.Message
		}

		bytes, _ := json.Marshal(utils.VueResponse{
			Res:  entity.Status == 1,
			Info: logInData,
		})
		this.Ctx.WriteString(string(bytes))
	}
}

// @Title 解析pToken
// @Summary  解析pToken
// @Description
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /jone/parseToken [post]
func (this *XiaopController) ParseToken() {
	var params map[string]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &params)
	pToken := params["pToken"]

	verifyURL := "https://login.sogou-inc.com/backend/verify/"
	defer func() {
		if r := recover(); r != nil {
			models.Logger.Error("%v", r)
			bytes, _ := json.Marshal(utils.VueResponse{
				Res:  false,
				Info: fmt.Sprintf("%v", r),
			})
			this.Ctx.WriteString(string(bytes))
		}
	}()

	//client := http.Client{Timeout: 30 * time.Second}
	formData := url.Values{"appid": {"1812"}, "ptoken": {pToken}}
	resp, err := http.PostForm(verifyURL, formData)

	if err != nil {
		panic(err.Error())
	} else {
		entity := LogInEntity{}
		body, _ := ioutil.ReadAll(resp.Body)
		_ = json.Unmarshal(body, &entity)

		models.Logger.Info("SSO登录验证结果:%v", entity)

		logInData := "接口token验证失败"
		if entity.Status == 1 {
			bytes,_ := json.Marshal(entity.Data)
			logInData = string(bytes)
		} else {
			logInData = logInData + ":" + entity.Message
		}

		bytes, _ := json.Marshal(utils.VueResponse{
			Res:  entity.Status == 1,
			Info: logInData,
		})
		this.Ctx.WriteString(string(bytes))
	}
}

type LogInEntity struct {
	Status  int       `json:"status"`
	Message string    `json:"message"`
	Data    LogInData `json:"data"`
}

type LogInData struct {
	Uid     string `json:"uid"`
	Name    string `json:"name"`
	Unumber string `json:"uno"`
}
