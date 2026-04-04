import { ElNotification } from 'element-plus'

export const notifyError = (message: string) => {
  ElNotification.error({
    title: 'Ошибка',
    position: 'bottom-right',
    message,
  })
}
