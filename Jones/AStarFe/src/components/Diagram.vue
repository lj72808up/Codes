<template>
  <div
    class="page-content"
    @mousedown="startNodesBus($event)"
    @mousemove="moveNodesBus($event)"
    @mouseup="endNodesBus($event)"
  >
    <el-button type="primary" size="small" class="elbutton" @mousedown.native="dragIt('new-node')">拖拽新建Node</el-button>
    <el-button type="primary" size="small" style="position: absolute;top:10px;right: 10px"
               v-if="isCheck"
               @click="onUpdateDAG">提交修改
    </el-button>
    <el-button type="primary" size="small" style="position: absolute;top:10px;right: 10px"
               @click="onSubmitDAG"
               v-if="!isCheck">提交创建
    </el-button>
    <DAGBoard
      :DataAll="DataAll"
      @updateDAG="updateDAG"
      @editNodeDetails="editNodeDetails"
      @doSthPersonal="doSthPersonal"
    ></DAGBoard>
    <!-- DAG-Diagram主体 -->

    <!-- 用来模拟拖拽添加的元素 -->
    <node-bus
      v-if="dragBus"
      :value="busValue.value"
      :pos_x="busValue.pos_x"
      :pos_y="busValue.pos_y"
    />
    <!-- 右侧JSOn展示,你忽略就行了 -->
    <editor
      ref="myEditor"
      class="json-editor"
      v-model="jsonEditor"
      :options="options"
      @init="editorInit"
      lang="json"
      theme="chrome"
      width="400"
      v-show="false"
    ></editor>

    <!--右侧弹框-->
    <!--:before-close="handleClose"-->
    <el-drawer
      title="详情"
      :with-header="false"
      :visible.sync="drawerVisible"
      direction="rtl"
      custom-class="demo-drawer"
      ref="drawer"
      size="50%"
      :before-close="handleClose"
    >
      <div class="allDiv">
        <div class="content">
          <div class="content-inside" style="max-height: 720px ; overflow:auto">
            <el-form :model="nodeForm" ref="nodeForm" label-width="130px">
              <el-form-item label="任务名称">
                <el-col :span="21">
                  <el-input size="small" v-model="nodeForm.tmpName"/>
                </el-col>
                <el-col :span="2" :offset="1">
                  <el-tooltip content="任务名称只能包含a-z,0-9和中划线" placement="bottom" effect="light" v-if="!checkJobName">
                    <el-button type="danger" icon="el-icon-close" v-if="!checkJobName" circle size="mini"/>
                  </el-tooltip>
                </el-col>
              </el-form-item>
              <el-form-item label="任务描述">
                <el-input size="small" v-model="nodeForm.taskDesc"/>
              </el-form-item>
              <el-form-item label="镜像">
                <el-select size="small" v-model="nodeForm.dockerImg" filterable clearable @change="onChangeImg"
                           style="width: 100%">
                  <el-option
                    v-for="item in imgOptions"
                    :key="item.Img"
                    :label="item.Img"
                    :value="item.Id">
                  </el-option>
                </el-select>
              </el-form-item>
              <el-form-item label="命令">
                <el-input size="small" v-model="nodeForm.dockerCmd" :disabled="true"/>
              </el-form-item>
              <el-form-item label="参数" v-show="false">
                <el-input size="small" v-model="nodeForm.param"/>
              </el-form-item>
              <el-divider content-position="center">输入</el-divider>
              <el-form-item label="任务类型">
                <el-select size="small" v-model="nodeForm.taskType" filterable @change="taskTypeChange">
                  <el-option
                    v-for="t in getInputTypeOptionByImg"
                    :key="t.id"
                    :label="t.label"
                    :value="t.id">
                  </el-option>
                </el-select>
              </el-form-item>

              <el-form-item label="任务数据源" v-show="nodeForm.taskType==='mysql' || nodeForm.taskType==='lvpi'">
                <el-select size="small" v-model="nodeForm.taskTypeInfo.dsid" filterable clearable>
                  <el-option
                    v-for="opt in taskTypeDsidOptions"
                    :key="opt.DatasourceName"
                    :label="opt.DatasourceName"
                    :value="opt.Dsid">
                  </el-option>
                </el-select>
              </el-form-item>

              <el-form-item label="文件路径" v-show="nodeForm.taskType==='hdfs'">
                <el-input size="small" v-model="nodeForm.taskTypeInfo.hdfsLocation"/>
              </el-form-item>

              <el-form-item label="文件编码" v-show="nodeForm.taskType==='hdfs'">
                <el-select size="small" v-model="nodeForm.taskTypeInfo.fileEncode">
                  <el-option label="GBK" value="GBK"></el-option>
                  <el-option label="UTF-8" value="UTF-8"></el-option>
                </el-select>
              </el-form-item>

              <el-row>
                <el-col :span="8">
                  <el-form-item label="文件分隔符" v-show="nodeForm.taskType==='hdfs'">
                    <el-input size="small" v-model="nodeForm.taskTypeInfo.splitter"/>
                  </el-form-item>
                </el-col>

                <el-col :span="8">
                  <el-form-item label="字段个数" v-show="nodeForm.taskType==='hdfs'">
                    <el-input size="small" v-model="nodeForm.taskTypeInfo.validFieldCount"/>
                  </el-form-item>
                </el-col>
              </el-row>

              <el-form-item label="字段约束" v-show="nodeForm.taskType==='hdfs'">
                <el-row v-for="(extendItem,idx) in nodeForm.taskTypeInfo.schema">
                  <el-col :span="9">
                    <el-input size="small" v-model="extendItem.id">
                      <template slot="prepend">文件中的第几列</template>
                    </el-input>
                  </el-col>
                  <el-col :span="1" style="text-align: center;">-</el-col>
                  <el-col :span="9">
                    <el-input size="small" v-model="extendItem.name">
                      <template slot="prepend">字段名</template>
                    </el-input>
                  </el-col>
                  <el-col :span="1" :offset="1">
                    <el-button type="text" icon="el-icon-circle-plus-outline" @click="addExtendSchema(idx,'schema')"
                               style="font-size: 18px"/>
                  </el-col>
                  <el-col :span="1">
                    <el-button type="text" icon="el-icon-remove-outline" @click="delExtendSchema(idx,'schema')"
                               style="font-size: 18px; color: red; "/>
                  </el-col>
                </el-row>
              </el-form-item>

              <el-row>
                <el-col :span="8">
                  <el-form-item label="检测done" v-show="nodeForm.taskType!=='custom'">
                    <el-radio-group v-model="nodeForm.taskTypeInfo.isCheckDone && nodeForm.taskType!=='custom'">
                      <el-radio :label="true">是</el-radio>
                      <el-radio :label="false">否</el-radio>
                    </el-radio-group>
                  </el-form-item>
                </el-col>

                <div v-show="nodeForm.taskType==='custom'">
                  <el-form-item label="程序包HDFS地址:">
                    <el-input size="small" v-model="nodeForm.taskTypeInfo.downloadHdfsPath"/>
                  </el-form-item>
                  <el-form-item label="执行命令:">
                    <el-input size="small" v-model="nodeForm.taskTypeInfo.cmd"/>
                  </el-form-item>
                  <el-form-item label="HDFS授权账号:">
                    <el-input size="small" v-model="nodeForm.taskTypeInfo.hdfsAuthUser"/>
                  </el-form-item>
                  <el-form-item label="HDFS授权密码:">
                    <el-input size="small" v-model="nodeForm.taskTypeInfo.hdfsAuthPassword"/>
                  </el-form-item>
                </div>

                <el-col :span="16">
                  <el-form-item label="done位置:"
                                v-show="nodeForm.taskType==='hdfs' && nodeForm.taskTypeInfo.isCheckDone">
                    <el-input size="small" v-model="nodeForm.taskTypeInfo.hdfsDoneLocation"/>
                  </el-form-item>
                </el-col>
              </el-row>

              <el-form-item label="扩展字段" v-show="nodeForm.taskType==='hdfs'">
                <el-radio-group v-model="nodeForm.taskTypeInfo.hasExtend">
                  <el-radio :label="true">包含</el-radio>
                  <el-radio :label="false">不包含</el-radio>
                </el-radio-group>


                <el-row v-for="(extendItem,idx) in nodeForm.taskTypeInfo.extendSchema"
                        v-show="nodeForm.taskType==='hdfs' && nodeForm.taskTypeInfo.hasExtend">
                  <el-col :span="9">
                    <el-input size="small" v-model="extendItem.id">
                      <template slot="prepend">文件中的第几列</template>
                    </el-input>
                  </el-col>
                  <el-col :span="1" style="text-align: center;">-</el-col>
                  <el-col :span="9">
                    <el-input size="small" v-model="extendItem.name">
                      <template slot="prepend">字段名</template>
                    </el-input>
                  </el-col>
                  <el-col :span="1" :offset="1">
                    <el-button type="text" icon="el-icon-circle-plus-outline"
                               @click="addExtendSchema(idx,'extendSchema')"
                               style="font-size: 18px"/>
                  </el-col>
                  <el-col :span="1">
                    <el-button type="text" icon="el-icon-remove-outline" @click="delExtendSchema(idx,'extendSchema')"
                               style="font-size: 18px; color: red; "/>
                  </el-col>
                </el-row>
              </el-form-item>


              <el-form-item label="sql语句" v-show="nodeForm.taskType!=='hdfs'&&nodeForm.taskType!=='custom'">
                <el-tooltip :content="helpContent" placement="top">
                  <el-link type="warning" underline="false" icon="el-icon-info">DAG时间变量请使用内置变量</el-link>
                </el-tooltip>
                <el-input size="small" type="textarea" rows=5 v-model="nodeForm.originalSql" :readonly="false"/>
                <el-button type="primary" size="small" @click="CreateSQLByTemplate()">使用SQL模板</el-button>
              </el-form-item>
              <el-form-item label="sql参数" v-show="false">
                <el-input size="small" type="textarea" rows=4 v-model="nodeForm.sqlParam" :readonly="false"/>
              </el-form-item>

              <div v-show="nodeForm.taskType!=='custom'">
                <el-divider content-position="center">输出</el-divider>
                <el-form-item label="输出类型">
                  <el-select size="small" v-model="nodeForm.outTableType" filterable clearable
                             @change="getConnection">
                    <el-option
                      v-for="t in outType"
                      :key="t.id"
                      :label="t.label"
                      :value="t.id">
                    </el-option>
                  </el-select>
                </el-form-item>
                <el-form-item label="" v-show="getOutputType!=='mail'">
                  <el-radio-group v-model="nodeForm.newOrOldOutTable">
                    <el-radio label="newTable">新创建</el-radio>
                    <el-radio label="oldTable">选择</el-radio>
                  </el-radio-group>
                </el-form-item>
                <!--配置新的输出表 start-->
                <el-form-item label="数据源"
                              v-show="getOutputType==='mysql' || getOutputType==='clickhouse' || getOutputType==='lvpi' ||  getOutputType==='hive'">
                  <!--                <el-input size="small" v-model="nodeForm.lvpi.dsid"/>-->
                  <el-select v-model="nodeForm.lvpi.dsid" clearable filterable
                             placeholder="请选择" size="small" @change="getOutTableNameOptions">
                    <el-option
                      v-for="conn in dsOptions"
                      :key="conn.Dsid"
                      :label="conn.DatasourceName"
                      :value="conn.Dsid">
                    </el-option>
                  </el-select>
                </el-form-item>
                <el-form-item label="输出表名" v-show="getNewOrOldOutTable==='newTable'">
                  <el-col :span="10">
                    <el-input size="small" v-model="nodeForm.outTableName"/>
                  </el-col>
                  <el-col :span="10" :offset="1">
                    <!--                  <el-button @click="parseSql" size="small" type="primary">配置表</el-button>-->
                    <el-button type="primary" @click="parseSql" v-show="nodeForm.taskType!=='hdfs'"
                               :icon="parseButtonIcon" :disabled="onParsing">{{ parseButtonText }}
                    </el-button>
                  </el-col>
                </el-form-item>

                <!--hive-->
                <el-form-item label="存储格式" v-show="getOutputType==='hive' && getNewOrOldOutTable==='newTable'">
                  <el-select v-model="nodeForm.storageEngine" placeholder="请选择" size="small">
                    <el-option
                      v-for="engine in nodeForm.storageEngineOptions"
                      :key="engine"
                      :label="engine"
                      :value="engine">
                    </el-option>
                  </el-select>
                </el-form-item>
                <el-form-item label="输出表描述" v-show="getOutputType!=='mail' && getNewOrOldOutTable==='newTable'">
                  <el-input size="small" v-model="nodeForm.outTableDesc"/>
                </el-form-item>
                <el-form-item label="输出字段"
                              v-show="getOutputType!=='mail' && getNewOrOldOutTable==='newTable' &&
                            this.nodeForm.taskType !== 'hdfs'">
                  <vxe-table
                    border="inner"
                    show-overflow
                    keep-source
                    ref="fieldTable"
                    stripe="true"
                    :data="nodeForm.sqlFields"
                    :edit-config="{trigger: 'dblclick', mode: 'cell', showStatus: true, icon: 'fa fa-pencil-square-o'}">
                    <!--显示字段-->
                    <vxe-table-column field="fieldName" title="字段名"/>
                    <vxe-table-column field="description" title="描述"
                                      :edit-render="{name: 'input', immediate: true, attrs: {type: 'text'}}"/>
                  </vxe-table>
                </el-form-item>

                <!--clickhouse-->
                <el-form-item label="partitions"
                              v-show="getOutputType==='clickhouse' &&  getNewOrOldOutTable==='newTable'">
                  <el-input size="small" v-model="nodeForm.clickhouse.partitionStr"/>
                </el-form-item>
                <el-form-item label="indices"
                              v-show="getOutputType==='clickhouse' &&  getNewOrOldOutTable==='newTable'">
                  <el-input size="small" v-model="nodeForm.clickhouse.indexStr"/>
                </el-form-item>

                <!--mysql-->
                <!--<el-form-item label="jdbcString" v-show="false && getOutputType==='mysql' &&  getNewOrOldOutTable==='newTable'">
                  <el-input size="small" v-model="nodeForm.lvpi.jdbcStr"/>
                </el-form-item>-->


                <!--配置新的输出表 end-->
                <!--选择旧输出表 start-->
                <el-form-item label="输出表名" v-show="getNewOrOldOutTable==='oldTable'">
                  <el-col :span="12">
                    <el-select v-model="nodeForm.outTableName" clearable filterable
                               placeholder="请选择" size="small">
                      <el-option
                        v-for="tb in nodeForm.outTableNameOptions"
                        :key="tb"
                        :label="tb"
                        :value="tb">
                      </el-option>
                    </el-select>
                  </el-col>
                  <el-col :span="12">
                    <el-button @click="checkMetaData" type="primary" size="small">校验输出表</el-button>
                    <!--                  <span>aaa</span>-->
                    <el-tooltip content="输出表schema和sql匹配" placement="bottom" effect="light"
                                v-if="checkOutputMeta.success">
                      <el-button type="success" icon="el-icon-check" v-if="checkOutputMeta.success" circle size="mini"
                                 :readonly="isCheck"/>
                    </el-tooltip>
                    <el-tooltip content="输出表schema和sql不匹配,或未配置sql参数 请检查" placement="bottom" effect="light"
                                v-if="checkOutputMeta.fail">
                      <el-button type="danger" icon="el-icon-close" v-if="checkOutputMeta.fail" circle size="mini"
                                 :readonly="isCheck"/>
                    </el-tooltip>
                    <i class="el-icon-loading" v-if="checkOutputMeta.start"/>
                  </el-col>
                </el-form-item>
                <!--选择旧输出表 end-->

                <!--输出类型是email-->
                <div v-show="getOutputType==='mail'">
                  <el-form-item label="邮件标题">
                    <el-col :span="18">
                      <el-input placeholder="" size="small" v-model="nodeForm.email.title"/>
                    </el-col>
                    <el-col :span="5" :offset="1">
                      <el-button size="small" type="primary" @click="exploreEmail">预览表格</el-button>
                    </el-col>
                  </el-form-item>
                  <el-form-item label="收件人">
                    <el-input placeholder="多个用户用逗号隔开(不用配置邮箱后缀)" size="small" v-model="nodeForm.email.receiver"/>
                  </el-form-item>
                  <el-form-item label="附件发送">
                    <el-switch v-model="nodeForm.email.withAttachment" active-color="#13ce66" inactive-color="#ff4949">
                    </el-switch>
                  </el-form-item>
                  <div id="emailExplore"></div>
                </div>
                <el-divider content-position="center"/>
              </div>
            </el-form>
          </div>
        </div>
        <!--        <div class="footer">-->
        <!--          <el-button type="primary" @click="handleClose">确定</el-button>-->
        <!--        </div>-->
      </div>
    </el-drawer>

    <!--右侧弹框-->
    <!--:before-close="handleClose"-->
    <el-drawer
      title="详情"
      :with-header="false"
      :visible.sync="drawerDAGVisible"
      direction="rtl"
      custom-class="demo-drawer"
      ref="drawer"
      :before-close="handleDAGClose"
    >
      <div class="allDiv">
        <div class="content">
          <div class="content-inside" style="max-height: 800px ; overflow:auto">
            <el-form :model="dagForm" label-width="120px">

              <el-form-item label="DAG ID">
                <el-col :span="22">

                  <el-input v-model="dagForm.name" size="mini" v-if="!isCheck" @blur="checkDuplicateName"
                            :readonly="isCheck"/>
                  <el-input v-model="dagForm.name" size="mini" v-if="isCheck" :readonly="isCheck"/>

                </el-col>
                <el-col :span="1" :offset="1">

                  <el-tooltip content="ID可用" placement="bottom" effect="light" v-if="checkTaskName.success">
                    <el-button type="success" icon="el-icon-check" v-if="checkTaskName.success" circle size="mini"
                               :readonly="isCheck"/>
                  </el-tooltip>
                  <el-tooltip content="ID重复或包含非法字符(只能是a-z, 0-9和中划线)" placement="bottom" effect="light"
                              v-if="checkTaskName.fail">
                    <el-button type="danger" icon="el-icon-close" v-if="checkTaskName.fail" circle size="mini"
                               :readonly="isCheck"/>
                  </el-tooltip>
                  <i class="el-icon-loading" v-if="checkTaskName.start"/>
                </el-col>
              </el-form-item>

              <el-form-item label="DAG名称">
                <el-input v-model="dagForm.exhibitName" size="mini" :readonly="isCheck"/>
              </el-form-item>

              <el-form-item label="开始时间">
                <el-date-picker
                  size="mini"
                  v-model="dagForm.startDate"
                  type="date"
                  placeholder="选择任务日期"
                  format="yyyy-MM-dd"
                  value-format="yyyy-MM-dd"
                  style="width: 80%;"/>
              </el-form-item>
              <el-form-item label="邮件列表">
                <el-input v-model="dagForm.mailList" size="mini"/>
              </el-form-item>
              <el-form-item label="">
                <el-radio-group v-model="dagForm.dagType" @change="onChangeDAGType">
                  <el-radio label="day">天级任务</el-radio>
                  <el-radio label="hour">小时级任务</el-radio>
                </el-radio-group>
              </el-form-item>
              <el-form-item label="延时">
                <el-input-number v-model="dagForm.delay" :min="0" size="mini"/>
                <span style="color: rgb(189, 189, 189); font-size: 12px;">(填1, 表示任务执行前1小时或前1天的数据)</span>
              </el-form-item>
              <el-form-item label="cron表达式">
                <el-input :placeholder="cronPlaceHolder" v-model="dagForm.cronExpression" size="mini"/>
              </el-form-item>
              <el-form-item label="重试次数">
                <el-input-number v-model="dagForm.retries" :min="0" size="mini"/>
              </el-form-item>
              <el-form-item label="重试延时(分钟)">
                <el-input-number v-model="dagForm.retryDelay" :min="0" size="mini"/>
              </el-form-item>
            </el-form>
          </div>
        </div>
      </div>
    </el-drawer>

    <!--提交的dialog-->
    <el-dialog title="请等待DAG生成完毕" :visible.sync="submitDialogVisible" width="40%" :before-close="handleCloseSubmit">
      <div style="margin-top: 25px">
        <el-timeline>
          <el-timeline-item
            v-for="(activity, index) in submitActivities"
            :key="index"
            :icon="activity.icon"
            :type="activity.type"
            :color="activity.color"
            size="large"
            :timestamp="activity.timestamp">
            {{ activity.content }}
          </el-timeline-item>
        </el-timeline>
      </div>
      <span slot="footer" class="dialog-footer">
        <el-button type="primary" size="small" @click="afterSubmitDag()"
                   :icon="buttonIcon" :disabled="onProcessing">{{ buttonText }}</el-button>
      </span>
    </el-dialog>


    <!--生成SQL的 dialog-->
    <el-dialog title="SQL生成" :visible.sync="sqlDialogVisible" width="80%">
      <div style="height: 460px;width: 100%;">
        <el-col :span=5>
          <div id="treeDiv" style="height: 460px;overflow:auto">
            <span>库中的模板</span>
            <el-input
              placeholder="输入关键字过滤"
              v-model="filterText">
            </el-input>
            <!--同步加载树-->
            <el-tree :data="treeData" :props="defaultSqlProps" :filter-node-method="filterNode" ref="tree" accordion
                     @node-click="qyeryInfo"/>
          </div>
        </el-col>

        <el-col :span=13 :offset=1>
          <el-form ref="form" :model="form">
            <el-row>
              <el-input v-model="description"
                        v-show="true"
                        placeholder="工单描述"
                        :rows="1"
                        :max-rows="3"
                        :readonly="true">
                <template slot="prepend">工单备注:</template>
              </el-input>
            </el-row>

            <el-form-item prop="sql">
              <codemirror v-model="form.sql" :options="cmOptions" ref="cmEditor"/>
            </el-form-item>

            <el-form-item>
              <el-col :span=16>
                <el-row type="flex" justify="start">
                  <el-button @click="onParam('form')">生成参数列表</el-button>
                </el-row>
              </el-col>
              <el-col :span=8>
                <el-row type="flex" justify="end">
                  <el-button @click="CommitSQL()" type="primary">确定
                  </el-button>
                </el-row>
              </el-col>

            </el-form-item>
          </el-form>
        </el-col>

        <el-col :span=4 :offset=1>
          <div>
            <el-row>
              <label>DAG内置变量:</label>
              <el-tooltip content="表示时间间隔的变量, 后缀可以加上 .start 和 .end" placement="top">
                <i class="el-icon-question"/>
              </el-tooltip>
              <div style="overflow:auto; max-height: 150px; margin-left: 6px; line-height: 25px">
                <li>时间间隔:</li>
                <span style="margin-left: 18px; color: #c0ccda">{{ intervalVars }}</span>
                <li>时间点:</li>
                <span style="margin-left: 18px; color: #c0ccda">{{ momentVars }}</span>
              </div>
            </el-row>
            <el-row style="margin-top: 15px">
              <label>参数配置</label>
              <el-tooltip class="item" effect="dark" content="参数用于提交的sql中" placement="top-start">
                <i class="el-icon-question" @click="makeParam()"/>
              </el-tooltip>
              <el-button @click="makeParam()" v-show="false">替换参数模板</el-button>
            </el-row>
            <div style="overflow:auto; max-height: 230px">
              <template v-for="(item,index) in paramNames">
                <el-row>
                  <span>参数名: {{ item }}</span>
                </el-row>
                <el-row>
                  <el-select v-model="paramTypes[index]" placeholder="请选择参数类型" @change="onChangeParamType(index)"
                             size="mini">
                    <el-option
                      v-for="item in paramOptions"
                      :key="item.value"
                      :label="item.label"
                      :value="item.value">
                    </el-option>
                  </el-select>
                </el-row>
                <el-row>
                  <el-input v-model="paramValues[index]" placeholder="请输入"
                            v-show="paramTypes[index]==='string'" size="mini"/>
                  <el-date-picker
                    size="mini"
                    v-model="paramValues[index]"
                    v-show="paramTypes[index]==='date'"
                    type="date"
                    placeholder="选择任务日期"
                    format="yyyyMMdd"
                    value-format="yyyyMMdd"
                  />
                </el-row>
              </template>
            </div>
          </div>
        </el-col>
      </div>
    </el-dialog>


  </div>
</template>

<script>
import router from '../router'
import 'codemirror/mode/sql/sql.js'  // 语言模式
import 'codemirror/lib/codemirror.js'
import 'codemirror/mode/clike/clike.js'
import 'codemirror/addon/display/autorefresh.js'
import 'codemirror/addon/edit/matchbrackets.js'
import 'codemirror/addon/selection/active-line.js'
import 'codemirror/addon/display/fullscreen.js'
import 'codemirror/addon/hint/show-hint.js'
import 'codemirror/addon/hint/sql-hint.js'
// theme css
import 'codemirror/addon/hint/show-hint.css'
import 'codemirror/addon/display/fullscreen.css'
import sqlFormatter from 'sql-formatter'

const {TopologicalSort} = require('topological-sort')
const axios = require('axios')

import {checkChinese, formatDateLine} from '../conf/utils'
//https://github.com/murongqimiao/DAG-diagram
export default {
  components: {
    editor: require('vue2-ace-editor')
  },
  data() {
    return {
      internalVars: [],
      intervalVars: ['Last7days', 'LastWeek', 'ThisWeek', 'LastMonth', 'ThisMonth'],
      momentVars: ['yyyyMMdd', 'HH'],
      checkOutputMeta: {
        start: false,
        success: false,
        fail: false,
      },
      cmOptions: {
        // codemirror options
        tabSize: 4,
        mode: 'text/x-mysql',   // https://codemirror.net/mode/index.html
        theme: 'default',
        lineNumbers: true,
        dragDrop: false,
        line: true,
        hintOptions: {
          tables: {},
          databases: {}
        },
        extraKeys: {
          'Tab': 'autocomplete',
          'F12': function (cm) {
            console.log(cm.getOption('fullScreen'))
            cm.setOption('fullScreen', !cm.getOption('fullScreen'))
          },
          'Esc': function (cm) {
            if (cm.getOption('fullScreen')) cm.setOption('fullScreen', false)
          }
        }
      },
      dagId: 0,
      helpContent: "自动数仓中和时间相关的变量请使用内置变量. "
        + "表示时间间隔的变量, 可以用.start和.end表示一段时间的开始和结束:[ \"Last7days\", \"LastWeek\", \"ThisWeek\", \"LastMonth\", \"ThisMonth\" ]\n" +
        "表示时间点的变量, 分别用 \"yyyyMMdd\" 表示天; \"HH\" 表示小时]",
      taskType: [{'label': 'hive', 'id': 'hive'},
        {'label': 'mysql', 'id': 'mysql'},
        {'label': 'hdfs', 'id': 'hdfs'},
        {'label': '绿皮', 'id': 'lvpi'},//['mysql', 'hive', 'clickhouse'],
        {'label': '自定义类型', 'id': 'custom'}],
      taskTypeDsidOptions: [],
      outType: [{'label': 'hive', 'id': 'hive'},
        {'label': 'mail', 'id': 'mail'},
        {'label': 'clickhouse', 'id': 'clickhouse'},
        {'label': 'mysql', 'id': 'mysql'},
        {'label': '绿皮', 'id': 'lvpi'},],
      imgOptions: [],
      dsOptions: [],
      options: {
        enableBasicAutocompletion: true,
        enableSnippets: true,
        enableLiveAutocompletion: true/* 自动补全 */
      },
      jsonshow: false,
      dragBus: false,
      busValue: {
        value: 'name',
        pos_x: 100,
        pos_y: 100
      },
      defaultProps: {
        children: 'children',
        label: 'label'
      },
      DataAll: {'edges': [], 'nodes': []},
      jsonEditor: JSON.stringify('//change JSON click left button show the change'),
      drawerVisible: false,
      drawerDAGVisible: false,
      drawerDAGLoading: false,
      fieldDescription: false,
      parsing: false,
      // newOrOldOutTable: '',
      nodeForm: {
        tmpName: '',
        taskDesc: '',
        tmpNodeId: '',
        taskType: '',
        taskTypeInfo: {
          dsid: 0,
          hdfsLocation: '',
          splitter: '',
          validFieldCount: '0',
          schema: [{'id': '', 'name': ''}],
          extendSchema: [{'id': '', 'name': ''}],
          hasExtend: false,
          hdfsDoneLocation: '',
          fileEncode: 'GBK',
          isCheckDone: false,
          downloadHdfsPath:"",
          cmd:"",
          hdfsAuthUser:"",
          hdfsAuthPassword:"",
        },
        dockerImg: '',
        dockerCmd: '',
        param: '',
        originalSql: '',
        sqlFields: [],  // [{"fieldName":"", description:""}]
        outTableType: '',
        outTableName: '',
        outTableDesc: '',
        outTableNameOptions: [],
        envs: [], //[{env_key: '', env_val: '',}],
        storageEngine: '',
        storageEngineOptions: ['orc/snappy', 'tsv'],
        newOrOldOutTable: '',
        email: {
          title: '',
          receiver: '',
          withAttachment: false,
        },
        clickhouse: {
          partitionStr: '',
          indexStr: ''
        },
        lvpi: {
          jdbcStr: '',
          dsid: 0
        },
      },
      dagForm: {
        name: '',
        exhibitName: '',
        startDate: formatDateLine(new Date()),
        mailList: '',
        cronExpression: '',
        dagType: 'day',
        delay: '',
        retries: '',
        retryDelay: '',
      },
      cronPlaceHolder: '',
      checkTaskName: {
        start: false,
        success: false,
        fail: false,
      },
      isCheck: false,
      submitDialogVisible: false,
      submitActivities: [{
        content: '创建DAG任务输出表',
        timestamp: '',
        type: 'primary',
        icon: '',
      }, {
        content: '连通airflow',
        timestamp: '',
        type: '',
        icon: '',
      },
        {
          content: '生成表层级血缘关系',
          timestamp: '',
          type: '',
          icon: '',
        }, {
          content: '等待完成',
          timestamp: '',
          type: '',
          icon: '',
        }],

      // SQL param
      filterText: '',
      treeData: '',
      sqlDialogVisible: false,
      defaultSqlProps: {
        children: 'children',
        label: 'label',
        templateId: 'templateId',
      },
      form: {
        name: '',
        sql: '',
      },
      paramNames: [],
      paramTypes: [],
      paramValues: [],
      jsonParams: '',
      description: '',
      paramOptions: [{
        label: '一般类型',
        value: 'string',
        defaultValue: ''
      }, {
        label: '日期类型',
        value: 'date',
        defaultValue: '19990101'
      },
        //{
        //  label: '查询词文件分区',
        //  value: 'qw_partition',
        //  defaultValue: ''
        //}
      ],
    }
  },
  computed: {
    getInputTypeOptionByImg: function () {
      console.log("调用计算了")
      let targetImgName = ""
      this.imgOptions.forEach(img => {
        if (img.Id == this.nodeForm.dockerImg) {
          targetImgName = img.Img
        }
      })
      if ((targetImgName + "").indexOf('custom') === -1) {
        console.log("不包含 custom")
        this.taskType = [{'label': 'hive', 'id': 'hive'},
          {'label': 'mysql', 'id': 'mysql'},
          {'label': 'hdfs', 'id': 'hdfs'},
          {'label': '绿皮', 'id': 'lvpi'}]
      } else {
        console.log("包含 custom")
        this.taskType = [{'label': '自定义类型', 'id': 'custom'}]
      }
      try {
        this.nodeForm.taskType = this.taskType[0].id
      } catch (e) {
        console.error(e)
      }
      return this.taskType
    },
    onParsing: function () {
      return this.parsing
    },
    parseButtonIcon: function () {
      if (this.parsing) {
        return 'el-icon-loading'
      } else {
        return ''
      }
    },
    parseButtonText: function () {
      if (this.parsing) {
        return '字段分析中'
      } else {
        return '配置表'
      }
    },
    checkJobName() {
      // 只支持 a-z, 0-9, 中横线 的任务名和dag名
      let re = new RegExp('^[0-9a-z\-]+$')
      return re.test(this.nodeForm.tmpName)
      // return (!checkChinese(this.nodeForm.tmpName))&&(this.nodeForm.tmpName.indexOf('_') === -1)
    },
    codemirror() {
      console.log('计算了')
      return this.$refs.cmEditor.codemirror
    },
    onProcessing: function () {
      return this.isProcessing()
    },
    buttonText: function () {
      if (this.isProcessing()) {
        return '提交中'
      } else if (!this.isSubmitSuccess()) {
        return '保存下次修改'
      } else {
        return '确定'
      }
    },
    buttonIcon: function () {
      if (this.isProcessing()) {
        return 'el-icon-loading'
      } else {
        return ''
      }
    },
    getOutputType: function () {
      return this.nodeForm.outTableType
    },
    getNewOrOldOutTable: function () {
      let isMail = this.nodeForm.outTableType === 'mail'
      // let isHdfsInput = this.nodeForm.taskType === 'hdfs'
      if (isMail) {
        return ''
      }
      return this.nodeForm.newOrOldOutTable
    }
  },
  watch: {
    filterText(val) {
      this.$refs.tree.filter(val)
    },
  },
  methods: {
    delExtendSchema(idx, type) {
      if (type === 'extendSchema') {
        if (this.nodeForm.taskTypeInfo.extendSchema.length > 1) {
          this.nodeForm.taskTypeInfo.extendSchema.splice(idx, 1)
        }
      } else {
        if (this.nodeForm.taskTypeInfo.schema.length > 1) {
          this.nodeForm.taskTypeInfo.schema.splice(idx, 1)
        }
      }

    },
    addExtendSchema(idx, type) {
      let newItem = {'id': '', 'name': ''}
      if (type === 'extendSchema') {
        this.nodeForm.taskTypeInfo.extendSchema.splice(idx + 1, 0, newItem)
      } else {
        this.nodeForm.taskTypeInfo.schema.splice(idx + 1, 0, newItem)
      }
    },
    exploreEmail() {
      this.parsing = true
      let tableHtml = ''
      let param = {
        'sql': this.nodeForm.originalSql,
        'params': this.nodeForm.sqlParam,
        'taskType': this.nodeForm.taskType,
      }
      let self = this
      axios.post('/airflow/getSqlFields', param, {responseType: 'json'})
        .then(function (response) {
          let fields = JSON.parse(response.data.Info)
          console.log(fields)
          // let fieldArray = []
          if (fields.length === 0) {
            self.$notify({
              type: 'info',
              duration: 5000,
              message: '未能静态解析出字段, 请确保配置了所有参数. 或是使用了特殊语法'
            })
          } else {
            let headerHtml = ''
            let columnHtml = ''
            fields.forEach(function (f) {
              headerHtml = headerHtml + '\n<th style="border: 1px solid black;"><span style="color: white">' + f + '</span></th>'
              columnHtml = columnHtml + '\n<td style="border: 1px solid black;"> - </td>'
            })
            tableHtml = '<table style="width:100%; border: 1px solid black; border-collapse: collapse; font-size: 11px; table-layout:fixed;">' +
              '<tr bgcolor="#2561CF" margin="1" align="center" >' + headerHtml +
              '</tr>' +
              '<tr align="center">' + columnHtml +
              '</tr>' +
              '</table>'
            console.log(tableHtml)
            document.getElementById('emailExplore').innerHTML = tableHtml
          }
          self.parsing = false
        }).catch(function (error) {
        console.log(error)
        self.parsing = false
      })
    },
    afterSubmitDag() {
      console.log('是否创建成功:' + this.isSubmitSuccess())
      let self = this
      if (!this.isSubmitSuccess()) {  // 创建失败, 保存失败的dag
        axios.post('/airflow/saveDagWhenPostFail', self.DataAll, {responseType: 'json'})
          .then(function (response) {
            let res = response.data
            if (res.Res) {
              self.submitDialogVisible = false
              self.$notify({
                message: '保存成功',
                type: 'success'
              })
              self.isCheck = true //保存失败的dag后切换到查看/修改模式
            } else {
              self.$notify({
                message: res.Info,
                type: 'error'
              })
            }
          })
          .catch(function (error) {
            console.log(error)
            self.checkOutputMeta.start = false
          })
      } else {  // 创建成功
        this.submitDialogVisible = false
      }

    },
    // 配置输出表时, 如果选择的旧表, 要校验元数据是否和dataFrame匹配
    checkMetaData() {
      if (this.nodeForm.outTableName === '') {
        return
      }
      let self = this
      let param = {
        'sql': this.nodeForm.originalSql,
        'params': this.nodeForm.sqlParam,
        'outTableName': this.nodeForm.outTableName,
        'outTableType': this.nodeForm.outTableType,
        'taskType': this.nodeForm.taskType,
      }
      // console.log('(*********)')
      // console.log(param)
      self.checkOutputMeta.start = true
      axios.post('/airflow/checkOutTable', param, {responseType: 'json'})
        .then(function (response) {
          self.checkOutputMeta.start = false
          let res = response.data
          if (res.Res) {
            self.checkOutputMeta.success = true
            self.checkOutputMeta.fail = false
          } else {
            self.checkOutputMeta.fail = true
            self.checkOutputMeta.success = false
            self.$notify({
              message: res.Info,
              type: 'error'
            })
          }
        })
        .catch(function (error) {
          console.log(error)
          self.checkOutputMeta.start = false
        })
    },
    // 获取输出表选项
    getOutTableNameOptions() {
      let param = {
        'tableType': this.nodeForm.outTableType,
        'dsId': this.nodeForm.lvpi.dsid + '',
      }
      let self = this
      axios.post('/airflow/getOutTables', param, {responseType: 'json'})
        .then(function (response) {
          self.nodeForm.outTableNameOptions = response.data
        })
        .catch(function (error) {
          console.log(error)
        })
    },
    echoSqlParam(sqlParamStr) {
      console.log('sqlParam: ' + sqlParamStr)
      let self = this
      this.paramNames = []
      this.paramTypes = []
      this.paramValues = []
      if (sqlParamStr !== undefined && sqlParamStr !== '') {
        let paramArr = JSON.parse(sqlParamStr)
        paramArr.forEach(function (ele) {
          self.paramNames.push(ele.name)
          self.paramTypes.push(ele.type)
          self.paramValues.push(ele.value)
        })
      }
    },
    parseSql() {
      this.parsing = true
      let sql = this.nodeForm.originalSql
      let sqlParam = this.nodeForm.sqlParam
      this.echoSqlParam(sqlParam)
      let param = {
        'params': sqlParam,
        'sql': sql,
        'taskType': this.nodeForm.taskType,
      }
      console.log(param)
      let self = this
      axios.post('/airflow/getSqlFields', param, {responseType: 'json'})
        .then(function (response) {
          let data = response.data
          let fields = []
          if (data.Res) {
            fields = JSON.parse(data.Info)
            console.log(fields)
            let fieldArray = []
            if (fields.length === 0) {
              self.$notify({
                type: 'info',
                duration: 5000,
                message: '未能静态解析出字段, 请确保配置了所有参数. 或是使用了特殊语法'
              })
            } else {
              fields.forEach(function (f) {
                fieldArray.push({
                  'fieldName': f,
                  'description': ''
                })
              })
              self.nodeForm.sqlFields = fieldArray
            }
          } else {   // 接口访问出错
            self.$notify({
              type: 'error',
              duration: 5000,
              message: data.Info
            })
            self.nodeForm.sqlFields = []
          }
          self.parsing = false
        }).catch(function (error) {
        self.parsing = false
      })
    },
    isProcessing() {
      let lastIndex = this.submitActivities.length - 1
      return this.submitActivities[lastIndex].type === ''
    },
    isSubmitSuccess() {
      let lastIndex = this.submitActivities.length - 1
      return this.submitActivities[lastIndex].type !== 'danger'
    }
    ,
    handleCloseSubmit(done) {
      if (this.isProcessing()) {
        this.$confirm('DAG在生成中,确认关闭吗?')
          .then(_ => {
            done()
          })
          .catch(_ => {
          })
      } else {
        done()
      }
    },

    closeFields(done) {
      console.log(this.nodeForm.sqlFields)
      done()
    },
    getImageById(id) {
      if (this.imgOptions === null) {
        return
      }
      for (let i = 0; i < this.imgOptions.length; i++) {
        if (this.imgOptions[i].Id === id) {
          return this.imgOptions[i]
        }
      }
    },
    getImageIdByImageName(imgName){
      for (let i = 0; i < this.imgOptions.length; i++) {
        if (this.imgOptions[i].Img === imgName) {
          return this.imgOptions[i].Id
        }
      }
    },
    onChangeImg(row) {
      console.log(this.imgOptions)
      console.log(row)
      let img = this.getImageById(row)
      this.nodeForm.dockerCmd = img.Command
    },
    checkDuplicateName() {
      let re = new RegExp('^[0-9a-z\-]+$')
      // let hasChinese = checkChinese(this.dagForm.name)   // 中文检测
      let hasIllegal = !re.test(this.dagForm.name)
      if (hasIllegal) {
        this.checkTaskName.start = false
        this.checkTaskName.success = false
        this.checkTaskName.fail = true
        return
      }
      this.checkTaskName = {
        start: true,
        success: false,
        fail: false,
      }
      // alert('访问接口, 检测重名')
      let self = this
      axios.get('/airflow/checkDagName/' + this.dagForm.name).then(function (response) {
        let res = response.data.Res
        self.checkTaskName.start = false
        if (!res) {  //没重名
          self.checkTaskName.success = true
          self.checkTaskName.fail = false
        } else {
          self.checkTaskName.success = false
          self.checkTaskName.fail = true
        }
      }).catch(function (error) {
        self.checkTaskName.start = false
      })
    },
    markActivityStatus(index, typeAttr, iconAttr, timeAttr) {
      this.submitActivities[index].type = typeAttr
      this.submitActivities[index].icon = iconAttr
      this.submitActivities[index].timestamp = timeAttr

      if (typeAttr === 'danger' && index !== 2) {
        let lastIndex = this.submitActivities.length - 1
        this.submitActivities[lastIndex].type = 'danger'
        this.submitActivities[lastIndex].icon = 'el-icon-error'
        this.submitActivities[lastIndex].timestamp = 'DAG创建失败'
      }
    },
    resetActivityStatus() {
      this.submitActivities = [{
        content: '创建DAG任务输出表',
        timestamp: '',
        type: '',
        icon: '',
      }, {
        content: '连通airflow',
        timestamp: '',
        type: '',
        icon: '',
      }, {
        content: '生成表层级血缘关系',
        timestamp: '',
        type: '',
        icon: '',
      }, {
        content: '等待完成',
        timestamp: '',
        type: '',
        icon: '',
      }]
    },
    onUpdateDAG() {
      let self = this
      if (this.checkCycle()) {
        this.$alert('该DAG中存在环路, 请调整后提交', '注意', {
          confirmButtonText: '确定'
        })
        return
      }
      // console.log(this.DataAll)
      this.DataAll.dagConfig = this.dagForm
      this.submitDialogVisible = true
      // step1. 创建输出表
      this.resetActivityStatus()
      this.markActivityStatus(0, 'primary', 'el-icon-loading', '')
      axios.post('/airflow/preAddDag', self.DataAll, {responseType: 'json'})
        .then(function (response) {
          let res = response.data
          if (res.Res) {  // 第一步成功
            self.markActivityStatus(0, 'success', 'el-icon-success', 'success!')
            // step2. 上传文件
            axios.post('/airflow/updateDag?dagName=' + self.dagId, self.DataAll, {responseType: 'json'})
              .then(function (response) {
                let res = response.data
                if (res.Res) {  // 第二步成功
                  self.markActivityStatus(1, 'success', 'el-icon-success', 'success!')
                  axios.post('/airflow/afterPostDag', self.DataAll, {responseType: 'json'})
                    .then(function (response) {
                      if (response.data.Res) { // 第三步开始
                        self.isCheck = true // 创建成功后,切换到查看模式
                        self.markActivityStatus(2, 'success', 'el-icon-success', 'success!')
                      } else {
                        self.markActivityStatus(2, 'danger', 'el-icon-error', response.data.Info)
                      }
                      // 这一步不影响最终的结果
                      self.markActivityStatus(3, 'success', 'el-icon-success', 'DAG创建成功!')
                    })
                    .catch(function (error) {
                      console.error(error)
                      self.markActivityStatus(2, 'danger', 'el-icon-error', error)
                      // 这一步不影响最终的结果
                      self.markActivityStatus(3, 'success', 'el-icon-success', 'DAG创建成功!')
                    })
                } else {
                  self.markActivityStatus(1, 'danger', 'el-icon-error', res.Info)
                }
                // router.push({path: 'airflow'})
              }).catch(function (error) {
              console.error(error)
              self.markActivityStatus(1, 'danger', 'el-icon-error', error)
            })
          } else {
            self.markActivityStatus(0, 'danger', 'el-icon-error', res.Info)
          }
        }).catch(function (error) {
        self.markActivityStatus(0, 'danger', 'el-icon-error', error)
        console.error(error)
      })
    },
    onSubmitDAG() {
      let self = this
      if (!this.checkTaskName.success || checkChinese(this.dagForm.name) || this.dagForm.name.indexOf('_') !== -1) {
        this.$alert('DAG名称重复或含有中文,下划线, 请更改', '注意', {
          confirmButtonText: '确定'
        })
        return
      }
      let nodes = this.DataAll.nodes
      for (let i = 0; i < nodes.length; ++i) {
        if (checkChinese(nodes[i].name)) {
          this.$alert('DAG任务的名称中含有中文, 请修改', '注意', {
            confirmButtonText: '确定'
          })
          return
        }
      }

      if (this.dagForm.name === '') {
        this.$alert('请先配置DAG调度信息', '注意', {
          confirmButtonText: '确定'
        })
        return
      }
      if (!this.checkCycle()) {
        this.DataAll.dagConfig = this.dagForm
        console.log(this.DataAll)
        self.submitDialogVisible = true
        // step1. 创建输出表
        this.resetActivityStatus()
        this.markActivityStatus(0, 'primary', 'el-icon-loading', '')
        axios.post('/airflow/preAddDag', self.DataAll, {responseType: 'json'})
          .then(function (response) {
            let res = response.data
            if (res.Res) {  // 第一步成功
              self.markActivityStatus(0, 'success', 'el-icon-success', 'success!')
              // step2. 上传文件
              axios.post('/airflow/addDag', self.DataAll, {responseType: 'json'})
                .then(function (response) {
                  let res = response.data
                  if (res.Res) {  // 第二步成功
                    self.markActivityStatus(1, 'success', 'el-icon-success', 'success!')
                    axios.post('/airflow/afterPostDag', self.DataAll, {responseType: 'json'})
                      .then(function (response) {
                        if (response.data.Res) { // 第三步开始
                          self.isCheck = true // 创建成功后,切换到查看模式
                          self.markActivityStatus(2, 'success', 'el-icon-success', 'success!')
                        } else {
                          self.markActivityStatus(2, 'danger', 'el-icon-error', response.data.Info)
                        }
                        // 这一步不影响提交结果
                        self.markActivityStatus(3, 'success', 'el-icon-success', 'DAG创建成功!')
                      })
                      .catch(function (error) {
                        console.error(error)
                        self.markActivityStatus(2, 'danger', 'el-icon-error', error)
                        // 这一步不影响最终的结果
                        self.markActivityStatus(3, 'success', 'el-icon-success', 'DAG创建成功!')

                      })
                  } else {
                    self.markActivityStatus(1, 'danger', 'el-icon-error', res.Info)
                  }
                  // router.push({path: 'airflow'})
                }).catch(function (error) {
                console.error(error)
                self.markActivityStatus(1, 'danger', 'el-icon-error', error)
              })
            } else {
              self.markActivityStatus(0, 'danger', 'el-icon-error', res.Info)
            }
          }).catch(function (error) {
          self.markActivityStatus(0, 'danger', 'el-icon-error', error)
          console.error(error)
        })
      } else {
        this.$alert('该DAG中存在环路, 请调整后提交', '注意', {
          confirmButtonText: '确定'
        })
      }
    },
    checkCycle() {
      let flag = false
      let nodes = this.DataAll.nodes
      let edges = this.DataAll.edges

      const nodeMap = new Map()
      for (let i = 0; i < nodes.length; i++) {
        nodeMap.set(nodes[i].id + '', nodes[i].id)
      }
      const sortOp = new TopologicalSort(nodeMap)
      for (let i = 0; i < edges.length; i++) {
        sortOp.addEdge(edges[i].src_node_id + '', edges[i].dst_node_id + '')
      }
      try {
        let sorted = sortOp.sort()
        console.log('没环')
        flag = false
      } catch (e) {
        console.log(e)
        console.log('有环')
        flag = true
      }
      return flag
    },
    resetForm() {
      this.nodeForm.tmpName = 'new-node'
      this.nodeForm.tmpNodeId = ''
    },
    /**
     * 保存配置, 更新DataAll
     */
    findIdAndUpdateJson() {
      let id = this.nodeForm.tmpNodeId
      let self = this
      for (let i = 0; i < this.DataAll.nodes.length; i++) {
        if (this.DataAll.nodes[i].id === id) {
          this.DataAll.nodes[i].name = this.nodeForm.tmpName   // 节点名作为taskId和taskName
          this.DataAll.nodes[i].taskDesc = this.nodeForm.taskDesc
          let dockerImg = this.getImageById(this.nodeForm.dockerImg)
          console.log('img: ' + dockerImg)
          if (dockerImg !== undefined) {
            this.DataAll.nodes[i].image = dockerImg.Img
          }
          this.DataAll.nodes[i].command = this.nodeForm.dockerCmd
          this.DataAll.nodes[i].envs = this.nodeForm.envs
          this.DataAll.nodes[i].param = this.nodeForm.param
          this.DataAll.nodes[i].outTableType = this.nodeForm.outTableType
          this.DataAll.nodes[i].outTableName = this.nodeForm.outTableName
          this.DataAll.nodes[i].outTableDesc = this.nodeForm.outTableDesc
          this.DataAll.nodes[i].taskType = this.nodeForm.taskType
          this.DataAll.nodes[i].taskTypeInfo = this.nodeForm.taskTypeInfo
          this.DataAll.nodes[i].originalSql = this.nodeForm.originalSql
          this.DataAll.nodes[i].sqlParam = this.nodeForm.sqlParam
          this.DataAll.nodes[i].sqlFields = this.nodeForm.sqlFields
          this.DataAll.nodes[i].storageEngine = this.nodeForm.storageEngine
          this.DataAll.nodes[i].newOrOldOutTable = this.nodeForm.newOrOldOutTable
          this.DataAll.nodes[i].email = this.nodeForm.email

          this.DataAll.nodes[i].lvpi = this.nodeForm.lvpi
          this.DataAll.nodes[i].clickhouse = this.nodeForm.clickhouse
          // console.log('当前节点已更新为:' + JSON.stringify(this.DataAll.nodes[i], null, 4))
        }
      }
      this.jsonEditor = JSON.stringify(this.DataAll, null, 4)
      this.saveChange()
    },
    handleDAGClose(done) {
      this.drawerDAGVisible = false
    },
    handleClose(done) {
      // 从jsonEditor中找到当前的id节点, 并更新name属性
      this.findIdAndUpdateJson()
      this.cancelForm()
      console.log('json已更新')
      console.log(this.DataAll)
    },
    cancelForm() {
      this.drawerVisible = false
      this.nodeForm = {
        tmpName: '',
        taskDesc: '',
        tmpNodeId: '',
        taskType: '',
        taskTypeInfo: {
          dsid: 0,
          hdfsLocation: '',
          splitter: '',
          validFieldCount: '0',
          schema: [{'id': '', 'name': ''}],
          extendSchema: [{'id': '', 'name': ''}],
          hasExtend: false,
          hdfsDoneLocation: '',
          isCheckDone: false,
        },
        dockerImg: '',
        dockerCmd: '',
        param: '',
        originalSql: '',
        sqlFields: [],  // [{"fieldName":"", description:""}]
        outTableType: '',
        outTableName: '',
        outTableDesc: '',
        outTableNameOptions: [],
        envs: [], //[{env_key: '', env_val: '',}],
        storageEngine: '',
        storageEngineOptions: ['orc/snappy', 'tsv'],
        newOrOldOutTable: '',
        email: {
          title: '',
          receiver: '',
          withAttachment: false,
        },
        clickhouse: {
          partitionStr: '',
          indexStr: ''
        },
        lvpi: {
          jdbcStr: '',
          dsid: 0
        }
      }
      this.resetForm()
    },
    dragIt(val) {
      sessionStorage['dragDes'] = JSON.stringify({
        drag: true,
        name: val
      })
    },
    updateDAG(data, actionName, id) {
      // DAG-Board更新 data是最新的数据,  actionName是事件的名称
      this.DataAll = data
      this.jsonEditor = JSON.stringify(this.DataAll, null, 4)
      // console.log('actionName:', actionName, 'id:', id, 'data:', data)
    },
    saveChange() {
      this.DataAll = JSON.parse(this.jsonEditor)
    },
    startNodesBus(e) {
      /**
       *  别的组件调用时, 先放入缓存
       * dragDes: {
       *    drag: true,
       *    name: 组件名称
       *    type: 组件类型
       *    model_id: 跟后台交互使用
       * }
       **/
      let dragDes = null
      if (sessionStorage['dragDes']) {
        dragDes = JSON.parse(sessionStorage['dragDes'])
      }
      if (dragDes && dragDes.drag) {
        const x = e.pageX
        const y = e.pageY
        this.busValue = Object.assign({}, this.busValue, {
          pos_x: x,
          pos_y: y,
          value: dragDes.name
        })
        this.dragBus = false
      }
    },
    moveNodesBus(e) {
      // 移动模拟节点
      if (this.dragBus) {
        const x = e.pageX
        const y = e.pageY
        this.busValue = Object.assign({}, this.busValue, {
          pos_x: x,
          pos_y: y
        })
      }
    },
    endNodesBus(e) {
      // 节点放入svg
      let dragDes = null
      if (sessionStorage['dragDes']) {
        dragDes = JSON.parse(sessionStorage['dragDes'])
      }
      if (dragDes && dragDes.drag && e.toElement.id === 'svgContent') {
        const {model_id, type} = dragDes
        const pos_x =
          (e.offsetX - 15 - (sessionStorage['svg_left'] || 0)) /
          (sessionStorage['svgScale'] || 1) // 参数修正
        const pos_y =
          (e.offsetY - 15 - (sessionStorage['svg_top'] || 0)) /
          (sessionStorage['svgScale'] || 1) // 参数修正
        const params = {
          model_id: sessionStorage['newGraph'],
          desp: {
            type,
            pos_x,
            pos_y,
            name: this.busValue.value,
            rightClickEvent: [
              {
                'label': '配置DAG调度',
                'eventName': 'editDAG'
              },
              {
                'label': '配置任务',
                'eventName': 'editTask'
              }
            ],
          }
        }
        this.DataAll.nodes.push({
          ...params.desp,
          id: this.DataAll.nodes.length + 100,
          in_ports: [0],
          out_ports: [0]
        })
      }
      window.sessionStorage['dragDes'] = null
      this.dragBus = false
    },
    handleNodeClick(data) {
      console.log('handleNodeClick')
      const {value, url} = data
      if (url) {
        this.$router.push(url)
        return false
      }
      if (value) {
        localStorage['currentExample'] = JSON.stringify(value)
        this.DataAll = value
        this.jsonEditor = JSON.stringify(this.DataAll, null, '\t')
        this.jsonshow = true
        this.svgScale = sessionStorage['svgScale'] = 1.5
      }
    },
    editorInit: function (editor) { // 右侧JSON相关可以忽略
      require('brace/ext/language_tools') // language extension prerequsite...
      require('brace/mode/html')
      require('brace/mode/javascript') // language
      require('brace/mode/less')
      require('brace/theme/chrome')
      require('brace/snippets/javascript') // snippet
      require('brace/mode/json')
      require('brace/theme/tomorrow')
    },
    editNodeDetails(value) {
      alert(`edit id ${value.id} , info : ${JSON.stringify(value.detail, null, 4)} `)
    },
    doSthPersonal(eventName, id) {
      console.log(eventName + ':' + id)
      // 如果是单选事件, 则打开编辑面板
      if (eventName === 'editTask') {
        this.nodeForm.tmpNodeId = id
        this.getNodeInfoById(id)
        this.drawerVisible = true
        console.log(this.nodeForm)
      } else if (eventName === 'editDAG') {
        this.nodeForm.tmpNodeId = id
        // this.getNodeInfoById(id)
        this.drawerDAGVisible = true
      }
    },
    /**
     * 回显dag节点
     * @param id
     */
    getNodeInfoById(id) {
      let self = this
      for (let i = 0; i < this.DataAll.nodes.length; i++) {
        if (this.DataAll.nodes[i].id === id) {
          console.log('当前节点数据:')
          console.log(this.DataAll.nodes[i])
          this.nodeForm.tmpName = this.DataAll.nodes[i].name
          this.nodeForm.taskDesc = this.DataAll.nodes[i].taskDesc
          this.nodeForm.taskType = this.DataAll.nodes[i].taskType
          axios.get('/airflow/GetConnByPrivilege/' + this.nodeForm.taskType)
            .then(function (response) {
              self.taskTypeDsidOptions = response.data
            })
            .catch(function (error) {
              console.log(error)
            })
          // 查找数据源
          this.nodeForm.dockerImg = this.getImageIdByImageName(this.DataAll.nodes[i].image)

          this.nodeForm.dockerCmd = this.DataAll.nodes[i].command
          this.nodeForm.envs = this.DataAll.nodes[i].envs === undefined ? [] : this.DataAll.nodes[i].envs
          this.nodeForm.param = this.DataAll.nodes[i].param
          this.nodeForm.outTableType = this.DataAll.nodes[i].outTableType
          this.nodeForm.outTableName = this.DataAll.nodes[i].outTableName
          this.nodeForm.outTableDesc = this.DataAll.nodes[i].outTableDesc
          this.nodeForm.originalSql = this.DataAll.nodes[i].originalSql
          this.nodeForm.sqlParam = this.DataAll.nodes[i].sqlParam
          this.nodeForm.sqlFields = this.DataAll.nodes[i].sqlFields
          this.nodeForm.storageEngine = this.DataAll.nodes[i].storageEngine
          this.nodeForm.storageEngineOptions = ['orc/snappy', 'tsv']
          this.nodeForm.newOrOldOutTable = this.DataAll.nodes[i].newOrOldOutTable
          if (this.DataAll.nodes[i].email !== undefined) {
            this.nodeForm.email = {
              title: this.DataAll.nodes[i].email.title === undefined ? '' : this.DataAll.nodes[i].email.title,
              receiver: this.DataAll.nodes[i].email.receiver === undefined ? '' : this.DataAll.nodes[i].email.receiver,
              withAttachment: this.DataAll.nodes[i].email.withAttachment === undefined ? false : this.DataAll.nodes[i].email.withAttachment,
            }
          }
          console.log(this.DataAll.nodes[i].clickhouse)
          if (this.DataAll.nodes[i].clickhouse !== undefined) {
            this.nodeForm.clickhouse = {
              partitionStr: this.DataAll.nodes[i].clickhouse.partitionStr === undefined ? '' : this.DataAll.nodes[i].clickhouse.partitionStr,
              indexStr: this.DataAll.nodes[i].clickhouse.indexStr === undefined ? '' : this.DataAll.nodes[i].clickhouse.indexStr,
            }
          }
          if (this.DataAll.nodes[i].lvpi !== undefined) {
            this.nodeForm.lvpi = {
              jdbcStr: this.DataAll.nodes[i].lvpi.jdbcStr === undefined ? '' : this.DataAll.nodes[i].lvpi.jdbcStr,
              dsid: this.DataAll.nodes[i].lvpi.dsid === undefined ? 0 : this.DataAll.nodes[i].lvpi.dsid,
            }
          }
          if (this.DataAll.nodes[i].taskTypeInfo !== undefined) {
            this.nodeForm.taskTypeInfo = this.DataAll.nodes[i].taskTypeInfo
            /*this.nodeForm.taskTypeInfo = {
              dsid: this.DataAll.nodes[i].taskTypeInfo.dsid === undefined ? '' : this.DataAll.nodes[i].taskTypeInfo.dsid,
              hdfsLocation: this.DataAll.nodes[i].taskTypeInfo.hdfsLocation === undefined ? '' : this.DataAll.nodes[i].taskTypeInfo.hdfsLocation,
            }*/
          }

          //预加载SQL模板相关数据
          this.form.name = ''
          this.paramNames = []
          this.paramTypes = []
          this.paramValues = []
          this.form.sql = this.nodeForm.originalSql == null ? '' : this.nodeForm.originalSql
          this.jsonParams = this.nodeForm.sqlParam == null ? '' : this.nodeForm.sqlParam

          if (this.jsonParams != '') {
            let paramList = JSON.parse(this.jsonParams)
            for (let i = 0; i < paramList.length; i++) {
              this.paramNames.push(paramList[i]['name'])
              this.paramTypes.push(paramList[i]['type'])
              this.paramValues.push(paramList[i]['value'])
            }
          }

          this.getConnection()   // 获取数据源选项
        }
      }
      console.log()
    },
    onAddParams() {
      this.nodeForm.envs.push({
        env_key: '',  // airflow 任务参数键
        env_val: '', // airflow 任务参数值
      })
    },
    removeNewParam(item) {
      let index = this.nodeForm.envs.indexOf(item)
      if (index !== -1) {
        this.nodeForm.envs.splice(index, 1)
      }
    },
    onChangeDAGType(selectValue) {
      if (selectValue === 'hour') {
        this.cronPlaceHolder = '示例:55 * * * *'
      } else if (selectValue === 'day') {
        this.cronPlaceHolder = '示例:00 15 * * *'
      }
    },
    loadDagData(dagId) {
      let self = this
      if (dagId == null) {
        return
      }
      axios.get('/airflow/getDag?dagName=' + dagId).then(function (response) {
        console.log(response.data)
        let dataJson = response.data
        self.dagForm = dataJson.dagConfig  // dag调度的描述
        delete dataJson.dagConfig
        self.DataAll = dataJson            // task调度的描述

        self.isCheck = true
      }).catch(function (error) {
        console.log(error)
      })

    },
    CreateSQLByTemplate() {
      let self = this
      this.form.sql = this.nodeForm.originalSql === undefined ? '' : this.nodeForm.originalSql
      let queryType = self.nodeForm.taskType

      axios.get('/airflow/template/getAllTemplate/' + queryType)
        .then(function (response) {
          self.treeData = response.data
        })
        .catch(function (error) {
          console.log(error)
        })

      self.sqlDialogVisible = true
      // console.log('=================')
      console.log(self.nodeForm)
    },
    filterNode(value, data) {
      if (!value) return true
      return data.label.indexOf(value) !== -1
    },
    resetParam() {
      this.paramTypes = []
      this.paramNames = []
      this.paramValues = []
      this.jsonParams = ''
    },
    qyeryInfo(data) {
      let self = this
      if (data.pid != '0') {
        this.resetParam()
        self.curTmplateId = data.templateId
        console.log('templateId:' + data.templateId)
        console.log('label:' + data.label)
        axios.get('/airflow/template/getById/' + data.templateId)
          .then(function (response) {
            console.log(response.data.Template)
            self.form.name = response.data.Title
            self.form.sql = response.data.Template
            self.description = response.data.Description
            // 触发参数解析
            self.onParam('form')
          }).catch(function (error) {
          console.log(error)
        })
      }
    },
    onParam(formName) {
      this.resetParam()
      let sql = this.form.sql
      let expr = /\{\{(.+?)\}\}/g   // 匹配双括号中的参数名
      let originParams = sql.match(expr)
      console.log('originParams:' + originParams)

      originParams.forEach(p => {
        let attr = p.substring(2, p.length - 2)
        console.log(attr)
        console.log(this.internalVars.indexOf(attr))
        if (this.internalVars.indexOf(attr) === -1) {  // 不在内部变量里, 才显示
          console.log('找到了:' + attr)
          this.paramNames.push(attr)
        }
      })
      console.log('paramNames:' + this.paramNames)
      for (let step = 0; step < this.paramNames.length; step++) {
        this.paramValues.push('')
        this.paramTypes.push('')
      }
    },
    makeParam() {
      let params = []
      for (let i = 0; i < this.paramNames.length; i++) {
        params.push({
          'name': this.paramNames[i],
          'type': this.paramTypes[i],
          'value': this.paramValues[i]
        })
      }
      console.log('paramsName:' + this.paramNames + '\nparamTypes:' + this.paramTypes + '\nparamValues:' + this.paramValues)
      this.jsonParams = JSON.stringify(params)
      console.log('jsonParam:' + this.jsonParams)
    },
    CommitSQL() {
      let self = this
      this.makeParam()

      this.nodeForm.originalSql = this.form.sql
      this.nodeForm.sqlParam = this.jsonParams
      this.sqlDialogVisible = false
    },
    onChangeParamType(index) {
      this.paramValues[index] = ''  // 赋一个默认值
      console.log(this.paramTypes)
    },
    taskTypeChange() {
      this.nodeForm.taskTypeInfo.dsid = 0
      // this.nodeForm.taskTypeInfo.hdfsLocation = ''
      let curTaskType = this.nodeForm.taskType
      let self = this
      axios.get('/airflow/GetConnByPrivilege/' + curTaskType)
        .then(function (response) {
          // console.log(response.data)
          self.taskTypeDsidOptions = response.data
        })
        .catch(function (error) {
          console.log(error)
        })
    },
    initInternalVars() {
      let intervalEndPoints = []
      this.intervalVars.forEach(x => {
        intervalEndPoints.push(x + '.start')
        intervalEndPoints.push(x + '.end')
      })
      this.momentVars.forEach(x => {
        intervalEndPoints.push(x)
      })
      this.internalVars = intervalEndPoints
      console.log(this.internalVars)
    },
    getConnection() {
      let self = this
      let outType = this.nodeForm.outTableType
      if (outType === 'lvpi') {
        outType = 'mysql'
      }
      axios.get('/connMeta/get?connType=' + outType)
        .then(function (response) {
          self.dsOptions = response.data
        })
    }
  },
  created() {
  },
  mounted() {
    let self = this
    let dagId = sessionStorage.getItem('jones_airflow_dagId')
    console.log('需要回显的dag为: ' + dagId)
    this.dagId = dagId

    axios.get('/airflow/getAllDagImgs').then(function (response) {
      self.imgOptions = response.data
    }).catch(function (error) {
      console.log(error)
    })
    this.loadDagData(dagId)
    sessionStorage.removeItem('jones_airflow_dagId')
    this.initInternalVars()
    this.onChangeDAGType('day')  // 默认是天级别任务
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped type='text/css'>

.title {
  font-size: 12px;
  COLOR: #FFFFFF;
  font-family: Tahoma;
}

.el-divider__text {
  color: #C0C0C0;
}

.content {
  min-height: 100%;
  height: 100%;
}

.content-inside {
  padding: 20px;
  margin-right: 20px;
  margin-top: 20px;
  padding-bottom: 100px;
}

.footer {
  margin-left: 20px;
  height: 50px;
  margin-top: -100px;
  margin-right: 30px;
  float: right;
}

.allDiv {
  height: 100%;
}

@keyframes line_success {
  0% {
    stroke-dashoffset: 100%;
    stroke-dasharray: 100%;
  }
  50% {
    stroke-dashoffset: 0%;
    stroke-dasharray: 100%;
  }
  99% {
    stroke-dashoffset: 0%;
    stroke-dasharray: 100%;
  }
  100% {
    stroke-dashoffset: 100%;
    stroke-dasharray: 0px;
  }
}

@keyframes fadeInLeft {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
    transform: none;
  }
}

.page-content {
  position: relative;
  left: 0;
  top: 0;
  right: 0;
  bottom: 0;
  height: 100%;
  background-size: 50px 50px;
  background-image: linear-gradient(0deg, transparent 24%, rgba(255, 255, 255, 0.05) 25%, rgba(255, 255, 255, 0.05) 26%, transparent 27%, transparent 74%, rgba(255, 255, 255, 0.05) 75%, rgba(255, 255, 255, 0.05) 76%, transparent 77%, transparent), linear-gradient(90deg, transparent 24%, rgba(255, 255, 255, 0.05) 25%, rgba(255, 255, 255, 0.05) 26%, transparent 27%, transparent 74%, rgba(255, 255, 255, 0.05) 75%, rgba(255, 255, 255, 0.05) 76%, transparent 77%, transparent);
  background-color: rgb(60, 60, 60) !important;
  padding-left: 10px;
}

.json-editor {
  position: absolute;
  bottom: 0;
  right: 0;
}

.elbutton {
  position: absolute;
  top: 10px;
}

.mainContent {
  width: 100%;
  position: relative;
  left: 0;
  top: 70px;
  bottom: 0;
  height: 100%;
  text-align: left;
}

.mainContent #helpContent {
  margin-top: 100px;
  text-align: left;
  text-indent: 20px;
  line-height: 20px;
}

.mainContent .nav {
  width: 300px;
  background: #212528;
  position: relative;
  left: 0;
  bottom: 0;
  top: 0;
}

.mainContent .nav div {
  color: #ffffff;
  height: 50px;
  line-height: 50px;
  width: 100%;
  text-align: center;
}

.DAG-content {
  position: relative;
  left: 0px;
  top: 0;
  bottom: 0;
  right: 0;
}

.nodes_bus {
  user-select: none;
  text-align: center;
}

.nodes_bus span {
  display: block;
  border: 1px white solid;
  height: 50px;
  width: 200px;
  margin-left: 50%;
  transform: translateX(-50%);
  margin-bottom: 30px;
  cursor: move;
  border-radius: 50px;
}

.CodeMirror {
  border: 1px solid #eee;
  /*height: 500px;*/
  font-size: 15px;
  line-height: 120%;
}
</style>
