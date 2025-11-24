// URL Shortener page logic

let currentShortURL = null;

// Initialize page
document.addEventListener('DOMContentLoaded', () => {
  // Check if viewing a short URL (e.g., /shorten.html?code=abc123)
  const urlParams = new URLSearchParams(window.location.search);
  const code = urlParams.get('code');
  
  if (code) {
    loadShortURL(code);
  }
});

// Handle short URL creation
async function handleCreateShortURL(event) {
  event.preventDefault();
  
  const submitBtn = document.getElementById('submitBtn');
  const form = document.getElementById('shortenForm');
  const formData = new FormData(form);
  
  // Prepare request data
  const longURL = formData.get('long_url');
  const customCode = formData.get('custom_code');
  const expiryDays = parseInt(formData.get('expiry_days'));
  
  // Validate
  if (!longURL || longURL.trim().length === 0) {
    app.showToast('URL cannot be empty', 'error');
    return;
  }
  
  // Show loading state
  submitBtn.disabled = true;
  submitBtn.innerHTML = '<span class="loading-spinner"></span> Shortening...';
  
  try {
    const shortURL = await app.api.post('/shorten', {
      long_url: longURL,
      custom_code: customCode || undefined,
      expiry_days: expiryDays === 0 ? null : expiryDays,
    });
    
    app.showToast('Short URL created successfully!');
    
    // Redirect to view
    window.history.pushState({}, '', `/shorten.html?code=${shortURL.short_code}`);
    currentShortURL = shortURL;
    showShortURLView(shortURL);
    
  } catch (error) {
    app.showToast(error.message || 'Failed to create short URL', 'error');
  } finally {
    submitBtn.disabled = false;
    submitBtn.innerHTML = '<span>✨</span> Shorten URL';
  }
}

// Load short URL by code
async function loadShortURL(code) {
  const createView = document.getElementById('createView');
  const viewView = document.getElementById('viewView');
  
  try {
    // Show loading
    createView.style.display = 'none';
    viewView.style.display = 'block';
    
    const shortURL = await app.api.get(`/shorten/${code}`);
    currentShortURL = shortURL;
    showShortURLView(shortURL);
    
  } catch (error) {
    viewView.style.display = 'none';
    createView.style.display = 'block';
    
    if (error.message.includes('not found')) {
      app.showToast('Short URL not found', 'error');
    } else if (error.message.includes('expired')) {
      app.showToast('This short URL has expired', 'error');
    } else {
      app.showToast('Failed to load short URL', 'error');
    }
  }
}

// Show short URL view
function showShortURLView(shortURL) {
  console.log('showShortURLView called with:', shortURL);
  
  const createView = document.getElementById('createView');
  const viewView = document.getElementById('viewView');
  
  // IMPORTANT: Show the view FIRST before accessing child elements
  createView.style.display = 'none';
  viewView.style.display = 'block';

  requestAnimationFrame(() => {
    requestAnimationFrame(() => {
        console.log(viewView.offsetHeight);
        // Now get the elements (after view is visible)
        const urlMeta = document.getElementById('urlMeta');
        const shortURLLink = document.getElementById('shortURLLink');
        const longURLText = document.getElementById('longURLText');
        
        console.log('Elements found:', {
            urlMeta: !!urlMeta,
            shortURLLink: !!shortURLLink,
            longURLText: !!longURLText
        });
        
        if (!shortURLLink || !longURLText || !urlMeta) {
            console.error('Required elements not found!');
            return;
        }

        // Build short URL
        const shortURLFull = `${window.location.origin}/s/${shortURL.short_code}`;
        console.log('Short URL:', shortURLFull);
        
        // Update display
        try {
            shortURLLink.href = shortURLFull;
            shortURLLink.textContent = shortURLFull;
            longURLText.textContent = shortURL.long_url;
            console.log('URLs set successfully');
        } catch (error) {
            console.error('Error setting URLs:', error);
            return;
        }
        
        // Update metadata
        let created = 'Unknown';
        try {
            const createdDate = new Date(shortURL.created_at);
            if (!isNaN(createdDate.getTime())) {
            created = createdDate.toLocaleString();
            }
        } catch (e) {
            console.error('Error parsing created_at:', e);
        }
        
        let expiryText = 'Never expires';
        if (shortURL.expires_at) {
            try {
            const expiresDate = new Date(shortURL.expires_at);
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
        
        urlMeta.innerHTML = `
        <div style="display:flex; flex-direction:row; align-items:center; gap:1.5rem">
            <div style="display:flex; flex-direction:row; align-items:center; gap:0.5rem;"><span class="material-symbols-outlined btn-sized">calendar_clock</span><span>Created: ${created}</span></div>
            <div style="display:flex; flex-direction:row; align-items:center; gap:0.5rem;"><span class="material-symbols-outlined btn-sized">free_cancellation</span><span>${expiryText}</span></div>
            <div style="display:flex; flex-direction:row; align-items:center; gap:0.5rem;"><span class="material-symbols-outlined btn-sized">bar_chart</span><span>Views: ${shortURL.views}</span></div>
            <div style="display:flex; flex-direction:row; align-items:center; gap:0.5rem;"><span class="material-symbols-outlined btn-sized">id_card</span><span>ID: ${shortURL.short_code || 'Unknown'}</span></div>
        </div>`;
    })
  })
}

// Copy short URL to clipboard
async function copyShortURL() {
  if (currentShortURL) {
    const shortURLFull = `${window.location.origin}/s/${currentShortURL.short_code}`;
    await app.copyToClipboard(shortURLFull);
    app.showToast('Short URL copied!');
  }
}

// Copy long URL to clipboard
async function copyLongURL() {
  if (currentShortURL) {
    await app.copyToClipboard(currentShortURL.long_url);
    app.showToast('Original URL copied!');
  }
}

// Delete short URL
async function deleteShortURL() {
  if (!currentShortURL) return;
  
  if (!confirm('Are you sure you want to delete this short URL? This action cannot be undone.')) {
    return;
  }
  
  try {
    await app.api.delete(`/shorten/${currentShortURL.id}`);
    app.showToast('Short URL deleted successfully');
    
    // Redirect to create view
    window.history.pushState({}, '', '/shorten.html');
    createNewShortURL();
    
  } catch (error) {
    console.error('Delete error:', error);
    app.showToast(error.message || 'Failed to delete short URL', 'error');
  }
}

// Create new short URL (show create form)
function createNewShortURL() {
  const createView = document.getElementById('createView');
  const viewView = document.getElementById('viewView');
  
  createView.style.display = 'block';
  viewView.style.display = 'none';
  currentShortURL = null;
  
  // Clear form
  clearForm();
  
  // Update URL
  window.history.pushState({}, '', '/shorten.html');
}

// Clear form
function clearForm() {
  document.getElementById('shortenForm').reset();
}

// Show QR Code Modal
function showQRCodeModal() {
  if (!currentShortURL) return;
  
  const modal = document.getElementById('qrModal');
  modal.classList.add('show');
  
  // Prevent body scroll
  document.body.style.overflow = 'hidden';
  
  // Generate QR code
  generateQRCode();
}

// Close QR Code Modal
function closeQRModal(event) {
  const modal = document.getElementById('qrModal');
  modal.classList.remove('show');
  
  // Restore body scroll
  document.body.style.overflow = '';
}

// Generate QR Code
function generateQRCode() {
  if (!currentShortURL) return;
  
  const qrElement = document.getElementById('qrcode');
  const qrURLElement = document.getElementById('qrCodeURL');
  const shortURLFull = `${window.location.origin}/s/${currentShortURL.short_code}`;
  
  // Clear previous QR code
  qrElement.innerHTML = '';
  
  // Generate new QR code
  qrCodeInstance = new QRCode(qrElement, shortURLFull);
  
  // Show URL
  qrURLElement.textContent = shortURLFull;
}

// Download QR Code as PNG
function downloadQRCode() {
  if (!currentShortURL) return;
  
  const canvas = document.querySelector('#qrcode canvas');
  if (!canvas) {
    app.showToast('QR code not generated yet', 'error');
    return;
  }
  
  // Convert canvas to blob and download
  canvas.toBlob((blob) => {
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `short-url-${currentShortURL.short_code}-qrcode.png`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
    
    app.showToast('QR code downloaded!');
  });
}

// Close modal with Escape key
document.addEventListener('keydown', (e) => {
  if (e.key === 'Escape') {
    const modal = document.getElementById('qrModal');
    if (modal.classList.contains('show')) {
      closeQRModal();
    }
  }
});

// Handle browser back/forward
window.addEventListener('popstate', () => {
  const urlParams = new URLSearchParams(window.location.search);
  const code = urlParams.get('code');
  
  if (code) {
    loadShortURL(code);
  } else {
    createNewShortURL();
  }
});