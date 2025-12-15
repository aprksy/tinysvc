// JSON Bin page logic

let currentJSON = null;
let currentView = 'formatted'; // 'formatted' or 'raw'
let currentExample = 'curl';

// Initialize page
document.addEventListener('DOMContentLoaded', () => {
  // Configure highlight.js
  hljs.highlightAll();

  // Character counter
  const textarea = document.getElementById('jsonContent');
  const charCount = document.getElementById('charCount');
  
  if (textarea) {
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
  }

  // Check if viewing a JSON bin (e.g., /json.html?id=abc123)
  const urlParams = new URLSearchParams(window.location.search);
  const id = urlParams.get('id');
  
  if (id) {
    loadJSON(id);
  }
});

// Handle JSON bin creation
async function handleCreateJSON(event) {
  event.preventDefault();
  
  const submitBtn = document.getElementById('submitBtn');
  const form = document.getElementById('jsonForm');
  const formData = new FormData(form);
  
  // Prepare request data
  const content = formData.get('content');
  const name = formData.get('name');
  const expiryDays = parseInt(formData.get('expiry_days'));
  
  // Validate JSON
  try {
    JSON.parse(content);
  } catch (e) {
    showValidation('Invalid JSON: ' + e.message, false);
    return;
  }
  
  // Show loading state
  submitBtn.disabled = true;
  submitBtn.innerHTML = '<span class="loading-spinner"></span> Creating...';
  
  try {
    const jsonBin = await app.api.post('/json', {
      content: JSON.parse(content),
      name: name || undefined,
      is_public: true,
      expiry_days: expiryDays === 0 ? null : expiryDays,
    });
    
    app.showToast('JSON bin created successfully!');
    
    // Redirect to view
    window.history.pushState({}, '', `/json.html?id=${jsonBin.id}`);
    currentJSON = jsonBin;
    showJSONView(jsonBin);
    
  } catch (error) {
    app.showToast(error.message || 'Failed to create JSON bin', 'error');
  } finally {
    submitBtn.disabled = false;
    submitBtn.innerHTML = '<span>✨</span> Create JSON Bin';
  }
}

// Load JSON bin by ID
async function loadJSON(id) {
  const createView = document.getElementById('createView');
  const viewView = document.getElementById('viewView');
  
  try {
    createView.style.display = 'none';
    viewView.style.display = 'block';
    
    const jsonBin = await app.api.get(`/json/${id}`);
    currentJSON = jsonBin;
    showJSONView(jsonBin);
    
  } catch (error) {
    viewView.style.display = 'none';
    createView.style.display = 'block';
    
    if (error.message.includes('not found')) {
      app.showToast('JSON bin not found', 'error');
    } else if (error.message.includes('expired')) {
      app.showToast('This JSON bin has expired', 'error');
    } else {
      app.showToast('Failed to load JSON bin', 'error');
    }
  }
}

// Show JSON bin view
function showJSONView(jsonBin) {
  console.log('showJSONView called with:', jsonBin);
  
  const createView = document.getElementById('createView');
  const viewView = document.getElementById('viewView');
  
  createView.style.display = 'none';
  viewView.style.display = 'block';
  
  setTimeout(() => {
    // Update title
    const title = jsonBin.name || 'JSON Bin';
    document.getElementById('jsonTitle').textContent = title;
    
    // Update metadata
    let created = 'Unknown';
    try {
      const createdDate = new Date(jsonBin.created_at);
      if (!isNaN(createdDate.getTime())) {
        created = createdDate.toLocaleString();
      }
    } catch (e) {
      console.error('Error parsing created_at:', e);
    }
    
    let expiryText = 'Never expires';
    if (jsonBin.expires_at) {
      try {
        const expiresDate = new Date(jsonBin.expires_at);
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
    
    const metadataHTML = `
  <div style="display:flex; flex-direction:row; align-items:center; gap:1.5rem">
    <div style="display:flex; flex-direction:row; align-items:center; gap:0.5rem;"><span class="material-symbols-outlined btn-sized">calendar_clock</span><span>Created: ${created}</span></div>
    <div style="display:flex; flex-direction:row; align-items:center; gap:0.5rem;"><span class="material-symbols-outlined btn-sized">free_cancellation</span><span>${expiryText}</span></div>
    <div style="display:flex; flex-direction:row; align-items:center; gap:0.5rem;"><span class="material-symbols-outlined btn-sized">visibility</span><span>${jsonBin.views}</span></div>
    <div style="display:flex; flex-direction:row; align-items:center; gap:0.5rem;"><span class="material-symbols-outlined btn-sized">id_card</span><span>ID: ${jsonBin.id || 'Unknown'}</span></div>
  </div>`;
    
    const jsonMetaElement = document.getElementById('jsonMeta');
    if (jsonMetaElement) {
      jsonMetaElement.innerHTML = metadataHTML;
    }
    
    // Display JSON with syntax highlighting
    const jsonCode = document.getElementById('jsonCode');
    const formattedJSON = JSON.stringify(jsonBin.content, null, 2);
    jsonCode.textContent = formattedJSON;
    hljs.highlightElement(jsonCode);
    
    // Set raw URL
    const rawURL = `${window.location.origin}/api/v1/json/${jsonBin.id}/raw`;
    document.getElementById('rawURL').value = rawURL;
    
    // Update code examples
    updateCodeExamples();
    
  }, 100);
}

// Format JSON in textarea
function formatJSON() {
  const textarea = document.getElementById('jsonContent');
  try {
    const parsed = JSON.parse(textarea.value);
    textarea.value = JSON.stringify(parsed, null, 2);
    showValidation('✓ JSON formatted successfully', true);
  } catch (e) {
    showValidation('Invalid JSON: ' + e.message, false);
  }
}

// Minify JSON in textarea
function minifyJSON() {
  const textarea = document.getElementById('jsonContent');
  try {
    const parsed = JSON.parse(textarea.value);
    textarea.value = JSON.stringify(parsed);
    showValidation('✓ JSON minified successfully', true);
  } catch (e) {
    showValidation('Invalid JSON: ' + e.message, false);
  }
}

// Validate JSON in textarea
function validateJSON() {
  const textarea = document.getElementById('jsonContent');
  try {
    JSON.parse(textarea.value);
    showValidation('✓ Valid JSON', true);
  } catch (e) {
    showValidation('Invalid JSON: ' + e.message, false);
  }
}

// Show validation message
function showValidation(message, isSuccess) {
  const validationMsg = document.getElementById('validationMessage');
  validationMsg.textContent = message;
  validationMsg.className = 'validation-message ' + (isSuccess ? 'success' : 'error');
  
  if (isSuccess) {
    setTimeout(() => {
      validationMsg.style.display = 'none';
    }, 3000);
  }
}

// Toggle between formatted and raw view
function toggleView() {
  const viewToggleText = document.getElementById('viewToggleText');
  const jsonCode = document.getElementById('jsonCode');
  
  if (currentView === 'formatted') {
    // Show minified
    currentView = 'raw';
    viewToggleText.textContent = 'Formatted';
    const minified = JSON.stringify(currentJSON.content);
    jsonCode.textContent = minified;
  } else {
    // Show formatted
    currentView = 'formatted';
    viewToggleText.textContent = 'Raw';
    const formatted = JSON.stringify(currentJSON.content, null, 2);
    jsonCode.textContent = formatted;
  }
  
  hljs.highlightElement(jsonCode);
}

// Copy JSON content to clipboard
async function copyJSONContent() {
  if (currentJSON) {
    const content = JSON.stringify(currentJSON.content, null, 2);
    await app.copyToClipboard(content);
  }
}

// Copy raw URL to clipboard
async function copyRawURL() {
  const rawURL = document.getElementById('rawURL').value;
  await app.copyToClipboard(rawURL);
}

// Copy JSON bin URL to clipboard
async function copyJSONURL() {
  if (currentJSON) {
    const url = `${window.location.origin}/json.html?id=${currentJSON.id}`;
    await app.copyToClipboard(url);
  }
}

// Download JSON as file
function downloadJSON() {
  if (!currentJSON) return;
  
  const content = JSON.stringify(currentJSON.content, null, 2);
  const filename = currentJSON.name ? 
    `${currentJSON.name.replace(/[^a-z0-9]/gi, '_')}.json` : 
    `json-${currentJSON.id}.json`;
  
  const blob = new Blob([content], { type: 'application/json' });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = filename;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  URL.revokeObjectURL(url);
  
  app.showToast('JSON downloaded!');
}

// Show update modal
function showUpdateModal() {
  if (!currentJSON) return;
  
  const modal = document.getElementById('updateModal');
  const updateContent = document.getElementById('updateContent');
  
  updateContent.value = JSON.stringify(currentJSON.content, null, 2);
  modal.classList.add('show');
  document.body.style.overflow = 'hidden';
}

// Close update modal
function closeUpdateModal() {
  const modal = document.getElementById('updateModal');
  modal.classList.remove('show');
  document.body.style.overflow = '';
}

// Format JSON in update modal
function formatUpdateJSON() {
  const textarea = document.getElementById('updateContent');
  const validationMsg = document.getElementById('updateValidationMessage');
  
  try {
    const parsed = JSON.parse(textarea.value);
    textarea.value = JSON.stringify(parsed, null, 2);
    validationMsg.textContent = '✓ JSON formatted successfully';
    validationMsg.className = 'validation-message success';
    setTimeout(() => validationMsg.style.display = 'none', 3000);
  } catch (e) {
    validationMsg.textContent = 'Invalid JSON: ' + e.message;
    validationMsg.className = 'validation-message error';
  }
}

// Validate JSON in update modal
function validateUpdateJSON() {
  const textarea = document.getElementById('updateContent');
  const validationMsg = document.getElementById('updateValidationMessage');
  
  try {
    JSON.parse(textarea.value);
    validationMsg.textContent = '✓ Valid JSON';
    validationMsg.className = 'validation-message success';
    setTimeout(() => validationMsg.style.display = 'none', 3000);
  } catch (e) {
    validationMsg.textContent = 'Invalid JSON: ' + e.message;
    validationMsg.className = 'validation-message error';
  }
}

// Save update
async function saveUpdate() {
  if (!currentJSON) return;
  
  const updateContent = document.getElementById('updateContent');
  
  try {
    const parsed = JSON.parse(updateContent.value);
    
    const updated = await app.api.request(`/json/${currentJSON.id}`, {
      method: 'PUT',
      body: JSON.stringify({ content: parsed }),
    });
    
    app.showToast('JSON bin updated successfully!');
    closeUpdateModal();
    
    currentJSON = updated;
    showJSONView(updated);
    
  } catch (error) {
    const validationMsg = document.getElementById('updateValidationMessage');
    validationMsg.textContent = error.message || 'Failed to update';
    validationMsg.className = 'validation-message error';
  }
}

// Delete JSON bin
async function deleteJSON() {
  if (!currentJSON) return;
  
  if (!confirm('Are you sure you want to delete this JSON bin? This action cannot be undone.')) {
    return;
  }
  
  try {
    await app.api.delete(`/json/${currentJSON.id}`);
    app.showToast('JSON bin deleted successfully');
    
    window.history.pushState({}, '', '/json.html');
    createNewJSON();
    
  } catch (error) {
    console.error('Delete error:', error);
    app.showToast(error.message || 'Failed to delete JSON bin', 'error');
  }
}

// Create new JSON bin (show create form)
function createNewJSON() {
  const createView = document.getElementById('createView');
  const viewView = document.getElementById('viewView');
  
  createView.style.display = 'block';
  viewView.style.display = 'none';
  currentJSON = null;
  
  clearForm();
  window.history.pushState({}, '', '/json.html');
}

// Clear form
function clearForm() {
  document.getElementById('jsonForm').reset();
  document.getElementById('charCount').textContent = '0';
  document.getElementById('charCount').style.color = 'var(--text-muted)';
  document.getElementById('validationMessage').style.display = 'none';
}

// Update code examples
function updateCodeExamples() {
  if (!currentJSON) return;
  
  const rawURL = `${window.location.origin}/api/v1/json/${currentJSON.id}/raw`;
  showExample(currentExample);
}

// Show code example
function showExample(type) {
  currentExample = type;
  
  if (!currentJSON) return;
  
  const rawURL = `${window.location.origin}/api/v1/json/${currentJSON.id}/raw`;
  const exampleCode = document.getElementById('exampleCode');
  
  // Update active tab
  document.querySelectorAll('.tab-btn').forEach(btn => btn.classList.remove('active'));
  event?.target?.classList.add('active');
  
  let code = '';
  
  switch (type) {
    case 'curl':
      code = `# GET request\ncurl ${rawURL}\n\n# POST/UPDATE request\ncurl -X PUT ${window.location.origin}/api/v1/json/${currentJSON.id} \\\n  -H "Content-Type: application/json" \\\n  -d '{"key": "value"}'`;
      exampleCode.className = 'language-bash';
      break;
    case 'fetch':
      code = `// GET request\nfetch('${rawURL}')\n  .then(res => res.json())\n  .then(data => console.log(data));\n\n// PUT request\nfetch('${window.location.origin}/api/v1/json/${currentJSON.id}', {\n  method: 'PUT',\n  headers: { 'Content-Type': 'application/json' },\n  body: JSON.stringify({ content: { key: 'value' } })\n});`;
      exampleCode.className = 'language-javascript';
      break;
    case 'axios':
      code = `// GET request\naxios.get('${rawURL}')\n  .then(res => console.log(res.data));\n\n// PUT request\naxios.put('${window.location.origin}/api/v1/json/${currentJSON.id}', {\n  content: { key: 'value' }\n});`;
      exampleCode.className = 'language-javascript';
      break;
  }
  
  exampleCode.textContent = code;
  hljs.highlightElement(exampleCode);
}

// Copy example code
function copyExample() {
  const exampleCode = document.getElementById('exampleCode');
  app.copyToClipboard(exampleCode.textContent);
}

// Close modal with Escape key
document.addEventListener('keydown', (e) => {
  if (e.key === 'Escape') {
    const modal = document.getElementById('updateModal');
    if (modal.classList.contains('show')) {
      closeUpdateModal();
    }
  }
});

// Handle browser back/forward
window.addEventListener('popstate', () => {
  const urlParams = new URLSearchParams(window.location.search);
  const id = urlParams.get('id');
  
  if (id) {
    loadJSON(id);
  } else {
    createNewJSON();
  }
});