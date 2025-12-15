// Navigation dropdown and mobile menu functionality

document.addEventListener('DOMContentLoaded', () => {
  // Dropdown functionality
  const dropdown = document.querySelector('.dropdown');
  const dropdownToggle = document.querySelector('.dropdown-toggle');
  
  if (dropdownToggle) {
    dropdownToggle.addEventListener('click', (e) => {
      e.stopPropagation();
      dropdown.classList.toggle('active');
    });
  }

  // Close dropdown when clicking outside
  document.addEventListener('click', (e) => {
    if (dropdown && !dropdown.contains(e.target)) {
      dropdown.classList.remove('active');
    }
  });

  // Mobile menu toggle
  const mobileMenuToggle = document.querySelector('.mobile-menu-toggle');
  const nav = document.querySelector('nav');
  
  if (mobileMenuToggle) {
    mobileMenuToggle.addEventListener('click', () => {
      nav.classList.toggle('mobile-active');
    });
  }

  // Close mobile menu when clicking a link
  const navLinks = document.querySelectorAll('nav a');
  navLinks.forEach(link => {
    link.addEventListener('click', () => {
      if (nav) {
        nav.classList.remove('mobile-active');
      }
    });
  });
});