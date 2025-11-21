// Paste page logic

let currentPaste = null;
let isRawView = false;

// Initialize page
document.addEventListener('DOMContentLoaded', () => {
  // Configure marked.js
  marked.setOptions({
    highlight: function(code, lang) {
      if (lang && hljs.getLanguage(lang)) {
        try {
          return hljs.highlight(code, { language: lang }).value;
        } catch (err) {}
      }
      return hljs.highlightAuto(code).value;
    },
    breaks: true,
    gfm: true,
  });

  // Character counter
  const textarea = document.getElementById('pasteContent');
  const charCount = document.getElementById('charCount');
  
  textarea.addEventListener('input', () => {
    const count = textarea.value.length;
    charCount.textContent = count.toLocaleString();
    
    // Warn if approaching limit
    const maxSize = 10 * 1024 * 1024; // 10MB
    if (count > maxSize * 0.9) {
      charCount.style.color = 'var(--error)';
    } else if (count > maxSize * 0.7) {
      charCount.style.color = 'var(--accent)';
    } else {
      charCount.style.color = 'var(--text-muted)';
    }
  });

  // Check if viewing a paste (e.g., /paste.html?id=abc123)
  const urlParams = new URLSearchParams(window.location.search);
  const pasteId = urlParams.get('id');
  
  if (pasteId) {
    loadPaste(pasteId);
  }
});

// Handle paste creation
async function handleCreatePaste(event) {
  event.preventDefault();
  
  const submitBtn = document.getElementById('submitBtn');
  const form = document.getElementById('pasteForm');
  const formData = new FormData(form);
  
  // Prepare request data
  const content = formData.get('content');
  const isMarkdown = document.getElementById('isMarkdown').checked;
  const expiryDays = parseInt(formData.get('expiry_days'));
  
  // Validate
  if (!content || content.trim().length === 0) {
    app.showToast('Content cannot be empty', 'error');
    return;
  }
  
  if (content.length > 10 * 1024 * 1024) {
    app.showToast('Content exceeds 10MB limit', 'error');
    return;
  }
  
  // Show loading state
  submitBtn.disabled = true;
  submitBtn.innerHTML = '<span class="loading-spinner"></span> Creating...';
  
  try {
    const paste = await app.api.post('/paste', {
      content: content,
      is_markdown: isMarkdown,
      expiry_days: expiryDays === 0 ? null : expiryDays,
    });
    
    app.showToast('Paste created successfully!');
    
    // Redirect to view paste
    window.history.pushState({}, '', `/paste.html?id=${paste.id}`);
    currentPaste = paste;
    showPasteView(paste);
    
  } catch (error) {
    app.showToast(error.message || 'Failed to create paste', 'error');
  } finally {
    submitBtn.disabled = false;
    submitBtn.innerHTML = '<span>✨</span> Create Paste';
  }
}

// Load paste by ID
async function loadPaste(id) {
  const createView = document.getElementById('createView');
  const viewView = document.getElementById('viewView');
  
  try {
    // Show loading
    createView.style.display = 'none';
    viewView.style.display = 'block';
    document.getElementById('pasteContentRendered').innerHTML = '<div style="text-align: center; padding: 2rem; color: var(--text-muted);">Loading paste...</div>';
    
    const paste = await app.api.get(`/paste/${id}`);
    currentPaste = paste;
    showPasteView(paste);
    
  } catch (error) {
    viewView.style.display = 'none';
    createView.style.display = 'block';
    
    if (error.message.includes('not found')) {
      app.showToast('Paste not found', 'error');
    } else if (error.message.includes('expired')) {
      app.showToast('This paste has expired', 'error');
    } else {
      app.showToast('Failed to load paste', 'error');
    }
  }
}

// Show paste view
// Show paste view
function showPasteView(paste) {
  const createView = document.getElementById('createView');
  const viewView = document.getElementById('viewView');
  const renderedContent = document.getElementById('pasteContentRendered');
  const rawContent = document.getElementById('rawContent');
  const pasteMeta = document.getElementById('pasteMeta');
  
  console.log('Paste data:', paste); // Debug log
  
  // Hide create, show view
  createView.style.display = 'none';
  viewView.style.display = 'block';
  
  // Update title
  const title = paste.is_markdown ? 'Markdown Paste' : 'Text Paste';
  document.getElementById('pasteTitle').textContent = title;
  
  // Update metadata
  let created = 'Unknown';
  try {
    const createdDate = new Date(paste.created_at);
    if (!isNaN(createdDate.getTime())) {
      created = createdDate.toLocaleString();
    }
  } catch (e) {
    console.error('Error parsing created_at:', e);
  }
  
  let expiryText = 'Never expires';
  if (paste.expires_at) {
    try {
      const expiresDate = new Date(paste.expires_at);
      if (!isNaN(expiresDate.getTime())) {
        const now = new Date();
        const daysLeft = Math.ceil((expiresDate - now) / (1000 * 60 * 60 * 24));
        
        if (daysLeft > 1) {
          expiryText = `Expires in ${daysLeft} days`;
        } else if (daysLeft === 1) {
          expiryText = 'Expires in 1 day';
        } else if (daysLeft === 0) {
          const hoursLeft = Math.ceil((expiresDate - now) / (1000 * 60 * 60));
          expiryText = `Expires in ${hoursLeft} hour${hoursLeft === 1 ? '' : 's'}`;
        } else {
          expiryText = 'Expires soon';
        }
      }
    } catch (e) {
      console.error('Error parsing expires_at:', e);
    }
  }
  
  pasteMeta.innerHTML = `
    <span>📅 Created: ${created}</span>
    <span>⏰ ${expiryText}</span>
    <span>📝 ${paste.is_markdown ? 'Markdown' : 'Plain Text'}</span>
    <span>🆔 ID: ${paste.id || 'Unknown'}</span>
  `;
  
  // Render content
  if (paste.is_markdown) {
    try {
      renderedContent.innerHTML = marked.parse(paste.content);
      // Highlight code blocks
      renderedContent.querySelectorAll('pre code').forEach((block) => {
        hljs.highlightElement(block);
      });
    } catch (e) {
      console.error('Error rendering markdown:', e);
      renderedContent.innerHTML = `<pre style="margin: 0; white-space: pre-wrap; word-wrap: break-word;">${escapeHtml(paste.content)}</pre>`;
    }
  } else {
    renderedContent.innerHTML = `<pre style="margin: 0; white-space: pre-wrap; word-wrap: break-word;">${escapeHtml(paste.content)}</pre>`;
  }
  
  // Set raw content
  rawContent.textContent = paste.content;
  
  // Reset view state
  isRawView = false;
  document.getElementById('pasteContentRendered').style.display = 'block';
  document.getElementById('pasteContentRaw').style.display = 'none';
  document.getElementById('rawToggle').innerHTML = '<span>👁️</span> Raw';
}

// Toggle between rendered and raw view
function toggleRawView() {
  const renderedContent = document.getElementById('pasteContentRendered');
  const rawContent = document.getElementById('pasteContentRaw');
  const rawToggle = document.getElementById('rawToggle');
  
  isRawView = !isRawView;
  
  if (isRawView) {
    renderedContent.style.display = 'none';
    rawContent.style.display = 'block';
    rawToggle.innerHTML = '<span>🎨</span> Rendered';
  } else {
    renderedContent.style.display = 'block';
    rawContent.style.display = 'none';
    rawToggle.innerHTML = '<span>👁️</span> Raw';
  }
}

// Copy paste content to clipboard
async function copyPasteContent() {
  if (currentPaste) {
    await app.copyToClipboard(currentPaste.content);
  }
}

// Copy paste URL to clipboard
async function copyPasteURL() {
  if (currentPaste) {
    const url = `${window.location.origin}/paste.html?id=${currentPaste.id}`;
    await app.copyToClipboard(url);
    app.showToast('Paste URL copied!');
  }
}

// Download paste as file
function downloadPaste() {
  if (!currentPaste) return;
  
  const content = currentPaste.content;
  const filename = `paste-${currentPaste.id}.${currentPaste.is_markdown ? 'md' : 'txt'}`;
  
  const blob = new Blob([content], { type: 'text/plain' });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = filename;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  URL.revokeObjectURL(url);
  
  app.showToast('Paste downloaded!');
}

// Delete paste
async function deletePaste() {
  if (!currentPaste) return;
  
  if (!confirm('Are you sure you want to delete this paste? This action cannot be undone.')) {
    return;
  }
  
  try {
    // DELETE endpoint returns 204 No Content (no body)
    const result = await app.api.delete(`/paste/${currentPaste.id}`);
    
    // Success - result will be null for 204 responses
    app.showToast('Paste deleted successfully');
    
    // Redirect to create view
    window.history.pushState({}, '', '/paste.html');
    createNewPaste();
    
  } catch (error) {
    console.error('Delete error:', error);
    app.showToast(error.message || 'Failed to delete paste', 'error');
  }
}

// Create new paste (show create form)
function createNewPaste() {
  const createView = document.getElementById('createView');
  const viewView = document.getElementById('viewView');
  
  createView.style.display = 'block';
  viewView.style.display = 'none';
  currentPaste = null;
  
  // Clear form
  clearForm();
  
  // Update URL
  window.history.pushState({}, '', '/paste.html');
}

// Clear form
function clearForm() {
  document.getElementById('pasteForm').reset();
  document.getElementById('charCount').textContent = '0';
  document.getElementById('charCount').style.color = 'var(--text-muted)';
}

// Utility: Escape HTML
function escapeHtml(text) {
  const map = {
    '&': '&amp;',
    '<': '&lt;',
    '>': '&gt;',
    '"': '&quot;',
    "'": '&#039;'
  };
  return text.replace(/[&<>"']/g, m => map[m]);
}

// Handle browser back/forward
window.addEventListener('popstate', () => {
  const urlParams = new URLSearchParams(window.location.search);
  const pasteId = urlParams.get('id');
  
  if (pasteId) {
    loadPaste(pasteId);
  } else {
    createNewPaste();
  }
});