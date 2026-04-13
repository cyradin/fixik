<template>
  <el-card shadow="hover">
    <template #header>
      <span>{{ title }}</span>
    </template>

    <el-form :model="formModel" :rules="rules" ref="formRef" show-message inline-message>
      <el-table :data="rows" style="width: 100%">
        <el-table-column
          v-for="col in columns"
          :key="col.key"
          :label="col.label"
          :width="col.width"
        >
          <template #default="{ row }">
            <el-form-item
              v-if="row._editing"
              :prop="`rows.${row._uid}.${col.key}`"
              style="margin-bottom: 0"
            >
              <component
                :is="col.editor"
                v-model="formModel.rows[row._uid][col.key]"
                v-bind="col.editorProps"
                @keyup.enter="save(row)"
              >
                <template v-if="col.editor === 'el-select'">
                  <el-option
                    v-for="opt in col.options"
                    :key="opt.value"
                    :label="opt.label"
                    :value="opt.value"
                  />
                </template>
              </component>
            </el-form-item>

            <span v-else @click="edit(row)">
              {{ formatValue(row[col.key], col) }}
            </span>
          </template>
        </el-table-column>

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
    </el-form>

    <div style="margin-top: 10px">
      <el-button @click="addRow"> + Добавить </el-button>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { reactive, watch, computed, ref, ComputedRef } from 'vue'

interface Column {
  key: string
  label: string
  width?: string
  editor: any
  editorProps?: Record<string, any>
  required?: boolean
  formatter?: (value: any) => string
  options?: { label: string; value: any }[]

  validation?: {
    min?: number
    max?: number
    pattern?: RegExp
    message?: string
  }
}

interface Row {
  id: number
  _uid: number
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
  remove: (id: number) => Promise<void>
  defaultRow?: () => any
  onEditRow?: (row: any, formRow: any) => void
}>()

const emit = defineEmits(['refresh'])

const activeEditUid = ref<number | null>(null)

const rows = reactive<Row[]>([])
const original = new Map<number, any>()

const formRef = ref()

const formModel = reactive({
  rows: {} as Record<number, any>,
})

watch(
  () => props.items,
  (items) => {
    rows.splice(
      0,
      rows.length,
      ...items.map((i) => ({
        ...i,
        _uid: i.id,
      })),
    )

    rows.forEach((row) => {
      formModel.rows[row._uid] = { ...row }
    })
  },
  { immediate: true },
)

const rules = computed(() => {
  const r: any = {}

  rows.forEach((row) => {
    props.columns.forEach((col) => {
      const path = `rows.${row._uid}.${col.key}`

      const fieldRules: any[] = []

      if (col.required) {
        fieldRules.push({
          required: true,
          message: `${col.label} обязательно`,
          trigger: ['blur', 'change'],
        })
      }

      if (col.validation?.min !== undefined) {
        fieldRules.push({
          min: col.validation.min,
          message: col.validation.message || `${col.label} минимум ${col.validation.min}`,
          trigger: ['blur', 'change'],
        })
      }

      if (col.validation?.max !== undefined) {
        fieldRules.push({
          max: col.validation.max,
          message: col.validation.message || `${col.label} максимум ${col.validation.max}`,
          trigger: ['blur', 'change'],
        })
      }

      if (col.validation?.pattern) {
        fieldRules.push({
          pattern: col.validation.pattern,
          message: col.validation.message || `${col.label} некорректный формат`,
          trigger: ['blur'],
        })
      }

      if (fieldRules.length) {
        r[path] = fieldRules
      }
    })
  })

  return r
})

const formatValue = (value: any, col: Column) => {
  if (col.formatter) return col.formatter(value)

  if (col.options?.length) {
    const found = col.options.find((o) => o.value === value)
    return found ? found.label : value
  }

  return value
}

const edit = (row: Row) => {
  closeActiveEdit()

  original.set(row._uid, { ...row })
  row._editing = true

  formModel.rows[row._uid] = { ...row }

  activeEditUid.value = row._uid

  props.onEditRow?.(row, formModel.rows[row._uid])
}

const cancel = (row: Row) => {
  if (row._isNew) {
    const index = rows.indexOf(row)
    if (index !== -1) rows.splice(index, 1)
  } else {
    const orig = original.get(row._uid)
    if (orig) {
      Object.assign(row, orig)
      formModel.rows[row._uid] = { ...orig }
    }
  }

  row._editing = false

  if (activeEditUid.value === row._uid) {
    activeEditUid.value = null
  }
}

const closeActiveEdit = () => {
  if (activeEditUid.value !== null) {
    const prev = rows.find((r) => r._uid === activeEditUid.value)
    if (prev) cancel(prev)
  }
}

const save = async (row: Row) => {
  try {
    await formRef.value.validate()

    const payload = { ...formModel.rows[row._uid] }
    delete payload._editing
    delete payload._isNew
    delete payload._uid

    if (row._isNew) {
      const res = await props.create(payload)
      Object.assign(row, res)
      original.set(row._uid, { ...row })
    } else {
      await props.update(row.id, payload)

      original.set(row._uid, { ...formModel.rows[row._uid] })
    }

    row._editing = false
    row._isNew = false

    emit('refresh')
  } catch (e) {
    throw e
  }
}

const remove = async (row: Row) => {
  if (!row.id) return

  if (!confirm('Удалить?')) return

  try {
    await props.remove(row.id)
    emit('refresh')
  } catch (e) {
    console.error(e)
  }
}

const addRow = () => {
  closeActiveEdit()

  const uid = Date.now()

  const base = props.defaultRow ? props.defaultRow() : {}

  const empty: Row = {
    id: 0,
    _uid: uid,
    _editing: true,
    _isNew: true,
    ...base,
  }

  rows.unshift(empty)
  formModel.rows[uid] = { ...empty }

  activeEditUid.value = uid
}
</script>
