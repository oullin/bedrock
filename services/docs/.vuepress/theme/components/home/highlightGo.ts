// Tiny Go syntax highlighter. Input is a static constant (see HERO_TABS),
// so regex order is chosen for clarity, not resilience.
export function highlightGo(src: string): string {
  const esc = (s: string) => s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')
  let t = esc(src)
  t = t.replace(/(\/\/[^\n]*)/g, '<span class="c">$1</span>')
  t = t.replace(/(&quot;|")([^&"\n]*?)(&quot;|")/g, '<span class="s">"$2"</span>')
  t = t.replace(/\b(package|import|func|return|if|else|for|range|var|const|type|struct|interface|map|chan|go|defer|switch|case|break|continue|nil|true|false)\b/g, '<span class="k">$1</span>')
  t = t.replace(/\b([A-Z][A-Za-z0-9]*)\b/g, '<span class="t">$1</span>')
  t = t.replace(/([a-zA-Z_][a-zA-Z0-9_]*)(\()/g, '<span class="fn">$1</span>$2')
  return t
}
