export const meta = {
    name: 'todo-scanner',
    description: '扫描指定目录下所有 Python 文件的 TODO 注释，并生成汇总报告',
    phases: [
      { title: 'Scan',    detail: '并行扫描每个 Python 文件的 TODO 注释' },
      { title: 'Report',  detail: '汇总所有发现，生成优先级报告' },
    ],
  }
  
  const TARGET_DIR = (args && args.dir) || '.'
  
  const FINDING_SCHEMA = {
    type: 'object',
    required: ['file', 'todos'],
    properties: {
      file: { type: 'string' },
      todos: {
        type: 'array',
        items: {
          type: 'object',
          required: ['line', 'text', 'priority'],
          properties: {
            line:     { type: 'number' },
            text:     { type: 'string' },
            priority: { type: 'string', enum: ['high', 'medium', 'low'] },
          },
        },
      },
    },
  }
  
  // ── Phase 1: 并行扫描 ──────────────────────────────────────────
  phase('Scan')
  log(`扫描目录：${TARGET_DIR}`)
  
  // 先用一个 agent 列出所有 .py 文件
  const fileList = await agent(
    `列出 ${TARGET_DIR} 目录下所有 .py 文件的绝对路径（递归），` +
    `每行一个，不要其他任何内容。`,
    { label: 'list-files', phase: 'Scan' }
  )
  
  const files = fileList.trim().split('\n').filter(Boolean)
  log(`找到 ${files.length} 个 Python 文件，开始并行扫描`)
  
  // 对每个文件并行启动一个 agent
  const findings = await parallel(
    files.map((f) => () =>
      agent(
        `读取文件 ${f}，找出所有 TODO 注释（# TODO: 或 # todo:）。` +
        `对每条 TODO 判断优先级：含 FIXME/urgent/critical 的为 high，` +
        `含 improve/optimize 的为 medium，其余为 low。` +
        `若文件没有任何 TODO，返回 todos 字段为空数组。`,
        { label: `scan:${f.split('/').pop()}`, phase: 'Scan', schema: FINDING_SCHEMA }
      )
    )
  )
  
  // ── Phase 2: 汇总报告 ─────────────────────────────────────────
  phase('Report')
  
  const nonEmpty = findings.filter((f) => f && f.todos && f.todos.length > 0)
  log(`共发现 ${nonEmpty.reduce((s, f) => s + f.todos.length, 0)} 条 TODO`)
  
  const report = await agent(
    `以下是从 Python 项目中扫描到的所有 TODO 注释（JSON）：\n` +
    JSON.stringify(nonEmpty, null, 2) +
    `\n\n请生成一份中文 Markdown 报告，结构如下：\n` +
    `1. 总览（high/medium/low 各多少条）\n` +
    `2. 按优先级分组列出每条 TODO（文件名 + 行号 + 内容）\n` +
    `3. 建议优先处理的前 3 条，附理由`,
    { label: 'generate-report', phase: 'Report' }
  )
  
  return { total: findings.reduce((s, f) => s + (f?.todos?.length || 0), 0), report }