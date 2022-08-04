<template>
  <div class="table-box">

    <h2 class="title">CRUD</h2>

    <div class="query-box">
      <el-input
          class="el-inp"
          v-model="inputQuery"
          placeholder="请搜索姓名"
      />
      <div>
        <el-button
            type="primary"
            @click="handleAdd('add')">
          增加
        </el-button>
        <el-button
            type="danger"
            @click="handleMultipleDelete()"
            v-if="batchDeleteList.length > 0">
          删除选中
        </el-button>
      </div>
    </div>

    <el-table
        border
        :data="tableData"
        style="width: 100%"
        @selection-change="handleSelectionChange"
        ref="tableRef">

      <el-table-column type="selection" width="55"/>
      <el-table-column prop="name" label="姓名" width="120"/>
      <el-table-column prop="phone" label="电话" width="120"/>
      <el-table-column prop="email" label="邮箱" width="200"/>
      <el-table-column prop="state" label="状态" width="120"/>
      <el-table-column prop="address" label="地址" width="300"/>
      <el-table-column fixed="right" label="操作" width="120" align="center">
        <template #default="scope">
          <el-button
              link
              type="primary"
              style="color: #E47470"
              size="small"
              @click="handleDelete(scope.row)">
            删除
          </el-button>
          <el-button
              link
              type="primary"
              size="small"
              @click="handleEdit('edit',scope.row)">
            编辑
          </el-button>
        </template>
      </el-table-column>

    </el-table>

    <!-- 增加和编辑的dialog -->
    <el-dialog
        v-model="dialogFormVisible"
        draggable
        width="40%"
    >
      <template #header>
        {{ dialogTitle }}
      </template>
      <el-form :model="dialogForm">
        <el-form-item label="姓名">
          <el-input v-model="dialogForm.name" autocomplete="off"/>
        </el-form-item>
        <el-form-item label="电话">
          <el-input v-model="dialogForm.phone" autocomplete="off"/>
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="dialogForm.email" autocomplete="off"/>
        </el-form-item>
        <el-form-item label="状态">
          <el-input v-model="dialogForm.state" autocomplete="off"/>
        </el-form-item>
        <el-form-item label="地址">
          <el-input v-model="dialogForm.address" autocomplete="off"/>
        </el-form-item>
      </el-form>

      <template #footer>
        <div class="dialog-footer">
          <el-button @click="handleDialogCancel">取消</el-button>
          <el-button type="primary" @click="handleDialogConfirm">确认</el-button>
        </div>
      </template>

    </el-dialog>
  </div>
</template>

<script setup>

import {ref, watch} from "vue";

// ------------------- 数据 -------------------

// 注意点 : 需要 watch 监听, 则不能使用 $ref
let inputQuery = ref('') // 搜索框数据
let tableData = $ref([
  {
    id: 1,
    name: 'Tom1',
    state: 'California',
    phone: '13800138000',
    email: '13800138000@qq.com',
    address: 'No. 189, Grove St, Los Angeles',
  },
  {
    id: 2,
    name: 'Tom2',
    state: 'California',
    phone: '13800138000',
    email: '13800138000@qq.com',
    address: 'No. 189, Grove St, Los Angeles',
  },
  {
    id: 3,
    name: 'Tom3',
    state: 'California',
    phone: '13800138000',
    email: '13800138000@qq.com',
    address: 'No. 189, Grove St, Los Angeles',
  },
  {
    id: 4,
    name: 'Tom4',
    state: 'California',
    phone: '13800138000',
    email: '13800138000@qq.com',
    address: 'No. 189, Grove St, Los Angeles',
  },
]) // 表格数据
let copyTableData = Object.assign(tableData) // 浅拷贝表格数据
let dialogFormVisible = $ref(false) // 弹窗显示状态
let dialogType = $ref('add')  // 弹窗状态
let dialogTitle = $ref('新增') // 弹窗名称
let dialogForm = $ref({}) // 弹窗内容信息
let batchDeleteList = $ref([]) // 多选的需要批量删除的列表

// ------------------- 选中元素 -------------------

const tableRef = ref() // ref 选中元素

// ------------------- 方法 -------------------

// 监听搜索框
watch(inputQuery, val => {
  if (val.length > 0) {
    // 过滤自己的name然后使用match正则匹配输入的name
    // toLowerCase 是为了统一名称都是小写, 方便检索
    tableData = tableData.filter(item => item.name.toLowerCase().match(val.toLowerCase()))
  } else {
    tableData = copyTableData
  }
})

// 表格多选
const handleSelectionChange = (row) => {
  batchDeleteList = []
  // 循环遍历数据 , 获取里面的所有 id 放到数组中
  row.forEach(item => {
    batchDeleteList.push(item.id)
  })
}

// 多选删除
const handleMultipleDelete = () => {
  // 循环遍历多个数据
  batchDeleteList.forEach(id => {
    handleDelete({id})
  })
  // 清空选中的需要批量删除的数组
  batchDeleteList = []
  // 清除目前选中的选择框
  tableRef.value.clearSelection()
}

// 删除
const handleDelete = (row) => {
  // 1. 查找索引
  let index = tableData.findIndex(item => item.id === row.id)
  // 2. 根据索引删除对应选项
  tableData.splice(index, 1)
}

// 编辑
const handleEdit = (type, row) => {
  dialogStateModify(type)
  // 把目前所选中的数据全部放到 form 表单中
  dialogForm = {...row}
}

// 增加
const handleAdd = (type) => {
  dialogStateModify(type)
  // 表单清空
  dialogForm = {}
}

// 弹窗状态修改
const dialogStateModify = (type) => {
  dialogFormVisible = true
  dialogType = type
  if (dialogType === 'add') {
    dialogTitle = '新增'
  } else {
    dialogTitle = '编辑'
  }
}

// 弹窗取消
const handleDialogCancel = () => {
  dialogFormVisible = false
}

// 弹窗确认
const handleDialogConfirm = () => {
  dialogFormVisible = false

  // 1. 判断是新增还是编辑
  if (dialogType === 'add') {
    // 1.1. 添加数据
    tableData.push({
      id: tableData.length + 1,
      ...dialogForm
    })
  } else if (dialogType === 'edit') {
    // 2.1. 获取form中的数据的id , 找出对应的索引值
    // 2.2. 编辑数据
    let index = tableData.findIndex(item => dialogForm.id === item.id)
    tableData[index] = dialogForm
  }
}

</script>

<style>
.table-box {
  width: 800px;
  margin: 200px auto;
}

.title {
  text-align: center;
}

.query-box {
  margin-bottom: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.el-inp {
  width: 200px;
}
</style>
