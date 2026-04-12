<template>
  <el-card shadow="hover">
    <template #header>
      <span>{{ title }}</span>
    </template>

    <el-table :data="rows" style="width: 100%">
      <el-table-column v-for="col in columns" :key="col.key" :label="col.label" :width="col.width">
        <template #default="{ row }">
          <component
            v-if="row._editing"
            :is="col.editor"
            v-model="row[col.key]"
            v-bind="col.editorProps"
          />

          <span v-else @click="edit(row)">
            {{ formatValue(row[col.key], col) }}
          </span>
        </template>
      </el-table-column>

      <!-- ACTIONS -->
      <el-table-column width="160">
        <template #default="{ row }">
          <template v-if="row._editing">
            <el-button size="small" @click="save(row)">💾</el-button>
            <el-button size="small" @click="cancel(row)">❌</el-button>
          </template>

          <template v-else>
            <el-button size="small" @click="edit(row)">✏️</el-button>
            <el-button size="small" type="danger" @click="remove(row)">🗑</el-button>
          </template>
        </template>
      </el-table-column>
    </el-table>

    <div style="margin-top: 10px">
      <el-button @click="addRow"> + Добавить </el-button>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { reactive, watch } from 'vue'

interface Column {
  key: string
  label: string
  width?: string
  editor: any
  editorProps?: Record<string, any>
  formatter?: (value: any) => string
}

interface Row {
  id: number
  _editing?: boolean
  _isNew?: boolean
  [key: string]: any
}

const props = defineProps<{
  title: string
  items: any[]
  columns: Column[]
  create: (data: any) => Promise<any>
  update: (id: number, data: any) => Promise<void>
  delete: (id: number) => Promise<void>
}>()

const emit = defineEmits(['refresh'])

const rows = reactive<Row[]>([])
const original = new Map<number, any>()

// sync with props
watch(
  () => props.items,
  (items) => {
    rows.splice(0, rows.length, ...items.map((i) => ({ ...i })))
  },
  { immediate: true },
)

const formatValue = (value: any, col: Column) => {
  if (col.formatter) return col.formatter(value)
  return value
}

const edit = (row: Row) => {
  original.set(row.id, { ...row })
  row._editing = true
}

const cancel = (row: Row) => {
  if (row._isNew) {
    const index = rows.indexOf(row)
    if (index !== -1) rows.splice(index, 1)
    return
  }

  const orig = original.get(row.id)
  if (orig) Object.assign(row, orig)

  row._editing = false
}

const save = async (row: Row) => {
  try {
    const payload = { ...row }
    delete payload._editing
    delete payload._isNew

    if (row._isNew) {
      await props.create(payload)
    } else {
      await props.update(row.id, payload)
    }

    row._editing = false
    row._isNew = false

    emit('refresh')
  } catch (e) {
    console.error(e)
  }
}

const remove = async (row: Row) => {
  if (!row.id) return

  try {
    await props.delete(row.id)
    emit('refresh')
  } catch (e) {
    console.error(e)
  }
}

const addRow = () => {
  const empty: Row = {
    id: Date.now(),
    _editing: true,
    _isNew: true,
  }

  props.columns.forEach((c) => {
    empty[c.key] = null
  })

  rows.unshift(empty)
}
</script>
