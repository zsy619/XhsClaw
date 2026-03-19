import { http, type ApiResponse } from './request'

// 仪表盘统计数据接口
export interface DashboardStats {
  total_contents: number
  published_count: number
  draft_count: number
  pending_count: number
  failed_count: number
  today_contents: number
  today_published: number
  today_views: number
  weekly_trend: DailyStats[]
  status_distribution: StatusCount[]
  avg_generation_time: number
  success_rate: number
  total_tokens: number
  today_tokens: number
  total_cost: number
  last_activity_time: string
  updated_at: string
}

// 每日统计数据
export interface DailyStats {
  date: string
  contents: number
  published: number
  views: number
  tokens: number
}

// 状态数量
export interface StatusCount {
  status: number
  label: string
  count: number
  color: string
}

// 用户活动
export interface UserActivity {
  id: number
  user_id: number
  type: string
  title: string
  status: string
  time: string
  time_ago: string
}

// 内容趋势
export interface ContentTrend {
  date: string
  generate: number
  publish: number
}

// 仪表盘完整响应
export interface DashboardData {
  stats: DashboardStats
  activities: UserActivity[]
  trends: ContentTrend[]
}

// 获取仪表盘统计数据
export const getDashboardStats = () => {
  return http.get<DashboardStats>('/dashboard/stats')
}

// 获取完整仪表盘数据
export const getDashboardData = () => {
  return http.get<DashboardData>('/dashboard/data')
}

// 获取用户最近活动
export const getUserActivities = () => {
  return http.get<UserActivity[]>('/dashboard/activities')
}

// 获取内容趋势
export const getContentTrends = () => {
  return http.get<ContentTrend[]>('/dashboard/trends')
}