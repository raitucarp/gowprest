(function () {
  const STORAGE_KEY = 'gowprest-theme';

  function getInitialTheme() {
    const saved = localStorage.getItem(STORAGE_KEY);
    if (saved) return saved;
    return window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches
      ? 'dark'
      : 'light';
  }

  function applyTheme(theme) {
    document.documentElement.setAttribute('data-theme', theme);
    localStorage.setItem(STORAGE_KEY, theme);
    
    // Update theme-color meta tag for mobile browsers
    const metaTheme = document.querySelector('meta[name="theme-color"]');
    if (metaTheme) {
      metaTheme.setAttribute('content', theme === 'dark' ? '#090f19' : '#0073aa');
    }

    const toggleBtn = document.getElementById('theme-toggle');
    if (toggleBtn) {
      const icon = toggleBtn.querySelector('.theme-icon');
      const text = toggleBtn.querySelector('.theme-text');
      if (icon) icon.textContent = theme === 'dark' ? '☀️' : '🌙';
      if (text) text.textContent = theme === 'dark' ? 'Light' : 'Dark';
    }
  }

  // Apply immediately to prevent flash of wrong theme
  applyTheme(getInitialTheme());

  document.addEventListener('DOMContentLoaded', () => {
    applyTheme(getInitialTheme());

    // 1. Theme toggle
    const toggleBtn = document.getElementById('theme-toggle');
    if (toggleBtn) {
      toggleBtn.addEventListener('click', () => {
        const current = document.documentElement.getAttribute('data-theme') || 'light';
        const next = current === 'dark' ? 'light' : 'dark';
        applyTheme(next);
      });
    }

    // 2. Mobile Drawer Navigation
    const mobileBtn = document.getElementById('mobile-menu-btn');
    const sidebar = document.getElementById('site-sidebar');
    const backdrop = document.getElementById('sidebar-backdrop');
    const closeBtn = document.getElementById('sidebar-close-btn');

    function openDrawer() {
      if (!sidebar) return;
      sidebar.classList.add('open');
      if (backdrop) backdrop.classList.add('active');
      document.body.classList.add('drawer-open');
      if (mobileBtn) mobileBtn.setAttribute('aria-expanded', 'true');
    }

    function closeDrawer() {
      if (!sidebar) return;
      sidebar.classList.remove('open');
      if (backdrop) backdrop.classList.remove('active');
      document.body.classList.remove('drawer-open');
      if (mobileBtn) mobileBtn.setAttribute('aria-expanded', 'false');
    }

    if (mobileBtn) {
      mobileBtn.addEventListener('click', (e) => {
        e.stopPropagation();
        if (sidebar && sidebar.classList.contains('open')) {
          closeDrawer();
        } else {
          openDrawer();
        }
      });
    }

    if (closeBtn) {
      closeBtn.addEventListener('click', (e) => {
        e.stopPropagation();
        closeDrawer();
      });
    }

    if (backdrop) {
      backdrop.addEventListener('click', () => {
        closeDrawer();
      });
    }

    // Close on Escape key
    document.addEventListener('keydown', (e) => {
      if (e.key === 'Escape') {
        closeDrawer();
      }
    });

    // Auto-close drawer on link click in mobile view
    if (sidebar) {
      sidebar.querySelectorAll('a').forEach((link) => {
        link.addEventListener('click', () => {
          if (window.innerWidth <= 900) {
            closeDrawer();
          }
        });
      });
    }

    // 3. Enhance Code Blocks (Copy Button & Language Badge)
    enhanceCodeBlocks();

    // 4. Responsive Table Wrapper
    wrapTables();
  });

  function enhanceCodeBlocks() {
    // Select Hugo's .highlight divs or raw <pre> blocks in the content
    const codeContainers = document.querySelectorAll('.content-article .highlight, .content-article pre:not(.highlight pre)');

    codeContainers.forEach((container) => {
      // Don't re-enhance if already processed
      if (container.closest('.code-block-wrapper')) return;

      const pre = container.tagName === 'PRE' ? container : container.querySelector('pre');
      if (!pre) return;

      const code = pre.querySelector('code');

      // Detect language
      let lang = 'CODE';
      const classList = [
        ...(container.className ? container.className.split(' ') : []),
        ...(code && code.className ? code.className.split(' ') : [])
      ];

      for (const cls of classList) {
        if (cls.startsWith('language-')) {
          lang = cls.replace('language-', '').toUpperCase();
          break;
        }
        if (['go', 'bash', 'sh', 'json', 'yaml', 'toml', 'html', 'http'].includes(cls.toLowerCase())) {
          lang = cls.toUpperCase();
          break;
        }
      }

      // Format known labels
      if (lang === 'SH') lang = 'BASH';

      // Create code-block-wrapper
      const wrapper = document.createElement('div');
      wrapper.className = 'code-block-wrapper';

      // Create header bar
      const header = document.createElement('div');
      header.className = 'code-header';

      const langSpan = document.createElement('span');
      langSpan.className = 'code-lang';
      langSpan.textContent = lang;

      const copyBtn = document.createElement('button');
      copyBtn.className = 'code-copy-btn';
      copyBtn.setAttribute('type', 'button');
      copyBtn.setAttribute('aria-label', 'Copy code to clipboard');
      copyBtn.innerHTML = `
        <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
          <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
        </svg>
        <span>Copy</span>
      `;

      copyBtn.addEventListener('click', async () => {
        const textToCopy = (code ? code.innerText : pre.innerText).replace(/\n$/, '');
        try {
          if (navigator.clipboard && navigator.clipboard.writeText) {
            await navigator.clipboard.writeText(textToCopy);
          } else {
            // Fallback for older mobile webviews
            const textarea = document.createElement('textarea');
            textarea.value = textToCopy;
            textarea.style.position = 'fixed';
            textarea.style.opacity = '0';
            document.body.appendChild(textarea);
            textarea.select();
            document.execCommand('copy');
            document.body.removeChild(textarea);
          }

          copyBtn.classList.add('copied');
          copyBtn.innerHTML = `
            <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
              <polyline points="20 6 9 17 4 12"></polyline>
            </svg>
            <span>Copied!</span>
          `;

          setTimeout(() => {
            copyBtn.classList.remove('copied');
            copyBtn.innerHTML = `
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
                <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
              </svg>
              <span>Copy</span>
            `;
          }, 2000);
        } catch (err) {
          console.error('Failed to copy code: ', err);
        }
      });

      header.appendChild(langSpan);
      header.appendChild(copyBtn);

      // Insert wrapper in place of container
      container.parentNode.insertBefore(wrapper, container);
      wrapper.appendChild(header);
      wrapper.appendChild(container);
    });
  }

  function wrapTables() {
    const tables = document.querySelectorAll('.content-article table');
    tables.forEach((table) => {
      if (table.parentElement.classList.contains('table-responsive')) return;
      const wrapper = document.createElement('div');
      wrapper.className = 'table-responsive';
      table.parentNode.insertBefore(wrapper, table);
      wrapper.appendChild(table);
    });
  }
})();
