// IP page logic

let currentIP = null;

// Initialize page
document.addEventListener('DOMContentLoaded', () => {
  fetchIP();
});

// Fetch IP from API
async function fetchIP() {
  const ipDisplay = document.getElementById('ipDisplay').querySelector('.ip-address');
  const copyBtn = document.getElementById('copyBtn');
  const ipInfo = document.getElementById('ipInfo');
  
  try {
    ipDisplay.textContent = 'Loading...';
    ipDisplay.className = 'ip-address loading';
    copyBtn.disabled = true;
    
    const data = await app.api.get('/ip');
    
    if (data && data.ip) {
      currentIP = data.ip;
      ipDisplay.textContent = data.ip;
      ipDisplay.className = 'ip-address';
      copyBtn.disabled = false;
      
      // Show additional info
      updateIPInfo(data.ip);
      ipInfo.style.display = 'grid';
    } else {
      throw new Error('Invalid response');
    }
  } catch (error) {
    ipDisplay.textContent = 'Failed to load';
    ipDisplay.className = 'ip-address';
    ipDisplay.style.color = 'var(--error)';
    app.showToast('Failed to fetch IP address', 'error');
    console.error('Error fetching IP:', error);
  }
}

// Update IP information display
function updateIPInfo(ip) {
  const ipVersion = document.getElementById('ipVersion');
  const ipType = document.getElementById('ipType');
  const lastChecked = document.getElementById('lastChecked');
  
  // Determine IP version
  if (ip.includes(':')) {
    ipVersion.textContent = 'IPv6';
    ipType.textContent = 'Public';
  } else if (ip.includes('.')) {
    ipVersion.textContent = 'IPv4';
    
    // Check if private IP
    const parts = ip.split('.').map(Number);
    if (
      parts[0] === 10 ||
      (parts[0] === 172 && parts[1] >= 16 && parts[1] <= 31) ||
      (parts[0] === 192 && parts[1] === 168)
    ) {
      ipType.textContent = 'Private';
    } else {
      ipType.textContent = 'Public';
    }
  } else {
    ipVersion.textContent = 'Unknown';
    ipType.textContent = 'Unknown';
  }
  
  // Update last checked time
  const now = new Date();
  lastChecked.textContent = now.toLocaleTimeString();
}

// Copy IP to clipboard
async function copyIP() {
  if (currentIP) {
    await app.copyToClipboard(currentIP);
  }
}

// Refresh IP
function refreshIP() {
  fetchIP();
}