import { ElNotification } from 'element-plus'

export const notifyError = (message: string, title: string = 'Ошибка') => {
  ElNotification.error({
    title: title,
    position: 'bottom-right',
    message,
  })
}

export const notifySuccess = (message: string, title: string = 'Успешно') => {
  ElNotification.success({
    title: title,
    position: 'bottom-right',
    message,
  })
}
