export const formatDateTime = (dateStr: string): string => {
  const d = new Date(dateStr)

  const pad = (n: number) => String(n).padStart(2, '0')

  return (
    `${pad(d.getDate())}.${pad(d.getMonth() + 1)}.${d.getFullYear()} ` +
    `${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`
  )
}
