<template>
  <div class="markdown-editor">
    <!-- Tab Navigation -->
    <ul class="nav nav-tabs" role="tablist">
      <li class="nav-item" role="presentation">
        <button
          class="nav-link"
          :class="{ active: activeTab === 'write' }"
          @click="activeTab = 'write'"
          type="button"
        >
          <i class="bi bi-pencil-square me-1"></i>
          Write
        </button>
      </li>
      <li class="nav-item" role="presentation">
        <button
          class="nav-link"
          :class="{ active: activeTab === 'preview' }"
          @click="activeTab = 'preview'"
          type="button"
        >
          <i class="bi bi-eye me-1"></i>
          Preview
        </button>
      </li>
    </ul>

    <!-- Tab Content -->
    <div class="tab-content markdown-editor-content">
      <!-- Write Tab -->
      <div
        v-show="activeTab === 'write'"
        class="tab-pane fade show active"
      >
        <div class="markdown-toolbar">
          <button
            type="button"
            class="btn btn-sm btn-outline-secondary me-1"
            @click="insertMarkdown('**', '**', 'bold text')"
            title="Bold"
          >
            <i class="bi bi-type-bold"></i>
          </button>
          <button
            type="button"
            class="btn btn-sm btn-outline-secondary me-1"
            @click="insertMarkdown('*', '*', 'italic text')"
            title="Italic"
          >
            <i class="bi bi-type-italic"></i>
          </button>
          <button
            type="button"
            class="btn btn-sm btn-outline-secondary me-1"
            @click="insertMarkdown('`', '`', 'code')"
            title="Inline Code"
          >
            <i class="bi bi-code"></i>
          </button>
          <button
            type="button"
            class="btn btn-sm btn-outline-secondary me-1"
            @click="insertMarkdown('[', '](url)', 'link text')"
            title="Link"
          >
            <i class="bi bi-link-45deg"></i>
          </button>
          <button
            type="button"
            class="btn btn-sm btn-outline-secondary me-1"
            @click="insertMarkdown('![', '](image-url)', 'alt text')"
            title="Image"
          >
            <i class="bi bi-image"></i>
          </button>
          <div class="btn-group me-1" role="group">
            <button
              type="button"
              class="btn btn-sm btn-outline-secondary"
              @click="insertHeading(1)"
              title="Heading 1"
            >
              H1
            </button>
            <button
              type="button"
              class="btn btn-sm btn-outline-secondary"
              @click="insertHeading(2)"
              title="Heading 2"
            >
              H2
            </button>
            <button
              type="button"
              class="btn btn-sm btn-outline-secondary"
              @click="insertHeading(3)"
              title="Heading 3"
            >
              H3
            </button>
          </div>
          <button
            type="button"
            class="btn btn-sm btn-outline-secondary me-1"
            @click="insertList('-')"
            title="Bullet List"
          >
            <i class="bi bi-list-ul"></i>
          </button>
          <button
            type="button"
            class="btn btn-sm btn-outline-secondary me-1"
            @click="insertList('1.')"
            title="Numbered List"
          >
            <i class="bi bi-list-ol"></i>
          </button>
          <button
            type="button"
            class="btn btn-sm btn-outline-secondary me-1"
            @click="insertBlockquote()"
            title="Quote"
          >
            <i class="bi bi-quote"></i>
          </button>
        </div>
        
        <textarea
          ref="textareaRef"
          v-model="localValue"
          :class="typeof textareaClass === 'string' ? textareaClass : textareaClass"
          :rows="rows"
          :placeholder="placeholder"
          @input="$emit('update:modelValue', ($event.target as HTMLTextAreaElement)?.value || '')"
        ></textarea>
      </div>

      <!-- Preview Tab -->
      <div
        v-show="activeTab === 'preview'"
        class="tab-pane fade show active"
      >
        <div 
          class="markdown-preview"
          :style="{ minHeight: previewHeight }"
          v-html="renderedMarkdown"
        ></div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, nextTick, onMounted } from 'vue'
import { marked } from 'marked'
import hljs from 'highlight.js/lib/core'
import javascript from 'highlight.js/lib/languages/javascript'
import typescript from 'highlight.js/lib/languages/typescript'
import python from 'highlight.js/lib/languages/python'
import java from 'highlight.js/lib/languages/java'
import go from 'highlight.js/lib/languages/go'
import cpp from 'highlight.js/lib/languages/cpp'
import sql from 'highlight.js/lib/languages/sql'
import json from 'highlight.js/lib/languages/json'
import xml from 'highlight.js/lib/languages/xml'
import css from 'highlight.js/lib/languages/css'
import bash from 'highlight.js/lib/languages/bash'

// Register languages
hljs.registerLanguage('javascript', javascript)
hljs.registerLanguage('typescript', typescript)
hljs.registerLanguage('python', python)
hljs.registerLanguage('java', java)
hljs.registerLanguage('go', go)
hljs.registerLanguage('cpp', cpp)
hljs.registerLanguage('sql', sql)
hljs.registerLanguage('json', json)
hljs.registerLanguage('xml', xml)
hljs.registerLanguage('css', css)
hljs.registerLanguage('bash', bash)

interface Props {
  modelValue: string
  rows?: number
  placeholder?: string
  textareaClass?: string | Record<string, boolean>
}

const props = withDefaults(defineProps<Props>(), {
  rows: 12,
  placeholder: 'Write your content in Markdown...',
  textareaClass: 'form-control'
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const activeTab = ref<'write' | 'preview'>('write')
const textareaRef = ref<HTMLTextAreaElement>()

const localValue = computed({
  get: () => props.modelValue,
  set: (value: string) => emit('update:modelValue', value)
})

const previewHeight = computed(() => {
  return `${props.rows * 1.5}em`
})

// Configure marked with simple options
marked.setOptions({
  breaks: true,
  gfm: true
})

// Helper function to highlight code blocks
function highlightCode(code: string, lang?: string): string {
  if (lang && hljs.getLanguage(lang)) {
    try {
      return hljs.highlight(code, { language: lang }).value
    } catch (err) {
      console.warn('Highlight.js error:', err)
    }
  }
  return hljs.highlightAuto(code).value
}

const renderedMarkdown = computed(() => {
  if (!props.modelValue.trim()) {
    return '<p class="text-muted"><em>Nothing to preview yet. Write some markdown content!</em></p>'
  }
  
  try {
    const html = marked.parse(props.modelValue) as string
    
    // Post-process to add syntax highlighting to code blocks
    return html.replace(/<pre><code class="language-(\w+)">([\s\S]*?)<\/code><\/pre>/g, (_match: any, lang: any, code: any) => {
      const decodedCode = code.replace(/&lt;/g, '<').replace(/&gt;/g, '>').replace(/&amp;/g, '&')
      const highlightedCode = highlightCode(decodedCode, lang)
      return `<pre><code class="hljs language-${lang}">${highlightedCode}</code></pre>`
    }).replace(/<pre><code>([\s\S]*?)<\/code><\/pre>/g, (_match: any, code: any) => {
      const decodedCode = code.replace(/&lt;/g, '<').replace(/&gt;/g, '>').replace(/&amp;/g, '&')
      const highlightedCode = highlightCode(decodedCode)
      return `<pre><code class="hljs">${highlightedCode}</code></pre>`
    })
  } catch (error) {
    console.error('Markdown parsing error:', error)
    return '<p class="text-danger"><em>Error parsing markdown content</em></p>'
  }
})

// Toolbar functions
function insertMarkdown(before: string, after: string, placeholder: string) {
  if (!textareaRef.value) return
  
  const textarea = textareaRef.value
  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  const selectedText = textarea.value.substring(start, end)
  const textToInsert = selectedText || placeholder
  const newText = before + textToInsert + after
  
  const newValue = 
    textarea.value.substring(0, start) + 
    newText + 
    textarea.value.substring(end)
  
  emit('update:modelValue', newValue)
  
  nextTick(() => {
    textarea.focus()
    if (selectedText) {
      textarea.setSelectionRange(start, start + newText.length)
    } else {
      textarea.setSelectionRange(start + before.length, start + before.length + textToInsert.length)
    }
  })
}

function insertHeading(level: number) {
  if (!textareaRef.value) return
  
  const textarea = textareaRef.value
  const start = textarea.selectionStart
  const lineStart = textarea.value.lastIndexOf('\n', start - 1) + 1
  const prefix = '#'.repeat(level) + ' '
  
  const newValue = 
    textarea.value.substring(0, lineStart) + 
    prefix + 
    textarea.value.substring(lineStart)
  
  emit('update:modelValue', newValue)
  
  nextTick(() => {
    textarea.focus()
    textarea.setSelectionRange(start + prefix.length, start + prefix.length)
  })
}

function insertList(marker: string) {
  if (!textareaRef.value) return
  
  const textarea = textareaRef.value
  const start = textarea.selectionStart
  const lineStart = textarea.value.lastIndexOf('\n', start - 1) + 1
  const prefix = marker + ' '
  
  const newValue = 
    textarea.value.substring(0, lineStart) + 
    prefix + 
    textarea.value.substring(lineStart)
  
  emit('update:modelValue', newValue)
  
  nextTick(() => {
    textarea.focus()
    textarea.setSelectionRange(start + prefix.length, start + prefix.length)
  })
}

function insertBlockquote() {
  if (!textareaRef.value) return
  
  const textarea = textareaRef.value
  const start = textarea.selectionStart
  const lineStart = textarea.value.lastIndexOf('\n', start - 1) + 1
  const prefix = '> '
  
  const newValue = 
    textarea.value.substring(0, lineStart) + 
    prefix + 
    textarea.value.substring(lineStart)
  
  emit('update:modelValue', newValue)
  
  nextTick(() => {
    textarea.focus()
    textarea.setSelectionRange(start + prefix.length, start + prefix.length)
  })
}
</script>

<style scoped>
.markdown-editor {
  border: 1px solid #dee2e6;
  border-radius: 0.375rem;
  overflow: hidden;
}

.nav-tabs {
  border-bottom: 1px solid #dee2e6;
  background-color: #f8f9fa;
  margin-bottom: 0;
}

.nav-tabs .nav-link {
  border: none;
  border-radius: 0;
  color: #6c757d;
  padding: 0.75rem 1rem;
  font-size: 0.875rem;
}

.nav-tabs .nav-link:hover {
  border: none;
  background-color: #e9ecef;
}

.nav-tabs .nav-link.active {
  background-color: #fff;
  color: #0d6efd;
  border: none;
  border-bottom: 2px solid #0d6efd;
}

.markdown-editor-content {
  background-color: #fff;
}

.markdown-toolbar {
  padding: 0.5rem;
  border-bottom: 1px solid #e9ecef;
  background-color: #f8f9fa;
  display: flex;
  flex-wrap: wrap;
  gap: 0.25rem;
}

.markdown-toolbar .btn {
  padding: 0.25rem 0.5rem;
  font-size: 0.75rem;
  line-height: 1.5;
}

.markdown-editor textarea {
  border: none;
  border-radius: 0;
  resize: vertical;
  font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, Courier, monospace;
  font-size: 0.875rem;
  line-height: 1.5;
  padding: 1rem;
}

.markdown-editor textarea:focus {
  box-shadow: none;
  border: none;
}

.markdown-preview {
  padding: 1rem;
  overflow-y: auto;
}

/* Markdown content styling */
.markdown-preview :deep(h1),
.markdown-preview :deep(h2),
.markdown-preview :deep(h3),
.markdown-preview :deep(h4),
.markdown-preview :deep(h5),
.markdown-preview :deep(h6) {
  margin-top: 1.5rem;
  margin-bottom: 0.5rem;
  font-weight: 600;
}

.markdown-preview :deep(h1) { font-size: 2rem; }
.markdown-preview :deep(h2) { font-size: 1.5rem; }
.markdown-preview :deep(h3) { font-size: 1.25rem; }
.markdown-preview :deep(h4) { font-size: 1.1rem; }
.markdown-preview :deep(h5) { font-size: 1rem; }
.markdown-preview :deep(h6) { font-size: 0.875rem; }

.markdown-preview :deep(p) {
  margin-bottom: 1rem;
  line-height: 1.6;
}

.markdown-preview :deep(ul),
.markdown-preview :deep(ol) {
  margin-bottom: 1rem;
  padding-left: 2rem;
}

.markdown-preview :deep(li) {
  margin-bottom: 0.25rem;
}

.markdown-preview :deep(blockquote) {
  margin: 1rem 0;
  padding: 0.5rem 1rem;
  border-left: 4px solid #dee2e6;
  background-color: #f8f9fa;
  color: #6c757d;
}

.markdown-preview :deep(code) {
  background-color: #f8f9fa;
  padding: 0.125rem 0.25rem;
  border-radius: 0.25rem;
  font-size: 0.875em;
  color: #d63384;
}

.markdown-preview :deep(pre) {
  background-color: #f8f9fa;
  padding: 1rem;
  border-radius: 0.375rem;
  overflow-x: auto;
  margin-bottom: 1rem;
}

.markdown-preview :deep(pre code) {
  background-color: transparent;
  padding: 0;
  color: inherit;
}

.markdown-preview :deep(table) {
  width: 100%;
  margin-bottom: 1rem;
  border-collapse: collapse;
}

.markdown-preview :deep(th),
.markdown-preview :deep(td) {
  padding: 0.5rem;
  border: 1px solid #dee2e6;
}

.markdown-preview :deep(th) {
  background-color: #f8f9fa;
  font-weight: 600;
}

.markdown-preview :deep(img) {
  max-width: 100%;
  height: auto;
  border-radius: 0.375rem;
}

.markdown-preview :deep(a) {
  color: #0d6efd;
  text-decoration: none;
}

.markdown-preview :deep(a:hover) {
  text-decoration: underline;
}

.markdown-preview :deep(hr) {
  margin: 2rem 0;
  border: 0;
  border-top: 1px solid #dee2e6;
}
</style>
