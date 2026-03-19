// 系统诊断和自我校验工具
import axios from 'axios'
import { ElMessage } from 'element-plus'

// 诊断结果接口
interface DiagnosticResult {
  success: boolean
  message: string
  details?: string[]
  suggestions?: string[]
}

// 校验项接口
interface CheckItem {
  name: string
  check: () => Promise<boolean>
  errorMessage: string
  fixSuggestion: string
}

/**
 * 系统诊断类
 */
export class SystemDiagnostic {
  private baseURL: string
  private checkResults: { [key: string]: boolean } = {}

  constructor(baseURL?: string) {
    this.baseURL = baseURL || import.meta.env.VITE_API_BASE_URL || 'http://localhost:8000/api/v1'
  }

  /**
   * 执行完整的系统诊断
   */
  async runFullDiagnostic(): Promise<DiagnosticResult> {
    const details: string[] = []
    const suggestions: string[] = []
    let allPassed = true

    console.log('🧪 开始系统诊断...')

    // 定义所有检查项
    const checks: CheckItem[] = [
      {
        name: '网络连接检查',
        check: () => this.checkNetworkConnection(),
        errorMessage: '无法连接到网络',
        fixSuggestion: '请检查您的网络连接'
      },
      {
        name: '后端服务器健康检查',
        check: () => this.checkBackendHealth(),
        errorMessage: '后端服务器未响应',
        fixSuggestion: '请确保后端服务器正在运行'
      },
      {
        name: 'API路由检查',
        check: () => this.checkAPIRoutes(),
        errorMessage: 'API路由访问失败',
        fixSuggestion: '请检查后端路由配置'
      },
      {
        name: 'CORS配置检查',
        check: () => this.checkCORS(),
        errorMessage: 'CORS配置有问题',
        fixSuggestion: '请检查后端CORS配置'
      }
    ]

    // 执行所有检查
    for (const check of checks) {
      try {
        console.log(`🔍 检查: ${check.name}...`)
        const passed = await check.check()
        this.checkResults[check.name] = passed
        
        if (passed) {
          details.push(`✓ ${check.name}: 通过`)
          console.log(`✅ ${check.name}: 通过`)
        } else {
          details.push(`✗ ${check.name}: 失败 - ${check.errorMessage}`)
          suggestions.push(check.fixSuggestion)
          console.log(`❌ ${check.name}: 失败 - ${check.errorMessage}`)
          allPassed = false
        }
      } catch (error) {
        details.push(`✗ ${check.name}: 异常 - ${error}`)
        suggestions.push(check.fixSuggestion)
        console.log(`💥 ${check.name}: 异常`, error)
        allPassed = false
      }
    }

    // 添加快速修复建议
    if (!allPassed) {
      suggestions.push('💡 快速修复：运行项目根目录的 start.sh 脚本启动服务')
      suggestions.push('💡 或者手动运行：cd backend && go run cmd/server/main.go')
    }

    return {
      success: allPassed,
      message: allPassed 
        ? '🎉 所有检查通过！系统运行正常。'
        : '⚠️ 部分检查失败，请查看详细信息。',
      details,
      suggestions
    }
  }

  /**
   * 检查网络连接
   */
  private async checkNetworkConnection(): Promise<boolean> {
    try {
      // 尝试连接一个公共地址来验证网络
      await axios.get('https://www.baidu.com', { timeout: 5000 })
      return true
    } catch {
      // 如果无法连接外网，至少检查localhost是否可达
      try {
        await axios.get('http://localhost:8000/health', { timeout: 3000 })
        return true
      } catch {
        return false
      }
    }
  }

  /**
   * 检查后端健康状态
   */
  private async checkBackendHealth(): Promise<boolean> {
    try {
      const healthURL = this.baseURL.replace('/api/v1', '/health')
      const response = await axios.get(healthURL, { timeout: 5000 })
      return response.status === 200
    } catch (error) {
      console.log('健康检查失败:', error)
      return false
    }
  }

  /**
   * 检查API路由
   */
  private async checkAPIRoutes(): Promise<boolean> {
    try {
      // 尝试访问一个不需要认证的公共API
      const response = await axios.get(`${this.baseURL}/health`, { timeout: 5000 })
      return response.status === 200
    } catch (error) {
      console.log('API路由检查失败:', error)
      return false
    }
  }

  /**
   * 检查CORS配置
   */
  private async checkCORS(): Promise<boolean> {
    try {
      // 尝试发送OPTIONS预检请求
      const response = await axios.options(this.baseURL.replace('/api/v1', '/health'), {
        timeout: 5000
      })
      // 检查是否有CORS相关的响应头
      const hasCorsHeaders = response.headers['access-control-allow-origin'] !== undefined
      return hasCorsHeaders || response.status === 200 || response.status === 204
    } catch (error) {
      // 如果OPTIONS请求失败，尝试GET请求
      try {
        await axios.get(this.baseURL.replace('/api/v1', '/health'), { timeout: 5000 })
        return true
      } catch {
        return false
      }
    }
  }

  /**
   * 获取单个检查结果
   */
  getCheckResult(checkName: string): boolean | undefined {
    return this.checkResults[checkName]
  }

  /**
   * 显示诊断结果
   */
  showDiagnosticResult(result: DiagnosticResult): void {
    if (result.success) {
      ElMessage.success(result.message)
    } else {
      ElMessage.warning(result.message)
    }

    // 在控制台输出详细信息
    console.log('📊 诊断结果:')
    result.details?.forEach(detail => console.log(detail))
    
    if (result.suggestions && result.suggestions.length > 0) {
      console.log('💡 建议:')
      result.suggestions.forEach(suggestion => console.log(suggestion))
    }
  }
}

/**
 * 快速诊断函数
 */
export async function quickDiagnostic(): Promise<DiagnosticResult> {
  const diagnostic = new SystemDiagnostic()
  const result = await diagnostic.runFullDiagnostic()
  diagnostic.showDiagnosticResult(result)
  return result
}

/**
 * 一键修复建议
 */
export function getQuickFixSteps(): string[] {
  return [
    '1. 打开终端，进入项目根目录',
    '2. 运行: ./start.sh (或 sh start.sh)',
    '3. 等待服务启动完成',
    '4. 刷新浏览器页面',
    '',
    '或者手动启动:',
    '1. 后端: cd backend && go run cmd/server/main.go',
    '2. 前端: cd frontend && npm run dev'
  ]
}

// 导出默认实例
export default new SystemDiagnostic()
