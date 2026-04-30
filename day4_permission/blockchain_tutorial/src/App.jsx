import { useState } from 'react'
import './App.css'

const BLOCKS = [
  { id: 0, hash: '0000000000000000', prevHash: 'Genesis', data: '创世区块', nonce: 0 },
  { id: 1, hash: 'a3f8c2d1e9b47065', prevHash: '0000000000000000', data: 'Alice → Bob: 5 BTC', nonce: 42831 },
  { id: 2, hash: 'b7e1a09c3f254d68', prevHash: 'a3f8c2d1e9b47065', data: 'Bob → Carol: 2 BTC', nonce: 19274 },
  { id: 3, hash: 'c2d4e5f0a1b38c97', prevHash: 'b7e1a09c3f254d68', data: 'Carol → Dave: 1 BTC', nonce: 83651 },
]

const CONCEPTS = [
  { icon: '🔗', title: '哈希函数', desc: '将任意数据转换为固定长度的唯一摘要，数据微小改动会导致哈希值完全不同。' },
  { icon: '⛏️', title: '工作量证明', desc: '矿工需要反复尝试不同的 Nonce 值，直到找到满足难度要求的哈希值。' },
  { icon: '📡', title: '去中心化网络', desc: '全网节点共同维护账本副本，无单一故障点，无法被单方面篡改。' },
  { icon: '🔐', title: '非对称加密', desc: '用私钥签名交易，任何人可用公钥验证签名，确保交易真实性。' },
]

function Block({ block, index }) {
  const [open, setOpen] = useState(false)
  const isGenesis = index === 0

  return (
    <div className={`block ${isGenesis ? 'genesis' : ''}`}>
      <div className="block-header" onClick={() => setOpen(o => !o)}>
        <span className="block-index">#{block.id}</span>
        <span className="block-hash">{block.hash}</span>
        <span className="block-toggle">{open ? '▲' : '▼'}</span>
      </div>
      {open && (
        <div className="block-body">
          <div className="block-row"><span>前块哈希</span><code>{block.prevHash}</code></div>
          <div className="block-row"><span>数据</span><code>{block.data}</code></div>
          <div className="block-row"><span>Nonce</span><code>{block.nonce}</code></div>
          <div className="block-row"><span>哈希</span><code className="accent">{block.hash}</code></div>
        </div>
      )}
    </div>
  )
}

function ChainView() {
  return (
    <section className="chain-section">
      <h2 className="section-title">区块链结构可视化</h2>
      <p className="section-desc">点击区块查看详细字段，每个区块通过前块哈希链接在一起。</p>
      <div className="chain">
        {BLOCKS.map((block, i) => (
          <div key={block.id} className="chain-item">
            <Block block={block} index={i} />
            {i < BLOCKS.length - 1 && <div className="chain-arrow">→</div>}
          </div>
        ))}
      </div>
    </section>
  )
}

function ConceptCard({ icon, title, desc }) {
  return (
    <div className="concept-card">
      <div className="concept-icon">{icon}</div>
      <h3>{title}</h3>
      <p>{desc}</p>
    </div>
  )
}

function ConceptsView() {
  return (
    <section className="concepts-section">
      <h2 className="section-title">核心概念</h2>
      <div className="concepts-grid">
        {CONCEPTS.map(c => <ConceptCard key={c.title} {...c} />)}
      </div>
    </section>
  )
}

const TABS = [
  { id: 'chain', label: '⛓ 区块链结构' },
  { id: 'concepts', label: '📚 核心概念' },
]

export default function App() {
  const [tab, setTab] = useState('chain')

  return (
    <div className="app">
      <header className="header">
        <div className="logo">⛓</div>
        <div>
          <h1>区块链教学平台</h1>
          <p className="subtitle">Blockchain Interactive Learning</p>
        </div>
      </header>

      <nav className="tabs">
        {TABS.map(t => (
          <button
            key={t.id}
            className={`tab-btn ${tab === t.id ? 'active' : ''}`}
            onClick={() => setTab(t.id)}
          >
            {t.label}
          </button>
        ))}
      </nav>

      <main className="main">
        {tab === 'chain' && <ChainView />}
        {tab === 'concepts' && <ConceptsView />}
      </main>

      <footer className="footer">
        <span>区块链教学演示 · Blockchain Demo · NTU 2026</span>
      </footer>
    </div>
  )
}
