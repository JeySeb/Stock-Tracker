function escapeHtml(input: string): string {
  return input
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#039;')
}

export function formatWhatsAppText(rawText: string): string {
  if (!rawText) return ''

  // 1) Extract code blocks (``` ... ```) before escaping
  const codeBlocks: string[] = []
  let text = rawText.replace(/```([\s\S]*?)```/g, (_match, codeContent: string) => {
    const escaped = escapeHtml(codeContent.trim())
    const html = `<pre><code>${escaped}</code></pre>`
    const idx = codeBlocks.push(html) - 1
    return `__CODEBLOCK_${idx}__`
  })

  // 2) Escape everything to prevent HTML injection
  text = escapeHtml(text)

  // 3) Extract inline code spans (`...`) as placeholders (now working on escaped content)
  const codeSpans: string[] = []
  text = text.replace(/`([^`]+?)`/g, (_match, codeInline: string) => {
    const html = `<code>${codeInline}</code>`
    const idx = codeSpans.push(html) - 1
    return `__CODESPAN_${idx}__`
  })

  // 4) Process block structures: quotes, lists
  const lines = text.split(/\r?\n/)
  const htmlParts: string[] = []
  let i = 0

  while (i < lines.length) {
    const line = lines[i]

    // Blockquote lines: start with '&gt; '
    if (/^&gt;\s+/.test(line)) {
      const quoteLines: string[] = []
      while (i < lines.length && /^&gt;\s+/.test(lines[i])) {
        quoteLines.push(lines[i].replace(/^&gt;\s+/, ''))
        i++
      }
      htmlParts.push(`<blockquote>${quoteLines.join('<br>')}</blockquote>`)
      continue
    }

    // Unordered list: lines starting with '* ' or '- '
    if (/^(\*|-)\s+/.test(line)) {
      const items: string[] = []
      // capture items allowing blank lines between bullets
      while (i < lines.length) {
        if (/^(\*|-)\s+/.test(lines[i])) {
          items.push(lines[i].replace(/^(\*|-)\s+/, ''))
          i++
          // skip subsequent blank lines inside the list
          while (i < lines.length && /^\s*$/.test(lines[i])) i++
          continue
        }
        // allow blank lines before the next bullet if the next non-blank is a bullet
        if (/^\s*$/.test(lines[i])) {
          let j = i + 1
          while (j < lines.length && /^\s*$/.test(lines[j])) j++
          if (j < lines.length && /^(\*|-)\s+/.test(lines[j])) {
            i = j
            continue
          }
        }
        break
      }
      {
        const listHtml = '<ul>' + items.map((it) => '<li>' + it + '</li>').join('') + '</ul>'
        htmlParts.push(listHtml)
      }
      continue
    }

    // Ordered list: lines starting with '1. ', '2. ', etc. (allow blank lines between items)
    if (/^\d+\.\s+/.test(line)) {
      const items: string[] = []
      // detect starting number from the first matched line
      const firstMatch = line.match(/^(\d+)\.\s+/)
      const startNum = firstMatch ? Number(firstMatch[1]) : 1

      while (i < lines.length) {
        if (/^\d+\.\s+/.test(lines[i])) {
          items.push(lines[i].replace(/^\d+\.\s+/, ''))
          i++
          // skip subsequent blank lines inside the list
          while (i < lines.length && /^\s*$/.test(lines[i])) i++
          continue
        }
        // allow blank lines before the next numeric line if the next non-blank is numeric
        if (/^\s*$/.test(lines[i])) {
          let j = i + 1
          while (j < lines.length && /^\s*$/.test(lines[j])) j++
          if (j < lines.length && /^\d+\.\s+/.test(lines[j])) {
            i = j
            continue
          }
        }
        break
      }
      {
        const startAttr = startNum !== 1 ? ' start="' + String(startNum) + '"' : ''
        const listHtml = '<ol' + startAttr + '>' + items.map((it) => '<li>' + it + '</li>').join('') + '</ol>'
        htmlParts.push(listHtml)
      }
      continue
    }

    // Normal line → keep; we'll add <br> where appropriate
    htmlParts.push(line)
    i++
  }

  let html = htmlParts.join('<br>')

  // 5) Inline formatting (avoid matching list bullets: ensure not followed by space)
  // Bold: *text*
  html = html.replace(/\*(?!\s)([\s\S]*?[^\s])\*/g, '<strong>$1</strong>')

  // Italic: _text_
  html = html.replace(/_(?!\s)([\s\S]*?[^\s])_/g, '<em>$1</em>')

  // Strikethrough: ~text~
  html = html.replace(/~(?!\s)([\s\S]*?[^\s])~/g, '<del>$1</del>')

  // 6) Restore code spans and code blocks
  html = html.replace(/__CODESPAN_(\d+)__/g, (_m, idx: string) => codeSpans[Number(idx)])
  html = html.replace(/__CODEBLOCK_(\d+)__/g, (_m, idx: string) => codeBlocks[Number(idx)])

  return html
}

export default formatWhatsAppText

