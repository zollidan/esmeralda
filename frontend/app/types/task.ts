export type Task = {
  id: string
  date: string
  status: string
  created_at: string
}

export type ProgressBar = {
  task_id: string
  percent: number
  message: string
}

export type SportRuAPIStatus = {
  status: string
  latency_ms: number
}
