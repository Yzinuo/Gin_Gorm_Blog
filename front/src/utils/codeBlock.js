/**
 * 增强 Markdown 渲染后的代码块：
 * 添加仿 Mac/VS Code 终端顶栏（红黄绿三点）、语言标识及一键复制代码按钮
 */
export function enhanceCodeBlocks(container) {
  if (!container)
    return

  const preElements = container.querySelectorAll('pre')
  preElements.forEach((pre) => {
    // 避免重复包装
    if (pre.parentElement?.classList.contains('code-block-wrapper'))
      return

    const code = pre.querySelector('code')
    if (!code)
      return

    // 提取语言名称
    let lang = 'CODE'
    for (const cls of code.classList) {
      if (cls.startsWith('language-')) {
        lang = cls.replace('language-', '').toUpperCase()
        break
      }
    }

    // 创建最外层容器
    const wrapper = document.createElement('div')
    wrapper.className = 'code-block-wrapper'

    // 创建顶栏
    const header = document.createElement('div')
    header.className = 'code-block-header'
    header.innerHTML = `
      <div class="code-lang-tag">
        <span class="code-dot red"></span>
        <span class="code-dot yellow"></span>
        <span class="code-dot green"></span>
        <span class="code-lang-name">${lang}</span>
      </div>
      <button type="button" class="code-copy-btn" aria-label="复制代码">
        <span class="copy-icon">⎘</span>
        <span class="copy-text">复制</span>
      </button>
    `

    const copyBtn = header.querySelector('.code-copy-btn')
    copyBtn.addEventListener('click', async () => {
      try {
        await navigator.clipboard.writeText(code.innerText)
        copyBtn.classList.add('copied')
        const textEl = copyBtn.querySelector('.copy-text')
        const iconEl = copyBtn.querySelector('.copy-icon')
        if (textEl)
          textEl.textContent = '已复制'
        if (iconEl)
          iconEl.textContent = '✓'

        setTimeout(() => {
          copyBtn.classList.remove('copied')
          if (textEl)
            textEl.textContent = '复制'
          if (iconEl)
            iconEl.textContent = '⎘'
        }, 2000)
      }
      catch (err) {
        console.error('Failed to copy code', err)
      }
    })

    // 插入 wrapper 并将 pre 移入其中
    pre.parentNode.insertBefore(wrapper, pre)
    wrapper.appendChild(header)
    wrapper.appendChild(pre)
  })
}
